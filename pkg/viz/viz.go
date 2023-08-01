package viz

import (
	"errors"
	"fmt"
	"os"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/vmware-tanzu/cartographer/pkg/apis/v1alpha1"
	"golang.org/x/exp/maps"

	"github.com/nickjameswebb/cartoviz/pkg/util"
)

// TODO: return a byte buffer from this function, instead of rendering directly in function
// TODO: pass a logger to this function
// TODO: pass a context to this function
// TODO: hide visualization behind an interface so that we can test
// TODO: configurable filename, format
// TODO: wrap all errors
func GraphSupplyChain(supplyChain *v1alpha1.ClusterSupplyChain) error {
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

	err = graphSupplyChain(graph, supplyChain)
	if err != nil {
		return err
	}

	return g.RenderFilename(graph, graphviz.PNG, "graph.png")
}

type nodeSource struct {
	Node     *cgraph.Node
	Resource *v1alpha1.SupplyChainResource
}

func graphSupplyChain(graph *cgraph.Graph, supplyChain *v1alpha1.ClusterSupplyChain) error {
	nodesources := []*nodeSource{}
	for _, resource := range supplyChain.Spec.Resources {
		nodeName := fmt.Sprintf("%s/%s", resource.TemplateRef.Kind, resource.Name)
		node, err := graph.CreateNode(nodeName)
		if err != nil {
			return err
		}

		nodesource := &nodeSource{
			Node:     node,
			Resource: resource.DeepCopy(),
		}

		nodesources = append(nodesources, nodesource)
	}

	for _, nodesource := range nodesources {
		depResourceRefs := append(nodesource.Resource.Images, nodesource.Resource.Configs...)
		depResourceRefs = append(depResourceRefs, nodesource.Resource.Sources...)

		for _, depRef := range depResourceRefs {
			var dep *nodeSource
			for _, depNodesource := range nodesources {
				if depNodesource.Resource.Name == depRef.Resource {
					dep = depNodesource
				}
			}
			if dep == nil {
				return errors.New("resource ref has no corresponding resource")
			}

			edgeName, err := edgeNameFromResource(dep.Resource)
			if err != nil {
				return fmt.Errorf("error getting edge name for resource: %w", err)
			}

			edge, err := graph.CreateEdge(edgeName, dep.Node, nodesource.Node)
			if err != nil {
				return err
			}
			edge.SetLabel(edgeName)
		}
	}
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

	if !util.Contains(maps.Keys(templateKindToEdgeNameMapping), templateKind) {
		return "", errors.New("invalid template kind")
	}

	return templateKindToEdgeNameMapping[templateKind], nil
}
