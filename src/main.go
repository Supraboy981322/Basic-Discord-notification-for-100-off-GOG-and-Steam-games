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
	"time"
	"io/ioutil"
	"log"

	"github.com/tidwall/gjson"
)

type webhooksStruct struct {
	storeName string `json:"store"`
	discordWebhooks []string `json:"webhooks"`
}

var (
	stores = []string {
		"Steam",
		"GOG",
	}
)

const (
	//will replace with parsing the settings.json file
	steamSearchURL = "https://store.steampowered.com/search/?maxprice=free&category1=998%2C997%2C993%2C996%2C994&specials=1&ndl=1"
	gogSearchURL = "https://www.gog.com/en/games?priceRange=0,0&discounted=true"

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

//the purpose of this fn is so I can print to stdout
//  and a log file in one fn call
func log_(text string, eror error) {
	var content string
	//if there was no err passed to the fn
	if eror != nil {
		//this is formatting for the log file 
		content = fmt.Sprintf(
			"%s  ;  err:  %s\n%s  ;  %s\n\n",
			time.Now(), eror,
			time.Now(), text)
		//this mirrors to stdout 
		fmt.Println(content)
	} else {
		//this is also formatting for log file
		content = fmt.Sprintf(
			"%s  ; %s\n\n", 
			time.Now(), text)
		//this prints just the text passed to
		//  the fn, without the formatting 
		fmt.Println(text)
	}

	//open the log file 
	fi, err := os.OpenFile(
		logFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644)
	if err != nil {
		log.Fatalf("err openning log file:  %s\n", err)
	}

	//prevent it from being closed until fn finishes
	defer fi.Close()
	
	//write the formatted string to the log file
	_, err = fi.WriteString(content)
	if err != nil {
		log.Fatalf("err writing to log file:  %s\n", err)
	}
}


func getWebhook(store string, which int) string {
	//read the settings.json file
	storesJSONbyte, err := ioutil.ReadFile("settings.json")
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

	if index != -1 {
		webhookPath := fmt.Sprintf(
			"%d.webhooks.%d",
			index, which)
		
		log_(fmt.Sprintf(
				"webhookPath := %s\n",
				webhookPath), nil)

		webhook := gjson.Get(
			storesJSON,
			webhookPath).String()

		return webhook
	} else {
		return " "
	}
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


func main() {
	//iterate through each of the stores
	for i := 0; i < len(stores); i++ {
		var numGames int
		var storeURL string
		var gamesData [][]string

		//call the scraper, err if unsupported
		switch stores[i] {
		case "Steam":
			storeURL = steamSearchURL
			numGames, gamesData = scrapeSteam(storeURL)
		case "GOG":
			storeURL = gogSearchURL
			numGames, gamesData = scrapeGOG(storeURL)
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
		sendAmountToDiscord(
			strconv.Itoa(numGames),
			stores[i])

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
				3447003, //Discord embed card color
				currDiscordURL)

			printGame(-1, gamesData[i])
	
			currDiscordURL++
			if currDiscordURL == 3 {
				currDiscordURL = 0
			}
		}
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

func sendAmountToDiscord(amount string, platform string) {
	webhookURL := getWebhook(platform, 0)
	
	//make sure to user the proper grammatical number
	verbageISare := "are"
	verbageGame := "games"
	if amount == "1" {
		log_("Only one game? sad.\n", nil)
		verbageISare = "is"
		verbageGame = "game"
	}

	//payload struct
	payload := map[string]interface{}{
		"content": fmt.Sprintf(
				"%s currently has %s %s that %s 100%% off.",
				platform,
				amount,
				verbageGame,
				verbageISare),
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
