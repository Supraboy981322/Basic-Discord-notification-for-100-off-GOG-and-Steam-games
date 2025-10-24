package main

import (
    "fmt"
    "strings"
    "github.com/gocolly/colly/v2"
    "net/http"
    "os"
    "bytes"
    "strconv"
    "encoding/json"
//    "regexp"
    "github.com/tidwall/gjson"
    "io/ioutil"
    "log"
)

type webhooksStruct struct {
    storeName string `json:"store"`
    discordWebhooks []string `json:"webhooks"`
}

const (
  /*****************************************/
  /**Replace these with your webhook URLs **/
  /******************************************/
  // for testing

  steamSearchURL = "https://store.steampowered.com/search/?maxprice=free&category1=998%2C997%2C993%2C996%2C994&specials=1&ndl=1"

  //color codes
  /*(change these for Discord)*/
  ColorReset   = "\033[0m"
  ColorRed     = "\033[31m"
  ColorGreen   = "\033[32m"
  ColorYellow  = "\033[33m"
  ColorBlue    = "\033[34m"
  ColorMagenta = "\033[35m"
  ColorCyan    = "\033[36m"
)


func getWebhook(store string, which int) string {
    //read the servers.json file
    storesJSONbyte, err := ioutil.ReadFile("servers.json")
    if err != nil {
        log.Fatalf("err reading servers.json:  ", err)
    }

    //convert the byte[] from json file to string
    storesJSON := string(storesJSONbyte)

    index := -1
    gjson.Parse(storesJSON).ForEach(func(i, v gjson.Result) bool {
        if v.Get("store").String() == store {
            //get the index of the store in json
            index = int(i.Int())

            //exit the current iteration of the loop
            return false
        }
        //move to the next iteration of the loop
        return true
    })

    if index != -1 {
        webhookPath := fmt.Sprintf("%d.webhooks.%d", index, which)

        fmt.Printf("webhookPath: %s\n", webhookPath)
        webhook := gjson.Get(storesJSON, webhookPath).String()
        
        return webhook
    } else {
        return " "
    }
}


/*********************/
/**  main function  **/
/*********************/
func main() {
  //create default collector
  c := colly.NewCollector(
    //only visit these domains
    colly.AllowedDomains("store.steampowered.com", "gog.com"),
  )

  //return error if there are any
  c.OnError(func(r *colly.Response, err error) {
    fmt.Printf("Request URL: %s failed with response: %v", r.Request.URL, err)
  })
  
  numSteamGames, steamGames := scrapeSteam(steamSearchURL)

  //make sure the url is set correctly
  fmt.Printf(steamSearchURL)

  //scrape Steam
  fmt.Printf("games data: %+v\n", steamGames)

  fmt.Printf("\033[H\033[2J")
  fmt.Printf("number of games: %d\n", numSteamGames)

  //send the amount of steam games to discord
  sendAmountToDiscord(strconv.Itoa(numSteamGames), "Steam")

  //sets which discord url to start with
  currentDiscordURL := 1
  for i := 0; i < numSteamGames; i++ {
    fmt.Printf("game %d\n", i + 1)
    fmt.Printf("  %sgame link:%s\n    %s\n", ColorRed, ColorReset, steamGames[i][0])
    fmt.Printf("  %sgame name:%s\n    %s\n", ColorGreen, ColorReset, steamGames[i][1])
    fmt.Printf("  %sgame img:%s\n    %s\n", ColorBlue, ColorReset, steamGames[i][2])
    fmt.Printf("  %sgame tags:%s\n    %s\n", ColorYellow, ColorReset, steamGames[i][3])
    fmt.Printf("  %sgame desc:%s\n    %s\n", ColorMagenta, ColorReset, steamGames[i][4])
    fmt.Printf("\n\n")

    //what the list values mean
    gameName := steamGames[i][1]
    gameLink := steamGames[i][0]
    gameIMG := steamGames[i][2]
    gameTags := steamGames[i][3]
    gameDesc := steamGames[i][4]
    sendGamesToDiscord(gameName, gameDesc, gameTags, gameLink, gameIMG, 3447003, currentDiscordURL)

    //move to next discord webhook url
    if currentDiscordURL == 0 {
      currentDiscordURL = 1
    } else if currentDiscordURL == 1 {
      currentDiscordURL = 2
    } else if currentDiscordURL == 2 {
      currentDiscordURL = 0
    }
  }
  
  //scrape GOG
  /*C.Visit("gogSearchURL")*/  
}


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

/*    //get the game's capsule image
    capsuleIMG := e.ChildAttr("img[src^='https://shared.fastly.steamstatic.com/store_item_assets/steam/apps/']", "src")
    capsuleIMGregex := regexp.MustCompile(`capsule_[0-9]x[0-9]`)

    //swap the file in the url for the header (for a higher-resolution image)
    rawHeaderIMG := capsuleIMGregex.ReplaceAllString(capsuleIMG, "header")

    //cleanup the url (remove args)
    cleanedHeaderIMG := strings.Split(rawHeaderIMG, "?")[0]

    //check if the cleanedHeaderIMG contains a string from some capsule
    // IMGs that point to a subfolder which doesn't have the header img
    // and fix it
    invalidHeaderFilePath := `apps/[a-zA-Z0-9]+/[a-zA-Z0-9]+/header\.jpg` //apps/[appid]/[random string]/header.jpg
    containsInvalidPath, err := regexp.MatchString(invalidHeaderFilePath, cleanedHeaderIMG)
    //error handling
    if err != nil {
        fmt.Println("Error matching regex:", err)
        return
    }
    if containsInvalidPath {
      //compile the string that invalidates the url
      invalidPath := regexp.MustCompile(`[a-zA-Z0-9]+/header\.jpg`)
      fixedPath := invalidPath.ReplaceAllString(cleanedHeaderIMG, "header.jpg")
      //replace with valid url
      cleanedHeaderIMG = fixedPath
    }

    //save the IMG url to the current data array
//    gatheredData = append(gatheredData, cleanedHeaderIMG)
    gatheredData = append(gatheredData, capsuleIMG)*/

    //go to the game's store page to get the description and tags
    gameTags, gameDesc, gameIMG := goToGamePage(cleanedLink)
    
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
func goToGamePage(gameURL string) (string, string, string) {
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


func sendGamesToDiscord(gameName string, gameDesc string, gameTags string, gameURL string, gameIMG string, gameColor int, currentDiscordURL int) {
  var webHookURL string
    if strings.Contains(gameURL, "steampowered") {
        webHookURL = getWebhook("Steam", currentDiscordURL)
    } else if strings.Contains(gameURL, "gog.com") {
        webHookURL = getWebhook("GOG", currentDiscordURL)
    } else if strings.Contains(gameURL, "epic.com") {
        webHookURL = getWebhook("Epic", currentDiscordURL)
    } else {
        fmt.Println("TODO: Itch.io")
    }
  //payload struct
	payload := map[string]interface{}{
		"embeds": []map[string]interface{}{
			{
				"title": gameName,
				"description": "**Link:**\n" + gameURL + "\n\n**Description:**\n" + gameDesc + "\n\n**Tags:**\n" + gameTags,
        "image": map[string]string{
          "url": gameIMG,
        },
				"color": gameColor,
			},
		},
	}

	//convert payload struct to json
	data, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("error converting payload struct to json:", err)
		os.Exit(1)
	}
	
	//create the request
	resp, err := http.Post(webHookURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("error sendig request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	//print the response status
	fmt.Println("response status:", resp.Status)
}

func sendAmountToDiscord(amount string, platform string) {
    webHookURL := getWebhook(platform, 0)

  //make sure to use the proper grammatical number
  verbageISare := "are"
  verbageGame := " games"
  if amount == "1" {
    fmt.Printf("only one game? sad.\n")
    verbageISare = "is"
    verbageGame = " game"
  }


  //payload struct
	payload := map[string]interface{}{
    "content": platform + " currently has " + amount + verbageGame + " that " + verbageISare + " 100% off.",
	}

	//convert payload struct to json
	data, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("error converting payload struct to json:", err)
		os.Exit(1)
	}
	
	//create the request
	resp, err := http.Post(webHookURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("error sendig request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	//print the response status
	fmt.Println("response status:", resp.Status)
}
