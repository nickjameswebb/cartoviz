package main

import (
	"os"

	"github.com/nickjameswebb/cartoviz/pkg/cmd"
	"github.com/nickjameswebb/cartoviz/pkg/types"
	"github.com/spf13/pflag"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func main() {
	flags := pflag.NewFlagSet("cartoviz", pflag.ExitOnError)
	pflag.CommandLine = flags

	scheme := runtime.NewScheme()
	types.AddToScheme(scheme)

	// TODO: better way to pass scheme
	root := cmd.NewCmdViz(&cmd.CmdVizOptions{
		IOStreams: genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr},
		Scheme:    scheme,
	})
	if err := root.Execute(); err != nil {
		// error printed by cobra, no need to do so here
		os.Exit(1)
	}
}
