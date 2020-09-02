/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/spf13/cobra"
	"golang.org/x/net/html"
)

type MovieDetail struct {
	Rank               string `json:"rank"`
	Title              string `json:"title"`
	Movie_release_year string `json:"movie_release_year"`
	IMDB_rating        string `json:"imdb_rating"`
}

func sendForMarshal(jsonForm MovieDetail) {
	js, _ := json.Marshal(jsonForm)
	fmt.Println(string(js))
}

func collectText(n *html.Node, buf *bytes.Buffer) {
	if n.Type == html.TextNode {
		buf.WriteString(n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		collectText(c, buf)
	}
}

// imdbChartFetcherCmd represents the imdbChartFetcher command
var imdbChartFetcherCmd = &cobra.Command{
	Use:   "imdbChartFetcher",
	Short: "This command will pull movie details from IMDB list",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("imdbChartFetcher called")
		if len(args) != 2 {
			fmt.Println("Invalid arguments")
			//panic("Exiting")
			os.Exit(0)
		}
		var fetchURL = args[0]
		var movieCountLimit, _ = strconv.Atoi(args[1])
		//fmt.Println(fetchURL, movieCountLimit)
		count := 0
		resp, err := http.Get(fetchURL)
		if err != nil {
			log.Println(err)
		}
		defer resp.Body.Close()

		doc, err := html.Parse(resp.Body)
		if err != nil {
			log.Println(err)
		}
		var jsonForm MovieDetail
		// Recursively visit nodes in the parse tree
		var f func(*html.Node)
		f = func(n *html.Node) {

			if n.Type == html.ElementNode && n.Data == "td" {
				text := &bytes.Buffer{}
				for _, a := range n.Attr {
					var titleAndYear, imdbRatings string
					if a.Key == "class" && a.Val == "titleColumn" {
						collectText(n, text)
						if text.String() != " " {
							str := text.String()

							titleAndYear = strings.TrimFunc(str, func(r rune) bool {
								return !unicode.IsLetter(r) && !unicode.IsNumber(r)
							})
							titleAndYearSlice := strings.Split(titleAndYear, "\n")

							jsonForm.Rank = strings.Replace(titleAndYearSlice[0], ".", "", -1)
							jsonForm.Title = strings.Trim(titleAndYearSlice[1], " ")
							jsonForm.Movie_release_year = strings.Replace(strings.Trim(titleAndYearSlice[2], " "), "(", "", -1)
							//jsonForm = jsonBuilder(titleAndYear)
							//fmt.Println(jsonForm)
						}

					}

					if a.Key == "class" && a.Val == "ratingColumn imdbRating" {
						collectText(n, text)
						if text.String() != " " {
							str := text.String()

							imdbRatings = strings.TrimFunc(str, func(r rune) bool {
								return !unicode.IsLetter(r) && !unicode.IsNumber(r)
							})
							//jsonForm.IMDB_rating = imdbRatings
							//jsonBuilder(str2)
							//fmt.Println("Moive Ratings : ", imdbRatings)
							jsonForm.IMDB_rating = imdbRatings
							if count < movieCountLimit {
								count++
								sendForMarshal(jsonForm)
							} else {
								break
							}

						}
					}

					//jsonbuilt = titleAndYear //+ "\n" + imdbRatings + "\n"
					//fmt.Println(jsonbuilt)
					//jsonBuilder(jsonbuilt)

				}

			}

			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}

		}

		f(doc)

	},
}

func init() {
	rootCmd.AddCommand(imdbChartFetcherCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// imdbChartFetcherCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// imdbChartFetcherCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
