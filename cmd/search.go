package cmd

import (
	"fmt"
	"genread-api/pkg/api"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"log"
	"strconv"
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
		books, err := api.SearchBookByTitle(args[0], 5)
		if err != nil {
			log.Fatalln(err)
		}

		if err != nil {
			log.Fatal(err)
			return
		}

		var titles []string

		for _, book := range books {
			titles = append(titles, book.Title+"."+book.Extension+" (by: "+book.Author+")")
		}

		prompt := promptui.Select{
			Label: "Select Day",
			Items: titles,
		}

		_, result, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		fmt.Printf("You choose %q\n", result)

		selection, _ := strconv.Atoi(result)

		api.DownloadSelection(books[selection])
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
