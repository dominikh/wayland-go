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

func eventTypeName(iface elInterface, event elEvent) string {
	return fmt.Sprintf("%sEvent%s", typeName(iface.Name), exportedGoIdentifier(event.Name))
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
	for i, req := range iface.Requests {
		printDocs(req.Description)

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
				// XXX
				typ = "int32"
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
					// a new_id without an interface turns into "sun", i.e. interface name, interface version, id.
					fmt.Printf("%s.Interface().Name, version, %s,", ctor.Name, ctor.Name)
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
	events := make([]string, len(iface.Events))
	for i, event := range iface.Events {
		events[i] = fmt.Sprintf("(*%s)(nil)", eventTypeName(iface, event))
	}

	fmt.Printf(`
var %s = &wayland.Interface{
  Name: "%s",
  Version: %s,
  Events: []wayland.Event{%s},
}
`, ifaceName(iface.Name), iface.Name, iface.Version, strings.Join(events, ","))
}

func printEvents(iface elInterface) {
	for _, event := range iface.Events {
		fmt.Printf("type %s struct {\n", eventTypeName(iface, event))
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
				// XXX
				typ = "uintptr"
			case "fd":
				typ = "uintptr"
			}
			printArgDocs(arg)
			fmt.Printf("%s %s", exportedGoIdentifier(arg.Name), typ)
			if arg.Type == "new_id" {
				fmt.Print(" `wl:\"new_id\"`")
			}
			fmt.Println()
		}
		fmt.Printf("}\n\n")
	}
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
		printEvents(iface)

		printDocs(iface.Description)
		fmt.Printf("type %s struct { wayland.Proxy }\n", typeName(iface.Name))

		fmt.Printf("func (*%s) Interface() *wayland.Interface { return %s }\n", typeName(iface.Name), ifaceName(iface.Name))

		printRequests(iface)
	}
}
