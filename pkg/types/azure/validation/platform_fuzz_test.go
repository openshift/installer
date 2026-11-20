package validation

import (
	"fmt"
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/util/validation/field"
)

const (
	tagKeyInvalidMsg   = "key is invalid or contains invalid characters: key can have a maximum of 128 characters, cannot be empty and must begin with a letter, end with a letter, number or underscore, and must contain only alphanumeric characters and the following special characters `_ . -`"
	tagValueInvalidMsg = "value is invalid or contains invalid characters: value can have a maximum of 256 characters, cannot be empty and must contain only alphanumeric characters and the following special characters `_ + , - . / : ; < = > ? @`"
	tagPrefixMsg       = "key contains restricted prefix"

	// tagValueSpecials is the value character set named by validateTag, aside from ASCII letters and digits.
	tagValueSpecials = "_+,.-/:;<=>?@"
)

func FuzzValidateTag(f *testing.F) {
	f.Add("key1", "value1")
	f.Add("", "")
	f.Add(strings.Repeat("a", 129), "value")
	f.Add(strings.Repeat("a", 128), strings.Repeat("b", 256))
	f.Add("a", strings.Repeat("b", 257))
	f.Add("name", "value")
	f.Add("name123", "value")
	f.Add("kubernetes.io_cluster", "value")
	f.Add("azure", "microsoft")
	f.Add("for_openshift.io", "azure")
	f.Add("key@", "1value")
	f.Add("key", "")
	f.Add("key", "value*^%")
	f.Add("_key", "value")
	f.Add("key_", "value_2")
	f.Add("Key.-_7", "value@7")

	f.Fuzz(func(t *testing.T, key, value string) {
		err := validateTag(key, value)
		want := expectedTagError(key, value)
		if want == "" {
			if err != nil {
				t.Fatalf("validateTag(%q, %q) = %v, want accept", key, value, err)
			}
			return
		}
		if err == nil {
			t.Fatalf("validateTag(%q, %q) accepted, want %q", key, value, want)
		}
		if err.Error() != want {
			t.Fatalf("validateTag(%q, %q) = %q, want %q", key, value, err.Error(), want)
		}
	})
}

func FuzzFindDuplicateTagKeys(f *testing.F) {
	f.Add("environment", "Environment", "key")
	f.Add("key1", "key2", "key3")
	f.Add("Key", "key", "KEY")
	f.Add("", "a", "A")
	f.Add("name", "name123", "resourcename")

	f.Fuzz(func(t *testing.T, k1, k2, k3 string) {
		tags := map[string]string{k1: "v1", k2: "v2", k3: "v3"}
		assertDuplicateKeys(t, tags, findDuplicateTagKeys(tags))
	})
}

func FuzzValidateUserTags(f *testing.F) {
	f.Add("key1", "value1", "key_2", "value_2", "key.3", "value.3", false)
	f.Add("", "value", "key", "", "1key", "1value", false)
	f.Add("name", "value", "name123", "value", "kubernetes.io_cluster", "value", false)
	f.Add("azure", "microsoft", "environment", "test", "Environment", "lab", false)
	f.Add(strings.Repeat("a", 129), "value", "key", strings.Repeat("b", 257), "key@", "value*^%", false)
	f.Add("key", "value", "key2", "value2", "key3", "value3", true)

	f.Fuzz(func(t *testing.T, k1, v1, k2, v2, k3, v3 string, exceedLimit bool) {
		fldPath := field.NewPath("userTags")

		if errs := validateUserTags(nil, fldPath); len(errs) != 0 {
			t.Fatalf("nil tags: got %d errors, want 0", len(errs))
		}
		if errs := validateUserTags(map[string]string{}, fldPath); len(errs) != 0 {
			t.Fatalf("empty tags: got %d errors, want 0", len(errs))
		}

		tags := map[string]string{k1: v1, k2: v2, k3: v3}
		if exceedLimit {
			for i := 0; len(tags) <= maxUserTagLimit; i++ {
				tags[fmt.Sprintf("fuzzpad%04d", i)] = "v"
			}
		}
		assertUserTagErrors(t, tags, validateUserTags(tags, fldPath))
	})
}

// tagKeyShapeOK reports whether key meets the length and character rules, ignoring restricted prefixes.
func tagKeyShapeOK(key string) bool {
	if len(key) < 1 || len(key) > 128 {
		return false
	}
	if !isASCIILetter(key[0]) {
		return false
	}
	if len(key) == 1 {
		return true
	}
	if !isTagKeyTail(key[len(key)-1]) {
		return false
	}
	for i := 1; i < len(key)-1; i++ {
		if !isTagKeyMiddle(key[i]) {
			return false
		}
	}
	return true
}

func tagKeyRestricted(key string) bool {
	if strings.EqualFold(key, "name") {
		return true
	}
	folded := strings.ToLower(key)
	for _, prefix := range []string{"kubernetes.io", "openshift.io", "microsoft", "azure", "windows"} {
		if strings.HasPrefix(folded, prefix) {
			return true
		}
	}
	return false
}

func tagValueAccepted(value string) bool {
	if len(value) < 1 || len(value) > 256 {
		return false
	}
	for i := 0; i < len(value); i++ {
		c := value[i]
		if isASCIILetter(c) || isDigit(c) || strings.ContainsRune(tagValueSpecials, rune(c)) {
			continue
		}
		return false
	}
	return true
}

func expectedTagError(key, value string) string {
	if !tagKeyShapeOK(key) {
		return tagKeyInvalidMsg
	}
	if !tagValueAccepted(value) {
		return tagValueInvalidMsg
	}
	if tagKeyRestricted(key) {
		return tagPrefixMsg
	}
	return ""
}

func assertDuplicateKeys(t *testing.T, tags map[string]string, err error) {
	t.Helper()
	if err != nil && !strings.Contains(err.Error(), "found duplicate tag keys") {
		t.Fatalf("duplicate error %q does not name duplicate tag keys", err)
	}

	var shaped []string
	allShaped := true
	for key := range tags {
		if !tagKeyShapeOK(key) {
			allShaped = false
			continue
		}
		shaped = append(shaped, key)
	}
	foldDup := hasEqualFoldPair(shaped)
	if foldDup && err == nil {
		t.Fatalf("tags %v have case-insensitive duplicate keys, got nil error", tags)
	}
	if allShaped && !foldDup && err != nil {
		t.Fatalf("tags %v have unique keys, got %v", tags, err)
	}
}

func assertUserTagErrors(t *testing.T, tags map[string]string, errs field.ErrorList) {
	t.Helper()

	var tooMany, forbidden, invalid int
	for _, err := range errs {
		if err == nil {
			t.Fatal("validateUserTags returned a nil error")
		}
		switch err.Type {
		case field.ErrorTypeTooMany:
			tooMany++
		case field.ErrorTypeForbidden:
			forbidden++
		case field.ErrorTypeInvalid:
			invalid++
		default:
			t.Fatalf("unexpected error type %q: %v", err.Type, err)
		}
	}
	if tooMany+forbidden+invalid != len(errs) {
		t.Fatalf("error count %d does not match classified errors", len(errs))
	}
	if tooMany > 1 || forbidden > 1 || invalid > len(tags) {
		t.Fatalf("errors exceed bounds: tooMany=%d forbidden=%d invalid=%d tags=%d", tooMany, forbidden, invalid, len(tags))
	}

	wantTooMany := 0
	if len(tags) > maxUserTagLimit {
		wantTooMany = 1
	}
	if tooMany != wantTooMany {
		t.Fatalf("tooMany=%d, want %d for %d tags", tooMany, wantTooMany, len(tags))
	}

	var invalidTags int
	shaped := make([]string, 0, len(tags))
	allShaped := true
	for key, value := range tags {
		if expectedTagError(key, value) != "" {
			invalidTags++
		}
		if !tagKeyShapeOK(key) {
			allShaped = false
			continue
		}
		shaped = append(shaped, key)
	}
	if invalid != invalidTags {
		t.Fatalf("invalid=%d, want %d for tags %v (%v)", invalid, invalidTags, tags, errs)
	}

	foldDup := hasEqualFoldPair(shaped)
	if foldDup && forbidden != 1 {
		t.Fatalf("tags %v have case-insensitive duplicate keys, forbidden=%d", tags, forbidden)
	}
	if allShaped && !foldDup && forbidden != 0 {
		t.Fatalf("tags %v have unique keys, forbidden=%d (%v)", tags, forbidden, errs)
	}
}

func hasEqualFoldPair(keys []string) bool {
	for i := range keys {
		for j := i + 1; j < len(keys); j++ {
			if strings.EqualFold(keys[i], keys[j]) {
				return true
			}
		}
	}
	return false
}

func isASCIILetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isTagKeyTail(c byte) bool {
	return isASCIILetter(c) || isDigit(c) || c == '_'
}

func isTagKeyMiddle(c byte) bool {
	return isTagKeyTail(c) || c == '.' || c == '-'
}
