package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/zoh/betfair"
	"log"
	"os"
)

var confFile = flag.String("conf", "", "A json configuration file")

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()
	if *confFile == "" {
		log.Fatal("You must specify a json configuration file (-conf=CONF_FILE)")
	}
	file, err := os.Open(*confFile)
	checkErr(err)
	defer file.Close()

	dec := json.NewDecoder(file)
	config := new(betfair.Config)
	dec.Decode(&config)

	s, err := betfair.NewSession(config)
	checkErr(err)
	s.Live = true

	loginErr := s.LoginNonInteractive()
	checkErr(loginErr)
	defer s.Logout()

	marketId := "1.122459972"

	marketBooks, err := s.ListMarketBook([]string{marketId})

	for _, runner := range marketBooks[0].Runners {
		fmt.Println(runner.Ex.AvailableToLay[0].Price)
	}
}
