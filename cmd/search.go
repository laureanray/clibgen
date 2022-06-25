package cmd

import (
	"fmt"
	"github.com/laureanray/clibgen/pkg/api"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"log"
	"strconv"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "search for a book, paper or article",
	Long: `search for a book, paper or article

example: clibgen search "Elon Musk"
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please enter a search query!")
			return
		}
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
			titles = append(titles, book.Title+"(by: "+book.Author+") ["+book.Extension+"]")
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
