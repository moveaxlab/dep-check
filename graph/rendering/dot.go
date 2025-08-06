package rendering

import (
	"github.com/moveaxlab/dep-check/graph"
	"fmt"
)

const (
	mermaidFileName    = "graph.mmd"
	mermaidGraphHeader = "flowchart LR"
	mermaidGraphFooter = ""
)

func DotNode(n *graph.Node) string {
	return fmt.Sprintf("\t\"%s\"", n.Name)
}

func DotEdge(e *graph.Edge) string {
	return fmt.Sprintf("\t\"%s\" -> \"%s\"", e.From, e.To)
}
