package writers

import "io"
import "fmt"
import "strings"
import "github.com/dadleyy/dannylister/tree"

func marshalNode(node tree.Node, depth int) (output string) {
	output = node.Name()
	var nested []string


	if node.Link != "" {
		output += fmt.Sprintf("* (%s)", node.Link)
	} else if node.IsDir() {
		output += "/"
	}

	output += "\n"

	ident := ""

	for i := -1; i < depth; i++ {
		ident += "  "
	}

	for _, child := range node.Children {
		nested = append(nested, ident + marshalNode(child, depth+1))
	}

	output += strings.Join(nested, "")
	return
}

func Text(output io.Writer, root tree.Node) (int, error) {
	buffer := []byte(marshalNode(root, 0))
	return output.Write(buffer)
}
