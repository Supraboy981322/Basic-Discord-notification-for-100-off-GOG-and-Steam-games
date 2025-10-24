package main

import (
	"fmt"
  	"strings" 
	
    "github.com/tidwall/gjson"
	"github.com/gocolly/colly/v2"
)

func scrapeGOG(searchURL string) (int, [][]string) {
	//create the games array
	var games [][]string
	var gogJSON string
	var gogAmount int
	
	//set the default gocolly collector
	c := colly.NewCollector(
		//only visit these domains
		colly.AllowedDomains("www.gog.com"),
	)

	//log the search url
	log_(fmt.Sprintf(
			"%s\n\n",
			searchURL), nil)
	
	//get the search results json data 
    c.OnHTML(`#gogcom-store-state`, func(e *colly.HTMLElement) {
        gogJSON = e.Text
    })

	c.Visit(searchURL)
	
	//iterate through the json data 
	gjson.Parse(gogJSON).ForEach(func(i, v gjson.Result) bool {
		if v.Get("*.products").Exists() {
			//get the indec of the body
			gogAmount = int(v.Get("*.productCount").Int())
			
			if gogAmount != 0 {
				for in := 0; in < gogAmount; in++ {
					var gatheredData []string
 
					currentPath := fmt.Sprintf("*.products.%d", in)

					//add the game's page url
					gatheredData = append(
						gatheredData,
						v.Get(currentPath + ".storeLink").String())
					//add the game's title
					gatheredData = append(
						gatheredData,
						v.Get(currentPath + ".title").String())
					//add the game's image
                    gatheredData = append(
						gatheredData,
						v.Get(currentPath + ".coverHorizontal").String())

					//get the tags array
                    tagsJSONarray := v.Get(currentPath + ".tags").String()
                    
					//iterate through each of the tags and add them to an array
					var tagsArray []string
                    gjson.Parse(tagsJSONarray).ForEach(func(ind, va gjson.Result) bool {
                        tagsArray = append(
							tagsArray,
							va.Get("name").String())
                        
						//move to next iteration of loop
						return true
                    })

					//iterate through each of the tags to add them to a
					//  formatted string so an unformatted array isn't
					//  sent to Discord (the first item is added outside
					//  the because it simplifies the code by not parsing
					//  the finished string to remove a comma and space
					//  at the beginning of the string)
                    tags := fmt.Sprintf(
							"'%s'",
							tagsArray[0])
					//iterate through each of the tags
                    for ind := 1; ind < len(tagsArray); ind++ {
                        tags = fmt.Sprintf(
							"%s, '%s'",
							tags,
							tagsArray[ind])
                    }

					//add the formatted tags string to the data array
                    gatheredData = append(
						gatheredData,
						tags)

					//get the game's description, since it's the only thing
					//  not in the json
                	desc := goToGOGgamePage(
						v.Get(currentPath + ".storeLink").String())

					//add the description to the array
                    gatheredData = append(
						gatheredData,
						desc)

					//add the whole data array to the games array as an item
                    games = append(
						games,
						gatheredData)
                }
            }
        }

        //move to the next iteration of loop
        return true
    })

    //return error if there are any
    c.OnError(func(r *colly.Response, err error) {
        log_(fmt.Sprintf(
			"Request URL: %s failed",
			r.Request.URL), err)
    })
  
    //return the data when done
    return gogAmount, games
}


func goToGOGgamePage(gameURL string) (string) {
	var desc string
	
	//create default collector
	c := colly.NewCollector(
		//only visit these domains
		colly.AllowedDomains("www.gog.com"),
	)

	log_("checking gog game page", nil)

	//look for script element
	c.OnHTML("script", func(e *colly.HTMLElement) {
		//if the element contains the comment that GOG puts in a game's
		//  page data JSON
        if strings.Contains(e.Text, "// define global namespace object") {
			//find the description object
    		startString := `},"description":"`
        	if strings.Contains(e.Text, startString) {
                start := strings.Index(e.Text, startString)
                start += len(startString)
                end := strings.Index(e.Text[start:], "\\n")
				//use the value as the description
                desc = e.Text[start : start+end]
            }
        }
    })

    //return error if colly has one 
    c.OnError(func(r *colly.Response, err error) {
        log_(fmt.Sprintf("Request URL: %s failed", r.Request.URL), err)
    })

    //go to the game's page
    c.Visit(gameURL)

    //if the game's description is blank...
    if desc == "" {
        desc = "no desc found"
    }

    return desc
}
