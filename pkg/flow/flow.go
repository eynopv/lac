package flow

import (
	"path/filepath"

	"github.com/eynopv/lac/pkg/expectation"
	"github.com/eynopv/lac/pkg/utils"
)

type Flow struct {
	FilePath string
	Items    []FlowItem
}

type FlowItem struct {
	Id        string                   `json:"id" yaml:"id"`
	Uses      string                   `json:"uses" yaml:"uses"`
	Variables map[string]string        `json:"variables" yaml:"variables"`
	Expect    *expectation.Expectation `json:"expect" yaml:"expect"`
}

func LoadFlow(flowPath string) (*Flow, error) {
	flow := Flow{
		FilePath: flowPath,
	}
	if err := utils.LoadItem(flowPath, &flow.Items); err != nil {
		return nil, err
	}
	flow.ResolveItemPaths()
	return &flow, nil
}

func (flow *Flow) ResolveItemPaths() {
	for i := range flow.Items {
		usesPath := filepath.Join(filepath.Dir(flow.FilePath), flow.Items[i].Uses)
		flow.Items[i].Uses = usesPath
	}
}
