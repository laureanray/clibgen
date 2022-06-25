package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "clibgen",
	Short: "Library Genesis command line / terminal client",
	Long: `
Clibgen is a CLI application to search and download epubs, pdfs, from library genesis. 
Useful if you are lazy to open up a browser to download e-books/resources.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
