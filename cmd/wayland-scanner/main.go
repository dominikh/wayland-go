package main

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/shurcooL/graphql/ident"
)

type elProtocol struct {
	Name string `xml:"name,attr"`

	Copyright   elCopyright   `xml:"copyright"`
	Description elDescription `xml:"description"`
	Interfaces  []elInterface `xml:"interface"`
}

type elCopyright struct {
	Text string `xml:",cdata"`
}

type elInterface struct {
	Name    string `xml:"name,attr"`
	Version string `xml:"version,attr"`

	Description elDescription `xml:"description"`
	Requests    []elRequest   `xml:"request"`
	Events      []elEvent     `xml:"event"`
	Enums       []elEnum      `xml:"enum"`
}

type elRequest struct {
	Name  string `xml:"name,attr"`
	Type  string `xml:"type,attr"`
	Since string `xml:"since,attr"`

	Description elDescription `xml:"description"`
	Args        []elArg       `xml:"arg"`
}

type elEvent struct {
	Name  string `xml:"name,attr"`
	Since string `xml:"since,attr"`

	Description elDescription `xml:"description"`
	Args        []elArg       `xml:"arg"`
}

type elEnum struct {
	Name     string `xml:"name,attr"`
	Since    string `xml:"since,attr"`
	Bitfield string `xml:"bitfield,attr"`

	Description elDescription `xml:"description"`
	Entries     []elEntry     `xml:"entry"`
}

type elEntry struct {
	Name    string `xml:"name,attr"`
	Value   string `xml:"value,attr"`
	Summary string `xml:"summary,attr"`
	Since   string `xml:"since,attr"`

	Description elDescription `xml:"description"`
}

type elArg struct {
	Name      string `xml:"name,attr"`
	Type      string `xml:"type,attr"`
	Summary   string `xml:"summary,attr"`
	Interface string `xml:"interface,attr"`
	AllowNull string `xml:"allow-null,attr"`
	Enum      string `xml:"enum,attr"`

	Description elDescription `xml:"description"`
}

type elDescription struct {
	Summary string `xml:"summary,attr"`

	Text string `xml:",cdata"`
}

func typeName(name string) string {
	name = strings.TrimPrefix(name, "wl_")
	if len(name) == 0 {
		// XXX
	}
	return exportedGoIdentifier(name)
}

func eventTypeName(iface elInterface, event elEvent) string {
	return fmt.Sprintf("%sEvent%s", typeName(iface.Name), exportedGoIdentifier(event.Name))
}

func eventsTypeName(iface elInterface) string {
	return fmt.Sprintf("%sEvents", typeName(iface.Name))
}

// XXX check for reserved names

func exportedGoIdentifier(name string) string {
	name = strings.TrimSuffix(name, "_")
	return ident.Name(strings.Split(name, "_")).ToMixedCaps()
}

func goIdentifier(name string) string {
	name = strings.TrimSuffix(name, "_")
	name = ident.Name(strings.Split(name, "_")).ToLowerCamelCase()
	return mapReserved(name)
}

func mapReserved(name string) string {
	switch name {
	case "interface":
		return "interface_"
	default:
		return name
	}
}

func printRequests(iface elInterface) {
	hasDestroy := false
	for i, req := range iface.Requests {
		if req.Name == "destroy" {
			hasDestroy = true
		}

		var ctor elArg
		if ctor.Name != "" {
			if ctor.Interface != "" {
				fmt.Printf("_ret := &%s{}; obj.Conn().NewProxy(0, _ret, obj.Queue());\n", typeName(ctor.Interface))
			}
		}

		fmt.Printf("obj.Conn().SendRequest(obj, %d, ", i)
		for _, arg := range req.Args {
			switch arg.Type {
			case "new_id":
				if ctor.Interface != "" {
					fmt.Print("_ret,")
				} else {
					// a new_id without an interface turns into "sun", i.e. interface name, interface version, id.
					fmt.Printf("%s.Interface().Name, version, %s,", ctor.Name, ctor.Name)
				}
			default:
				fmt.Printf("%s,", goIdentifier(arg.Name))
			}
		}
		fmt.Println(")")

		if req.Type == "destructor" {
			fmt.Println("obj.Conn().Destroy(obj)")
		}

		if ctor.Interface != "" {
			fmt.Println("return _ret")
		}

		fmt.Println("}")
	}

	if !hasDestroy {
		fmt.Printf("\nfunc (obj *%s) Destroy() { obj.Conn().Destroy(obj) }\n", typeName(iface.Name))
	}
}

func wlprotoInterfaceName(iface elInterface) string {
	name := strings.TrimPrefix(iface.Name, "wl_")
	if len(name) == 0 {
		// XXX
	}
	return goIdentifier(name) + "Interface"
}

func wlprotoArg(arg elArg) string {
	var typ string
	switch arg.Type {
	case "int":
		typ = "ArgTypeInt"
	case "uint":
		typ = "ArgTypeUint"
	case "fixed":
		typ = "ArgTypeFixed"
	case "string":
		typ = "ArgTypeString"
	case "object":
		typ = "ArgTypeObject"
	case "new_id":
		typ = "ArgTypeNewID"
	case "array":
		typ = "ArgTypeArray"
	case "fd":
		typ = "ArgTypeFd"
	default:
		panic("XXX")
	}
	if arg.Interface == "" {
		return fmt.Sprintf("{Type: wlproto.%s}", typ)
	} else {
		return fmt.Sprintf("{Type: wlproto.%s, Aux: reflect.TypeOf((*%s)(nil))}", typ, typeName(arg.Interface))
	}
}

func goTypeFromWlType(typ string, iface string) string {
	switch typ {
	case "int":
		return "int32"
	case "uint":
		return "uint32"
	case "fixed":
		return "wayland.Fixed"
	case "string":
		return "string"
	case "object":
		if iface != "" {
			return "*" + typeName(iface)
		} else {
			return "wayland.Object"
		}
	case "new_id":
		if iface != "" {
			return "*" + typeName(iface)
		} else {
			return "wayland.Object"
		}
	case "array":
		return "[]byte"
	case "fd":
		return "uintptr"
	default:
		return "UNKNOWN_TYPE_" + typ
	}
}

func docString(docs elDescription) string {
	text := docs.Text
	if text == "" {
		text = docs.Summary
		if text == "" {
			return ""
		}
	}

	text = strings.TrimSpace(text)
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = "// " + strings.TrimSpace(line)
	}
	return strings.Join(lines, "\n")
}

func printDocs(docs elDescription) {
	text := docs.Text
	if text == "" {
		text = docs.Summary
		if text == "" {
			return
		}
	}

	text = strings.TrimSpace(text)
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		fmt.Println("//", strings.TrimSpace(line))
	}
}

func enumEntryDocString(entry elEntry) string {
	if entry.Description.Text != "" || entry.Description.Summary != "" {
		return docString(entry.Description)
	}
	if entry.Summary == "" {
		return ""
	}
	return "// " + entry.Summary
}

var pkgTmpl = `// Code generated by wayland-scanner; DO NOT EDIT.
//
// Package pkg contains generated definitions of Wayland protocols.
//
// It was generated from the following files:
{{- range .InputFiles }}
// 	- {{ . }}
{{- end }}
package pkg

import (
{{- range .Imports }}
	"{{ . }}"
{{- end }}
)

var Interfaces = map[string]*wlproto.Interface{
{{- range .Specs }}
	{{- range .Interfaces }}
		"{{ .Name }}": {{ WlprotoInterfaceName . }},
	{{- end }}
{{ end }}
}

var Requests = map[string]*wlproto.Request{
{{- range $spec := .Specs }}
	{{- range $iface := .Interfaces }}
		{{- range $ireq, $req := .Requests }}
			"{{ $iface.Name }}_{{ $req.Name }}": &{{ WlprotoInterfaceName $iface }}.Requests[{{ $ireq }}],
		{{- end }}
	{{ end }}
{{ end }}
}

var Events = map[string]*wlproto.Event{
{{- range $spec := .Specs }}
	{{- range $iface := .Interfaces }}
		{{- range $iev, $ev := .Events }}
			"{{ $iface.Name }}_{{ $ev.Name }}": &{{ WlprotoInterfaceName $iface }}.Events[{{ $iev }}],
		{{- end }}
	{{ end }}
{{ end }}
}

{{ range $spec := .Specs }}
	{{ range $iface := $spec.Interfaces }}
		{{ range $enum := $iface.Enums }}
			{{ DocString .Description }}
			const (
			{{ $ename := print (TypeName $iface.Name) (ExportedGoIdentifier $enum.Name) }}
			{{ range $entry := $enum.Entries }}
				{{ EnumEntryDocString $entry }}
				{{ print $ename (ExportedGoIdentifier $entry.Name)  }} = {{ $entry.Value }}
			{{ end }}
			)
		{{ end }}

		var {{ WlprotoInterfaceName $iface }} = &wlproto.Interface{
			Name: "{{ $iface.Name }}",
			Version: {{ $iface.Version }},
			Requests: []wlproto.Request{
				{{ range $req := $iface.Requests }}
					{
						Name: "{{ $req.Name }}",
						Type: "{{ $req.Type }}",
						Since: {{ if $req.Since }} {{ $req.Since }} {{ else }} 1 {{ end }},
						Args: []wlproto.Arg{
							{{ range $arg := $req.Args }}
								{{ WlprotoArg $arg }},
							{{ end }}
						},
					},
				{{ end }}
			},
			Events: []wlproto.Event{
				{{ range $ev := $iface.Events }}
					{
						Name: "{{ $ev.Name }}",
						Since: {{ if $ev.Since }} {{ $ev.Since }} {{ else }} 1 {{ end }},
						Args: []wlproto.Arg{
							{{ range $arg := $ev.Args }}
								{{ WlprotoArg $arg }},
							{{ end }}
						},
					},
				{{ end }}
			},
		}

		{{ DocString .Description }}
		type {{ TypeName .Name }} struct { wayland.Proxy }

		func (*{{ TypeName .Name }}) Interface() *wlproto.Interface { return {{ WlprotoInterfaceName . }} }

		func (obj *{{ TypeName .Name }}) WithQueue(queue *wayland.EventQueue) *{{ TypeName .Name }} {
			wobj := &{{ TypeName .Name }}{}
			obj.Conn().NewWrapper(obj, wobj, queue)
			return wobj
		}

		type {{ EventsTypeName $iface }} struct {
			{{ range $event := $iface.Events }}
				{{ ExportedGoIdentifier $event.Name }} func(obj *{{TypeName $iface.Name}},
				{{- range $arg := $event.Args -}}
					{{ GoIdentifier $arg.Name }} {{ GoTypeFromWlType $arg.Type $arg.Interface }},
				{{- end -}}
				)
			{{ end }}
		}

		func (obj *{{ TypeName $iface.Name }}) AddListener(listeners {{ EventsTypeName $iface }}) {
			obj.Proxy.SetListeners(
				{{- range $event := $iface.Events -}}
					listeners.{{  ExportedGoIdentifier $event.Name}},
				{{- end -}}
			)
		}

		{{ $hasDestroy := false }}
		{{ range $ireq, $req := $iface.Requests }}
			{{ if eq $req.Name "destroy" }}
				{{ $hasDestroy = true }}
			{{ end }}
		    {{ $ctor := $.NoArg }}
			{{ DocString $req.Description }}
		    func (obj *{{ TypeName $iface.Name }}) {{ ExportedGoIdentifier $req.Name }}(
				{{- range $arg := $req.Args -}}
					{{- if eq $arg.Type "new_id" -}}
						{{- $ctor = $arg -}}
					{{- end -}}
					{{- if (or (ne $arg.Type "new_id") (eq $arg.Interface "")) -}}
						{{- GoIdentifier $arg.Name }} {{ GoTypeFromWlType $arg.Type $arg.Interface -}},
						{{- if eq $arg.Type "new_id" -}}
							version uint32,
						{{- end -}}
					{{- end -}}
				{{- end -}}
			) {{ if ne $ctor.Interface "" }} *{{ TypeName $ctor.Interface }} {{ end }} {
				{{- if (and (ne $ctor.Name "") (ne $ctor.Interface "")) }}
					_ret := &{{ TypeName $ctor.Interface }}{}
					obj.Conn().NewProxy(0, _ret, obj.Queue())
				{{ end }}
				obj.Conn().SendRequest(obj, {{ $ireq }},
					{{- range $arg := $req.Args -}}
						{{- if eq $arg.Type "new_id" -}}
							{{- if ne $ctor.Interface "" -}}
								_ret
							{{- else -}}
								{{- GoIdentifier $ctor.Name}}.Interface().Name, version, {{ GoIdentifier $ctor.Name -}}
							{{- end -}}
						{{- else -}}
							{{- GoIdentifier $arg.Name -}}
						{{- end }},
					{{- end -}}
				)
				{{ if eq $req.Type "destructor" }}
					obj.Conn().Destroy(obj)
				{{ end }}
				{{ if ne $ctor.Interface "" }}
					return _ret
				{{ end }}
			}
		{{ end }}
		{{ if not $hasDestroy }}
			func (obj *{{ TypeName $iface.Name }}) Destroy() { obj.Conn().Destroy(obj) }
		{{ end }}
	{{ end }}
{{ end }}
`

type tmplState struct {
	InputFiles []string
	Imports    []string
	Specs      []elProtocol

	NoArg elArg
}

func main() {
	var specs []elProtocol
	for _, arg := range os.Args[1:] {
		f, err := os.OpenFile(arg, os.O_RDONLY, 0)
		if err != nil {
			log.Fatal(err)
		}

		var spec elProtocol
		dec := xml.NewDecoder(f)
		if err := dec.Decode(&spec); err != nil {
			log.Fatal(err)
		}
		specs = append(specs, spec)
		f.Close()
	}

	var state tmplState
	for _, arg := range os.Args[1:] {
		state.InputFiles = append(state.InputFiles, filepath.Base(arg))
	}
	state.Imports = []string{
		"reflect",
		"honnef.co/go/wayland",
		"honnef.co/go/wayland/wlproto",
	}
	state.Specs = specs

	tmpl, err := template.New("pkg").Funcs(template.FuncMap{
		"TypeName":             typeName,
		"WlprotoInterfaceName": wlprotoInterfaceName,
		"WlprotoArg":           wlprotoArg,
		"DocString":            docString,
		"EnumEntryDocString":   enumEntryDocString,
		"ExportedGoIdentifier": exportedGoIdentifier,
		"GoIdentifier":         goIdentifier,
		"GoTypeFromWlType":     goTypeFromWlType,
		"EventsTypeName":       eventsTypeName,
	}).Parse(pkgTmpl)
	if err != nil {
		log.Fatal(err)
	}
	if err := tmpl.Execute(os.Stderr, state); err != nil {
		log.Fatal(err)
	}

	/*
		for _, iface := range spec.Interfaces {
			printRequests(iface)
		}
	*/
}
