package cmd

import (
	"fmt"
	"log"
	"robbykansas/another-novel-scraper/cmd/flags"
	"robbykansas/another-novel-scraper/cmd/novel"
	"robbykansas/another-novel-scraper/cmd/search"
	"robbykansas/another-novel-scraper/cmd/steps"
	"robbykansas/another-novel-scraper/cmd/ui/listInput"
	"robbykansas/another-novel-scraper/cmd/ui/textInput"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func init() {
	var flagWeb flags.Web
	rootCmd.AddCommand(titleCmd)

	titleCmd.Flags().StringP("title", "n", "", "title of the novel")
	titleCmd.Flags().VarP(&flagWeb, "web", "w", "available web")
}

type Options struct {
	Title *textInput.Output
	Web   *listInput.Selection
}

var titleCmd = &cobra.Command{
	Use:   "title",
	Short: "title novel",
	Long:  "",

	Run: func(cmd *cobra.Command, args []string) {
		var p *tea.Program

		flagTitle := cmd.Flag("title").Value.String()
		flagWeb := flags.Web(cmd.Flag("web").Value.String())

		options := Options{
			Title: &textInput.Output{},
			Web:   &listInput.Selection{},
		}

		novel := &novel.Novel{
			NovelTitle: flagTitle,
			Web:        flagWeb,
		}

		steps := steps.InitSteps(flagWeb)

		if novel.NovelTitle == "" {
			p = tea.NewProgram(textInput.InitialModel(options.Title, "Insert title novel?", novel))
			if _, err := p.Run(); err != nil {
				cobra.CheckErr(err)
			}

			novel.ExitCLI(p)

			novel.NovelTitle = options.Title.Output
			err := cmd.Flag("title").Value.Set(novel.NovelTitle)
			if err != nil {
				log.Fatal("failed to set title flag value", err)
			}
		}
		fmt.Println("test")
		// novelhall, _ := sources.NovelhallSearch(novel.NovelTitle)
		// firstKiss := sources.NewFirstKissNovel()
		searchTitle, _ := search.SearchTitle(novel.NovelTitle)
		// fmt.Printf("%+v\n", novelhall)

		if novel.Web == "" {
			step := steps.Steps["web"]
			step.ListTitle = searchTitle["Novelhall"]
			p = tea.NewProgram(listInput.InitialModelMulti(searchTitle["Novelhall"], options.Web, step.Headers, novel))
			if _, err := p.Run(); err != nil {
				cobra.CheckErr(err)
			}

			novel.ExitCLI(p)

			step.Field = options.Web.Choice

			novel.Web = flags.Web(options.Web.Choice)
			err := cmd.Flag("web").Value.Set(novel.Web.String())
			if err != nil {
				log.Fatal("failed to set web flag value", err)
			}
		}

		fmt.Println(novel.NovelTitle, novel.Web, "<<<<<<<<<<<<<<<<<<<< this")
	},
}
