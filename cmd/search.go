package cmd

import (
	"fmt"
	"log"

	"github.com/laureanray/clibgen/pkg/api"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	selectedMirror  string
	numberOfResults = 10

	searchCmd = &cobra.Command{
		Use:   "search",
		Short: "search for a book, paper or article",
		Long: `search for a book, paper or article
	example: clibgen search "Eloquent JavaScript"`,
		Run: func(_ *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Please enter a search query!")
				return
			}

			var libgenType = api.LibgenNew

			if selectedMirror == "old" {
				libgenType = api.LibgenOld
			} else if selectedMirror == "new" {
				libgenType = api.LibgenNew
			}

			books, err := api.SearchBookByTitle(args[0], numberOfResults, libgenType)
			if err != nil {
				log.Fatalln(err)
			}

			if err != nil {
				log.Fatal(err)
				return
			}

			var titles []string

			for _, book := range books {
				titles = append(titles, fmt.Sprintf("[%s] [%s] %s (%s)", book.FileSize, book.Extension, book.Title, book.Author))
			}

			prompt := promptui.Select{
				Label: "Select Title",
				Items: titles,
			}

			resultInt, _, err := prompt.Run()

			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}

			api.DownloadSelection(books[resultInt], libgenType)
		},
	}
)

func init() {
	searchCmd.
		PersistentFlags().
		StringVarP(&selectedMirror, "mirror", "m", "old", `select which mirror to use
		options: 
			"old" -> libgen.is
			"new" -> liggen.li 
	`)

	searchCmd.
		PersistentFlags().
		IntVarP(&numberOfResults, "number of results", "n", 10, `number of result(s) to be displayed maximum: 25`)

	rootCmd.AddCommand(searchCmd)
}
