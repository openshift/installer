package azure

import (
	"strconv"
	"strings"
	"testing"

	"golang.org/x/mod/semver"
)

func FuzzConvertToSemver(f *testing.F) {
	f.Add("4.17.2024101109")
	f.Add("418.94.20250122")
	f.Add("9.4.20250122")
	f.Add("10.0.1")
	f.Add("418")
	f.Add("418.94")
	f.Add("4.17")
	f.Add("9.4")
	f.Add("9.")
	f.Add("100.2.3")
	f.Add("010..10")
	f.Add("010..0.\n")
	f.Add("")

	expectedConversions := map[string]string{
		"4.17.2024101109": "v4.17.2024101109",
		"418.94.20250122": "v4.18.20250122",
		"9.4.20250122":    "v9.4.20250122",
		"10.0.1":          "v10.0.1",
		"418":             "",
		"418.94":          "",
		"4.17":            "v4.17",
		"9.4":             "v9.4",
		"9.":              "",
		"100.2.3":         "",
		"010..10":         "v0.10.10",
		"010..0.\n":       "v0.10.0",
		"":                "",
	}

	f.Fuzz(func(t *testing.T, ver string) {
		got := convertToSemver(ver)
		if got != "" && !semver.IsValid(got) {
			t.Fatalf("convertToSemver(%q) = %q, want a valid semantic version or empty", ver, got)
		}
		if want, ok := expectedConversions[ver]; ok && got != want {
			t.Fatalf("convertToSemver(%q) = %q, want %q", ver, got, want)
		}
	})
}

func FuzzParseAROSKUs(f *testing.F) {
	f.Add("4.22", x86)
	f.Add("4.22", arm64)
	f.Add("5.0", x86)
	f.Add("5.0", arm64)
	f.Add("5.0", "ppc64le")
	f.Add("not-a-version", x86)
	f.Add("", x86)

	f.Fuzz(func(t *testing.T, release, arch string) {
		gen1, gen2, err := parseAROSKUs(release, arch)
		major, minor, _ := strings.Cut(release, ".")
		majorInt, convErr := strconv.Atoi(major)
		if convErr != nil {
			if err == nil || gen1 != "" || gen2 != "" {
				t.Fatalf("parseAROSKUs(%q, %q) = %q, %q, %v", release, arch, gen1, gen2, err)
			}
			return
		}
		if err != nil {
			t.Fatalf("parseAROSKUs(%q, %q) unexpected error %v", release, arch, err)
		}
		xy := strings.ReplaceAll(release, ".", "")
		var want1, want2 string
		if majorInt >= 5 {
			switch arch {
			case x86:
				want2 = "aro_" + major + "-" + minor + "_x64_gen2"
			case arm64:
				want2 = "aro_" + major + "-" + minor + "_arm_gen2"
			}
		} else {
			switch arch {
			case x86:
				want1 = "aro_" + xy
				want2 = "aro_" + xy + "-v2"
			case arm64:
				want2 = "aro_" + xy + "-arm"
			}
		}
		if gen1 != want1 || gen2 != want2 {
			t.Fatalf("parseAROSKUs(%q, %q) = %q, %q, want %q, %q", release, arch, gen1, gen2, want1, want2)
		}
	})
}

func FuzzMatchesRHELMajor(f *testing.F) {
	f.Add("9.4.20250122", 9)
	f.Add("10.0.1", 10)
	f.Add("4.18.1", 9)
	f.Add("418.94.1", 9)
	f.Add("not-a-version", 9)
	f.Add("9.1", 10)

	f.Fuzz(func(t *testing.T, ver string, rhelMajor int) {
		got := matchesRHELMajor(ver, rhelMajor)
		major, _, _ := strings.Cut(ver, ".")
		verMajor, err := strconv.Atoi(major)
		want := rhelMajor == 9
		if err == nil && (verMajor == 9 || verMajor == 10) {
			want = verMajor == rhelMajor
		}
		if got != want {
			t.Fatalf("matchesRHELMajor(%q, %d) = %t, want %t", ver, rhelMajor, got, want)
		}
	})
}

func FuzzCheckIfNewer(f *testing.F) {
	f.Add("v4.18.1", "4.17")
	f.Add("v4.17.1", "4.18")
	f.Add("v9.4.1", "4.18")
	f.Add("not-semver", "4.18")
	f.Add("v4.18.1", "4")

	f.Fuzz(func(t *testing.T, candidate, release string) {
		got := checkIfNewer(candidate, release)
		if candidateMajorIsRHEL(candidate) {
			if got {
				t.Fatalf("checkIfNewer(%q, %q) = true, want false for a RHEL-numbered image", candidate, release)
			}
			return
		}
		if !strings.Contains(release, ".") && !got {
			t.Fatalf("checkIfNewer(%q, %q) = false, want true for a release without a minor version", candidate, release)
		}
	})
}

func candidateMajorIsRHEL(candidate string) bool {
	mm := semver.MajorMinor(candidate)
	if mm == "" {
		return false
	}
	major, err := strconv.Atoi(strings.TrimPrefix(strings.Split(mm, ".")[0], "v"))
	if err != nil {
		return false
	}
	return major == 9 || major == 10
}
