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
	Duration           string `json:"duration"`
	Genre              string `json:"genre"`
	Summary            string `json:"summary"`
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
	Long: `This command require 2 arguments,
	1.chart_url : URL of the Imdb Chart (chart_url is one of IMDb Top Indian Charts:
		● Top Rated Indian Movies ( https://www.imdb.com/india/top-rated-indian-movies )
		● Top Rated Tamil Movies ( https://www.imdb.com/india/top-rated-tamil-movies )
		● Top Rated Telugu Movies ( https://www.imdb.com/india/top-rated-telugu-movies )
	2.item_count :  .
	The script returns output a json string of the top items_count number of movie items (with
		attributes as listed below) in that particular IMDb chart.
		● rank
		● title
		● movie_release_year
		● imdb_rating
		● summary
		● duration
		● genre`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 2 {
			fmt.Println("Invalid arguments")

			os.Exit(0)
		}
		var fetchURL = args[0]
		var movieCountLimit, _ = strconv.Atoi(args[1])

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

			if n.Type == html.ElementNode && n.Data == "a" {
				for _, a := range n.Attr {
					if a.Key == "href" {
						path := strings.Split(a.Val, "/")

						if len(path) < 2 {
							break
						}
						if path[1] != "title" {
							break
						} else {
							if count < movieCountLimit {
								newURL := "https://www.imdb.com" + a.Val
								resp2, err := http.Get(newURL)

								if err != nil {
									log.Println(err)
								}
								defer resp.Body.Close()

								doc2, err := html.Parse(resp2.Body)
								if err != nil {
									log.Println(err)
								}
								var f2 func(*html.Node)
								f2 = func(n *html.Node) {
									if n.Type == html.ElementNode && n.Data == "div" {
										text := &bytes.Buffer{}
										for _, a := range n.Attr {
											var DurationAndGenre string
											if a.Key == "class" && a.Val == "subtext" {
												collectText(n, text)
												if text.String() != " " {
													str := text.String()

													DurationAndGenre = strings.TrimFunc(str, func(r rune) bool {
														return !unicode.IsNumber(r)
													})

													DurationAndGenre = strings.ReplaceAll(DurationAndGenre, string(' '), "")
													DurationAndGenre = strings.ReplaceAll(DurationAndGenre, "|", "")
													DurationAndGenre = strings.ReplaceAll(DurationAndGenre, "\t", "")
													t1 := strings.Split(DurationAndGenre, "\n")
													var time string
													var genre string
													time = t1[0]
													genre = t1[3]
													i := 4
													for i < len(t1) {
														if genre[len(genre)-1] != byte(',') {
															break
														} else {
															genre = genre + t1[i]
															i++
														}
													}

													jsonForm.Duration = time
													jsonForm.Genre = genre
												}

											}
											if a.Key == "class" && a.Val == "summary_text" {
												collectText(n, text)
												if text.String() != " " {
													str := text.String()

													summary := strings.TrimFunc(str, func(r rune) bool {
														return !unicode.IsLetter(r)
													})

													jsonForm.Summary = summary

												}

											}

										}

									}

									for c := n.FirstChild; c != nil; c = c.NextSibling {
										f2(c)
									}

								}

								f2(doc2)

							} else {
								break
							}

						}

					}
				}
			}
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
						}

					}

					if a.Key == "class" && a.Val == "ratingColumn imdbRating" {
						collectText(n, text)
						if text.String() != " " {
							str := text.String()

							imdbRatings = strings.TrimFunc(str, func(r rune) bool {
								return !unicode.IsLetter(r) && !unicode.IsNumber(r)
							})
							jsonForm.IMDB_rating = imdbRatings
							if count < movieCountLimit {
								count++
								sendForMarshal(jsonForm)
							} else {
								break
							}

						}
					}

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
