package writers

import "io"
import "encoding/json"
import "github.com/dadleyy/dannylister/tree"

func Json(output io.Writer, root tree.Node) (int, error) {
	buffer, err := json.MarshalIndent(root.Reduce(), "", "  ")

	if err != nil {
		return -1, err
	}

	return output.Write(buffer)
}

