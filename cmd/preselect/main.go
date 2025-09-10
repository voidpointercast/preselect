package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"strings"

	"preselect/api"
	"preselect/business"
	"preselect/data"
)

type extMapFlag map[string][]rune

func (e *extMapFlag) String() string {
	var parts []string
	for k, v := range *e {
		parts = append(parts, fmt.Sprintf("%s=%s", k, string(v)))
	}
	return strings.Join(parts, ",")
}

func (e *extMapFlag) Set(value string) error {
	parts := strings.SplitN(value, "=", 2)
	ext := strings.TrimPrefix(strings.ToLower(parts[0]), ".")
	var delims []rune
	if len(parts) == 2 {
		delims = []rune(parts[1])
	}
	if *e == nil {
		*e = make(map[string][]rune)
	}
	(*e)[ext] = delims
	return nil
}

func main() {
	var root string
	mappings := make(extMapFlag)

	flag.StringVar(&root, "dir", ".", "directory to scan")
	flag.Var(&mappings, "ext", "extension to delimiters mapping (e.g. -ext txt=,;)")
	flag.Parse()

	cfg := api.Config{Root: root, ExtMap: map[string]api.SourceFactory{}}
	for ext, delims := range mappings {
		d := delims
		cfg.ExtMap[ext] = func(r io.Reader) business.DataSource {
			return data.NewLoader(r, d)
		}
	}

	app := api.New(cfg)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
