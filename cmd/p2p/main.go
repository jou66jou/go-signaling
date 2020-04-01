package main

import (
	"fmt"
	"os"
	"p2p/cmd/apiserver"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "root",
	Short: "choose instance to run: apiserver",
	Long:  ``,
}

func main() {
	rootCmd.AddCommand(apiserver.ApiserverCmd)
	args := append([]string{"apiserver"}, os.Args[1:]...)
	rootCmd.SetArgs(args)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
