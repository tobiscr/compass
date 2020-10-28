package director

import (
	"context"
	"fmt"
	"reflect"

	"github.com/kyma-incubator/compass/components/director/pkg/graphql"

	gcli "github.com/machinebox/graphql"
	"github.com/pkg/errors"
)

type GenericPage struct {
	Data       interface{}       `json:"data"`
	PageInfo   *graphql.PageInfo `json:"pageInfo"`
	TotalCount int               `json:"totalCount"`
}

type GenericOutput struct {
	Result *GenericPage `json:"result"`
}

type Pager struct {
	QueryGenerator func(...interface{}) string
	PageSize       int
	PageToken      string
	Client         Client
	hasNext        bool
	currentDepth   int
	maxDepth       int
	prevParams     []interface{}
	levels         map[string][]Level
}

type Level struct {
	children map[string][]Level
}

func NewPager(queryGenerator func(...interface{}) string, currentDepth, maxDepth int, pageSize int, client Client, levels map[string][]Level, prevParams []interface{}) *Pager {
	return &Pager{
		QueryGenerator: queryGenerator,
		PageSize:       pageSize,
		Client:         client,
		hasNext:        true,
		currentDepth:   currentDepth,
		maxDepth:       maxDepth,
		prevParams:     prevParams,
		levels:         levels,
	}
}

func (p *Pager) Next(ctx context.Context, output interface{}) error {
	if !p.hasNext {
		return errors.New("no more pages")
	}
	var query string
	if p.currentDepth == 1 {
		params := make([]interface{}, 0)
		for i := 1; i < p.maxDepth; i++ {
			params = append(params, p.PageSize, "")
		}
		params = append([]interface{}{p.PageSize, p.PageToken}, params...)

		query = p.QueryGenerator(params...)
	} else if len(p.prevParams) > 0 {

		// TODO
	}

	req := gcli.NewRequest(query)

	response := GenericOutput{
		Result: &GenericPage{
			Data: output,
		},
	}

	err := p.Client.Do(ctx, req, &response)
	if err != nil {
		return errors.Wrap(err, "while getting page")
	}

	if response.Result == nil {
		return errors.New("unexpected empty response")
	}

	if !response.Result.PageInfo.HasNextPage {
		p.hasNext = false
		return nil
	}

	for _, children := range p.levels {
		for _, child := range children {
			innerPager := NewPager(p.QueryGenerator, p.currentDepth+1, p.maxDepth, p.PageSize, p.Client, child.children, []interface{}{
				p.PageSize,
				p.PageToken,
			})
			go func() {
				// TODO: Where to output?
				innerPager.ListAll(ctx, nil)
			}()
		}

	}

	p.PageToken = string(response.Result.PageInfo.EndCursor)

	return nil
}

func (p *Pager) HasNext() bool {
	return p.hasNext
}

func (p *Pager) ListAll(ctx context.Context, output interface{}) error {
	itemsType := reflect.TypeOf(output)
	if itemsType.Kind() != reflect.Ptr || itemsType.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("items should be a pointer to a slice, but got %v", itemsType)
	}

	allItems := reflect.MakeSlice(itemsType.Elem(), 0, 0)

	for p.HasNext() {
		pageSlice := reflect.New(itemsType.Elem())
		err := p.Next(ctx, pageSlice.Interface())
		if err != nil {
			return err
		}

		allItems = reflect.AppendSlice(allItems, pageSlice.Elem())
	}

	reflect.ValueOf(output).Elem().Set(allItems)
	return nil
}
