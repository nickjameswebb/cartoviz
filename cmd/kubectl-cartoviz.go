package main

import (
	"fmt"
	"os"

	"github.com/nickjameswebb/cartoviz/pkg/cmd"
	"github.com/nickjameswebb/cartoviz/pkg/types"
	"github.com/spf13/pflag"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes/scheme"
)

func main() {
	flags := pflag.NewFlagSet("cartoviz", pflag.ExitOnError)
	pflag.CommandLine = flags

	types.AddToScheme(scheme.Scheme)

	root := cmd.NewCmdViz(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	if err := root.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "cartoviz failed: %v\n", err)
		os.Exit(1)
	}
}
