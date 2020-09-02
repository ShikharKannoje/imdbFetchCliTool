"# imdbFetchCliTool" 

### IMDb Chart Fetcher

A command-line script written in Go(lang) that takes input the chart_url and items_count , where
chart_url is one of IMDb Top Indian Charts:

● Top Rated Indian Movies ( https://www.imdb.com/india/top-rated-indian-movies )
● Top Rated Tamil Movies ( https://www.imdb.com/india/top-rated-tamil-movies )
● Top Rated Telugu Movies ( https://www.imdb.com/india/top-rated-telugu-movies )

The script returns a json string of the top items_count number of movie items (with
attributes as listed below) in that particular IMDb chart.

Following attributes of each movie should be printed:
● rank
● title
● movie_release_year
● imdb_rating
● summary
● duration
● genre


Example:
> imdb_chart_fetcher imdbChartFetcher https://www.imdb.com/india/top-rated-indian-movies/ 2
{"rank":"1","title":"Pather Panchali","movie_release_year":"1955","imdb_rating":"8.5","duration":"2h5min","genre":"Drama","summary":"Impoverished priest Harihar Ray, dreaming of a better life for himself and his family, leaves his rural Bengal village in search of work"}
{"rank":"2","title":"Gol Maal","movie_release_year":"1979","imdb_rating":"8.5","duration":"2h24min","genre":"Comedy,Romance","summary":"A man's simple lie to secure his job escalates into more complex lies when his orthodox boss gets suspicious"}