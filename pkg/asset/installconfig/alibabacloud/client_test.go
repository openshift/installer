package alibabacloud

import (
	"fmt"
	"testing"
)

func TestDescribeRegions(t *testing.T) {
	client, err := NewClient(DefaultRegion)

	if err != nil {
		t.Errorf("Filed to create client")
	}

	_, err = client.DescribeRegions(DefaultRegion)
	if err != nil {
		t.Errorf(fmt.Sprintf("Filed to describe regions, Response: %s", _err))
	}
}

func TestListResourceGroups(t *testing.T) {
	client, err := NewClient(DefaultRegion)

	if err != nil {
		t.Errorf("Filed to create client")
	}

	_, err = client.ListResourceGroups()
	// fmt.Print(resp.ResourceGroups.ResourceGroup)
	if err != nil {
		t.Errorf(fmt.Sprintf("Filed to describe regions, Response: %s", _err))
	}
}
