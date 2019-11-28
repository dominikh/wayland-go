package main

import (
	"go/ast"
	"go/token"
	"log"

	"golang.org/x/tools/go/packages"
)

type Package struct {
	Path       string
	Name       string
	Interfaces map[string]string
}

func loadInterfaces(paths []string) map[string]*Package {
	cfg := &packages.Config{
		Mode:  packages.NeedName | packages.NeedSyntax,
		Tests: false,
	}
	pkgs, err := packages.Load(cfg, paths...)
	if err != nil {
		log.Fatal(err)
	}

	// map from package name (not path) to package
	m := map[string]*Package{}

invalidPackage:
	for _, pkg := range pkgs {
		if len(pkg.Errors) != 0 {
			log.Fatal(pkg.Errors)
		}

		if _, ok := m[pkg.Name]; ok {
			log.Fatalf("duplicate package name %q", pkg.Name)
		}

		ifaces :=map[string]string{}
		for _, f := range pkg.Syntax {
			for _, decl := range f.Decls {
				if decl, ok := decl.(*ast.GenDecl); ok && decl.Tok == token.VAR {
				wrongSpec:
					for _, spec := range decl.Specs {
						spec := spec.(*ast.ValueSpec)
						if len(spec.Names) != 1 {
							continue wrongSpec
						}
						if spec.Names[0].Name != "interfaceNames" {
							continue wrongSpec
						}
						if len(spec.Values) != 1 {
							continue invalidPackage
						}
						lit, ok := spec.Values[0].(*ast.CompositeLit)
						if !ok {
							continue invalidPackage
						}
						for _, elt := range lit.Elts {
							kv, ok := elt.(*ast.KeyValueExpr)
							if !ok {
								continue invalidPackage
							}
							key, ok := kv.Key.(*ast.BasicLit)
							if !ok || key.Kind != token.STRING {
								continue invalidPackage
							}
							value, ok := kv.Value.(*ast.BasicLit)
							if !ok || value.Kind != token.STRING {
								continue invalidPackage
							}
							ifaces[key.Value[1:len(key.Value)-1]] = value.Value[1:len(value.Value)-1]
						}
					}
				}
			}
		}

		m[pkg.Name] = &Package{
			Path:       pkg.PkgPath,
			Name:       pkg.Name,
			Interfaces: ifaces,
		}
	}

	out := map[string]*Package{}
	for _, pkg := range m {
		for iface := range pkg.Interfaces {
			if _, ok := out[iface]; ok {
				log.Fatalf("multiple packages define interface %q", iface)
			}
			out[iface] = pkg
		}
	}

	return out
}
