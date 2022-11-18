package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	mailHTMLDir := os.Args[1]
	templateDir := os.Args[2]
	htmlFiles, err := filepath.Glob(mailHTMLDir + "/*.html")
	if err != nil {
		panic(err)
	}

	for _, htmlFile := range htmlFiles {
		htmlFileName := filepath.Base(htmlFile)
		htmlFileName = htmlFileName[:len(htmlFileName)-len(filepath.Ext(htmlFileName))]

		htmlFileContent, err := os.ReadFile(htmlFile)
		if err != nil {
			panic(err)
		}
		fmt.Println(templateDir + "/" + htmlFileName + ".tmpl")
		templateContent, err := os.ReadFile(templateDir + "/" + htmlFileName + ".tmpl")
		if err != nil {
			panic(err)
		}

		newTemplateContent := "{{define \"htmlBody\"}}" + string(htmlFileContent) + "{{end}}"
		textContent := string(templateContent[:strings.Index(string(templateContent), "{{define \"htmlBody\"}}")])
		newContent := textContent + newTemplateContent
		err = os.WriteFile(templateDir+"\\"+htmlFileName+".tmpl", []byte(newContent), 0644)
		if err != nil {
			panic(err)
		}
	}

}
