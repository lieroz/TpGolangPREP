package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	// "log"
	"github.com/json-iterator/go"
	"bufio"
)

type User struct {
	Browsers []string `json:"browsers"`
	Company  string   `json:"-"`
	Country  string   `json:"-"`
	Email    string   `json:"email"`
	Job      string   `json:"-"`
	Name     string   `json:"name"`
	Phone    string   `json:"-"`
}

func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	seenBrowsers := []string{}
	uniqueBrowsers := 0
	foundUsers := ""
	i := 0

	for scanner.Scan() {
		user := new(User)
		// fmt.Printf("%v %v\n", err, line)
		iter := jsoniter.ConfigFastest.BorrowIterator(scanner.Bytes())
		iter.ReadVal(user)
		jsoniter.ConfigFastest.ReturnIterator(iter)
		if err != nil {
			panic(err)
		}

		isAndroid := false
		isMSIE := false

		for _, browser := range user.Browsers {
			check := false
			if strings.Contains(browser, "Android") {
				isAndroid = true
				check = true
			} else if strings.Contains(browser, "MSIE") {
				isMSIE = true
				check = true
			}
			if check {
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
						break
					}
				}
				if notSeenBefore {
					// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
					seenBrowsers = append(seenBrowsers, browser)
					uniqueBrowsers++
				}
			}
		}

		if !(isAndroid && isMSIE) {
			i++
			continue
		}

		// log.Println("Android and MSIE user:", user["name"], user["email"])
		email := strings.Replace(user.Email, "@", " [at] ", -1)
		foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, user.Name, email)
		i++
	}

	fmt.Fprintln(out, "found users:\n"+foundUsers)
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}
