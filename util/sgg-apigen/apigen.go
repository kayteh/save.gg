package main

////// SPECIFICATION
////
// https://regex101.com/r/iD9aU8/2
////
//
//   // <method> <pathSpec> [<version> [default]]
//  [// <middleware>[(<security flag>...)]...]
//   func <handlerName>
//
////
//
// 	// PATCH /api/user/:slug v1
// 	// RequireSession(All)
//   func patchUser(...) {
//
// becomes
//
//   meta.RegisterRoute("PATCH", "/api/user/:slug", mw.VR(VRMap{
//		"v1": mw.RequireSession(patchUser, &mw.SecurityFlags{All: true}),
//	 }))
//
////
//////

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sync"
	"text/template"
)

var re = regexp.MustCompile(`// ([A-Z]+) ([/a-zA-Z:_0-9-~]+) ?(?:(v[0-9a-zA-Z]+) ?(default)?)?\n(?:// ([a-zA-Z0-9\(\)\ ]+)\n)?func ([a-zA-Z0-9]+)`)
var middlewareRe = regexp.MustCompile(`([A-Za-z0-9]+)(?:\(([A-Za-z0-9]+[,\ ]?)+\))?`)

const tmpl = `{{block "main" .}}package api

//
// ATTENTION: This file is generated automagically.
// Do not touch it. Do not pass go. Do not collect $200.
// Instead run 'go generate' or 'make gen' to build this file.
//

import (
	mw "save.gg/sgg/cmd/sgg-api/run/middleware"
	"save.gg/sgg/meta"
)

func init() {
	{{range .}}
		meta.RegisterRoute("{{.Verb}}", "{{.URI}}", 
		{{- if .Versioned}}
			{{- template "versionedRoute" .}}
		{{else}}
			{{- template "regularRoute" .Versions.default}}
		{{end}})
	{{else}}
	// no routes
	{{end}}
}
{{end}}

{{block "regularRoute" .}}

	{{- range .Middleware}}mw.{{.Name}}({{end}}{{.FuncName}},
{{- range .Middleware}}
	{{- if .HasFlags}}
	&mw.SecurityFlags{
		{{- range .SecFlags}}
			{{.}}: true,
		{{end}}
	}{{end}}),
{{end}}

{{end}}

{{block "versionedRoute" .}}
	mw.VR(mw.VRMap{
		{{- range $k, $r := .Versions}}
		"{{$k}}": {{- template "regularRoute" $r}}
		{{end}}
	}),
{{end}}
`

type route struct {
	Verb      string
	URI       string
	Versions  map[string]routeVersion
	Versioned bool
}

type routeVersion struct {
	FuncName      string
	Middleware    []middleware
	Version       string
	MiddlewareRaw string
}

type middleware struct {
	Name     string
	SecFlags []string
	HasFlags bool
}

// map["VERB URI"]route
var routes map[string]route = map[string]route{}
var l sync.Mutex
var wg sync.WaitGroup

func main() {
	if len(os.Args) == 2 {
		os.Chdir(os.Args[1])
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	d, err := ioutil.ReadDir(cwd)
	if err != nil {
		log.Fatal(err)
	}

	if len(d) < 1 {
		log.Fatalf("[ERR] no files found in %s", cwd)
	}

	filepath.Walk(cwd, walker)

	wg.Wait()

	renderFile(cwd)

}

func walker(path string, info os.FileInfo, err error) error {
	if filepath.Ext(path) != ".go" || filepath.Base(path) == "api.generated.go" {
		log.Printf("%s skipped", path)
		return nil
	}

	wg.Add(1)
	go processFile(path)

	return nil
}

func processFile(path string) {
	defer func() {
		r := recover()
		if r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			log.Println("[ERR]", r)
		}
	}()

	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("[ERR] couldn't read file: %s", err)
	}
	m := re.FindAllStringSubmatch(string(b), -1)
	for _, match := range m {
		processMatches(match)
	}

	wg.Done()
	log.Printf("%s finished", path)
}

// 1 = verb
// 2 = uri
// 3 = version opt.
// 4 = version default opt.
// 5 = middleware opt.
// 6 = func name
func processMatches(m []string) {

	rv := routeVersion{
		FuncName:      m[6],
		Version:       m[3],
		MiddlewareRaw: m[5],
		Middleware:    parseMiddleware(m[5]),
	}
	_, exists := routes[fmt.Sprintf("%s %s", m[1], m[2])]
	l.Lock()
	defer l.Unlock()

	var rsr route
	if exists {

		rsr = routes[fmt.Sprintf("%s %s", m[1], m[2])]

		rsr.Versioned = true
		rsr.Versions[m[3]] = rv

		if m[4] != "" {
			rsr.Versions["default"] = rv
		}

		routes[fmt.Sprintf("%s %s", m[1], m[2])] = rsr

	} else {

		rsr = route{
			Verb:      m[1],
			URI:       m[2],
			Versioned: false,
			Versions:  map[string]routeVersion{"default": rv},
		}

		if m[3] != "" {
			rsr.Versioned = true
			rsr.Versions[m[3]] = rv

		}

		routes[fmt.Sprintf("%s %s", m[1], m[2])] = rsr

	}
}

func parseMiddleware(in string) []middleware {
	m := middlewareRe.FindAllStringSubmatch(in, -1)

	var mw []middleware

	for _, ms := range m {

		c := middleware{
			Name: ms[1],
		}

		if ms[2] != "" && len(ms) > 2 {
			c.SecFlags = ms[2:]
			c.HasFlags = len(c.SecFlags) != 0
		}

		mw = append(mw, c)
	}

	return mw
}

func renderFile(cwd string) {
	defer func() {
		r := recover()
		if r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			log.Println("[ERR]", r)
		}
	}()

	t := template.Must(template.New("tmpl").Parse(tmpl))

	f, err := os.OpenFile(cwd+"/api.generated.go", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("[ERR] couldn't write to output file: %s", err)
	}
	defer f.Close()

	err = t.ExecuteTemplate(f, "main", routes)
	if err != nil {
		log.Fatal(err)
	}

	goFmt(cwd + "/api.generated.go")

	log.Print("Finished!")
}

func goFmt(p string) {
	exec.Command("go", "fmt", p).Run()
}
