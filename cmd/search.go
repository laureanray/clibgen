package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/laureanray/clibgen/internal/book"
	"github.com/laureanray/clibgen/internal/libgen"
	"github.com/laureanray/clibgen/internal/mirror"
	"github.com/laureanray/clibgen/internal/utils"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

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
  selectedFilter  string
	outputDirectory string
	numberOfResults int

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

      m := mirror.NewLegacyMirror(libgen.IS)

      // Set Defaults
      m.SetNumberOfResults(numberOfResults)

			var books []book.Book

			switch selectedFilter {
			case libgen.AUTHOR:
				books, _ = m.SearchByAuthor(args[0])
      case libgen.ISBN:
        books, _ = m.SearchByTitle(args[0])
			default:
				books, _ = m.SearchByTitle(args[0])
			}

			if len(books) == 0 {
				return
			}

			var titles []string

			for _, book := range books {
				parsedTitle := utils.TruncateText(book.Title, 42)
				parsedAuthor := utils.TruncateText(book.Author, 24)
				parsedExt := getExtension(fmt.Sprintf("%-4s", book.Extension))
				titles = append(titles, fmt.Sprintf("%s %-6s | %-45s %s", parsedExt, book.FileSize, parsedTitle, parsedAuthor))
			}

			prompt := promptui.Select{
				Label: "Select Title",
				Items: titles,
        Size: 10,
			}

			resultInt, _, err := prompt.Run()

			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}

			m.DownloadSelection(books[resultInt], outputDirectory)
		},
	}
)

func init() {
	searchCmd.
		PersistentFlags().
		StringVarP(&selectedFilter, "filter", "f", "title", `search by [title, author, isbn]`)

	searchCmd.
		PersistentFlags().
		StringVarP(&outputDirectory, "output", "o", "./", `Output directory`)

	searchCmd.
		PersistentFlags().
		IntVarP(&numberOfResults, "number of results", "n", 10, `number of result(s) to be displayed maximum: 25`)

	rootCmd.AddCommand(searchCmd)
}
