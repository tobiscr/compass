package director

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"

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
	depth          int
	prevParams     []interface{}
	levels         []Level
	PageInfoPath   string
}

type Level struct {
	queryGenerator func(...interface{}) string
	PageInfoPath   string
	children       []Level
}

func NewPager(queryGenerator func(...interface{}) string, pageInfoPath string, depth int, pageSize int, client Client, levels []Level, prevParams []interface{}) *Pager {
	return &Pager{
		QueryGenerator: queryGenerator,
		PageSize:       pageSize,
		Client:         client,
		hasNext:        true,
		depth:          depth,
		prevParams:     prevParams,
		levels:         levels,
		PageInfoPath:   pageInfoPath,
	}
}

func (p *Pager) Next(ctx context.Context, wg *sync.WaitGroup, output interface{}) error {
	if !p.hasNext {
		return errors.New("no more pages")
	}
	var query string

	params := p.prevParams
	leftParams := p.depth - (len(params) / 2)
	for i := 0; i < leftParams; i++ {
		params = append(params, p.PageSize, p.PageToken)
	}
	query = p.QueryGenerator(params...)
	fmt.Println(">>>>>", query)

	req := gcli.NewRequest(query)

	// response := GenericOutput{
	// 	Result: &GenericPage{
	// 		Data: output,
	// 	},
	// }
	var response map[string]interface{}

	err := p.Client.Do(ctx, req, &response)
	if err != nil {
		return errors.Wrap(err, "while getting page")
	}

	p.processChildren(ctx, wg, p.levels)

	// if response.Result == nil {
	// 	return errors.New("unexpected empty response")
	// }

	pathSegements := strings.Split(p.PageInfoPath, ".")
	var pageInfo = response
	for _, ps := range pathSegements {
		pageInfo = pageInfo[ps].(map[string]interface{})
	}
	hasNextPage := pageInfo["hasNextPage"].(bool)

	if !hasNextPage {
		p.hasNext = false
		return nil
	}

	p.PageToken = pageInfo["endCursor"].(string)

	return nil
}

func (p *Pager) processChildren(ctx context.Context, wg *sync.WaitGroup, levels []Level) {
	for _, level := range levels {
		newParams := p.prevParams
		newParams = append(newParams, p.PageSize, p.PageToken)
		innerPager := NewPager(level.queryGenerator, level.PageInfoPath, p.depth+1, p.PageSize, p.Client, level.children, newParams)
		wg.Add(1)
		go func() {
			defer wg.Done()
			// TODO: Where to output?
			var result interface{}
			innerPager.ListAll(ctx, wg, &result)
		}()
	}
}

func (p *Pager) HasNext() bool {
	return p.hasNext
}

func (p *Pager) ListAll(ctx context.Context, wg *sync.WaitGroup, output interface{}) error {
	itemsType := reflect.TypeOf(output)
	if itemsType.Kind() != reflect.Ptr || itemsType.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("items should be a pointer to a slice, but got %v", itemsType)
	}

	allItems := reflect.MakeSlice(itemsType.Elem(), 0, 0)

	for p.HasNext() {
		pageSlice := reflect.New(itemsType.Elem())
		err := p.Next(ctx, wg, pageSlice.Interface())
		if err != nil {
			return err
		}

		allItems = reflect.AppendSlice(allItems, pageSlice.Elem())
	}

	reflect.ValueOf(output).Elem().Set(allItems)
	return nil
}
