package main

import "os"
import "fmt"
import "flag"
import "path/filepath"
import "encoding/json"

type optionsT struct {
	Help bool
	Recursive bool
	WorkingDir string
	Output string
}

type treeNode struct {
	os.FileInfo
	FullPath string
	Children []treeNode
	Link string
}

func (t *treeNode) Collect(recurse bool) error {
	var err error
	t.FileInfo, err = os.Stat(t.FullPath)

	if err != nil {
		return err
	}

	if t.IsDir() != true {
		return nil
	}

	file, err := os.Open(t.FullPath)

	if err != nil {
		return err
	}

	children, err := file.Readdir(0)

	if err != nil {
		return err
	}

	// now that we've collected all of the FileInfo objects of items in our directory,
	// we can loop over each creating new nodes and adding them into our children slice
	for _, info := range children {
		// build the full path to this child from our path and their name
		full := filepath.Join(t.FullPath, info.Name())

		// create the instance of the tree node
		child := treeNode{info, full, make([]treeNode, 0), ""}

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
		if recurse == true {
			err = child.Collect(true)
			if err != nil {
				return err
			}
		}

		// add this child to our list and move on
		t.Children = append(t.Children, child)
	}

	return nil
}

func main() {
	var options optionsT
	flag.BoolVar(&options.Help, "help", false, "display this help text")
	flag.BoolVar(&options.Recursive, "recursive", false, "when set, list files recursively. default: false")
	flag.StringVar(&options.WorkingDir, "path", "./", "specifies the working directory for the file scan")
	flag.StringVar(&options.Output, "output", "text", "specifies the type of format to output to stdout. options are: \"json\", \"yaml\" or \"text\". default is \"text\"")

	flag.Parse()

	if options.Help == true {
		flag.PrintDefaults()
		return
	}

	fmt.Printf("help[%t] recursive[%t] path[%s] format[%s]\n", options.Help, options.Recursive, options.WorkingDir, options.Output)

	info, err := os.Stat(options.WorkingDir)

	if err != nil {
		fmt.Printf("Error parsing working directory: %s\n", err.Error())
		return
	}

	if info.IsDir() != true {
		fmt.Printf("Path provided is not a directory\n")
		return
	}

	// create the start of the tree that we will build to assemble all of the files 
	// we encounter. we'll later send this into a writer to output the tree in the 
	// format prescribed by the user
	tree := treeNode{FullPath: options.WorkingDir}

	// launch our collection process
	err = tree.Collect(options.Recursive)

	if err != nil {
		fmt.Printf("Error while traversing the file sytem: %s\n", err.Error())
		return
	}

	output, err := json.Marshal(tree)

	if err != nil {
		fmt.Printf("Error while encoding tree as json: %s\n", err.Error())
		return
	}

	fmt.Println(string(output))
}
