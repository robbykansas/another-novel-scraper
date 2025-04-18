package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"robbykansas/another-novel-scraper/cmd/content"
	"robbykansas/another-novel-scraper/cmd/flags"
	"robbykansas/another-novel-scraper/cmd/models"
	"robbykansas/another-novel-scraper/cmd/novel"
	"robbykansas/another-novel-scraper/cmd/search"
	"robbykansas/another-novel-scraper/cmd/ui/listInput"
	"robbykansas/another-novel-scraper/cmd/ui/textInput"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/maps"
)

func init() {
	dir, errP := filepath.Abs(filepath.Dir(os.Args[0]))
	if errP != nil {
		log.Fatal(errP)
	}

	pathFile := fmt.Sprintf("%s/ans-config", dir)

	viper.SetConfigName("another-novel-scraper-config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(pathFile)

	err := viper.ReadInConfig()
	if err != nil {
		models.DefaultPath = pathFile
		os.Mkdir(pathFile, os.ModePerm)

		if _, err := os.Stat(pathFile); err != nil {
			if os.IsNotExist(err) {
				log.Fatal("Create directory failed")
			}
		}
	}

	var flagWeb flags.Web
	var flagChosenTitle flags.ChosenTitle
	rootCmd.AddCommand(titleCmd)

	titleCmd.Flags().StringP("title", "n", "", "title of the novel")
	titleCmd.Flags().VarP(&flagChosenTitle, "chosenTitle", "c", "chosen title")
	titleCmd.Flags().VarP(&flagWeb, "web", "w", "available web")
	titleCmd.Flags().StringP("folder", "f", "", "folder download")
}

type Options struct {
	Title       *textInput.Output
	ChosenTitle *listInput.Selection
	Web         *listInput.Selection
	Folder      *textInput.Output
}

var titleCmd = &cobra.Command{
	Use:   "title",
	Short: "title novel",
	Long:  "",

	Run: func(cmd *cobra.Command, args []string) {
		var p *tea.Program

		flagTitle := cmd.Flag("title").Value.String()
		flagChosenTitle := flags.ChosenTitle(cmd.Flag("chosenTitle").Value.String())
		flagWeb := flags.Web(cmd.Flag("web").Value.String())
		flagFolder := cmd.Flag("folder").Value.String()

		options := Options{
			Title:       &textInput.Output{},
			ChosenTitle: &listInput.Selection{},
			Web:         &listInput.Selection{},
			Folder:      &textInput.Output{},
		}

		novel := &novel.Novel{
			NovelTitle:  flagTitle,
			ChosenTitle: flagChosenTitle,
			Web:         flagWeb,
			Folder:      flagFolder,
		}

		if novel.NovelTitle == "" {
			state := textInput.TitleInput
			header := "Insert title novel?"
			placeholder := "Let This Grieving Soul Retire - Woe Is the Weakling Who Leads the Strongest Party"
			p = tea.NewProgram(textInput.InitialModel(options.Title, header, novel, placeholder, state))
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

		searchTitle, err := search.SearchTitle(novel.NovelTitle)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(1)
		}

		if novel.ChosenTitle == "" {
			headers := "Title Choices"
			var titleChoices []string
			titleChoices = append(titleChoices, maps.Keys(searchTitle)...)
			flags.AllowedTitle = titleChoices
			p = tea.NewProgram(listInput.InitialModelMulti(searchTitle, options.ChosenTitle, headers, novel, listInput.TitleView))
			if _, err := p.Run(); err != nil {
				cobra.CheckErr(err)
			}

			novel.ExitCLI(p)

			novel.ChosenTitle = flags.ChosenTitle(options.ChosenTitle.Choice)
			err := cmd.Flag("chosenTitle").Value.Set(novel.ChosenTitle.String())
			if err != nil {
				log.Fatal("failed to set chosen title")
			}
		}

		if novel.Web == "" {
			headers := "Where do you want the novel scraping from?"
			p = tea.NewProgram(listInput.InitialModelMulti(searchTitle, options.Web, headers, novel, listInput.WebView))
			if _, err := p.Run(); err != nil {
				cobra.CheckErr(err)
			}

			novel.ExitCLI(p)

			flags.AllowedWeb = append(flags.AllowedWeb, options.Web.Choice)
			novel.Web = flags.Web(options.Web.Choice)
			err := cmd.Flag("web").Value.Set(novel.Web.String())
			if err != nil {
				log.Fatal("failed to set web flag value", err)
			}
		}

		state := textInput.FolderInput
		header := "Download folder location?"
		placeholder := "/Users/yourname/Downloads"
		p = tea.NewProgram(textInput.InitialModel(options.Folder, header, novel, placeholder, state))
		if _, err := p.Run(); err != nil {
			cobra.CheckErr(err)
		}

		novel.ExitCLI(p)

		novel.Folder = options.Folder.Output
		errFlag := cmd.Flag("folder").Value.Set(novel.Folder)
		if errFlag != nil {
			log.Fatal("failed to set download folder flag value", errFlag)
		}

		viper.Set("downloadLocation", novel.Folder)
		loc := fmt.Sprintf("%s/another-novel-scraper-config.yaml", models.DefaultPath)
		errWrite := viper.WriteConfigAs(loc)
		if errWrite != nil {
			os.Exit(1)
		}

		content.GetContent(novel.Web.String(), novel.Folder, novel.ChosenTitle.String())

		fmt.Printf("Successfully downloaded novel at %s", novel.Folder)

		os.Exit(1)
	},
}
