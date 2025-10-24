package main

import (
    "fmt"
    "strings"
    "strconv"

    "github.com/gocolly/colly/v2"
)

/********************************/
/**  function to scrape Steam  **/
/********************************/
func scrapeSteam(searchURL string) (int, [][]string) {
  var numberOfGames int
  //create the games array
  var games [][]string
  
  //set the default gocolly collector
  c := colly.NewCollector(
    //only visit these domains
    colly.AllowedDomains("store.steampowered.com"),
  )

  //get the amount of games
  c.OnHTML("div.search_results_count", func(e *colly.HTMLElement) {
    //set the raw text
    rawText := e.Text

    //convert it to a list and chunck everything past the first word
    // (the number of 100% off games)
    amount := strings.Split(rawText, " ")[0]

    //convert the number games into an integer (for later use)
    numberOfGames, _ = strconv.Atoi(amount)
  })

  //on every <a> with Steam's search result item class
  c.OnHTML("a[class*='search_result_row']", func(e *colly.HTMLElement) {
    var gatheredData []string
    //get the raw, uncleaned, game store page url
    rawLink := e.Attr("href")

    //cleanup the game store page link (remove url args)
    cleanedLink := strings.Split(rawLink, "?")[0]

    //save the link to the current data array
    gatheredData = append(gatheredData, cleanedLink)

    //get the game's name
    name := e.ChildText(".responsive_search_name_combined > .search_name > .title")

    //save the name to the current data array
    gatheredData = append(gatheredData, name)

    //go to the game's store page to get the description and tags
    gameTags, gameDesc, gameIMG := goToSteamGamePage(cleanedLink)
    
	//add img url to the data list
    gatheredData = append(gatheredData, gameIMG)

    //add the tags to the data list
    gatheredData = append(gatheredData, gameTags)

    //add the description to the data list
    gatheredData = append(gatheredData, gameDesc)

    //save the gatheredData list to the games array
    games = append(games, gatheredData)
  })
  
  c.Visit(searchURL)

  //return error if there are any
  c.OnError(func(r *colly.Response, err error) {
    fmt.Printf("Request URL: %s failed with response: %v", r.Request.URL, err)
  })
  
  //return the data when done
  return numberOfGames, games
}


/************************************************************/
/** function to go to the game page to get the description **/
/************************************************************/
func goToSteamGamePage(gameURL string) (string, string, string) {
  var desc string
  var tags string
  var img string

  //create default collector
  c := colly.NewCollector(
    //only visit these domains
    colly.AllowedDomains("store.steampowered.com", "gog.com"),
  )

  //if the page contains the game's description snippet...
  c.OnHTML("div[class='game_description_snippet']", func(e *colly.HTMLElement) {
    //get the description and remove whitespace
    desc = strings.TrimSpace(e.Text)    
  })

  c.OnHTML("div[class~='glance_tags']", func(e *colly.HTMLElement) {
    var tagSlice []string
    e.ForEach(".app_tag", func(_ int, el *colly.HTMLElement) {
      tag := strings.TrimSpace(el.Text)
      if tag != "" {
        tagSlice = append(tagSlice, tag)
      }
    })
    //join the tags into a string list
    tagsJoined := strings.Join(tagSlice, "\", \"")
    
    //remove the random plus sign that steam puts in the tags for some reason
    tagsRemoveRandomPlusSign := strings.ReplaceAll(tagsJoined, "+", "")

    //remove the last `', '` (this was added when making it into a list,
    // but not removed immediately due to the weird plus signs)
    tagsCleaned := strings.TrimSuffix(tagsRemoveRandomPlusSign, ", \"")
    
    //
    tags = "\"" + tagsCleaned
  })

  //if the page contains the game's header
  c.OnHTML("img[class='game_header_image_full']", func(e *colly.HTMLElement) {
		//get the header
		rawIMG := e.Attr("src")
		img = strings.Split(rawIMG, "?")[0]
  })

  //return error if there are any
  c.OnError(func(r *colly.Response, err error) {
    fmt.Printf("Request URL: %s failed with response: %v", r.Request.URL, err)
  })

  //go to the game's page
  c.Visit(gameURL)

  return tags, desc, img
}
