package terrahelp

import (
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/hcl"
)

// Tfvars provides utility functions pertaining to the
// terraform.tfvars file
type Tfvars struct{}

// ExtractSensitiveVals returns a list of the sensitive values
// which were detected in the provided tfvars file
func (t *Tfvars) ExtractSensitiveVals(f string) ([]string, error) {
	// Read tfvars file
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	// Parse it
	astFile, err := hcl.ParseBytes(b)
	if err != nil {
		return nil, err
	}

	// Extract kv pairs
	var res map[string]string
	if err := hcl.DecodeObject(&res, astFile); err != nil {
		return nil, fmt.Errorf(
			"Error occured decoding Terraform vars file: %s\n\n"+
				"tfvars files are expected to only contain `key = \"value\"` type config.\n",
			err)
	}

	// Extract just the values
	var vals []string
	for _, v := range res {
		if v != "" {
			vals = append(vals, v)
		}
	}

	return vals, nil
}
