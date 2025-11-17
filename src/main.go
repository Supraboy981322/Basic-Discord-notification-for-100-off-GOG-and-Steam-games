package main

import (
	"fmt"
	"strings"
	"errors"
	"bytes"
	"net/http"
	"os"
	"strconv"
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/tidwall/gjson"
)

var (
	stores = []string {
		"Steam",
		"GOG",
	}

	binPath string
	settingsPath string
	steamSearchURL string
	gogSearchURL string 
)

const (
	//will replace with parsing the settings.json file

	//color codes
	ColorReset   = "\033[0m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"

	//the file that stdout is dumped into 
	logFile = "log.txt"
)

func init() {
	//get the path of the free-games-checker binary
	binaryPath, err := os.Executable()
	if err != nil {
		log.Fatalf("failed to get path of free-games-checker")
	}
	
	//get the directory of binary from path
	binPath = filepath.Dir(binaryPath)
	
	//construct settings.json path
	settingsPath = fmt.Sprintf("%s/settings.json", binPath)
	log_("using settings:  " + settingsPath, nil)

	steamSearchURL = getSettings("Steam", "url", false)
	gogSearchURL = getSettings("GOG", "url", false)
	log_(gogSearchURL, nil)
}

func getSettings(which string, what string, isCustom bool) (string) {
	storeJSON, err := getStoreSettings(which)
	if err != nil {
		log_("piss", nil)
		log.Fatal(err)
	}

	return gjson.Get(storeJSON.String(), what).String()
}

func getStoreSettings(store string) (gjson.Result, error) {
	//read the settings.json file
	storesJSONbyte, err := ioutil.ReadFile(settingsPath)
	if err != nil {
		log_("err reading settings.json", err)
	}
	
	//convert the byte[] from json file to a string
	storesJSON := string(storesJSONbyte)

	index := -1
	gjson.Parse(storesJSON).ForEach(
		func(i, v gjson.Result) bool {
			if v.Get("store").String() == store {
				//get the index of the store in json
				index = int(i.Int())
	
				//exit the current iteration of the loop
				return false
			}
			//move to the next iteration of the loop
			return true
	})

	json := gjson.Get(storesJSON, strconv.Itoa(index))

	if index > -1 {
		return json, nil
	} else {
		return json, errors.New("invalid index")
	}

	return json, errors.New("uncaught err")
}


func getWebhook(store string, which int) string {
	storeJSON, err := getStoreSettings(store)

	if err != nil {
		log_("getWebhook()", err)
	}

	webhookPath := fmt.Sprintf(
		"webhooks.%d", which)

	log_(fmt.Sprintf(
			"webhookPath := %s\n",
			webhookPath),
		nil)
	webhook := gjson.Get(
			storeJSON.String(), webhookPath).
		String()

  log_(webhook, nil)

	return webhook
}


func printGame(gameNum int, gameData []string) {
	var dataTypes = []string{
		"link",
		"name",
		"img",
		"tags",
		"desc",
	}

	if gameNum != -1 {
		//log the number of games
		log_(fmt.Sprintf("game %d", gameNum+1), nil)
	}

	//for each of the data types, log the data in
	//  a particular format
	for i := 0; i < len(dataTypes); i++ {
		var itemColor string
		switch dataTypes[i] {
		case "link":
			itemColor = ColorRed
		case "name":
			itemColor = ColorGreen
		case "img":
			itemColor = ColorBlue
		case "tags":
			itemColor = ColorYellow
		case "desc":
			itemColor = ColorMagenta
		default:
			err := errors.New("unsupported data type")
			log_(dataTypes[i], err)
		}
		log_(fmt.Sprintf(
				"  %sgame %s:%s\n    %s\n",
				itemColor, 
				dataTypes[i],
				ColorReset,
				gameData[i]), nil)
	}
}



func sendGamesToDiscord(
		gameName string,
		gameDesc string,
		gameTags string,
		gameURL string,
		gameIMG string,
		gameColor int,
		currentDiscordURL int) {
	log_("sendGamesToDiscord()", nil)
	var webHookURL string
	if strings.Contains(gameURL, "steampowered") {
		webHookURL = getWebhook("Steam", currentDiscordURL)
	} else if strings.Contains(gameURL, "gog.com") {
		webHookURL = getWebhook("GOG", currentDiscordURL)
	} else {
		log_("TODO: Itch.io and Epic", nil)
	}

	//payload struct
	payload := map[string]interface{}{
		"embeds": []map[string]interface{}{
			{
				"title": gameName,
				"description": fmt.Sprintf(
						"**Link:**\n%s\n\n**Description**\n%s\n\n**Tags**\n%s",
						gameURL,
						gameDesc,
						gameTags),
				"image": map[string]string{
					"url": gameIMG,
				},
				"color": gameColor,
			},
		},
	}

	//convert payload struct to json
	data, err  := json.Marshal(payload)
	if err != nil {
		log_("err converting payload struct for Discord to JSON:", err)
		os.Exit(1)
	}

	//create the request
	resp, err := http.Post(
		webHookURL, 
		"application/json",
		bytes.NewBuffer(data))
	if err != nil {
		log_("err sending request:", err)
		os.Exit(1)
	}

	//prevent body from being closed until fn done 
	defer resp.Body.Close()

	//log the reponse status
	log_(fmt.Sprintf(
		"response status:",
		resp.Status), nil)
}

func sendAmountToDiscord(amount int, platform string) {
	webhookURL := getWebhook(platform, 0)

	storeSettings, err := getStoreSettings(platform)
	if err != nil {
		log.Fatal(err)
	}

	msg := getSettings(platform, "numMsg", false)
	msgPar := gjson.Get(storeSettings.String(), "numMsgParams").Array()
	fmt.Println(msgPar)

	var input []any
	for i := 0; i < len(msgPar); i++ {
		switch (msgPar[i].String()) {
		case "$$num$$":
			input = append(input, amount)
		case "$$sNoS$$":
			if amount == 1 {
				input = append(input, "")
			} else {
				input = append(input, "s")
			}
		case "$$isAre$$":
			if amount == 1 {
				input = append(input, "is")
			} else {
				input = append(input, "are")
			}
		default:
			input = append(input, msgPar[i])
		}
	}

	fmt.Println(input)
	log_("webhookURL == " + webhookURL, nil)

	//payload struct
	payload := map[string]interface{}{
		"content": fmt.Sprintf(
				msg, input...),
	}
	
	//convert payload struct to json
	data, err := json.Marshal(payload)
	if err != nil {
		log_("err converting payload for Discord struct to json:", err)
		os.Exit(1)
	}

	//create the request
	resp, err := http.Post(
		webhookURL,
		"application/json",
		bytes.NewBuffer(data))

	if resp != nil {
		//prevent body from being closed until fn done
		defer resp.Body.Close()
	}
}

func main() {
	//iterate through each of the stores
	for i := 0; i < len(stores); i++ {
		var numGames int
		var storeURL string
		var gamesData [][]string
		var storeColor int

		//call the scraper, err if unsupported
		switch stores[i] {
		case "Steam":
			storeURL = steamSearchURL
			numGames, gamesData = scrapeSteam(storeURL)
			storeColor = 3447003
		case "GOG":
			storeURL = gogSearchURL
			numGames, gamesData = scrapeGOG(storeURL)
			storeColor = 10181046
		default:
			errMsg := "attempted to scrape unsupported store"
			err := errors.New(errMsg)
			log_(stores[i], err)
		}
	
		
		//log store's url
		log_(storeURL, nil)
	
		//log games' data array
		log_(fmt.Sprintf(
			"games data: %+v\n",
			gamesData), nil)

		//log the number of games
		log_(fmt.Sprintf(
				"number of games: %d\n",
				numGames), nil)

		//send the amount of games to Discord
		sendAmountToDiscord(numGames,	stores[i])

		//send the games to discord and log them
		currDiscordURL := 1
		for i := 0; i < numGames; i++ {	
			//so the next chunk is easier to read
	        gameLink := gamesData[i][0]
	        gameName := gamesData[i][1]
	        gameIMG := gamesData[i][2]
	        gameTags := gamesData[i][3]
	    	gameDesc := gamesData[i][4]

			sendGamesToDiscord(
				gameName,
				gameDesc,
				gameTags,
				gameLink,
				gameIMG,
				storeColor,
				currDiscordURL)

			printGame(-1, gamesData[i])
	
			currDiscordURL++
			if currDiscordURL == 3 {
				currDiscordURL = 0
			}
		}
	}
}
