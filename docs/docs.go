package docs

import (
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/rest-api/internal/setup"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type EP struct {
	setup.Endpoint
}

type Group struct {
	Name   string
	EPList []setup.Endpoint
}

type Key struct {
	Group   string
	Version string
}

func ParseTemplate() error {
	b, err := ioutil.ReadFile("./docs/source/index.html.tmpl")
	if err != nil {
		return fmt.Errorf("could not read index.html.tmpl file %w", err)
	}

	t := template.Must(template.New("api-docs").Parse(string(b))) // Parse template file.

	f, err := os.Create("./docs/source/index.html.md")
	if err != nil {
		return fmt.Errorf("could not create index.html.md file %w", err)
	}
	epg := EPGroups()
	err = t.Execute(f, epg)
	if err != nil {
		return fmt.Errorf("could not execute template  %w", err)
	}
	f.Close()
	return nil
}

func EPGroups() map[Key]Group {
	gs := make(map[Key]Group)
	// first organize the endpoints into groups
	for _, ep := range setup.EndpointsList() {
		k := Key{Group: ep.Group, Version: ep.Version}
		gs[k] = Group{}
	}

	// add endpoints to the correct groups
	for g, v := range gs {
		for _, ep := range setup.EndpointsList() {
			if g.Group == ep.Group && g.Version == ep.Version {
				v.EPList = append(v.EPList, ep)
			}
		}
		caser := cases.Title(language.English, cases.Compact)
		v.Name = caser.String(g.Group + " " + g.Version)
		if g.Group == "" && g.Version == "" {
			v.Name = "Domain"
		}
		gs[g] = v
	}

	return gs
}
