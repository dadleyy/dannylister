package tree

import "os"
import "testing"
import "path/filepath"

func TestCollectRecursive(t *testing.T) {
	dir, _ := os.Getwd()
	examples := filepath.Join(dir, "../example")

	start := Node{FullPath: examples}

	err := start.Collect(true)

	if err != nil {
		t.Fatal("unable to collect on example directory!")
		return
	}

	if start.Name() != "example" {
		t.Fatal("expected to have Name() return \"example\"")
		return
	}

	var linksource Node
	for _, child := range start.Children {
		if child.Name() == "linksource" {
			linksource = child
		}
	}

	if linksource.Name() != "linksource" {
		t.Fatal("unable to find linksource directory")
		return
	}

	if len(linksource.Children) == 0 {
		t.Fatal("linksource was missing it's children!")
		return
	}
}

func TestCollectNonRecursive(t *testing.T) {
	dir, _ := os.Getwd()
	examples := filepath.Join(dir, "../example")

	start := Node{FullPath: examples}

	err := start.Collect(false)

	if err != nil {
		t.Fatal("unable to collect on example directory!")
		return
	}

	if start.Name() != "example" {
		t.Fatal("expected to have Name() return \"example\"")
		return
	}

	var linksource Node
	for _, child := range start.Children {
		if child.Name() == "linksource" {
			linksource = child
		}
	}

	if linksource.Name() != "linksource" {
		t.Fatal("unable to find linksource directory")
		return
	}

	if len(linksource.Children) != 0 {
		t.Fatal("linksource had children!")
		return
	}
}
