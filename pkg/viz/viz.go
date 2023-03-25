package viz

import (
	"errors"
	"fmt"
	"os"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/vmware-tanzu/cartographer/pkg/apis/v1alpha1"
	"golang.org/x/exp/maps"
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

	err = visualize(graph, supplyChain)
	if err != nil {
		return err
	}

	return g.RenderFilename(graph, graphviz.PNG, "graph.png")
}

func visualize(graph *cgraph.Graph, supplyChain *v1alpha1.ClusterSupplyChain) error {
	var prevNode *cgraph.Node
	var prevResource *v1alpha1.SupplyChainResource

	for _, resource := range supplyChain.Spec.Resources {
		node, err := graph.CreateNode(resource.Name)
		if err != nil {
			return err
		}

		if prevNode != nil {
			edgeName, err := edgeNameFromResource(prevResource)
			if err != nil {
				return fmt.Errorf("no edge name for resource: %w", err)
			}
			edge, err := graph.CreateEdge(edgeName, prevNode, node)
			if err != nil {
				return err
			}
			edge.SetLabel(edgeName)
		}

		prevNode = node
		prevResource = resource.DeepCopy()
	}
	// e.SetLabel("e")
	return nil
}

func edgeNameFromResource(resource *v1alpha1.SupplyChainResource) (string, error) {
	if resource == nil {
		return "", errors.New("resource is nil")
	}

	templateKindToEdgeNameMapping := map[string]string{
		"ClusterImageTemplate":  "image",
		"ClusterSourceTemplate": "source",
		"ClusterConfigTemplate": "config",
		"ClusterTemplate":       "",
	}

	templateKind := resource.TemplateRef.Kind

	if !contains(maps.Keys(templateKindToEdgeNameMapping), templateKind) {
		return "", errors.New("invalid template kind")
	}

	return templateKindToEdgeNameMapping[templateKind], nil
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
