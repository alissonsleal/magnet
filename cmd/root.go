package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/andybalholm/cascadia"
	"github.com/atotto/clipboard"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"golang.org/x/net/html"
)

var rootCmd = &cobra.Command{
	Use:   "magnet",
	Short: "magnet is a tool to get magnet links from torrent galaxy",
	Long: `magnet is a tool to get magnet links from torrent galaxy. 

Usage:
magnet "For All Mankind S04E10"
magnet "For All Mankind" --season 4 --episode 10


`,
	Run: func(cmd *cobra.Command, args []string) {
		defaultTgUrl := "https://torrentgalaxy.to/torrents.php?"
		search := url.QueryEscape(args[0])
		season, _ := cmd.Flags().GetInt("season")
		episode, _ := cmd.Flags().GetInt("episode")

		parsedSearch := "&search=" + search

		if season != 0 && episode != 0 {
			parsedSearch += url.QueryEscape(fmt.Sprintf(" s%02de%02d", season, episode))
		}

		tgUrl := defaultTgUrl + "sort=seeder&order=desc" + parsedSearch

		spinner, _ := pterm.DefaultSpinner.Start("Searching for magnet links")
		spinner.RemoveWhenDone = true

		response, err := http.Get(tgUrl)
		if err != nil {
			spinner.UpdateText("Something went wrong")
			spinner.Fail()

			log.Fatalln(err)

		}

		spinner.UpdateText("Done")
		spinner.Info()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatalln(err)
		}

		str := string(body)

		ioReader := strings.NewReader(str)

		strHtml, _ := html.Parse(ioReader)

		rows := QueryAll(strHtml, ".tgxtablerow")

		records := []map[string]string{}

		for _, row := range rows {

			record := map[string]string{}
			record["name"] = AttrOr(Query(row, "a[title]"), "title", "no title")
			record["seeders"] = Query(row, "font[color='green']").FirstChild.FirstChild.Data
			record["leechers"] = Query(row, "font[color='#ff0000']").FirstChild.FirstChild.Data
			record["size"] = Query(row, "span.badge").FirstChild.Data
			record["magnet"] = Query(row, "a[role='button']").Attr[0].Val

			records = append(records, record)
		}

		data := pterm.TableData{{"Number", "Name", "Seeders", "Leechers", "Size"}}
		var options []string

		for index, record := range records {
			itemNumber := fmt.Sprintf("%d", index+1)

			data = append(data, []string{itemNumber, record["name"], record["seeders"], record["leechers"], record["size"]})
			options = append(options, itemNumber+" - "+record["name"])
		}

		pterm.DefaultTable.WithHasHeader().WithBoxed().WithData(data).Render()

		selectedOption, _ := pterm.DefaultInteractiveSelect.WithOptions(options).Show()

		selectedOptionString := strings.Split(selectedOption, " ")[0]

		selectedOptionNumber := 0

		fmt.Sscanf(selectedOptionString, "%d", &selectedOptionNumber)

		chosenMagnetLink := records[selectedOptionNumber-1]["magnet"]

		fmt.Sscanf(selectedOptionString, "%d", &selectedOptionNumber)

		pterm.Success.Println(pterm.Green(chosenMagnetLink))

		clipboard.WriteAll(chosenMagnetLink)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().IntP("season", "s", 0, "Season of the show")
	rootCmd.Flags().IntP("episode", "e", 0, "Episode of the show")
}

func Query(n *html.Node, query string) *html.Node {
	sel, err := cascadia.Parse(query)
	if err != nil {
		return &html.Node{}
	}
	return cascadia.Query(n, sel)
}

func QueryAll(n *html.Node, query string) []*html.Node {
	sel, err := cascadia.Parse(query)
	if err != nil {
		return []*html.Node{}
	}
	return cascadia.QueryAll(n, sel)
}

func AttrOr(n *html.Node, attrName, or string) string {
	for _, a := range n.Attr {
		if a.Key == attrName {
			return a.Val
		}
	}
	return or
}
