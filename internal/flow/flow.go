package flow

import (
	"path/filepath"

	"github.com/eynopv/lac/internal/expectation"
	"github.com/eynopv/lac/internal/utils"
)

type Flow struct {
	FilePath string
	Items    []FlowItem
}

type FlowItem struct {
	Id     string                   `json:"id" yaml:"id"`
	Uses   string                   `json:"uses" yaml:"uses"`
	With   map[string]string        `json:"with" yaml:"with"`
	Expect *expectation.Expectation `json:"expect" yaml:"expect"`
}

func LoadFlow(flowPath string) (*Flow, error) {
	flow := Flow{
		FilePath: flowPath,
	}
	if err := utils.LoadItem(flowPath, &flow.Items); err != nil {
		return nil, err
	}
	return &flow, nil
}

func (flow *Flow) ResolveItemPaths() {
	for i := range flow.Items {
		usesPath := filepath.Join(filepath.Dir(flow.FilePath), flow.Items[i].Uses)
		flow.Items[i].Uses = usesPath
	}
}
