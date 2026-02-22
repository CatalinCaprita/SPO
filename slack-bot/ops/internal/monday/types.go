package monday

import (
	"fmt"
	"strings"

	"github.com/shurcooL/graphql"
)

type WorkpacesQuery struct {
	Workspaces []WorkspaceListing
}

type WorkspaceListing struct {
	Id   graphql.ID
	Name graphql.String
	Kind graphql.String
}

type ListBoardsQuery struct {
	Boards []BoardListing `graphql:"boards(workspace_ids: [$wsId])"`
}

type BoardListing struct {
	Id          graphql.ID
	Name        graphql.String
	Description graphql.String
	BoardKind   graphql.String `graphql:"board_kind"`
	Columns     []Column
}

type Group struct {
	Id       graphql.ID
	Title    graphql.String
	Position graphql.String
}
type BoardWithGroups struct {
	Id          graphql.ID
	Name        graphql.String
	Description graphql.String
	BoardKind   graphql.String `graphql:"board_kind"`
	Groups      []Group
}

func (b BoardWithGroups) String() string {
	var builder = strings.Builder{}
	builder.WriteString(fmt.Sprintf("Board %s ID %s\nGroups:", b.Name, b.Id))
	for _, g := range b.Groups {
		builder.WriteString(fmt.Sprintf("(Name: %s ID: %s) ", g.Title, g.Id))
	}
	return builder.String()
}

type BoardWithGroupsByIdQuery struct {
	Boards []BoardWithGroups `graphql:"boards(ids: [$ids])"`
}

type EmailColumnValue struct {
	Email graphql.String `json:"email"`
	Text  graphql.String `json:"text"`
}

func NewEmailColumnValue(val string) EmailColumnValue {
	return EmailColumnValue{Email: graphql.String(val), Text: graphql.String(val)}
}

type PhoneColumnValue struct {
	Phone graphql.String `json:"phone"`
	Text  graphql.String `json:"text"`
}

func NewPhoneColumnValue(val string) PhoneColumnValue {
	return PhoneColumnValue{Phone: graphql.String(val), Text: graphql.String(val)}
}

type TextColumnValue struct {
	Text  graphql.String
	Value graphql.String
}
type ColumnValue struct {
	Id         graphql.ID
	Text       graphql.String
	Value      graphql.String
	TextValue  TextColumnValue  `graphql:"... on TextValue"`
	EmailValue EmailColumnValue `graphql:"... on EmailValue"`
	PhoneValue PhoneColumnValue `graphql:"... on PhoneValue"`
}

type Column struct {
	Id    graphql.ID
	Title graphql.String
	Type  graphql.String
}

type Item struct {
	Id           graphql.ID
	Name         graphql.String
	Group        Group
	ColumnValues []ColumnValue `graphql:"column_values"`
}

func (item Item) String() string {
	var email, phone = "", ""
	for _, c := range item.ColumnValues {
		if c.EmailValue.Email != "" {
			email = string(c.EmailValue.Email)
		}
		if c.PhoneValue.Phone != "" {
			phone = string(c.PhoneValue.Phone)
		}
	}
	return fmt.Sprintf("Name: %s, Email: %s, Phone: %s\n", item.Name, email, phone)
}

type ItemsPage struct {
	// Cursor graphql.String
	Items []Item `graphql:"items" json:"items"`
}
type BoardWithItemsPage struct {
	ItemsPage   ItemsPage `graphql:"items_page(limit: $limit query_params: $queryParams)" json:"items_page"`
	Id          graphql.ID
	Name        graphql.String
	Description graphql.String
}

type BoardByIdWithFilterItemsQuery struct {
	Boards []BoardWithItemsPage `graphql:"boards(ids: [$ids])"`
}

type ItemsQuery struct {
	Rules    []ItemsQueryRule   `graphql:"rules" json:"rules"`
	Operator ItemsQueryOperator `graphql:"operator" json:"operator"`
}

func (q *ItemsQuery) SetRules(rules []ItemsQueryRule) {
	q.Rules = rules
}

func (q *ItemsQuery) SetOperator(op ItemsQueryOperator) {
	q.Operator = op
}

func (q *ItemsQuery) AddRule(colId graphql.ID, colVar CompareValue, op ItemsQueryRuleOperator) {
	q.Rules = append(q.Rules, ItemsQueryRule{ColumnId: colId, CompareValue: colVar, Operator: op})
}

type ItemsQueryRule struct {
	ColumnId     graphql.ID             `json:"column_id"`
	CompareValue CompareValue           `json:"compare_value"`
	Operator     ItemsQueryRuleOperator `json:"operator"`
}

type ItemsQueryOperator string
type CompareValue graphql.String
type ItemsQueryRuleOperator string
type JSON string

const (
	ANY_OF                 ItemsQueryRuleOperator = "any_of"
	NOT_ANY_OF             ItemsQueryRuleOperator = "not_any_of"
	IS_EMPTY               ItemsQueryRuleOperator = "is_empty"
	IS_NOT_EMPTY           ItemsQueryRuleOperator = "is_not_empty"
	GREATER_THAN           ItemsQueryRuleOperator = "greater_than"
	GREATER_THAN_OR_EQUALS ItemsQueryRuleOperator = "greater_than_or_equals"
	LOWER_THAN             ItemsQueryRuleOperator = "lower_than"
	LOWER_THAN_OR_EQUAL    ItemsQueryRuleOperator = "lower_than_or_equal"
	BETWEEN                ItemsQueryRuleOperator = "between"
	NOT_CONTAINS_TEXT      ItemsQueryRuleOperator = "not_contains_text"
	CONTAINS_TEXT          ItemsQueryRuleOperator = "contains_text"
	CONTAINS_TERMS         ItemsQueryRuleOperator = "contains_terms"
	STARTS_WITH            ItemsQueryRuleOperator = "starts_with"
	ENDS_WITH              ItemsQueryRuleOperator = "ends_with"
	WITHIN_THE_NEXT        ItemsQueryRuleOperator = "within_the_next"
	WITHIN_THE_LAST        ItemsQueryRuleOperator = "within_the_last"
)

const (
	COLUMN_TYPE_STATUS string = "status"
)

type CreateItem struct {
	Id graphql.ID
}
type CreateItemMutation struct {
	CreateItem CreateItem `graphql:"create_item(board_id: $boardId group_id: $groupId item_name: $itemName column_values: $cols)"`
}

type CreateItemRequest struct {
	BoardName string
	GroupName string
	ItemName  string
	Name      string
	Email     string
	Phone     string
}
