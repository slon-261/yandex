// cd ./cmd/staticlint/
// go build main.go os_exit_analyzer.go
// cd ../../
// ./cmd/staticlint/main.exe ./...

package main

import (
	"encoding/json"
	"github.com/kisielk/errcheck/errcheck"
	"github.com/mdempsky/maligned/passes/maligned"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
)

// Config — имя файла конфигурации.
const Config = `config.json`

// ConfigData описывает структуру файла конфигурации.
type ConfigData struct {
	Staticcheck []string
	Stylecheck  []string
}

func main() {
	appfile, err := os.Executable()
	if err != nil {
		panic(err)
	}
	data, err := os.ReadFile(filepath.Join(filepath.Dir(appfile), Config))
	if err != nil {
		panic(err)
	}
	var cfg ConfigData
	if err = json.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}
	mychecks := []*analysis.Analyzer{
		OsExitAnalyzer, // свой анализатор
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		errcheck.Analyzer,
		maligned.Analyzer,
	}

	// добавляем анализаторы из staticcheck и stylecheck, которые указаны в файле конфигурации
	// ищем совпадение по первым симовлам или по всей строке
	for _, v := range staticcheck.Analyzers {
		for _, sc := range cfg.Staticcheck {
			if strings.HasPrefix(v.Analyzer.Name, sc) {
				mychecks = append(mychecks, v.Analyzer)
			}
		}
	}
	for _, v := range stylecheck.Analyzers {
		for _, sc := range cfg.Stylecheck {
			if strings.HasPrefix(v.Analyzer.Name, sc) {
				mychecks = append(mychecks, v.Analyzer)
			}
		}
	}

	multichecker.Main(
		mychecks...,
	)
}
