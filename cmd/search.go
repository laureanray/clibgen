package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/laureanray/clibgen/internal/libgen"
	"github.com/laureanray/clibgen/internal/mirror"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func truncateText(s string, max int) string {
	if max > len(s) {
		return s
	}
	return s[:strings.LastIndex(s[:max], " ")] + " ..."
}

func getExtension(s string) string {
	cyan := color.New(color.FgHiCyan).SprintFunc()
	magenta := color.New(color.FgHiMagenta).SprintFunc()
	blue := color.New(color.FgHiBlue).SprintFunc()

	switch s {
	case "pdf":
		return cyan(s)
	case "epub":
		return magenta(s)
	default:
		return blue(s)
	}
}

var (
	selectedSite    string
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

      var m mirror.Mirror

			//
			// var libgenType = api.LibgenNew
			//
			if selectedSite == "legacy" {
        m = mirror.NewLegacyMirror(libgen.IS);
			} else if selectedSite == "new" {
        m = mirror.
			}

      
      books, _ := m.SearchByTitle(args[0])
			//
			// books, siteUsed, err := api.SearchBookByTitle(args[0], numberOfResults, libgenType)
			// if err != nil {
			// 	log.Fatalln(err)
			// }
			//
			// if err != nil {
			// 	log.Fatal(err)
			// 	return
			// }
			//
			var titles []string
			
			for _, book := range books {
				parsedTitle := truncateText(book.Title, 42)
				parsedAuthor := truncateText(book.Author, 24)
				parsedExt := getExtension(fmt.Sprintf("%-4s", book.Extension))
				titles = append(titles, fmt.Sprintf("%s %-6s | %-45s %s", parsedExt, book.FileSize, parsedTitle, parsedAuthor))
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

      println(resultInt)
			//
			// api.DownloadSelection(books[resultInt], siteUsed)
		},
	}
)

func init() {
	searchCmd.
		PersistentFlags().
		StringVarP(&selectedSite, "site", "s", "legacy", `select which site to use
		options: 
			"legacy" 
			"new"
	`)

	searchCmd.
		PersistentFlags().
		IntVarP(&numberOfResults, "number of results", "n", 10, `number of result(s) to be displayed maximum: 25`)

	rootCmd.AddCommand(searchCmd)
}
