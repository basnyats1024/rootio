package rootio

import (
	"testing"
)

func TestFlatTree(t *testing.T) {
	f, err := Open("testdata/small.root")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer f.Close()

	obj, err := f.Get("tree")
	if err != nil {
		t.Fatalf("%v", err)
	}

	key := obj.(*Key)
	if key.Name() != "tree" {
		t.Fatalf("key.Name: expected [tree] (got=%v)", key.Name())
	}

	tree := key.Value().(*Tree)
	if tree.Name() != "tree" {
		t.Fatalf("tree.Name: expected [tree] (got=%v)", tree.Name())
	}

	entries := tree.Entries()
	if entries != 100 {
		t.Fatalf("tree.Entries: expected [100] (got=%v)", entries)
	}

	if tree.totbytes != 40506 {
		t.Fatalf("tree.totbytes: expected [40506] (got=%v)", tree.totbytes)
	}

	if tree.zipbytes != 4184 {
		t.Fatalf("tree.zipbytes: expected [4184] (got=%v)", tree.zipbytes)
	}
}

// EOF