package flow

import (
	"fmt"
	"testing"
)

func TestItemPathResolution(t *testing.T) {
	flow := Flow{
		FilePath: "path/to/flow/file.json",
		Items: []FlowItem{
			{
				Uses: "../item.json",
			},
		},
	}
	flow.ResolveItemPaths()
	if flow.Items[0].Uses != "path/to/item.json" {
		t.Fatalf("Expected item path to be path/to/item.json: " + fmt.Sprint(flow.Items[0].Uses))
	}
}
