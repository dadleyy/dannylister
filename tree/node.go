package tree

import "os"
import "path/filepath"

type Node struct {
	os.FileInfo
	FullPath string
	Children []Node
	Link string
}

func (t *Node) Collect(recurse bool) error {
	var err error
	t.FileInfo, err = os.Stat(t.FullPath)

	if err != nil {
		return err
	}

	// if this node is not a directory, there is nothing to collect.
	if t.IsDir() != true {
		return nil
	}

	// if it is a directory, we need to open it up to get all of the children
	file, err := os.Open(t.FullPath)

	if err != nil {
		return err
	}

	children, err := file.Readdir(0)


	if err != nil {
		return err
	}

	// cleanup the open file
	file.Close()

	// now that we've collected all of the FileInfo objects of items in our directory,
	// we can loop over each creating new nodes and adding them into our children slice
	for _, info := range children {
		// build the full path to this child from our path and their name
		full := filepath.Join(t.FullPath, info.Name())

		// create the instance of the tree node
		child := Node{info, full, make([]Node, 0), ""}

		// find out if this is a symbolic link
		if mode := child.Mode(); (mode & os.ModeSymlink) != 0 {
			dest, err := os.Readlink(full)

			// if we received an error here something strange happened; the os told us this
			// file had a symbolic link mode but was unable to resolve it to the destination
			if err != nil {
				return err
			}

			child.Link = dest

			t.Children = append(t.Children, child)
			continue
		}

		// if we were told to recurse, have the child collect all it's nodes as well
		if recurse != true || info.IsDir() != true {
			t.Children = append(t.Children, child)
			continue
		}

		err = child.Collect(true)
		if err != nil {
			return err
		}

		// add this child to our list and move on
		t.Children = append(t.Children, child)
	}

	return nil
}


