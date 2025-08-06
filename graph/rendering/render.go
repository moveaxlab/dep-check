package rendering

import (
	"github.com/moveaxlab/dep-check/graph"
	"fmt"
	"log"
	"os"
)

func RenderGraph(fileName, header, footer string, renderNode func(*graph.Node) string, renderEdge func(*graph.Edge) string) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalln("Error while create file :", file, err)
	}
	defer file.Close()

	_, err = fmt.Fprintln(file, header)
	if err != nil {
		log.Panicln("Error writing header :", header, "to file", file, err)
	}

	for _, n := range graph.Nodes {
		_, err = fmt.Fprintln(file, renderNode(n))
		if err != nil {
			log.Panicln("Error writing node :", renderNode(n), "to file", file, err)
		}
	}

	fmt.Fprintln(file)

	for _, e := range graph.Edges {
		_, err = fmt.Fprintln(file, renderEdge(e))
		if err != nil {
			log.Panicln("Error writing edge :", renderEdge(e), "to file", file, err)
		}
	}

	_, err = file.WriteString(footer)
	if err != nil {
		log.Panicln("Error writing footer :", footer, "to file", file, err)
	}
}
