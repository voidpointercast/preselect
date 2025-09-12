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

type extConfig struct {
	delims []rune
	quote  rune
}

type extMapFlag map[string]extConfig

func (e *extMapFlag) String() string {
	var parts []string
	for k, v := range *e {
		s := string(v.delims)
		if v.quote != 0 {
			s += ":" + string(v.quote)
		}
		parts = append(parts, fmt.Sprintf("%s=%s", k, s))
	}
	return strings.Join(parts, ",")
}

func (e *extMapFlag) Set(value string) error {
	parts := strings.SplitN(value, "=", 2)
	ext := strings.TrimPrefix(strings.ToLower(parts[0]), ".")
	var delims []rune
	var quote rune
	if len(parts) == 2 {
		dq := strings.SplitN(parts[1], ":", 2)
		delims = []rune(dq[0])
		if len(dq) == 2 {
			q := []rune(dq[1])
			if len(q) > 0 {
				quote = q[0]
			}
		}
	}
	if *e == nil {
		*e = make(map[string]extConfig)
	}
	(*e)[ext] = extConfig{delims: delims, quote: quote}
	return nil
}

func main() {
	var root string
	mappings := make(extMapFlag)

	flag.StringVar(&root, "dir", ".", "directory to scan")
	flag.Var(&mappings, "ext", "extension configuration (e.g. -ext txt=,; -ext csv=;:\" )")
	flag.Parse()

	cfg := api.Config{Root: root, ExtMap: map[string]api.SourceFactory{}}
	for ext, conf := range mappings {
		c := conf
		switch ext {
		case "csv":
			delim := ','
			if len(c.delims) > 0 {
				delim = c.delims[0]
			}
			quote := '"'
			if c.quote != 0 {
				quote = c.quote
			}
			cfg.ExtMap[ext] = func(r io.Reader) business.DataSource {
				return data.NewCSVLoader(r, delim, quote)
			}
		default:
			d := c.delims
			cfg.ExtMap[ext] = func(r io.Reader) business.DataSource {
				return data.NewLoader(r, d)
			}
		}
	}

	app := api.New(cfg)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
