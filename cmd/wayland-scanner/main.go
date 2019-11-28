package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"go/format"
	"io"
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

func (b *Builder) trimPrefix(name string) string {
	return strings.TrimPrefix(name, b.Prefix)
}

func (b *Builder) typeName(name string) string {
	name = b.trimPrefix(name)
	if len(name) == 0 {
		// XXX
	}
	return exportedGoIdentifier(name)
}

func (b *Builder) eventsTypeName(iface elInterface) string {
	return fmt.Sprintf("%sEvents", b.typeName(iface.Name))
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

func (b *Builder) wlprotoInterfaceName(iface elInterface) string {
	name := b.trimPrefix(iface.Name)
	if len(name) == 0 {
		// XXX
	}
	return goIdentifier(name) + "Interface"
}

func (b *Builder) wlprotoArg(arg elArg) string {
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
		return fmt.Sprintf("{Type: wlproto.%s, Aux: reflect.TypeOf((*%s)(nil))}", typ, b.qualifyTypeName(arg.Interface))
	}
}

func (b *Builder) goTypeFromWlType(typ string, iface string) string {
	switch typ {
	case "int":
		return "int32"
	case "uint":
		return "uint32"
	case "fixed":
		return "wlclient.Fixed"
	case "string":
		return "string"
	case "object":
		if iface != "" {
			return "*" + b.qualifyTypeName(iface)
		} else {
			return "wlclient.Object"
		}
	case "new_id":
		if iface != "" {
			return "*" + b.qualifyTypeName(iface)
		} else {
			return "wlclient.Object"
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

func enumEntryDocString(entry elEntry) string {
	if entry.Description.Text != "" || entry.Description.Summary != "" {
		return docString(entry.Description)
	}
	if entry.Summary == "" {
		return ""
	}
	return "// " + entry.Summary
}

func (b *Builder) qualifyTypeName(name string) string {
	pkg, ok := b.Interfaces[name]
	if !ok {
		panic(fmt.Sprintf("unknown interface %s", name))
	}
	if pkg != nil {
		b.Imports[pkg.Path] = true
		return pkg.Name + "." + pkg.Interfaces[name]
	}
	return b.typeName(name)
}

type Builder struct {
	Spec       elProtocol
	Imports    map[string]bool
	Interfaces map[string]*Package
	Code       bytes.Buffer
	Prefix     string
}

func (b *Builder) Write(data []byte) (int, error) {
	return b.Code.Write(data)
}

func (b *Builder) printSpecs(out io.Writer) {
	printPackage := func() {
		pkg := goIdentifier(b.Spec.Name)
		fmt.Fprintln(out, "// Code generated by wayland-scanner; DO NOT EDIT.")
		fmt.Fprintln(out)
		fmt.Fprintf(out, "// Package %s contains generated definitions of the %s Wayland protocol.\n", pkg, b.Spec.Name)
		if b.Spec.Description.Text != "" {
			fmt.Fprintln(out, "//")
			seenNonEmpty := false
			lines := strings.Split(b.Spec.Description.Text, "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" && !seenNonEmpty {
					continue
				}
				seenNonEmpty = true
				fmt.Fprintf(out, "// %s\n", line)
			}
		}
		fmt.Fprintf(out, "package %s\n", pkg)

	}

	printImports := func() {
		fmt.Fprintln(out, "import (")
		for imp := range b.Imports {
			fmt.Fprintf(out, "\t%q\n", imp)
		}
		fmt.Fprintln(out, ")")
	}

	printMaps := func() {
		fmt.Fprintln(b, "var interfaceNames = map[string]string{")
		for _, iface := range b.Spec.Interfaces {
			fmt.Fprintf(b, "\t%q: %q,\n", iface.Name, b.typeName(iface.Name))
		}
		fmt.Fprint(b, "}\n\n")

		fmt.Fprintln(b, "var Interfaces = map[string]*wlproto.Interface{")
		for _, iface := range b.Spec.Interfaces {
			fmt.Fprintf(b, "\t%q: %s,\n", iface.Name, b.wlprotoInterfaceName(iface))
		}
		fmt.Fprint(b, "}\n\n")

		fmt.Fprintln(b, "var Requests = map[string]*wlproto.Request{")
		for _, iface := range b.Spec.Interfaces {
			for ireq, req := range iface.Requests {
				fmt.Fprintf(b, "\t\"%s_%s\": &%s.Requests[%d],\n", iface.Name, req.Name, b.wlprotoInterfaceName(iface), ireq)
			}
		}
		fmt.Fprint(b, "}\n\n")

		fmt.Fprintln(b, "var Events = map[string]*wlproto.Event{")
		for _, iface := range b.Spec.Interfaces {
			for iev, ev := range iface.Events {
				fmt.Fprintf(b, "\t\"%s_%s\": &%s.Events[%d],\n", iface.Name, ev.Name, b.wlprotoInterfaceName(iface), iev)
			}
		}
		fmt.Fprintln(b, "}")
	}

	printInterface := func(iface elInterface) {
		printEnums := func() {
			for _, enum := range iface.Enums {
				fmt.Fprintln(b, docString(enum.Description))
				fmt.Fprintln(b, "const (")
				for _, entry := range enum.Entries {
					fmt.Fprintf(b, "\t%s\n", enumEntryDocString(entry))
					fmt.Fprintf(b, "\t%s%s%s = %s\n", b.typeName(iface.Name), exportedGoIdentifier(enum.Name), exportedGoIdentifier(entry.Name), entry.Value)
				}
				fmt.Fprintln(b, ")")
			}
		}

		printInterfaceVar := func() {
			fmt.Fprintf(b, "var %s = &wlproto.Interface{\n", b.wlprotoInterfaceName(iface))
			fmt.Fprintf(b, "\tName: %q,\n", iface.Name)
			fmt.Fprintf(b, "\tVersion: %s,\n", iface.Version)

			fmt.Fprintln(b, "\tRequests: []wlproto.Request{")
			for _, req := range iface.Requests {
				if req.Since == "" {
					req.Since = "1"
				}
				fmt.Fprintln(b, "\t\t{")
				fmt.Fprintf(b, "\t\t\tName: %q,\n", req.Name)
				fmt.Fprintf(b, "\t\t\tType: %q,\n", req.Type)
				fmt.Fprintf(b, "\t\t\tSince: %s,\n", req.Since)

				fmt.Fprintln(b, "\t\t\tArgs: []wlproto.Arg{")
				for _, arg := range req.Args {
					fmt.Fprintf(b, "\t\t\t\t%s,\n", b.wlprotoArg(arg))
				}
				fmt.Fprintln(b, "\t\t\t},")

				fmt.Fprintln(b, "\t\t},")
			}
			fmt.Fprintln(b, "\t},")

			fmt.Fprintln(b, "\tEvents: []wlproto.Event{")
			for _, ev := range iface.Events {
				if ev.Since == "" {
					ev.Since = "1"
				}
				fmt.Fprintln(b, "\t\t{")
				fmt.Fprintf(b, "\t\t\tName: %q,\n", ev.Name)
				fmt.Fprintf(b, "\t\t\tSince: %s,\n", ev.Since)

				fmt.Fprintln(b, "\t\t\tArgs: []wlproto.Arg{")
				for _, arg := range ev.Args {
					fmt.Fprintf(b, "\t\t\t\t%s,\n", b.wlprotoArg(arg))
				}
				fmt.Fprintln(b, "\t\t\t},")

				fmt.Fprintln(b, "\t\t},")
			}
			fmt.Fprintln(b, "\t},")

			fmt.Fprintln(b, "}")
		}

		printInterfaceEventsType := func() {
			fmt.Fprintf(b, "type %s struct {\n", b.eventsTypeName(iface))
			for _, ev := range iface.Events {
				fmt.Fprintf(b, "\t%s func(obj *%s,", exportedGoIdentifier(ev.Name), b.typeName(iface.Name))
				for _, arg := range ev.Args {
					fmt.Fprintf(b, "%s %s,", goIdentifier(arg.Name), b.goTypeFromWlType(arg.Type, arg.Interface))
				}
				fmt.Fprintln(b, ")")
			}
			fmt.Fprint(b, "}\n\n")
		}

		printInterfaceType := func() {
			fmt.Fprintln(b, docString(iface.Description))
			fmt.Fprintf(b, "type %s struct { wlclient.Proxy }\n\n", b.typeName(iface.Name))

			// basic methods
			fmt.Fprintf(b, "func (*%s) Interface() *wlproto.Interface { return %s }\n\n", b.typeName(iface.Name), b.wlprotoInterfaceName(iface))

			fmt.Fprintf(b, "func (obj *%[1]s) WithQueue(queue *wlclient.EventQueue) *%[1]s {\n", b.typeName(iface.Name))
			fmt.Fprintf(b, "\twobj := &%s{}\n", b.typeName(iface.Name))
			fmt.Fprintf(b, "\tobj.Conn().NewWrapper(obj, wobj, queue)\n")
			fmt.Fprintf(b, "\treturn wobj\n")
			fmt.Fprintf(b, "}\n\n")

			printInterfaceEventsType()
			fmt.Fprintf(b, "func (obj *%s) AddListener(listeners %s) {\n", b.typeName(iface.Name), b.eventsTypeName(iface))
			fmt.Fprint(b, "\tobj.Proxy.SetListeners(")
			for _, ev := range iface.Events {
				fmt.Fprintf(b, "listeners.%s,", exportedGoIdentifier(ev.Name))
			}
			fmt.Fprintln(b, ")")
			fmt.Fprint(b, "}\n\n")
		}

		printRequest := func(ireq int, req elRequest) {
			var ctor elArg
			fmt.Fprintln(b, docString(req.Description))
			fmt.Fprintf(b, "func (obj *%s) %s(", b.typeName(iface.Name), exportedGoIdentifier(req.Name))
			for _, arg := range req.Args {
				if arg.Type == "new_id" {
					ctor = arg
					if arg.Interface == "" {
						fmt.Fprintf(b, "%s %s, version uint32,", goIdentifier(arg.Name), b.goTypeFromWlType(arg.Type, arg.Interface))
					}
				} else {
					fmt.Fprintf(b, "%s %s,", goIdentifier(arg.Name), b.goTypeFromWlType(arg.Type, arg.Interface))
				}
			}
			fmt.Fprint(b, ")")
			if ctor.Interface != "" {
				fmt.Fprintf(b, "*%s", b.typeName(ctor.Interface))
			}

			fmt.Fprintln(b, "{")
			if ctor.Interface != "" {
				fmt.Fprintf(b, "\t_ret := &%s{}\n", b.typeName(ctor.Interface))
				fmt.Fprintln(b, "\tobj.Conn().NewProxy(0, _ret, obj.Queue())")
			}

			if req.Type == "destructor" {
				fmt.Fprintf(b, "\tobj.Conn().SendDestructor(obj, %d, ", ireq)
			} else {
				fmt.Fprintf(b, "\tobj.Conn().SendRequest(obj, %d, ", ireq)
			}
			for _, arg := range req.Args {
				if arg.Type == "new_id" {
					if ctor.Interface == "" {
						fmt.Fprintf(b, "%[1]s.Interface().Name, version, %[1]s,", goIdentifier(ctor.Name))
					} else {
						fmt.Fprint(b, "_ret,")
					}
				} else {
					fmt.Fprintf(b, "%s,", goIdentifier(arg.Name))
				}
			}
			fmt.Fprintln(b, ")")

			if ctor.Interface != "" {
				fmt.Fprintln(b, "\treturn _ret")
			}

			fmt.Fprint(b, "}\n\n")
		}

		printEnums()
		printInterfaceVar()
		printInterfaceType()

		hasDestroy := false
		for ireq, req := range iface.Requests {
			if req.Name == "destroy" {
				hasDestroy = true
			}
			printRequest(ireq, req)
		}
		if !hasDestroy {
			fmt.Fprintf(b, "func (obj *%s) Destroy() { obj.Conn().Destroy(obj) }\n\n", b.typeName(iface.Name))
		}
	}

	for _, iface := range b.Spec.Interfaces {
		b.Interfaces[iface.Name] = nil
	}

	printMaps()
	for _, iface := range b.Spec.Interfaces {
		printInterface(iface)
	}

	if b.Spec.Name == "wayland" {
		fmt.Fprintln(b, "func GetDisplay(conn *wlclient.Conn) *Display { _ret := &Display{}; conn.NewProxy(1, _ret, nil); return _ret }")
	}

	printPackage()
	printImports()

	out.Write(b.Code.Bytes())
}

type imports []string

func (imps *imports) String() string {
	return strings.Join([]string(*imps), ",")
}

func (imps *imports) Set(s string) error {
	*imps = append(*imps, s)
	return nil
}

func Build(file string, imports []string, prefix string, out io.Writer) {
	b := &Builder{
		Imports: map[string]bool{
			"reflect":                       true,
			"honnef.co/go/wayland/wlclient": true,
			"honnef.co/go/wayland/wlproto":  true,
		},
		Prefix: prefix,
	}

	if len(imports) > 0 {
		b.Interfaces = loadInterfaces(imports)
	} else {
		b.Interfaces = map[string]*Package{}
	}

	f, err := os.OpenFile(flag.Args()[0], os.O_RDONLY, 0)
	if err != nil {
		log.Fatal(err)
	}

	dec := xml.NewDecoder(f)
	if err := dec.Decode(&b.Spec); err != nil {
		log.Fatal(err)
	}
	f.Close()

	var buf bytes.Buffer
	b.printSpecs(&buf)
	d, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	out.Write(d)
}

func main() {
	var imps imports
	flag.Var(&imps, "i", "XXX")
	prefix := flag.String("prefix", "", "XXX")
	flag.Parse()

	Build(flag.Args()[0], imps, *prefix, os.Stdout)
}
