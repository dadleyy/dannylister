package main

import "os"
import "fmt"
import "flag"
import "time"

import "github.com/dadleyy/dannylister/tree"
import "github.com/dadleyy/dannylister/writers"

type optionsT struct {
	Recursive bool
	WorkingDir string
	Output string
}

func main() {
	var options optionsT
	start := time.Now()
	flag.BoolVar(&options.Recursive, "recursive", false, "when set, list files recursively. default: false")
	flag.StringVar(&options.WorkingDir, "path", "./", "specifies the working directory for the file scan")
	flag.StringVar(&options.Output, "output", "text", "specifies the type of format to output to stdout. options are: \"json\", \"yaml\" or \"text\". default is \"text\"")

	flag.Parse()


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
	root := tree.Node{FullPath: options.WorkingDir}

	// launch our collection process
	err = root.Collect(options.Recursive)

	if err != nil {
		fmt.Printf("Error while traversing the file sytem: %s\n", err.Error())
		return
	}

	switch options.Output {
	case "yaml":
		writers.Yaml(os.Stdout, root)
	case "json":
		writers.Json(os.Stdout, root)
	case "text":
		writers.Text(os.Stdout, root)
	default:
		fmt.Printf("invalid format\n")
	}

	fmt.Printf("\ntook %f miliseconds\n", time.Since(start).Seconds() * 1000)
}
