package scenarioassignment_test

import (
	"context"
	"database/sql/driver"
	"errors"

	"github.com/kyma-incubator/compass/components/director/pkg/graphql/externalschema"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/kyma-incubator/compass/components/director/internal/domain/scenarioassignment"
	"github.com/kyma-incubator/compass/components/director/internal/domain/tenant"
	"github.com/kyma-incubator/compass/components/director/internal/model"

	"github.com/kyma-incubator/compass/components/director/pkg/pagination"
)

const (
	tenantID         = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	externalTenantID = "eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"
	scenarioName     = "scenario-A"
	errMsg           = "some error"
)

func fixModel() model.AutomaticScenarioAssignment {
	return fixModelWithScenarioName(scenarioName)
}

func fixGQL() externalschema.AutomaticScenarioAssignment {
	return fixGQLWithScenarioName(scenarioName)
}

var testTableColumns = []string{"scenario", "tenant_id", "selector_key", "selector_value"}

func fixModelWithScenarioName(scenario string) model.AutomaticScenarioAssignment {
	return model.AutomaticScenarioAssignment{
		ScenarioName: scenario,
		Tenant:       tenantID,
		Selector: model.LabelSelector{
			Key:   "key",
			Value: "value",
		},
	}
}

func fixModelWithScenarioNameAndSelector(scenario string, selector model.LabelSelector) model.AutomaticScenarioAssignment {
	return model.AutomaticScenarioAssignment{
		ScenarioName: scenario,
		Tenant:       tenantID,
		Selector:     selector,
	}
}

func fixModelPage() model.AutomaticScenarioAssignmentPage {
	mod1 := fixModelWithScenarioName("foo")
	mod2 := fixModelWithScenarioName("bar")
	modItems := []*model.AutomaticScenarioAssignment{
		&mod1, &mod2,
	}
	return fixModelPageWithItems(modItems)
}

func fixModelPageWithItems(in []*model.AutomaticScenarioAssignment) model.AutomaticScenarioAssignmentPage {
	return model.AutomaticScenarioAssignmentPage{
		Data:       in,
		PageInfo:   &pagination.Page{},
		TotalCount: len(in),
	}
}

func fixGQLWithScenarioName(scenario string) externalschema.AutomaticScenarioAssignment {
	return externalschema.AutomaticScenarioAssignment{
		ScenarioName: scenario,
		Selector: &externalschema.Label{
			Key:   "key",
			Value: "value",
		},
	}
}

func fixGQLPage() externalschema.AutomaticScenarioAssignmentPage {
	gql1 := fixGQLWithScenarioName("foo")
	gql2 := fixGQLWithScenarioName("bar")
	gqlItems := []*externalschema.AutomaticScenarioAssignment{
		&gql1, &gql2,
	}
	return fixGQLPageWithItems(gqlItems)
}

func fixGQLPageWithItems(in []*externalschema.AutomaticScenarioAssignment) externalschema.AutomaticScenarioAssignmentPage {
	return externalschema.AutomaticScenarioAssignmentPage{
		Data:       in,
		PageInfo:   &externalschema.PageInfo{},
		TotalCount: len(in),
	}
}

func fixEntity() scenarioassignment.Entity {
	return scenarioassignment.Entity{
		Scenario:      scenarioName,
		TenantID:      tenantID,
		SelectorKey:   "key",
		SelectorValue: "value",
	}
}

func fixEntityWithScenarioName(scenario string) scenarioassignment.Entity {
	return scenarioassignment.Entity{
		Scenario:      scenario,
		TenantID:      tenantID,
		SelectorKey:   "key",
		SelectorValue: "value",
	}
}

func fixError() error {
	return errors.New(errMsg)
}

func fixCtxWithTenant() context.Context {
	return tenant.SaveToContext(context.TODO(), tenantID, externalTenantID)
}

func fixLabelSelector() model.LabelSelector {
	return model.LabelSelector{
		Key:   "key",
		Value: "value",
	}
}

type sqlRow struct {
	scenario      string
	tenantId      string
	selectorKey   string
	selectorValue string
}

func fixSQLRows(rows []sqlRow) *sqlmock.Rows {
	out := sqlmock.NewRows(testTableColumns)
	for _, row := range rows {
		out.AddRow(row.scenario, row.tenantId, row.selectorKey, row.selectorValue)
	}
	return out
}

func fixAutomaticScenarioAssignmentRow(scenarioName, tenantID string) []driver.Value {
	return []driver.Value{scenarioName, tenantID, "key", "value"}
}

func fixAutomaticScenarioAssignmentColumns() []string {
	return []string{"scenario", "tenant_id", "selector_key", "selector_value"}
}
