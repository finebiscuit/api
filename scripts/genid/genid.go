package main

import (
	"flag"
	"html/template"
	"log"
	"os"
)

func main() {
	pkgFlag := flag.String("pkg", "", "")
	flag.Parse()

	f, err := os.Create("id_gen.go")
	if err != nil {
		log.Fatal(err)
	}

	if pkgFlag == nil {
		log.Fatal("no pkg flag")
	}

	tpl.Execute(f, struct {
		PkgName string
	}{
		PkgName: *pkgFlag,
	})
}

var tpl = template.Must(template.New("").Parse(`// Code generated by go generate; DO NOT EDIT.
package {{ .PkgName }}

type ID string

const NoID ID = ""

func (id ID) String() string {
	return string(id)
}

func (id ID) Valid() bool {
	return id != NoID
}

func ParseID(s string) ID {
	return ID(s)
}
`))
