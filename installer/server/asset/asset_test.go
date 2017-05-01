package asset

import "testing"

func TestFind(t *testing.T) {
	assets := []Asset{
		asset{name: "A", data: []byte("foobar")},
		asset{name: "B", data: []byte("")},
	}

	// Empty list.
	_, err := Find([]Asset{}, "Null")
	if err == nil {
		t.Errorf("Find() should have failed on an empty list")
	}

	// Empty name.
	_, err = Find(assets, "")
	if err == nil {
		t.Errorf("Find() should have failed with an empty name")
	}

	// Non-existing Asset.
	_, err = Find(assets, "C")
	if err == nil {
		t.Errorf("Find() should have failed with a non-existing Asset")
	}

	// Existing Asset.
	asset, err := Find(assets, "A")
	if err != nil {
		t.Errorf("Find() shouldn't have failed with an existing Asset")
	}
	if asset.Name() != assets[0].Name() || string(asset.Data()) != string(assets[0].Data()) {
		t.Errorf("Find() didn't return the expected Asset")
	}
}

func TestReplace(t *testing.T) {
	assets := []Asset{asset{name: "A", data: []byte("foobar")}}

	// Replace a non-existing Asset should generate an error.
	b := asset{name: "B", data: []byte("foobar")}
	assetsA, err := Replace(assets, b)
	if err == nil || len(assetsA) == 2 {
		t.Errorf("Replace() on a non-existing Asset should generate an error")
	}

	// Typical Replace.
	a := asset{name: "A", data: []byte("foobar2")}
	assetsB, err := Replace(assets, a)
	if err != nil || len(assetsB) != 1 || string(assetsB[0].Data()) != string(a.data) {
		t.Errorf("Replace() should have replaced the Asset")
	}
}
