package rendering

import (
	"github.com/moveaxlab/dep-check/graph"
	"fmt"
)

const (
	dotFileName    = "dot_graph.txt"
	dotGraphHeader = "digraph dependency {"
	dotGraphFooter = "}"
)

func mermaidShortNode(n *graph.Node) string {
	return fmt.Sprintf("\t%s", n.Name)
}

func mermaidLongNode(n *graph.Node) string {
	return fmt.Sprintf("%s[\"`%s<br>*%s*`\"]", mermaidShortNode(n), n.Name, n.Path)
}

func mermaidEdge(e *graph.Edge) string {
	return fmt.Sprintf("\t%s;", e)
}
