package main

import (
	"go/token"
	"go/parser"
	"os"
	"log"
	"go/ast"
	"fmt"
	"reflect"
	"strings"
	"encoding/json"
	"strconv"
	"text/template"
	"bytes"
)

func generateImports(out *os.File, node *ast.File) {
	fmt.Fprintln(out, fmt.Sprintf(`package %s`, node.Name.Name))
	fmt.Fprintln(out)
	fmt.Fprintln(out, `import "net/http"`)
	fmt.Fprintln(out, `import "net/url"`)
	fmt.Fprintln(out, `import "strconv"`)
	fmt.Fprintln(out, `import "strings"`)
	fmt.Fprintln(out, `import "encoding/json"`)
	fmt.Fprintln(out, `import "io/ioutil"`)
	fmt.Fprintln(out, `import "io"`)
	fmt.Fprintln(out, `import "fmt"`)
}

var containsTemplate = template.Must(template.New("containsTemplate").Parse(`
func contains(list []string, item string) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}
`))

func generateListHelper(out *os.File, node *ast.File) {
	containsTemplate.Execute(out, nil)
}

var parseCrutchyBodyTemplate = template.Must(template.New("containsTemplate").Parse(`
func parseCrutchyBody(body io.ReadCloser) url.Values {
	b, _ := ioutil.ReadAll(body)
	defer body.Close()
	query := string(b)
	v, _ := url.ParseQuery(query)
	return v
}
`))

func generateParseQueryHelper(out *os.File, node *ast.File) {
	parseCrutchyBodyTemplate.Execute(out, nil)
	fmt.Fprintln(out)
}

func addFieldValidator(out *os.File, validatorArgs []string, field *ast.Field) {
	var (
		required             bool
		enum, regular, param string
		min, max             = -1, -1
	)
	fieldType, ok := field.Type.(*ast.Ident)
	if !ok {
		panic("unsupported interface")
	}
	for _, arg := range validatorArgs {
		if arg == "required" {
			required = true
		} else {
			if !strings.Contains(arg, "=") {
				panic("unsupported validator option")
			}
			value := arg[strings.Index(arg, "=")+1:]
			switch arg[:strings.Index(arg, "=")] {
			case "enum":
				enum = value
			case "default":
				regular = value
			case "min":
				min, _ = strconv.Atoi(value)
			case "max":
				max, _ = strconv.Atoi(value)
			case "paramname":
				param = value
			default:
				panic("unsupported validator option")
			}
		}
	}
	formatStr := "\tp." + field.Names[0].Name
	isInt := false
	switch fieldType.Name {
	case "int":
		isInt = true
		formatStr += `, err = strconv.Atoi(args.Get("%s"))`
	case "string":
		formatStr += ` = args.Get("%s")`
	default:
		panic("unsupported type")
	}
	if len(param) != 0 {
		fmt.Fprintln(out, fmt.Sprintf(formatStr, param))
	} else {
		fmt.Fprintln(out, fmt.Sprintf(formatStr, strings.ToLower(field.Names[0].Name)))
	}
	if isInt {
		fmt.Fprintln(out, "\tif err != nil {")
		fmt.Fprintln(out, "\t\t"+`return fmt.Errorf("`+strings.ToLower(field.Names[0].Name)+` must be int")`)
		fmt.Fprintln(out, "\t}")
	}
	if len(regular) != 0 {
		if isInt {
			fmt.Fprintln(out, "\tif p."+field.Names[0].Name+" == 0 {")
			fmt.Fprintln(out, "\tp."+field.Names[0].Name+" = "+regular)
		} else {
			fmt.Fprintln(out, "\tif len(p."+field.Names[0].Name+") == 0 {")
			fmt.Fprintln(out, "\t\tp."+field.Names[0].Name+` = "`+regular+`"`)
		}
		fmt.Fprintln(out, "\t}")
		fmt.Fprintln(out, "\t"+`if !contains(strings.Split("`+enum+`", "|"), p.`+field.Names[0].Name+`) {`)
		fmt.Fprint(out, "\t\t"+`return fmt.Errorf("`+strings.ToLower(field.Names[0].Name)+` must be one of [`)
		var e string
		for _, s := range strings.Split(enum, "|") {
			e += s + ", "
		}
		e = e[:len(e)-2]
		fmt.Fprintln(out, e+`]")`)
		fmt.Fprintln(out, "\t}")
	}
	if required {
		if isInt {
			fmt.Fprintln(out, "\tif p."+field.Names[0].Name+" == 0 {")
		} else {
			fmt.Fprintln(out, "\tif len(p."+field.Names[0].Name+") == 0 {")
		}
		fmt.Fprintln(out, "\t\t"+`return fmt.Errorf("`+strings.ToLower(field.Names[0].Name)+` must me not empty")`)
		fmt.Fprintln(out, "\t}")
	}
	if isInt {
		if min != -1 {
			fmt.Fprintln(out, "\tif p."+field.Names[0].Name+" <= "+strconv.Itoa(min)+" {")
			fmt.Fprintln(out, "\t\t"+`return fmt.Errorf("`+strings.ToLower(field.Names[0].Name)+` must be >= `+strconv.Itoa(min)+`")`)
			fmt.Fprintln(out, "\t}")
		}
		if max != -1 {
			fmt.Fprintln(out, "\tif p."+field.Names[0].Name+" >= "+strconv.Itoa(max)+" {")
			fmt.Fprintln(out, "\t\t"+`return fmt.Errorf("`+strings.ToLower(field.Names[0].Name)+` must be <= `+strconv.Itoa(max)+`")`)
			fmt.Fprintln(out, "\t}")
		}
	} else {
		fmt.Fprintln(out, "\tif len(p."+field.Names[0].Name+") <= "+strconv.Itoa(min)+" {")
		fmt.Fprintln(out, "\t\t"+`return fmt.Errorf("`+strings.ToLower(field.Names[0].Name)+` len must be >= `+strconv.Itoa(min)+`")`)
		fmt.Fprintln(out, "\t}")
	}
}

func generateValidators(out *os.File, node *ast.File) {
	for _, f := range node.Decls {
		g, ok := f.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range g.Specs {
			currType, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			currStruct, ok := currType.Type.(*ast.StructType)
			if !ok {
				continue
			}
			needValidator := false
			count := 0
			for _, field := range currStruct.Fields.List {
				if field.Tag != nil {
					tag := reflect.StructTag(field.Tag.Value[1: len(field.Tag.Value)-1])
					var text string
					text, needValidator = tag.Lookup("apivalidator")
					if needValidator && count == 0 {
						count++
						fmt.Fprintln(out, "func (p *"+currType.Name.Name+") validateAndFill"+currType.Name.Name+"(args url.Values) (err error) {")
					}
					if needValidator {
						addFieldValidator(out, strings.Split(text, ","), field)
					}
				}
			}
			if needValidator {
				fmt.Fprintln(out, "\treturn nil")
				fmt.Fprintln(out, "}")
				fmt.Fprintln(out)
			}
		}
	}
}

type structTpl struct {
	StructName string
	FuncName string
	Method string
	Auth string
	Cases string
}

var httpWrapperTemplate = template.Must(template.New("httpWrapperHeaderTemplate").Parse(`
func (srv *{{.StructName}}) handler{{.FuncName}}(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})
	resp["error"] = ""
	{{.Method}}
	{{.Auth}}
	var v url.Values
	switch r.Method {
	case "POST":
		v = parseCrutchyBody(r.Body)
	default:
		v = r.URL.Query()
	}
	var params {{.FuncName}}Params
	if err := params.validateAndFill{{.FuncName}}Params(v); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp["error"] = err.Error()
		body, _ := json.Marshal(resp)
		w.Write(body)
		return
	}
	user, err := srv.{{.FuncName}}(r.Context(), params)
	if err != nil {
		switch err.(type) {
		case ApiError:
			w.WriteHeader(err.(ApiError).HTTPStatus)
			resp["error"] = err.Error()
		default:
			w.WriteHeader(http.StatusInternalServerError)
			resp["error"] = "bad user"
		}
		body, _ := json.Marshal(resp)
		w.Write(body)
		return
	}
	resp["response"] = user
	body, _ := json.Marshal(resp)
	w.Write(body)
}
`))

type methodTpl struct {
	Method string
}

var methodTemplate = template.Must(template.New("methodTemplate").Parse(`
	if r.Method != "{{.Method}}" {
		w.WriteHeader(http.StatusNotAcceptable)
		resp["error"] = "bad method"
		body, _ := json.Marshal(resp)
		w.Write(body)
		return
	}`))

var authTemplate = template.Must(template.New("authTemplate").Parse(`
	if r.Header.Get("X-Auth") != "100500" {
		w.WriteHeader(http.StatusForbidden)
		resp["error"] = "unauthorized"
		body, _ := json.Marshal(resp)
		w.Write(body)
		return
	}`))

func generateHttpWrappers(out *os.File, file *ast.File) {
	for _, s := range structures {
		for _, node := range file.Decls {
			currFunc, ok := node.(*ast.FuncDecl)
			if !ok {
				continue
			}
			if currFunc.Doc != nil {
				methodSignature := new(MethodSignature)
				for _, doc := range currFunc.Doc.List {
					if strings.Contains(doc.Text, "apigen:api") {
						start := strings.Index(doc.Text, "{")
						end := strings.Index(doc.Text, "}") + 1
						json.Unmarshal([]byte(doc.Text[start:end]), methodSignature)

						tpl := structTpl{StructName: s, FuncName: currFunc.Name.Name}
						if len(methodSignature.Method) != 0 {
							output := new(bytes.Buffer)
							methodTemplate.Execute(output, methodTpl{methodSignature.Method})
							tpl.Method = output.String()
						}

						if methodSignature.Auth {
							output := new(bytes.Buffer)
							authTemplate.Execute(output, nil)
							tpl.Auth = output.String()
						}

						httpWrapperTemplate.Execute(out, tpl)
					}
				}
			}
		}
	}
}

type MethodSignature struct {
	Url    string `json:"url"`
	Auth   bool   `json:"auth"`
	Method string `json:"method"`
}

var serveHttpTemplate = template.Must(template.New("serveHttpTemplate").Parse(`
func (srv *{{.StructName}}) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})
	switch r.URL.Path {
	{{.Cases}}
	default:
		w.WriteHeader(http.StatusNotFound)
		resp["error"] = "unknown method"
		body, _ := json.Marshal(resp)
		w.Write(body)
	}
}
`))

type caseTpl struct {
	Path string
	Handler string
}

var caseTemplate = template.Must(template.New("caseTemplate").Parse(`
	case "{{.Path}}":
		srv.handler{{.Handler}}(w, r)
`))

func generateServeHTTP(out *os.File, file *ast.File) {
	for _, s := range structures {
		output := new(bytes.Buffer)
		for _, node := range file.Decls {
			currFunc, ok := node.(*ast.FuncDecl)
			if !ok {
				continue
			}
			if currFunc.Doc != nil {
				methodSignature := new(MethodSignature)
				for _, doc := range currFunc.Doc.List {
					if strings.Contains(doc.Text, "apigen:api") {
						start := strings.Index(doc.Text, "{")
						end := strings.Index(doc.Text, "}") + 1
						json.Unmarshal([]byte(doc.Text[start:end]), methodSignature)
						caseTemplate.Execute(output, caseTpl{methodSignature.Url, currFunc.Name.Name})
					}
				}
			}
		}
		serveHttpTemplate.Execute(out, structTpl{StructName: s, Cases: output.String()})
	}
}

var structures []string

func loopFunc(currFunc *ast.FuncDecl) {
	for _, doc := range currFunc.Doc.List {
		if strings.Contains(doc.Text, "apigen:api") {
		LOOP:
			for _, i := range currFunc.Recv.List {
				structName := i.Type.(*ast.StarExpr).X.(*ast.Ident)
				for _, i := range structures {
					if i == structName.Name {
						break LOOP
					}
				}
				structures = append(structures, structName.Name)
			}
		}
	}
}

func getNeededStructs(file *ast.File) {
	for _, node := range file.Decls {
		currFunc, ok := node.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if currFunc.Doc != nil {
			loopFunc(currFunc)
		}
	}
}

// это кодогенератор
func main() {
	fSet := token.NewFileSet()
	file, err := parser.ParseFile(fSet, os.Args[1], nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	out, _ := os.Create(os.Args[2])

	getNeededStructs(file)

	generateImports(out, file)
	generateListHelper(out, file)
	generateParseQueryHelper(out, file)
	generateValidators(out, file)
	generateHttpWrappers(out, file)
	generateServeHTTP(out, file)
}
