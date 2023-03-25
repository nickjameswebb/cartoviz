package viz

import (
	"errors"
	"fmt"
	"os"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/vmware-tanzu/cartographer/pkg/apis/v1alpha1"
)

// TODO: return a byte buffer from this function, instead of rendering directly in function
// TODO: pass a logger to this function
// TODO: pass a context to this function
// TODO: hide visualization behind an interface so that we can test
// TODO: configurable filename, format
func VisualizeSupplyChain(supplyChain *v1alpha1.ClusterSupplyChain) error {
	if supplyChain == nil {
		return errors.New("supply chain nil")
	}

	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		return err
	}
	defer func() {
		if err := graph.Close(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		g.Close()
	}()

	var prevNode *cgraph.Node
	for _, resource := range supplyChain.Spec.Resources {
		node, err := graph.CreateNode(resource.Name)
		if err != nil {
			return err
		}

		if prevNode != nil {
			_, err := graph.CreateEdge("e", prevNode, node)
			if err != nil {
				return err
			}
		}

		prevNode = node
	}

	// e.SetLabel("e")

	return g.RenderFilename(graph, graphviz.PNG, "graph.png")
}
