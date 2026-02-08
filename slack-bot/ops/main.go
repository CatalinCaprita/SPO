package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/CatalinCaprita/SPO/slack-bot/ops/internal/monday"
	"github.com/joho/godotenv"
)

const (
	MONDAY_TOKEN = "MONDAY_TOKEN"
	MONDAY_URL   = "https://api.monday.com/v2"
)

var (
	add    = flag.Bool("add", true, "Add a new item")
	search = flag.Bool("search", false, "Search an item by column type and a given value")
	column = flag.String("c", "", "Column after which to search")
	value  = flag.String("v", "", "Value to search in corresponding column")
	board  = flag.String("board", "", "Board Name to add")
	group  = flag.String("group", "", "Board Name to add")
	name   = flag.String("name", "", "Name to add")
	email  = flag.String("email", "", "Email to add")
	phone  = flag.String("phone", "", "Phone to add")
)

func main() {
	godotenv.Load("../.env")
	flag.Parse()

	if (!*add && !*search) || (*add && *search) {
		log.Fatal("Must use either -add or -search")
	}
	if *search && (*value == "" || *column == "") {
		log.Fatal("Use -c && -v to search")
	}
	client := monday.New(MONDAY_URL, os.Getenv(MONDAY_TOKEN))

	if *search {
		doSearch(client)
	} else {
		doAdd(client)
	}
}

func doSearch(client *monday.ApiClient) {
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

func doAdd(client *monday.ApiClient) {
	var request = monday.CreateItemRequest{
		BoardName: strings.ToLower(*board),
		GroupName: strings.ToLower(*group),
		ItemName:  "test",
		Name:      *name,
		Email:     *email,
		Phone:     *phone,
	}
	log.Println(request)
	if err := client.CreateItem(context.Background(), request); err != nil {
		log.Fatal(fmt.Errorf("Failed to create item: %w", err))
	}
}
