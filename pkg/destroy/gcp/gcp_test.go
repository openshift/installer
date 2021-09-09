package gcp

import "testing"

func TestGetNameFromURL(t *testing.T) {
	var testCases = []struct {
		item, url, expected string
	}{
		{
			item:     "zones",
			url:      "https://www.googleapis.com/compute/v1/projects/ci-op-lk2ifbjc/zones/us-central1-a",
			expected: "us-central1-a",
		},
		{
			item:     "networks",
			url:      "https://www.googleapis.com/compute/v1/projects/ci-op-lk2ifbjc/global/networks/ci-op-lk2ifbjc-15937-q68kj-network",
			expected: "ci-op-lk2ifbjc-15937-q68kj-network",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.item, func(t *testing.T) {
			if actual, expected := getNameFromURL(testCase.item, testCase.url), testCase.expected; actual != expected {
				t.Errorf("got incorrect name %s for item %s from url %s", actual, testCase.item, testCase.url)
			}
		})
	}
}

func TestGetRegionFromZone(t *testing.T) {
	if actual, expected := getRegionFromZone("us-central1-a"), "us-central1"; actual != expected {
		t.Errorf("got %s, not %s", actual, expected)
	}
}

func TestGetDiskLimit(t *testing.T) {
	var testCases = []struct {
		url, expected string
	}{
		{
			url:      "https://www.googleapis.com/compute/v1/projects/ci-op-lk2ifbjc/zones/us-central1-a/diskTypes/pd-standard",
			expected: "disks_total_storage",
		},
		{
			url:      "https://www.googleapis.com/compute/v1/projects/ci-op-lk2ifbjc/zones/us-central1-a/diskTypes/pd-ssd",
			expected: "ssd_total_storage",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.url, func(t *testing.T) {
			if actual, expected := getDiskLimit(testCase.url), testCase.expected; actual != expected {
				t.Errorf("got incorrect limit %s for url %s", actual, testCase.url)
			}
		})
	}
}
