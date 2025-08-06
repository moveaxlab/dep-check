package rendering

import (
	"github.com/moveaxlab/dep-check/graph"
	"fmt"
)

const (
	plantUmlFileName    = "plantUml_graph.txt"
	plantUmlGraphHeader = "@startuml"
	plantUmlGraphFooter = "@enduml"
)

func PlantUmlNode(n *graph.Node) string {
	return fmt.Sprintf("object %s", n.Name)
}

func PlantUmlEdge(e *graph.Edge) string {
	return e.String()
}
