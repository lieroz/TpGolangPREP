package main

import (
	"io"
	"os"
	"strings"
	// "log"
	"github.com/json-iterator/go"
	"bufio"
	"strconv"
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
	uniqueBrowsers ,i  := 0, 0
	user := new(User)
	out.Write([]byte("found users:\n"))

	for scanner.Scan() {
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
					seenBrowsers = append(seenBrowsers, browser)
					uniqueBrowsers++
				}
			}
		}

		if !(isAndroid && isMSIE) {
			i++
			continue
		}

		email := strings.Replace(user.Email, "@", " [at] ", -1)
		out.Write([]byte("[" + strconv.Itoa(i) + "] " + user.Name + " <" + email + ">\n"))
		i++
	}

	out.Write([]byte("\nTotal unique browsers " + strconv.Itoa(len(seenBrowsers)) + "\n"))
}
