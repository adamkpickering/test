package main

import (
	"fmt"
	"os"
	"text/template"
)

var templates *template.Template

func init() {
	templates = template.Must(template.ParseGlob("templates/*.gotmpl"))
}

func main() {
	for _, templateName := range []string{"template1.gotmpl", "template2.gotmpl"} {
		err := templates.ExecuteTemplate(os.Stdout, templateName, nil)
		if err != nil {
			fmt.Printf("failed to fill %s: %s\n", templateName, err)
			os.Exit(1)
		}
	}
}
