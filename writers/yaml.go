package writers

import "io"
import "github.com/go-yaml/yaml"
import "github.com/dadleyy/dannylister/tree"

func Yaml(output io.Writer, root tree.Node) (int, error) {
	buffer, err := yaml.Marshal(&root)

	if err != nil {
		return -1, err
	}

	return output.Write(buffer)
}
