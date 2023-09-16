package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version string

var rootCmd = &cobra.Command{
	Use:   "clibgen",
	Short: "Library Genesis command line / terminal client",
	Long: `
Clibgen is a CLI application to search and download epubs, pdfs, from library genesis. 
Useful if you are lazy to open up a browser to download e-books/resources.`,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Hugo",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
