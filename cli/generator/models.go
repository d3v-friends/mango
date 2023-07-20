package generator

import "github.com/spf13/cobra"

type (
	fmtModel struct {
		Path    fPath               `yaml:"path" json:"path"`
		Models  map[goStruct]fModel `yaml:"models" json:"models"`
		Imports []string            `yaml:"imports" json:"imports"`
	}

	fPath struct {
		Def string `yaml:"def" json:"def"`
		Out string `yaml:"out" json:"out"`
	}

	fModel struct {
		Model    string             `yaml:"model" json:"model"`
		Fields   map[goField]goType `yaml:"fields" json:"fields"`
		Datetime bool               `yaml:"datetime" json:"datetime"`
	}

	goStruct string
	goField  string
	goType   string
)

func NewModelCmd() (res *cobra.Command) {
	panic("not impl")
}
