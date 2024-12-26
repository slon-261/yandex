package main

import (
	"fmt"
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

// OsExitAnalyzer анализатор для проверки использования os.Exit в пакете main
var OsExitAnalyzer = &analysis.Analyzer{
	Name: "osexitanalyzer",
	Doc:  "Check os.Exit in main package",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.File:
				if x.Name.Name != "main" {
					return false
				}
			case *ast.SelectorExpr:
				if x.Sel.Name == "Exit" {
					fmt.Print(x)
					pass.Reportf(x.Pos(), "os.Exit in main package")
				}
			}
			return true
		})
	}
	return nil, nil
}
