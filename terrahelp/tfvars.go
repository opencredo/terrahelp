package terrahelp

import (
	"io/ioutil"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/token"
	"sort"
	"strings"
)

// Replaceables defines the values which should be replaced as part of
// various transformations actions
type Replaceables interface {
	Values() ([]string, error)
}

type DefaultReplaceables struct {
	Vals []string
}

func (d *DefaultReplaceables) Values() ([]string, error) {
	return d.Vals, nil
}

// Tfvars provides utility functions pertaining to the
// terraform.tfvars file
type Tfvars struct {
	filename string
	excludeWhitespaceOnly bool
}

// NewTfVars creates a new Tfvars holder based on the provided filename
func NewTfVars(f string, excl bool) *Tfvars {
	return &Tfvars{filename: f, excludeWhitespaceOnly: excl}
}

// ExtractSensitiveVals returns a list of the sensitive values
// which were detected in the provided tfvars file
func (t *Tfvars) Values() ([]string, error) {
	// Read tfvars file
	b, err := ioutil.ReadFile(t.filename)
	if err != nil {
		return nil, err
	}

	// Parse it
	astFile, err := hcl.ParseBytes(b)
	if err != nil {
		return nil, err
	}

	// Find sensitive values (all quoted value strings)
	var vals []string
	ast.Walk(astFile, func(node ast.Node) (ast.Node, bool) {
		if node == nil {
			return node, false
		}

		switch n := node.(type) {
		case *ast.LiteralType:
			switch n.Token.Type {
			case token.STRING:
				v :=  n.Token.Value().(string)
				if v != "" && t.excludeWhitespaceOnly && strings.TrimSpace(v) != "" {
					vals = append(vals, n.Token.Value().(string))
				}
				if v != "" && !t.excludeWhitespaceOnly {
					vals = append(vals, n.Token.Value().(string))
				}
			}
		}

		return node, true
	})

	// Reverse in case there are overlaps
	sort.Strings(vals)
	for i := len(vals)/2 - 1; i >= 0; i-- {
		opp := len(vals) - 1 - i
		vals[i], vals[opp] = vals[opp], vals[i]
	}

	return vals, nil
}
