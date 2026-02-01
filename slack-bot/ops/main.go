package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/CatalinCaprita/SPO/slack-bot/ops/internal/monday"
	"github.com/joho/godotenv"
)

const (
	MONDAY_TOKEN = "MONDAY_TOKEN"
	MONDAY_URL   = "https://api.monday.com/v2"
)

var (
	column = flag.String("column", "", "Column after which to search")
	value  = flag.String("value", "", "Value to search in corresponding column")
)

func main() {
	godotenv.Load("../.env")
	flag.Parse()

	client := monday.New(MONDAY_URL, os.Getenv(MONDAY_TOKEN))

	var params = monday.ItemsQuery{
		Rules: []monday.ItemsQueryRule{
			{
				ColumnId:     *column,
				CompareValue: monday.CompareValue(*value),
				Operator:     monday.CONTAINS_TEXT,
			},
		},
		Operator: "and",
	}
	items, err := client.GetItemsInAllBoards(context.Background(), params)
	if err != nil {
		panic(err)
	}
	for item := range items {
		log.Println("Found: ", item)
	}
}
