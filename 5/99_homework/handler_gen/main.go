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
	fmt.Fprintln(out)
}

func generateListHelper(out *os.File, node *ast.File) {
	fmt.Fprintln(out, "func contains(list []string, item string) bool {")
	fmt.Fprintln(out, "\tfor _, i := range list {")
	fmt.Fprintln(out, "\t\tif i == item {")
	fmt.Fprintln(out, "\t\t\treturn true")
	fmt.Fprintln(out, "\t\t}")
	fmt.Fprintln(out, "\t}")
	fmt.Fprintln(out, "\treturn false")
	fmt.Fprintln(out, "}")
	fmt.Fprintln(out)
}

func generateParseQueryHelper(out *os.File, node *ast.File) {
	fmt.Fprintln(out, "func parseCrutchyBody(body io.ReadCloser) url.Values {")
	fmt.Fprintln(out, "\tb, _ := ioutil.ReadAll(body)")
	fmt.Fprintln(out, "\tdefer body.Close()")
	fmt.Fprintln(out, "\tquery := string(b)")
	fmt.Fprintln(out, "\tv, _ := url.ParseQuery(query)")
	fmt.Fprintln(out, "\treturn v")
	fmt.Fprintln(out, "}")
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

func generateHttpWrappers(out *os.File, file *ast.File) {
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
					fmt.Fprintln(out, "func (srv *MyApi) handler"+currFunc.Name.Name+"(w http.ResponseWriter, r *http.Request) {")
					fmt.Fprintln(out, "\tresp := make(map[string]interface{})")
					fmt.Fprintln(out, "\t"+`resp["error"] = ""`)

					if len(methodSignature.Method) != 0 {
						fmt.Fprintln(out, "\t"+`if r.Method != "`+methodSignature.Method+`" {`)
						fmt.Fprintln(out, "\t\tw.WriteHeader(http.StatusNotAcceptable)")
						fmt.Fprintln(out, "\t\t"+`resp["error"] = "bad method"`)
						fmt.Fprintln(out, "\t\tbody, _ := json.Marshal(resp)")
						fmt.Fprintln(out, "\t\tw.Write(body)")
						fmt.Fprintln(out, "\t\treturn")
						fmt.Fprintln(out, "\t}")
					}

					if methodSignature.Auth {
						fmt.Fprintln(out, "\t"+`if r.Header.Get("X-Auth") != "100500" {`)
						fmt.Fprintln(out, "\t\tw.WriteHeader(http.StatusForbidden)")
						fmt.Fprintln(out, "\t\t"+`resp["error"] = "unauthorized"`)
						fmt.Fprintln(out, "\t\tbody, _ := json.Marshal(resp)")
						fmt.Fprintln(out, "\t\tw.Write(body)")
						fmt.Fprintln(out, "\t\treturn")
						fmt.Fprintln(out, "\t}")
					}

					fmt.Fprintln(out, "\tvar v url.Values")
					fmt.Fprintln(out, "\tswitch r.Method {")
					fmt.Fprintln(out, "\t"+`case "POST":`)
					fmt.Fprintln(out, "\t\tv = parseCrutchyBody(r.Body)")
					fmt.Fprintln(out, "\tdefault:")
					fmt.Fprintln(out, "\t\tv = r.URL.Query()")
					fmt.Fprintln(out, "\t}")

					fmt.Fprintln(out, "\tvar params "+currFunc.Name.Name+"Params")
					fmt.Fprintln(out, "\tif err := params.validateAndFill"+currFunc.Name.Name+"Params(v); err != nil {")
					fmt.Fprintln(out, "\t\tw.WriteHeader(http.StatusBadRequest)")
					fmt.Fprintln(out, "\t\t"+`resp["error"] = err.Error()`)
					fmt.Fprintln(out, "\t\tbody, _ := json.Marshal(resp)")
					fmt.Fprintln(out, "\t\tw.Write(body)")
					fmt.Fprintln(out, "\t\treturn")
					fmt.Fprintln(out, "\t}")

					fmt.Fprintln(out, "\tuser, err := srv."+currFunc.Name.Name+"(r.Context(), params)")

					fmt.Fprintln(out, "\tif err != nil {")
					fmt.Fprintln(out, "\t\tswitch err.(type) {")
					fmt.Fprintln(out, "\t\tcase ApiError:")
					fmt.Fprintln(out, "\t\t\tw.WriteHeader(err.(ApiError).HTTPStatus)")
					fmt.Fprintln(out, "\t\t\t"+`resp["error"] = err.Error()`)
					fmt.Fprintln(out, "\t\tdefault:")
					fmt.Fprintln(out, "\t\t\tw.WriteHeader(http.StatusInternalServerError)")
					fmt.Fprintln(out, "\t\t\t"+`resp["error"] = "bad user"`)
					fmt.Fprintln(out, "\t\t}")
					fmt.Fprintln(out, "\t\tbody, _ := json.Marshal(resp)")
					fmt.Fprintln(out, "\t\tw.Write(body)")
					fmt.Fprintln(out, "\t\treturn")
					fmt.Fprintln(out, "\t}")

					fmt.Fprintln(out, "\t"+`resp["response"] = user`)
					fmt.Fprintln(out, "\tbody, _ := json.Marshal(resp)")
					fmt.Fprintln(out, "\tw.Write(body)")
					fmt.Fprintln(out, "}")
					fmt.Fprintln(out)
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

func generateServeHTTP(out *os.File, file *ast.File) {
	fmt.Fprintln(out, "func (srv *MyApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {")
	fmt.Fprintln(out, "resp := make(map[string]interface{})")
	fmt.Fprintln(out, "\tswitch r.URL.Path {")
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
					fmt.Fprintln(out, "\t"+`case "`+methodSignature.Url+`":`)
					fmt.Fprintln(out, "\t\tsrv.handler"+currFunc.Name.Name+"(w, r)")
				}
			}
		}
	}
	fmt.Fprintln(out, "\tdefault:")
	fmt.Fprintln(out, "\t\tw.WriteHeader(http.StatusNotFound)")
	fmt.Fprintln(out, "\t\t"+`resp["error"] = "unknown method"`)
	fmt.Fprintln(out, "\t\tbody, _ := json.Marshal(resp)")
	fmt.Fprintln(out, "\t\tw.Write(body)")
	fmt.Fprintln(out, "\t}")
	fmt.Fprintln(out, "}")
}

// это кодогенератор
func main() {
	fSet := token.NewFileSet()
	file, err := parser.ParseFile(fSet, os.Args[1], nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	out, _ := os.Create(os.Args[2])

	generateImports(out, file)
	generateListHelper(out, file)
	generateParseQueryHelper(out, file)
	generateValidators(out, file)
	generateHttpWrappers(out, file)
	generateServeHTTP(out, file)
}
