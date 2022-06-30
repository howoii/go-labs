package main

import (
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"path/filepath"
)

func main() {
	var filename string
	flag.StringVar(&filename, "input", "src.go", "input file")
	flag.Parse()

	fset := token.NewFileSet()
	path, err := filepath.Abs(filepath.Join("src", filename))
	if err != nil {
		log.Fatal(err)
	}
	f, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}
	ast.Print(fset, f)
}
