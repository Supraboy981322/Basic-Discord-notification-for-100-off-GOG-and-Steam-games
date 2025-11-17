package main

import (
	"fmt"
	"os"
	"time"
	"log"
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
		log.Fatal(content)
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
