package cmd

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/andybalholm/cascadia"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"golang.org/x/net/html"
)

var rootCmd = &cobra.Command{
	Use:   "magnet",
	Short: "magnet is a tool to get magnet links from torrent galaxy",
	Long: `magnet is a tool to get magnet links from torrent galaxy. For example:

magnet "game of thrones" --season 1 --episode 10 --order-by seeders --limit 10 --page 1`,
	Run: func(cmd *cobra.Command, args []string) {
		defaultTgUrl := "https://torrentgalaxy.to/torrents.php?"
		search := url.QueryEscape(args[0])
		sortBy, _ := cmd.Flags().GetString("sort-by")
		season, _ := cmd.Flags().GetInt("season")
		episode, _ := cmd.Flags().GetInt("episode")

		parsedSortBy := "&sort=" + sortBy
		parsedSearch := "search=" + search

		if season != 0 && episode != 0 {
			parsedSearch += url.QueryEscape(fmt.Sprintf(" s%02de%02d", season, episode))
		}

		fmt.Println("getMagnets called")
		fmt.Println("sortBy:", parsedSortBy)
		fmt.Println("search:", parsedSearch)

		tgUrl := defaultTgUrl + parsedSearch + parsedSortBy + "&order=desc"

		fmt.Println(tgUrl)

		// response, err := http.Get("https://api.github.com/repos/elysiajs/elysia/issues")
		// if err != nil {
		// 	log.Fatalln(err)
		// }

		// body, err := io.ReadAll(response.Body)
		// if err != nil {
		// 	log.Fatalln(err)
		// }

		// sb := string(body)

		// fmt.Println(sb)

		str := GetTestString()
		ioReader := strings.NewReader(str)

		strHtml, _ := html.Parse(ioReader)

		rows := QueryAll(strHtml, ".tgxtablerow")

		fmt.Println(&rows)

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
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("sort-by", "o", "seeders", "Sort by seeders, leechers or time")
	rootCmd.Flags().IntP("limit", "l", 10, "Limit of results")
	rootCmd.Flags().IntP("page", "p", 1, "Page of results")
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
