package machine

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vincent-petithory/dataurl"
)

type fileAssertion struct {
	path       string
	data       string
	additional func(*testing.T, map[string]interface{}) bool
}

// assertFilesInIgnitionConfig asserts that the specified ignition config
// contains exactly the files enumerated in fileAssertions.
func assertFilesInIgnitionConfig(
	t *testing.T,
	ignitionConfig []byte,
	fileAssertions ...fileAssertion,
) bool {
	var ic map[string]interface{}
	if err := json.Unmarshal(ignitionConfig, &ic); err != nil {
		return assert.NoError(t, err, "unexpected error unmarshaling ignition config")
	}
	storage, ok := ic["storage"]
	if !assert.True(t, ok, "No storage in ignition config") {
		return false
	}
	files, ok := storage.(map[string]interface{})["files"]
	if !assert.True(t, ok, "No files in ignition config") {
		return false
	}
	expectedFilePaths := make([]string, len(fileAssertions))
	for i, a := range fileAssertions {
		expectedFilePaths[i] = a.path
	}
	filesList := files.([]interface{})
	actualFilePaths := make([]string, len(filesList))
	for i, f := range filesList {
		path, ok := f.(map[string]interface{})["path"]
		if !assert.True(t, ok, "file has no path: %+v", f) {
			return false
		}
		actualFilePaths[i] = path.(string)
	}
	if !assert.Equal(t, expectedFilePaths, actualFilePaths, "Unexpected file paths") {
		return false
	}
	for _, f := range filesList {
		file := f.(map[string]interface{})
		path := file["path"]
		var fa fileAssertion
		for _, a := range fileAssertions {
			if a.path != path {
				continue
			}
			fa = a
		}
		contents, ok := file["contents"]
		if !assert.True(t, ok, "file %q has no contents", path) {
			return false
		}
		source, ok := contents.(map[string]interface{})["source"]
		if !assert.True(t, ok, "file %q has no source", path) {
			return false
		}
		url, err := dataurl.DecodeString(source.(string))
		if !assert.NoError(t, err, "unexpected error decoding dataurl in file %q", path) {
			return false
		}
		if !assert.Equal(t, fa.data, string(url.Data), "unexpected data in file %q", path) {
			return false
		}
		if fa.additional != nil {
			if !fa.additional(t, file) {
				return false
			}
		}
	}
	return true
}

type systemdUnitAssertion struct {
	name       string
	dropinName string
	contents   string
	additional func(*testing.T, map[string]interface{}) bool
}

// assertSystemdUnitsInIgnitionConfig asserts that the specified ignition config
// contains exactly the systemd units enumerated in unitAssertions.
func assertSystemdUnitsInIgnitionConfig(
	t *testing.T,
	ignitionConfig []byte,
	unitAssertions ...systemdUnitAssertion,
) bool {
	var ic map[string]interface{}
	if err := json.Unmarshal(ignitionConfig, &ic); err != nil {
		return assert.NoError(t, err, "unexpected error unmarshaling ignition config")
	}
	systemd, ok := ic["systemd"]
	if !assert.True(t, ok, "No systemd in ignition config") {
		return false
	}
	units, ok := systemd.(map[string]interface{})["units"]
	if !assert.True(t, ok, "No units in ignition config") {
		return false
	}
	expectedUnitNames := make([]string, len(unitAssertions))
	for i, a := range unitAssertions {
		expectedUnitNames[i] = a.name
	}
	unitsList := units.([]interface{})
	actualUnitNames := make([]string, len(unitsList))
	for i, u := range unitsList {
		name, ok := u.(map[string]interface{})["name"]
		if !assert.True(t, ok, "unit has no name: %+v", u) {
			return false
		}
		actualUnitNames[i] = name.(string)
	}
	if !assert.Equal(t, expectedUnitNames, actualUnitNames, "Unexpected unit names") {
		return false
	}
	for _, u := range unitsList {
		unit := u.(map[string]interface{})
		name := unit["name"]
		var ua systemdUnitAssertion
		for _, a := range unitAssertions {
			if a.name != name {
				continue
			}
			ua = a
		}
		contentsParent := unit
		if ua.dropinName != "" {
			dropins, ok := unit["dropins"]
			if !assert.True(t, ok, "no dropins in systemd unit %q", name) {
				return false
			}
			dropinsList := dropins.([]interface{})
			if !assert.Equal(t, 1, len(dropinsList), "unexpected number of dropins in systemd unit %q", name) {
				return false
			}
			dropin := dropinsList[0].(map[string]interface{})
			dropinName, ok := dropin["name"]
			if !assert.True(t, ok, "no name in dropin in systemd unit %q", name) {
				return false
			}
			if !assert.Equal(t, ua.dropinName, dropinName.(string), "unexpected dropin name in systemd unit %q", name) {
				return false
			}
			contentsParent = dropin
		}
		contents, contentsOK := contentsParent["contents"]
		if ua.contents != "" {
			if !assert.True(t, contentsOK, "no contents in systemd unit %q", name) {
				return false
			}
			if !assert.Equal(t, ua.contents, contents.(string), "unexpected contents in systemd unit %q", name) {
				return false
			}
		} else {
			if !assert.False(t, contentsOK, "unexpected contents in systemd unit %q", name) {
				return false
			}
		}
		if ua.additional != nil {
			if !ua.additional(t, unit) {
				return false
			}
		}
	}
	return true
}

type userAssertion struct {
	name       string
	sshKey     string
	additional func(*testing.T, map[string]interface{}) bool
}

// assertUsersInIgnitionConfig asserts that the specified ignition config
// contains exactly the users enumerated in userAssertions.
func assertUsersInIgnitionConfig(
	t *testing.T,
	ignitionConfig []byte,
	userAssertions ...userAssertion,
) bool {
	var ic map[string]interface{}
	if err := json.Unmarshal(ignitionConfig, &ic); err != nil {
		return assert.NoError(t, err, "unexpected error unmarshaling ignition config")
	}
	passwd, ok := ic["passwd"]
	if !assert.True(t, ok, "No passwd in ignition config") {
		return false
	}
	users, ok := passwd.(map[string]interface{})["users"]
	if !assert.True(t, ok, "No users in ignition config") {
		return false
	}
	expectedUserNames := make([]string, len(userAssertions))
	for i, a := range userAssertions {
		expectedUserNames[i] = a.name
	}
	usersList := users.([]interface{})
	actualUserNames := make([]string, len(usersList))
	for i, u := range usersList {
		name, ok := u.(map[string]interface{})["name"]
		if !assert.True(t, ok, "user has no name: %+v", u) {
			return false
		}
		actualUserNames[i] = name.(string)
	}
	if !assert.Equal(t, expectedUserNames, actualUserNames, "Unexpected user names") {
		return false
	}
	for _, u := range usersList {
		user := u.(map[string]interface{})
		name := user["name"]
		var ua userAssertion
		for _, a := range userAssertions {
			if a.name != name {
				continue
			}
			ua = a
		}
		sshAuthorizedKeys, ok := user["sshAuthorizedKeys"]
		if !assert.True(t, ok, "no sshAuthorizedKeys in user %q", name) {
			return false
		}
		sshAuthorizedKeysList := sshAuthorizedKeys.([]interface{})
		if !assert.Equal(t, 1, len(sshAuthorizedKeysList), "unexpected number of sshAuthorizedKeys in user %q", name) {
			return false
		}
		sshAuthorizedKey := sshAuthorizedKeysList[0].(string)
		if !assert.Equal(t, ua.sshKey, sshAuthorizedKey, "unexpected ssh key in user %q", name) {
			return false
		}
		if ua.additional != nil {
			if !ua.additional(t, user) {
				return false
			}
		}
	}
	return true
}
