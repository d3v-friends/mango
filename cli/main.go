package main

import (
	"github.com/d3v-friends/mango/cli/generator"
	"github.com/spf13/cobra"
)

func main() {
	var err error
	cmd := &cobra.Command{
		Use: "mango",
	}

	cmd.AddCommand(generator.NewModelCmd())

	if err = cmd.Execute(); err != nil {
		panic(err)
	}
}
