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

			if n.Type == html.ElementNode && n.Data == "a" {
				for _, a := range n.Attr {
					if a.Key == "href" {
						path := strings.Split(a.Val, "/")
						//fmt.Println(path)
						if len(path) < 2 {
							break
						}
						if path[1] != "title" {
							break
						} else {
							if count < movieCountLimit {
								newURL := "https://www.imdb.com" + a.Val
								resp2, err := http.Get(newURL)
								//fmt.Println(newURL, "success")
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
													//fmt.Println(time, genre)
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
													//fmt.Println(jsonForm)
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
								//fmt.Println(jsonForm)
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
