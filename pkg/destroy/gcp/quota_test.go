package gcp

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openshift/installer/pkg/types/gcp"
)

func TestMergeAllUsage(t *testing.T) {
	var testCases = []struct {
		name     string
		into     []gcp.QuotaUsage
		update   []gcp.QuotaUsage
		expected []gcp.QuotaUsage
	}{
		{
			name:     "no previous data, one new",
			into:     nil,
			update:   []gcp.QuotaUsage{{Metric: &gcp.Metric{Service: "s", Limit: "l"}, Amount: 1}},
			expected: []gcp.QuotaUsage{{Metric: &gcp.Metric{Service: "s", Limit: "l"}, Amount: 1}},
		},
		{
			name: "no previous data, many new",
			into: nil,
			update: []gcp.QuotaUsage{
				{Metric: &gcp.Metric{Service: "s", Limit: "l"}, Amount: 1},
				{Metric: &gcp.Metric{Service: "S", Limit: "L"}, Amount: 1},
			},
			expected: []gcp.QuotaUsage{
				{Metric: &gcp.Metric{Service: "s", Limit: "l"}, Amount: 1},
				{Metric: &gcp.Metric{Service: "S", Limit: "L"}, Amount: 1},
			},
		},
		{
			name: "merge adds new entry",
			into: []gcp.QuotaUsage{
				{Metric: &gcp.Metric{Service: "s", Limit: "l"}, Amount: 1},
			},
			update: []gcp.QuotaUsage{
				{Metric: &gcp.Metric{Service: "S", Limit: "L"}, Amount: 1},
			},
			expected: []gcp.QuotaUsage{
				{Metric: &gcp.Metric{Service: "s", Limit: "l"}, Amount: 1},
				{Metric: &gcp.Metric{Service: "S", Limit: "L"}, Amount: 1},
			},
		},
		{
			name: "merge updates current entry",
			into: []gcp.QuotaUsage{
				{Metric: &gcp.Metric{Service: "s", Limit: "l"}, Amount: 1},
			},
			update: []gcp.QuotaUsage{
				{Metric: &gcp.Metric{Service: "s", Limit: "l"}, Amount: 1},
			},
			expected: []gcp.QuotaUsage{
				{Metric: &gcp.Metric{Service: "s", Limit: "l"}, Amount: 2},
			},
		},
		{
			name: "merge updates current entry only with full match",
			into: []gcp.QuotaUsage{
				{Metric: &gcp.Metric{Service: "s", Limit: "l", Dimensions: map[string]string{"a": "b"}}, Amount: 1},
			},
			update: []gcp.QuotaUsage{
				{Metric: &gcp.Metric{Service: "s", Limit: "l"}, Amount: 1},                                          // no dimensions
				{Metric: &gcp.Metric{Service: "s", Limit: "L", Dimensions: map[string]string{"a": "b"}}, Amount: 1}, // different limit
				{Metric: &gcp.Metric{Service: "S", Limit: "l", Dimensions: map[string]string{"a": "b"}}, Amount: 1}, // different service
				{Metric: &gcp.Metric{Service: "s", Limit: "l", Dimensions: map[string]string{"a": "B"}}, Amount: 1}, // different dimensions
				{Metric: &gcp.Metric{Service: "s", Limit: "l", Dimensions: map[string]string{"a": "b"}}, Amount: 1}, // match
			},
			expected: []gcp.QuotaUsage{
				{Metric: &gcp.Metric{Service: "s", Limit: "l", Dimensions: map[string]string{"a": "b"}}, Amount: 2},
				{Metric: &gcp.Metric{Service: "s", Limit: "l"}, Amount: 1},
				{Metric: &gcp.Metric{Service: "s", Limit: "L", Dimensions: map[string]string{"a": "b"}}, Amount: 1},
				{Metric: &gcp.Metric{Service: "S", Limit: "l", Dimensions: map[string]string{"a": "b"}}, Amount: 1},
				{Metric: &gcp.Metric{Service: "s", Limit: "l", Dimensions: map[string]string{"a": "B"}}, Amount: 1},
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if diff := cmp.Diff(mergeAllUsage(testCase.into, testCase.update), testCase.expected, cmp.AllowUnexported(gcp.QuotaUsage{}, gcp.Metric{})); diff != "" {
				t.Errorf("%s: got incorrect usage after merge: %v", testCase.name, diff)
			}
		})
	}
}

func TestQuotaIntegration(t *testing.T) {
	uninstaller := &ClusterUninstaller{
		pendingItemTracker: newPendingItemTracker(),
	}
	uninstaller.insertPendingItems("fake", []cloudResource{{
		key:   "first",
		quota: nil,
	}, {
		key: "second",
		quota: []gcp.QuotaUsage{
			{Metric: &gcp.Metric{Service: "s", Limit: "l"}, Amount: 1},
			{Metric: &gcp.Metric{Service: "S", Limit: "L"}, Amount: 1},
		},
	}, {
		key: "third",
		quota: []gcp.QuotaUsage{
			{Metric: &gcp.Metric{Service: "s", Limit: "l"}, Amount: 123},
			{Metric: &gcp.Metric{Service: "S", Limit: "L"}, Amount: 111},
		},
	}})
	for _, item := range uninstaller.getPendingItems("fake") {
		uninstaller.deletePendingItems("fake", []cloudResource{item})
	}
	if diff := cmp.Diff(uninstaller.pendingItemTracker.removedQuota, []gcp.QuotaUsage{
		{Metric: &gcp.Metric{Service: "s", Limit: "l"}, Amount: 124},
		{Metric: &gcp.Metric{Service: "S", Limit: "L"}, Amount: 112},
	}, cmp.AllowUnexported(gcp.QuotaUsage{}, gcp.Metric{})); diff != "" {
		t.Errorf("didn't get correct removed quota: %v", diff)
	}
}
