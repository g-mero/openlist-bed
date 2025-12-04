package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "openlist-bed",
	Short: "OpenList image bed CLI tool",
	Long:  `A command-line tool for uploading images to OpenList image bed service.`,
}

func main() {
	rootCmd.AddCommand(getUploadCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
