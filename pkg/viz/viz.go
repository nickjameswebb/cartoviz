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

// Node in the tree: contains resource, all dependent resources
// A resource can be dependent on multiple parents
// Graph all nodes
// Add edges between any node and all of it's dependent resources
// Grapher takes care of the rest

type nodeSource struct {
	Node     *cgraph.Node
	Resource *v1alpha1.SupplyChainResource
}

func graphSupplyChain(graph *cgraph.Graph, supplyChain *v1alpha1.ClusterSupplyChain) error {
	nodesources := []*nodeSource{}
	for _, resource := range supplyChain.Spec.Resources {
		node, err := graph.CreateNode(resource.Name)
		if err != nil {
			return err
		}

		nodesource := &nodeSource{
			Node:     node,
			Resource: resource.DeepCopy(), // TODO: do we need this?
		}

		nodesources = append(nodesources, nodesource)
	}

	for _, nodesource := range nodesources {
		// TODO: apply this to sources, configs
		for _, imageResourceRef := range nodesource.Resource.Images {
			// TODO: refactor this search into a function (generic?)
			var dependency *nodeSource
			for _, dnodesource := range nodesources {
				if dnodesource.Resource.Name == imageResourceRef.Resource {
					dependency = dnodesource
				}
			}
			if dependency == nil {
				return errors.New("image resource ref has no corresponding resource")
			}

			if dependency.Resource.Name == nodesource.Resource.Name {
				return errors.New("cyclic graph, something went wrong")
			}

			// draw edge from dependency -> nodesource
			edgeName, err := edgeNameFromResource(dependency.Resource)
			if err != nil {
				return fmt.Errorf("error getting edge name for resource: %w", err)
			}

			edge, err := graph.CreateEdge(edgeName, dependency.Node, nodesource.Node)
			if err != nil {
				return err
			}
			edge.SetLabel(edgeName)
		}

		for _, configResourceRef := range nodesource.Resource.Configs {
			// TODO: refactor this search into a function (generic?)
			var dependency *nodeSource
			for _, dnodesource := range nodesources {
				if dnodesource.Resource.Name == configResourceRef.Resource {
					dependency = dnodesource
				}
			}
			if dependency == nil {
				return errors.New("config resource ref has no corresponding resource")
			}

			if dependency.Resource.Name == nodesource.Resource.Name {
				return errors.New("cyclic graph, something went wrong")
			}

			// draw edge from dependency -> nodesource
			edgeName, err := edgeNameFromResource(dependency.Resource)
			if err != nil {
				return fmt.Errorf("error getting edge name for resource: %w", err)
			}

			edge, err := graph.CreateEdge(edgeName, dependency.Node, nodesource.Node)
			if err != nil {
				return err
			}
			edge.SetLabel(edgeName)
		}

		for _, sourceResourceRef := range nodesource.Resource.Sources {
			// TODO: refactor this search into a function (generic?)
			var dependency *nodeSource
			for _, dnodesource := range nodesources {
				if dnodesource.Resource.Name == sourceResourceRef.Resource {
					dependency = dnodesource
				}
			}
			if dependency == nil {
				return errors.New("source resource ref has no corresponding resource")
			}

			if dependency.Resource.Name == nodesource.Resource.Name {
				return errors.New("cyclic graph, something went wrong")
			}

			// draw edge from dependency -> nodesource
			edgeName, err := edgeNameFromResource(dependency.Resource)
			if err != nil {
				return fmt.Errorf("error getting edge name for resource: %w", err)
			}

			edge, err := graph.CreateEdge(edgeName, dependency.Node, nodesource.Node)
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
