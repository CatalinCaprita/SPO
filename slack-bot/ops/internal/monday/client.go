package monday

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"strings"
	"sync"

	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
)

type ApiClient struct {
	token  string
	client *graphql.Client
	url    string
}

func New(url, token string) *ApiClient {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	gqlClient := graphql.NewClient(url, httpClient)
	return &ApiClient{
		token:  token,
		client: gqlClient,
		url:    url,
	}
}

func (api *ApiClient) ListBoards(ctx context.Context) ([]BoardListing, error) {
	var simpleBoardsQuery = ListBoardsQuery{}
	select {
	case <-ctx.Done():
		log.Println("Context closed")
		return nil, fmt.Errorf("context closed")
	default:
	}
	if err := api.client.Query(ctx, &simpleBoardsQuery, nil); err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}
	return simpleBoardsQuery.Boards, nil
}

func (api *ApiClient) GetBoardWithGroups(ctx context.Context, id string) (*BoardWithGroups, error) {
	var byIdQuery = BoardWithGroupsByIdQuery{}
	var variables = map[string]any{
		"ids": id,
	}
	if err := api.client.Query(ctx, &byIdQuery, variables); err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}
	if len(byIdQuery.Boards) == 0 {
		return nil, fmt.Errorf("no board with id %s could be found", id)
	}
	return &byIdQuery.Boards[0], nil
}

func (api *ApiClient) GetBoardItemsFiltered(ctx context.Context, boardId graphql.ID, limit int, params ItemsQuery) ([]Item, error) {
	var query = BoardByIdWithFilterItemsQuery{}
	var variables = map[string]any{
		"limit":       graphql.Int(limit),
		"queryParams": params,
		"ids":         boardId,
	}
	if err := api.client.Query(ctx, &query, variables); err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}
	return query.Boards[0].ItemsPage.Items, nil

}

func (api *ApiClient) GetItemsInAllBoards(ctx context.Context, params ItemsQuery) (chan Item, error) {
	type Response struct {
		Items []Item
		Error error
	}

	boards, err := api.ListBoards(ctx)
	if err != nil {
		return nil, err
	}

	var listenChan = make(chan Response, len(boards))

	var wg = sync.WaitGroup{}
	log.Printf("Searching in %d boards", len(boards))

	wg.Add(len(boards))
	for _, board := range boards {
		go func() {
			defer wg.Done()
			// for each board, map the column_id to the column title
			var sb = &strings.Builder{}
			fmt.Fprintf(sb, "Available columns in board %s\n", board.Name)
			var columnNameToId = map[string]graphql.ID{}
			for _, col := range board.Columns {
				fmt.Fprintf(sb, "%s[%s], ", col.Title, col.Id)
				columnNameToId[strings.ToLower(string(col.Title))] = col.Id
			}
			var innerParams = &ItemsQuery{
				Rules: []ItemsQueryRule{},
			}
			innerParams.SetOperator(params.Operator)
			for _, rule := range params.Rules {
				var strId = rule.ColumnId.(string)
				var colId = columnNameToId[strings.ToLower(strId)]
				if colId == nil {
					return
				}
				slog.Debug("Replacing %s with %s in board '%s'", rule.ColumnId, colId, board.Name)
				innerParams.AddRule(colId, rule.CompareValue, rule.Operator)
			}
			items, err := api.GetBoardItemsFiltered(ctx, board.Id, 100, *innerParams)
			if err != nil {
				listenChan <- Response{Items: nil, Error: err}
			} else {
				if len(items) > 0 {
					slog.Debug("Found %d entries in Board: %s\n", string(len(items)), board.Name)
					listenChan <- Response{Items: items, Error: nil}
				}
			}
		}()
	}
	go func() { wg.Wait(); close(listenChan) }()

	var itemsChan = make(chan Item, 1)
	go func() {
		defer close(itemsChan)
		for resp := range listenChan {
			if resp.Error != nil {
				slog.Debug(fmt.Errorf("failed to query monday: %s", resp.Error).Error())
			}
			for _, i := range resp.Items {
				itemsChan <- i
			}
		}
	}()
	return itemsChan, nil
}
