/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log"
	"os"

	"github.com/CatalinCaprita/SPO/slack-bot/ops/internal/monday"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "monday",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		doRun()
	},
}

const (
	MONDAY_TOKEN = "MONDAY_TOKEN"
	MONDAY_URL   = "https://api.monday.com/v2"
)

var value string
var column string

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ops.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringVarP(&column, "column", "c", "name", "Column after wich to perform lookup")
	rootCmd.Flags().StringVarP(&value, "value", "v", "", "Text to look up")
}

func doRun() {
	godotenv.Load("../.env")
	client := monday.New(MONDAY_URL, os.Getenv(MONDAY_TOKEN))

	var params = monday.ItemsQuery{
		Rules: []monday.ItemsQueryRule{
			{
				ColumnId:     column,
				CompareValue: monday.CompareValue(value),
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
