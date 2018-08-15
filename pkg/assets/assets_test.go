package assets

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPutGet(t *testing.T) {
	ctx := context.Background()
	assets := &Assets{}

	assetA := Asset{
		Name: "a",
		Data: []byte("a-data"),
	}
	hashA, err := assets.Put(assetA)
	if err != nil {
		t.Fatal(err)
	}

	assetB := Asset{
		Name: "a",
		Data: []byte("b-data"),
	}
	hashB, err := assets.Put(assetB)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("GetByHash", func(t *testing.T) {
		for _, test := range []struct {
			hash  []byte
			asset *Asset
		}{
			{
				hash:  hashA,
				asset: &assetA,
			},
			{
				hash:  hashB,
				asset: &assetB,
			},
		} {
			t.Run(fmt.Sprintf("%x", test.hash), func(t *testing.T) {
				retrieved, err := assets.GetByHash(ctx, test.hash)
				if err != nil {
					t.Fatal(err)
				}

				if !reflect.DeepEqual(retrieved, *test.asset) {
					t.Errorf("%+v != %+v", retrieved, test.asset)
				}
			})
		}
	})

	t.Run("GetByName", func(t *testing.T) {
		retrieved, err := assets.GetByName(ctx, "a")
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(retrieved, assetB) {
			t.Errorf("%+v != %+v", retrieved, assetB)
		}
	})
}

func buildB(ctx context.Context, getByName GetByString) (asset *Asset, err error) {
	asset = &Asset{
		Name:          "b",
		RebuildHelper: buildB,
	}

	parents, err := asset.GetParents(ctx, getByName, "a/a")
	if err != nil {
		return nil, err
	}

	asset.Data = append(parents["a/a"].Data, []byte(", modified by b")...)
	return asset, nil
}

func buildC(ctx context.Context, getByName GetByString) (asset *Asset, err error) {
	asset = &Asset{
		Name:          "c",
		RebuildHelper: buildC,
	}

	parents, err := asset.GetParents(ctx, getByName, "b")
	if err != nil {
		return nil, err
	}

	asset.Data = append(parents["b"].Data, []byte(", modified by c")...)
	return asset, nil
}

func newAssets() *Assets {
	return &Assets{
		Root: Reference{
			Name: "c",
		},
		Rebuilders: map[string]Rebuild{
			"b": buildB,
			"c": buildC,
		},
	}
}

func defaultA(ctx context.Context, name string) (data []byte, err error) {
	if name == "a/a" {
		return []byte("a-data"), nil
	}
	return nil, os.ErrNotExist
}

func TestAssetsRead(t *testing.T) {
	ctx := context.Background()
	tempDir, err := ioutil.TempDir("", "openshift-install-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	t.Run("from scratch with a broken dependency", func(t *testing.T) {
		assets := newAssets()
		err = assets.Read(ctx, tempDir, nil, t)
		if err == nil {
			t.Fatal("unexpected success")
		}
		assert.Regexp(t, "^retrieve root: retrieve \"c\" by name: retrieve \"b\" by name: inject content for \"a/a\": cannot inject \"a/a\" without a file, default, or rebuilder: open .*a: no such file or directory$", err.Error())
	})

	var rootHash []byte
	t.Run("from scratch", func(t *testing.T) {
		assets := newAssets()
		err = assets.Read(ctx, tempDir, defaultA, t)
		if err != nil {
			t.Fatal(err)
		}

		rootHash = assets.Root.Hash
		assetC, err := assets.GetByHash(ctx, assets.Root.Hash)
		if err != nil {
			t.Fatal(err)
		}

		expected := "a-data, modified by b, modified by c"
		if string(assetC.Data) != expected {
			t.Fatalf("unexpected new asset C data: %q != %q", string(assetC.Data), expected)
		}

		refB := assetC.Parents[0]
		if refB.Name != "b" {
			t.Fatalf("asset %q has an unexpected parent name %q", assetC.Name, refB.Name)
		}

		assetB, err := assets.GetByHash(ctx, refB.Hash)
		if err != nil {
			t.Fatal(err)
		}

		expected = "a-data, modified by b"
		if string(assetB.Data) != expected {
			t.Fatalf("unexpected new asset %q data: %q != %q", refB.Name, string(assetB.Data), expected)
		}

		refA := assetB.Parents[0]
		if refA.Name != "a/a" {
			t.Fatalf("asset %q has an unexpected parent name %q", assetB.Name, refA.Name)
		}

		assetA, err := assets.GetByHash(ctx, refA.Hash)
		if err != nil {
			t.Fatal(err)
		}

		expected = "a-data"
		if string(assetA.Data) != expected {
			t.Fatalf("unexpected new asset %q data: %q != %q", refA.Name, string(assetA.Data), expected)
		}

		err = assets.Write(ctx, tempDir, false)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("read without edits", func(t *testing.T) {
		assets := newAssets()
		err = assets.Read(ctx, tempDir, defaultA, t)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, rootHash, assets.Root.Hash)
	})

	t.Run("read with leaf edit", func(t *testing.T) {
		err = ioutil.WriteFile(filepath.Join(tempDir, "a", "a"), []byte("edited a"), 0666)
		if err != nil {
			t.Fatal(err)
		}

		assets := newAssets()
		err = assets.Read(ctx, tempDir, defaultA, t)
		if err != nil {
			t.Fatal(err)
		}

		assert.NotEqual(t, rootHash, assets.Root.Hash)

		assetC, err := assets.GetByHash(ctx, assets.Root.Hash)
		if err != nil {
			t.Fatal(err)
		}

		expected := "edited a, modified by b, modified by c"
		if string(assetC.Data) != expected {
			t.Fatalf("unexpected new asset C data: %q != %q", string(assetC.Data), expected)
		}

		err = assets.Prune()
		if err != nil {
			t.Fatal(err)
		}

		_, err = assets.GetByHash(ctx, rootHash)
		if !os.IsNotExist(err) {
			t.Fatalf("can retrieve the original root asset by hash (%x) after pruning: %v", rootHash, err)
		}

		rootHash = assets.Root.Hash
		err = assets.Write(ctx, tempDir, false)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("reread after leaf edit", func(t *testing.T) {
		assets := newAssets()
		err = assets.Read(ctx, tempDir, defaultA, t)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, rootHash, assets.Root.Hash)
	})

	t.Run("read with branch edit", func(t *testing.T) {
		err = ioutil.WriteFile(filepath.Join(tempDir, "b"), []byte("edited b"), 0666)
		if err != nil {
			t.Fatal(err)
		}

		assets := newAssets()
		err = assets.Read(ctx, tempDir, defaultA, t)
		if err != nil {
			t.Fatal(err)
		}

		assetC, err := assets.GetByHash(ctx, assets.Root.Hash)
		if err != nil {
			t.Fatal(err)
		}

		expected := "edited b, modified by c"
		if string(assetC.Data) != expected {
			t.Fatalf("unexpected new asset C data: %q != %q", string(assetC.Data), expected)
		}

		assetB, err := assets.GetByHash(ctx, assetC.Parents[0].Hash)
		if err != nil {
			t.Fatal(err)
		}
		assert.Len(t, assetB.Parents, 0)

		rootHash = assets.Root.Hash
		err = assets.Write(ctx, tempDir, false)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("reread after branch edit", func(t *testing.T) {
		assets := newAssets()
		err = assets.Read(ctx, tempDir, defaultA, t)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, rootHash, assets.Root.Hash)
	})
}
