/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bufio"
	"fmt"
	"genread-api/pkg/api"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
	"strings"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		book, err := api.SearchBookByTitle(args[0], 5)
		if err != nil {
			log.Fatalln(err)
		}

		api.BookPrinter(book)

		fmt.Print("Select option: ")
		option, _ := reader.ReadString('\n')

		option = strings.TrimSuffix(option, "\n")

		intValue, err := strconv.Atoi(option)

		if err != nil {
			log.Fatal(err)
			return
		}

		api.DownloadSelection(book[intValue])

		//fmt.Println(text)
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
