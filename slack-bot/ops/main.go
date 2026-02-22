package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
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
	verbose       = flag.Bool("v", false, "verbose")
	searchFlagSet = flag.NewFlagSet("search", flag.ExitOnError)
	column        = searchFlagSet.String("col", "", "Column after which to search")
	value         = searchFlagSet.String("val", "", "Value to search in corresponding column")
	addFlagSet    = flag.NewFlagSet("add", flag.ExitOnError)
	board         = addFlagSet.String("board", "", "Board Name to add")
	group         = addFlagSet.String("group", "", "Board Name to add")
	name          = addFlagSet.String("name", "", "Name to add")
	email         = addFlagSet.String("email", "", "Email to add")
	phone         = addFlagSet.String("phone", "", "Phone to add")
)

func main() {
	godotenv.Load("../.env")
	parseFlags()
	if *verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
	client := monday.New(MONDAY_URL, os.Getenv(MONDAY_TOKEN))
	if searchFlagSet.Parsed() {
		doSearch(client)
	} else {
		doAdd(client)
	}
}

func parseFlags() {
	if len(os.Args) < 2 {
		fmt.Println("expected 'search' or 'add' subcommands")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "search":
		searchFlagSet.Parse(os.Args[2:])
		if *value == "" || *column == "" {
			log.Fatal("Use -c && -v to search")
		}

	case "add":
		addFlagSet.Parse(os.Args[2:])
	default:
		fmt.Println("expected 'search' or 'add' as subcommands")
		os.Exit(1)
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
		Name:      *name,
		Email:     *email,
		Phone:     *phone,
	}
	if err := client.CreateItem(context.Background(), request); err != nil {
		log.Fatal(fmt.Errorf("Failed to create item: %w", err))
	}
}
