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
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vmware-tanzu/cartographer/pkg/apis/v1alpha1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
)

func NewCmdViz(streams genericclioptions.IOStreams) *cobra.Command {
	configFlags := genericclioptions.NewConfigFlags(true)

	cmd := &cobra.Command{
		Use:          "cartoviz",
		Short:        "Visualize a Cartographer Supply Chain",
		Example:      "TODO",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			config, err := configFlags.ToRESTConfig()
			if err != nil {
				return fmt.Errorf("failed to get rest config: %w", err)
			}
			config.GroupVersion = &v1alpha1.SchemeGroupVersion
			config.APIPath = "/apis"
			config.ContentType = runtime.ContentTypeJSON
			config.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)

			client, err := rest.RESTClientFor(config)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			result := v1alpha1.ClusterSupplyChain{}
			err = client.
				Get().
				Resource("clustersupplychains").
				Name("source-to-url").
				Do(context.TODO()).
				Into(&result)
			if err != nil {
				return fmt.Errorf("failed to get supply chain: %w", err)
			}

			fmt.Println(result)

			return nil
			// if err := o.Complete(c, args); err != nil {
			// 	return err
			// }
			// if err := o.Validate(); err != nil {
			// 	return err
			// }
			// if err := o.Run(); err != nil {
			// 	return err
			// }

			// return nil
		},
	}

	configFlags.AddFlags(cmd.PersistentFlags())

	return cmd
}
