/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vmware-tanzu/cartographer/pkg/apis/v1alpha1"

	"github.com/nickjameswebb/cartoviz/pkg/viz"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func NewCmdViz(streams genericclioptions.IOStreams, scheme *runtime.Scheme) *cobra.Command {
	configFlags := genericclioptions.NewConfigFlags(true)

	cmd := &cobra.Command{
		Use:     "cartoviz",
		Short:   "Visualize a Cartographer Supply Chain",
		Example: "TODO",
		Args:    cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			supplyChainName := args[0]

			builder := resource.NewBuilder(configFlags)

			obj, err := builder.
				WithScheme(scheme, scheme.PrioritizedVersionsAllGroups()...).
				ResourceNames("clustersupplychain", supplyChainName).
				Do().
				Object()
			if err != nil {
				return err
			}

			supplyChain := obj.(*v1alpha1.ClusterSupplyChain)
			return viz.GraphSupplyChain(supplyChain)
		},
	}

	configFlags.AddFlags(cmd.PersistentFlags())

	return cmd
}
