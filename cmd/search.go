package cmd

import (
	"fmt"
	"github.com/laureanray/clibgen/pkg/api"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"log"
	"strconv"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "library gensis command line / terminal client",
	Long: `
clibgen is a CLI application to search and download epubs, pdfs, from library genesis. 
Useful if you are lazy to open up a browser to download e-books/resources.'`,
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
			Label: "Select Title",
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
