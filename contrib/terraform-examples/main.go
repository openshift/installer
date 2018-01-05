package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/token"
)

type Variable struct {
	Name        string
	Description string
	Default     string
}

type Index struct {
	Variables []Variable
}

type ByName []Variable

func (s ByName) Len() int           { return len(s) }
func (s ByName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByName) Less(i, j int) bool { return s[i].Name < s[j].Name }

type stateFunc func(n ast.Node, v *Variable) stateFunc

func newIndex(c []byte) (*Index, error) {
	file, err := hcl.ParseBytes(c)
	if err != nil {
		return nil, err
	}

	v := &Variable{}
	idx := &Index{}
	state := idx.content

	ast.Walk(file.Node, func(n ast.Node) (ast.Node, bool) {
		state = state(n, v)
		return n, true
	})

	idx.Variables = append(idx.Variables, *v)
	idx.Variables = idx.Variables[1:]

	return idx, nil
}

func (i *Index) content(n ast.Node, v *Variable) stateFunc {
	key, ok := n.(*ast.ObjectKey)
	if !ok {
		return i.content
	}

	if key.Token.Type != token.IDENT {
		return i.content
	}

	switch key.Token.Text {
	case "default":
		return i.defaultValue
	case "description":
		return i.description
	case "variable":
		i.Variables = append(i.Variables, *v)
		*v = Variable{}
		return i.variableName
	}

	return i.content
}

func (i *Index) variableName(n ast.Node, v *Variable) stateFunc {
	key, ok := n.(*ast.ObjectKey)
	if !ok {
		return i.variableName
	}

	if key.Token.Type != token.STRING {
		return i.variableName
	}

	v.Name = key.Token.Text

	return i.content
}

func (i *Index) defaultValue(n ast.Node, v *Variable) stateFunc {
	if _, ok := n.(*ast.ListType); ok {
		return i.content
	}

	if _, ok := n.(*ast.ObjectType); ok {
		return i.content
	}

	key, ok := n.(*ast.LiteralType)
	if !ok {
		return i.defaultValue
	}

	switch key.Token.Type {
	case token.BOOL:
		v.Default = key.Token.Text
		return i.content
	case token.STRING:
		v.Default = key.Token.Text
		return i.content
	}

	return i.defaultValue
}

func (i *Index) description(n ast.Node, v *Variable) stateFunc {
	key, ok := n.(*ast.LiteralType)
	if !ok {
		return i.description
	}

	switch key.Token.Type {
	case token.HEREDOC:
		v.Description = key.Token.Value().(string)
		return i.content
	case token.STRING:
		v.Description = key.Token.Value().(string)
		return i.content
	}

	return i.description
}

func contents(args []string) ([]byte, error) {
	if args[0] == "-" {
		if len(args) > 1 {
			return nil, errors.New("invalid stdin qualifier: multiple files provided")
		}

		return ioutil.ReadAll(os.Stdin)
	}

	var contents []byte

	for _, path := range flag.Args() {
		c, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}

		contents = append(contents, c...)
		contents = append(contents, []byte("\n")...)
	}

	return contents, nil
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: terraform-example [-|<FILE>...]\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	var c []byte
	c, err := contents(flag.Args())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	idx, err := newIndex(c)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	sort.Sort(ByName(idx.Variables))

	for _, v := range idx.Variables {
		if strings.HasPrefix(v.Description, "(internal)") {
			continue
		}

		fmt.Printf("\n")

		v.Name = strings.Trim(v.Name, `"`)

		if strings.HasPrefix(v.Description, "(optional)") {
			v.Name = `// ` + v.Name
		}

		for _, l := range strings.Split(strings.TrimSpace(v.Description), "\n") {
			fmt.Printf("// %s\n", l)
		}

		if v.Default == "" {
			fmt.Printf("%s = \"\"\n", v.Name)
		} else {
			fmt.Printf("%s = %s\n", v.Name, v.Default)
		}
	}
}
