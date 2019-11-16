package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
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

func ifaceName(name string) string {
	name = strings.TrimPrefix(name, "wl_")
	if len(name) == 0 {
		// XXX
	}
	return goIdentifier(name) + "Interface"
}

// XXX check for reserved names

func exportedGoIdentifier(name string) string {
	name = strings.TrimSuffix(name, "_")
	return ident.Name(strings.Split(name, "_")).ToMixedCaps()
}

func goIdentifier(name string) string {
	name = strings.TrimSuffix(name, "_")
	return ident.Name(strings.Split(name, "_")).ToLowerCamelCase()
}

func printEnums(iface elInterface) {
	typeName := typeName(iface.Name)

	for _, enum := range iface.Enums {
		// TODO emit enum.Description
		fmt.Println("const (")

		ename := typeName + exportedGoIdentifier(enum.Name)
		for _, entry := range enum.Entries {
			eename := ename + exportedGoIdentifier(entry.Name)
			// TODO emit entry.Summary or entry.Description
			fmt.Printf("%s = %s\n", eename, entry.Value)
		}

		fmt.Println(")")
	}
}

func printRequests(iface elInterface) {
	for i, req := range iface.Requests {
		reqName := exportedGoIdentifier(req.Name)
		// TODO emit req.Description
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
				// XXX interface
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
				// XXX
				typ = "int32"
			case "fd":
				// XXX
				typ = "int32"
			default:
				// XXX
				panic(fmt.Sprintf("unsupported type %s", arg.Type))
			}
			fmt.Printf("%s %s, ", goIdentifier(arg.Name), typ)
		}
		fmt.Printf(")")
		if ctor.Interface != "" {
			typ := "*" + typeName(ctor.Interface)
			fmt.Print(typ)
		}
		fmt.Println("{")
		fmt.Printf("const %s_%s = %d\n", iface.Name, req.Name, i)

		if ctor.Name != "" {
			if ctor.Interface != "" {
				fmt.Printf("_ret := &%s{}; obj.Conn().NewProxy(0, _ret);\n", typeName(ctor.Interface))
			}
		}

		fmt.Printf("obj.Conn().SendRequest(obj, %s_%s, ", iface.Name, req.Name)
		for _, arg := range req.Args {
			switch arg.Type {
			case "new_id":
				if ctor.Interface != "" {
					fmt.Print("_ret,")
				} else {
					fmt.Printf("%s,", ctor.Name)
				}
			default:
				fmt.Printf("%s,", goIdentifier(arg.Name))
			}
		}
		fmt.Println(")")

		if ctor.Interface != "" {
			fmt.Println("return _ret")
		}

		fmt.Println("}")
	}
}

func printInterface(iface elInterface) {
	/*
	requests := make([]string, len(iface.Requests))
	for i, req := range iface.Requests {
		args := make([]string, len(req.Args))
		for j, arg := range req.Args {
			args[j] = `"` + arg.Type + `"`
		}
		requests[i] = fmt.Sprintf(`
wayland.MessageRequest{
  Name: "%s",
  Types: []string{%s},
}`, req.Name, strings.Join(args, ","))
	}
*/
	
	events := make([]string, len(iface.Events))
	for i, event := range iface.Events {
		args := make([]string, len(event.Args))
		for j, arg := range event.Args {
			var typ string
			switch arg.Type {
			case "int":
				typ = "int32(0)"
			case "uint":
				typ = "uint32(0)"
			case "fixed":
				typ = "wayland.Fixed(0)"
			case "string":
				typ = `""`
			case "object":
				// when we receive an event with an object arg, then
				// we don't need to know the interface ahead of time,
				// because we'll look up the concrete object based on
				// its ID.
				typ = "nil"
			case "new_id":
				// when we receive an event with a new_id, we need to
				// create the appropriate proxy with the correct
				// interface.
				typ = ifaceName(arg.Interface)
			case "array":
				typ = `"XXX"`
			case "fd":
				typ = `"XXX"`
			default:
				// XXX
				panic(fmt.Sprintf("unsupported type %s", arg.Type))
			}
			args[j] = fmt.Sprintf("%s", typ)
		}

		events[i] = fmt.Sprintf(`
wayland.MessageEvent{
  Name: "%s",
  Types: []interface{}{%s},
}`, event.Name, strings.Join(args, ","))
	}

	fmt.Printf(`
var %s = &wayland.Interface{
  Name: "%s",
  Version: %s,
  Events: []wayland.MessageEvent{%s},
}
`, ifaceName(iface.Name), iface.Name, iface.Version,  strings.Join(events, ","))
}

func main() {
	f, err := os.OpenFile(os.Args[1], os.O_RDONLY, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var spec elProtocol
	dec := xml.NewDecoder(f)
	if err := dec.Decode(&spec); err != nil {
		log.Fatal(err)
	}

	fmt.Println(`package pkg; import "honnef.co/go/wayland";`)
	for _, iface := range spec.Interfaces {
		printEnums(iface)
		printInterface(iface)

		// TODO emit iface.Description
		fmt.Printf("type %s struct { wayland.Proxy }\n", typeName(iface.Name))
		fmt.Printf("func (*%s) Interface() *wayland.Interface { return %s }\n", typeName(iface.Name), ifaceName(iface.Name))

		printRequests(iface)
	}
}
