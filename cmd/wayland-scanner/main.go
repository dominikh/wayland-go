package main

import (
	"encoding/xml"
	"fmt"
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

func printEnums(iface elInterface) {
	typeName := typeName(iface.Name)

	for _, enum := range iface.Enums {
		printDocs(enum.Description)
		fmt.Println("const (")

		ename := typeName + exportedGoIdentifier(enum.Name)
		for _, entry := range enum.Entries {
			eename := ename + exportedGoIdentifier(entry.Name)
			printEnumEntryDocs(entry)
			fmt.Printf("%s = %s\n", eename, entry.Value)
		}

		fmt.Println(")")
	}
}

func printRequests(iface elInterface) {
	hasDestroy := false
	for i, req := range iface.Requests {
		printDocs(req.Description)

		if req.Name == "destroy" {
			hasDestroy = true
		}

		reqName := exportedGoIdentifier(req.Name)
		fmt.Printf("func (obj *%s) %s(", typeName(iface.Name), reqName)
		var ctor elArg
		for _, arg := range req.Args {
			var typ string
			switch arg.Type {
			case "int":
				typ = "int32"
			case "uint":
				typ = "uint32"
			case "fixed":
				typ = "wayland.Fixed"
			case "string":
				typ = "string"
			case "object":
				typ = "wayland.Object"
				if arg.Interface != "" {
					typ = "*" + typeName(arg.Interface)
				}
			case "new_id":
				ctor = arg
				if ctor.Interface != "" {
					continue
				} else {
					typ = "wayland.Object"
				}
			case "array":
				typ = "[]byte"
			case "fd":
				typ = "uintptr"
			default:
				// XXX
				panic(fmt.Sprintf("unsupported type %s", arg.Type))
			}
			fmt.Printf("%s %s, ", goIdentifier(arg.Name), typ)
			if ctor.Name != "" && ctor.Interface == "" {
				fmt.Printf("version uint32, ")
			}
		}
		fmt.Printf(")")
		if ctor.Interface != "" {
			typ := "*" + typeName(ctor.Interface)
			fmt.Print(typ)
		}
		fmt.Println("{")

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

func wlprotoInterface(iface elInterface) {
	fmt.Printf(`
var %s = &wlproto.Interface{
  Name: "%s",
  Version: %s,
  Requests: []wlproto.Request{
`, wlprotoInterfaceName(iface), iface.Name, iface.Version)
	for _, req := range iface.Requests {
		wlprotoRequest(req)
		fmt.Println(",")
	}
	fmt.Print(`
},
  Events: []wlproto.Event{
`)

	for _, ev := range iface.Events {
		wlprotoEvent(ev)
		fmt.Println(",")
	}
	fmt.Println(`
  },
}
`)
}

func wlprotoRequest(req elRequest) {
	if req.Since == "" {
		req.Since = "1"
	}
	fmt.Printf(`wlproto.Request{
		Name: %q,
		Type: %q,
		Since: %s,
		Args: []wlproto.Arg{
			`, req.Name, req.Type, req.Since)

	for _, arg := range req.Args {
		wlprotoArg(arg)
		fmt.Println(",")
	}
	fmt.Print(`
		},
	}`)
}

func wlprotoEvent(ev elEvent) {
	if ev.Since == "" {
		ev.Since = "1"
	}
	fmt.Printf(`wlproto.Event{
		Name: %q,
		Since: %s,
		Args: []wlproto.Arg{
			`, ev.Name, ev.Since)

	for _, arg := range ev.Args {
		wlprotoArg(arg)
		fmt.Println(",")
	}
	fmt.Print(`
		},
	}`)
}

func wlprotoArg(arg elArg) {
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
		fmt.Printf("{Type: wlproto.%s}", typ)
	} else {
		fmt.Printf("{Type: wlproto.%s, Aux: reflect.TypeOf((*%s)(nil))}", typ, typeName(arg.Interface))
	}
}

func printEvents(iface elInterface) {
	fmt.Printf("type %s struct {\n", eventsTypeName(iface))
	for _, event := range iface.Events {
		fmt.Printf("%s func(obj *%s, ", exportedGoIdentifier(event.Name), typeName(iface.Name))
		for _, arg := range event.Args {
			var typ string
			switch arg.Type {
			case "int":
				typ = "int32"
			case "uint":
				typ = "uint32"
			case "fixed":
				typ = "wayland.Fixed"
			case "string":
				typ = "string"
			case "object":
				if arg.Interface != "" {
					typ = "*" + typeName(arg.Interface)
				} else {
					typ = "wayland.Object"
				}
			case "new_id":
				typ = "*" + typeName(arg.Interface)
			case "array":
				typ = "[]byte"
			case "fd":
				typ = "uintptr"
			}
			fmt.Printf("%s %s,", goIdentifier(arg.Name), typ)
		}
		fmt.Println(")")
	}
	fmt.Println("}")

	fmt.Printf("func (obj *%s) AddListener(listeners %s) {\n", typeName(iface.Name), eventsTypeName(iface))
	fmt.Print("obj.Proxy.SetListeners(")
	for _, event := range iface.Events {
		fmt.Printf("listeners.%s,", exportedGoIdentifier(event.Name))
	}
	fmt.Println(")")
	fmt.Println("}")
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

func printEnumEntryDocs(entry elEntry) {
	if entry.Description.Text != "" || entry.Description.Summary != "" {
		printDocs(entry.Description)
		return
	}
	if entry.Summary == "" {
		return
	}
	fmt.Println("//", entry.Summary)
}

func printArgDocs(arg elArg) {
	if arg.Description.Text != "" || arg.Description.Summary != "" {
		printDocs(arg.Description)
		return
	}
	if arg.Summary == "" {
		return
	}
	fmt.Println("//", arg.Summary)
}

func printInterfacesMap(specs []elProtocol) {
	fmt.Println("var Interfaces = map[string]*wlproto.Interface{")
	for _, spec := range specs {
		for _, iface := range spec.Interfaces {
			fmt.Printf("\"%s\": %s,\n", iface.Name, wlprotoInterfaceName(iface))
		}
	}
	fmt.Println("}")
}

func printMethodsMap(specs []elProtocol) {
	fmt.Println("var Requests = map[string]*wlproto.Request{")
	for _, spec := range specs {
		for _, iface := range spec.Interfaces {
			for i, req := range iface.Requests {
				fmt.Printf("\"%s_%s\": &%s.Requests[%d],\n", iface.Name, req.Name, wlprotoInterfaceName(iface), i)
			}
		}
	}
	fmt.Println("}")

	fmt.Println("var Events = map[string]*wlproto.Event{")
	for _, spec := range specs {
		for _, iface := range spec.Interfaces {
			for i, ev := range iface.Events {
				fmt.Printf("\"%s_%s\": &%s.Events[%d],\n", iface.Name, ev.Name, wlprotoInterfaceName(iface), i)
			}
		}
	}
	fmt.Println("}")
}

func main() {
	fmt.Print("// Code generated by wayland-scanner; DO NOT EDIT.\n\n")
	imports := []string{
		"reflect",
		"honnef.co/go/wayland",
		"honnef.co/go/wayland/wlproto",
	}

	fmt.Println("// Package pkg contains generated definitions of Wayland protocols.")
	fmt.Println("//")
	fmt.Println("// It was generated from the following files:")
	for _, arg := range os.Args[1:] {
		base := filepath.Base(arg)
		fmt.Printf("// 	%s\n", base)
	}
	fmt.Println("package pkg")
	fmt.Println("import (")
	for _, imp := range imports {
		fmt.Printf("%q\n", imp)
	}
	fmt.Println(")")

	var specs []elProtocol
	for _, arg := range os.Args[1:] {
		f, err := os.OpenFile(arg, os.O_RDONLY, 0)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		var spec elProtocol
		dec := xml.NewDecoder(f)
		if err := dec.Decode(&spec); err != nil {
			log.Fatal(err)
		}
		specs = append(specs, spec)

		for _, iface := range spec.Interfaces {
			printEnums(iface)

			wlprotoInterface(iface)
			printEvents(iface)

			printDocs(iface.Description)
			fmt.Printf("type %s struct { wayland.Proxy }\n", typeName(iface.Name))

			fmt.Printf("func (*%s) Interface() *wlproto.Interface { return %s }\n",
				typeName(iface.Name), wlprotoInterfaceName(iface))

			fmt.Printf(`
func (obj *%s) WithQueue(queue *wayland.EventQueue) *%s {
  wobj := &%s{}
  obj.Conn().NewWrapper(obj, wobj, queue)
  return wobj
}
`, typeName(iface.Name), typeName(iface.Name), typeName(iface.Name))

			printRequests(iface)
		}
	}

	printInterfacesMap(specs)
	printMethodsMap(specs)
}
