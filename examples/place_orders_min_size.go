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

	marketId := "1.121189454"
	selectionId := 10400676

	instructions := []betfair.PlaceInstruction{
		betfair.PlaceInstruction{
			SelectionId: selectionId,
			Side:        betfair.BACK,
			OrderType:   betfair.LIMIT,
			Handicap:    0,
			LimitOrder: betfair.LimitOrder{
				PersistenceType: betfair.PERSIST,
				Size:            4,
				Price:           1000,
			},
		},
	}
	report, err := s.PlaceOrders(marketId, instructions, "")
	checkErr(err)
	fmt.Println(report)

	betId1 := report.InstructionReports[0].BetId

	instructions = []betfair.PlaceInstruction{
		betfair.PlaceInstruction{
			SelectionId: selectionId,
			Side:        betfair.BACK,
			OrderType:   betfair.LIMIT,
			Handicap:    0,
			LimitOrder: betfair.LimitOrder{
				PersistenceType: betfair.PERSIST,
				Size:            1,
				Price:           1000,
			},
		},
	}
	report, err = s.PlaceOrders(marketId, instructions, "")
	checkErr(err)
	fmt.Println(report)

	betId2 := report.InstructionReports[0].BetId

	cancelInstruction := []betfair.CancelInstruction{betfair.CancelInstruction{BetId: betId1}}
	report2, err := s.CancelOrders(marketId, cancelInstruction, "")
	checkErr(err)
	fmt.Println(report2)


	replaceInstruction := []betfair.ReplaceInstruction{
		betfair.ReplaceInstruction{betId2, 70},
	}
	report3, err := s.ReplaceOrders(marketId, replaceInstruction, "")
	checkErr(err)
	fmt.Println(report3)

	// place bet with 4 size (1)
	// place bet with 1 size (2)
	// remove (1) bet
	// update price (2) bet
}
