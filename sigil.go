package sigil

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/mgood/go-posix"
)

var (
	TemplatePath    []string
	TemplateDir     string
	PosixPreprocess bool
)

var fnMap = template.FuncMap{}

func Register(fm template.FuncMap) {
	for k, v := range fm {
		fnMap[k] = v
	}
}

func LookPath(file string) (string, error) {
	cwd, _ := os.Getwd()
	search := append([]string{cwd, TemplateDir}, TemplatePath...)
	for _, path := range search {
		filepath := filepath.Join(path, file)
		if _, err := os.Stat(filepath); err == nil {
			return filepath, nil
		}
	}
	return "", fmt.Errorf("Not found in path: %s", file)
}

func Execute(input string, vars map[string]string) (string, error) {
	var tmplVars string
	var err error
	for k, v := range vars {
		err := os.Setenv(k, v)
		if err != nil {
			return "", err
		}
		escaped := strings.Replace(v, "\"", "\\\"", -1)
		tmplVars = tmplVars + fmt.Sprintf("{{ $%s := \"%s\" }}", k, escaped)
	}
	if PosixPreprocess {
		input, err = posix.ExpandEnv(input)
		if err != nil {
			return "", err
		}
	}
	tmpl, err := template.New("template").Funcs(fnMap).Parse(tmplVars + input)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, nil)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
