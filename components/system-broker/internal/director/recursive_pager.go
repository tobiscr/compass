package director

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
	gcli "github.com/machinebox/graphql"
	"github.com/pkg/errors"
)

type RecursivePager struct {
	QueryGenerator func(pageSize int, page string) string
	PageSize       int
	PageToken      string
	Client         Client
	hasNext        bool
	PageInfoPath   string
}

func NewRecursivePager(queryGenerator func(pageSize int, page string) string, pageSize int, pageInfoPath string, client Client) *RecursivePager {
	return &RecursivePager{
		QueryGenerator: queryGenerator,
		PageSize:       pageSize,
		Client:         client,
		hasNext:        true,
		PageInfoPath:   pageInfoPath,
	}
}

func (p *RecursivePager) Next(ctx context.Context, output interface{}) error {
	if !p.hasNext {
		return errors.New("no more pages")
	}

	query := p.QueryGenerator(p.PageSize, p.PageToken)
	req := gcli.NewRequest(query)

	var response interface{}
	// response := GenericOutput{
	// 	Result: &GenericPage{
	// 		Data: output,
	// 	},
	// }

	err := p.Client.Do(ctx, req, &response)
	if err != nil {
		return errors.Wrap(err, "while getting page")
	}

	paths := strings.Split(p.PageInfoPath, ".")
	var r interface{}
	for _, level := range paths {
		r = response.(map[string]interface{})[level]
	}

	pInfo := r.(*graphql.PageInfo)

	// if response.Result == nil {
	// 	return errors.New("unexpected empty response")
	// }

	if !pInfo.HasNextPage {
		p.hasNext = false
		return nil
	}

	p.PageToken = string(pInfo.EndCursor)

	return nil
}

func (p *RecursivePager) HasNext() bool {
	return p.hasNext
}

func (p *RecursivePager) ListAll(ctx context.Context, output interface{}) error {
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
