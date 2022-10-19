package fwschema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type NestedBlock struct {
	Block
}

// ApplyTerraform5AttributePathStep allows Blocks to be walked using
// tftypes.Walk and tftypes.Transform.
func (b NestedBlock) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	a, ok := step.(tftypes.AttributeName)

	if !ok {
		return nil, fmt.Errorf("can't apply %T to block", step)
	}

	attrName := string(a)

	if attr, ok := b.Block.GetAttributes()[attrName]; ok {
		return attr, nil
	}

	if block, ok := b.Block.GetBlocks()[attrName]; ok {
		return block, nil
	}

	return nil, fmt.Errorf("no attribute %q on Attributes or Blocks", a)
}
