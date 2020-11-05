package catalog_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/kyma-incubator/compass/components/system-broker/internal/director"
	"github.com/kyma-incubator/compass/components/system-broker/internal/osb"
	"github.com/kyma-incubator/compass/components/system-broker/pkg/env"
	"github.com/kyma-incubator/compass/components/system-broker/tests/common"
	"github.com/pivotal-cf/brokerapi/v7/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestOSBCatalog(t *testing.T) {
	suite.Run(t, new(OSBCatalogTestSuite))
}

type OSBCatalogTestSuite struct {
	suite.Suite
	testContext *common.TestContext
	configURL   string
}

func (suite *OSBCatalogTestSuite) SetupSuite() {
	suite.testContext = common.NewTestContextBuilder().
		WithEnvExtensions(func(e env.Environment, servers map[string]common.FakeServer) {
			e.Set("director_gql.page_size", 1)
			e.Set("director_gql.page_concurrency", 1)
		}).Build(suite.T())
	suite.configURL = suite.testContext.Servers[common.DirectorServer].URL() + "/config"
}

func (suite *OSBCatalogTestSuite) SetupTest() {
	suite.testContext.HttpClient.Post(suite.configURL+"/reset", "application/json", nil)
}

func (suite *OSBCatalogTestSuite) TearDownSuite() {
	suite.testContext.CleanUp()
}

func (suite *OSBCatalogTestSuite) TeardownTest() {
	resp, err := suite.testContext.HttpClient.Get(suite.configURL + "/verify")
	assert.NoError(suite.T(), err)

	if resp.StatusCode == http.StatusInternalServerError {
		errorMsg, err := ioutil.ReadAll(resp.Body)
		assert.NoError(suite.T(), err)
		suite.Fail(string(errorMsg))
	}
	assert.Equal(suite.T(), resp.StatusCode, http.StatusOK)
}

func (suite *OSBCatalogTestSuite) TestEmptyResponse() {
	expectedCatalog := mockData(suite, 0, 0, 0, 0, 0)
	suite.testContext.SystemBroker.GET("/v2/catalog").WithHeader("X-Broker-API-Version", "2.15").Expect().
		Status(http.StatusOK).JSON().Equal(expectedCatalog)
}

func (suite *OSBCatalogTestSuite) TestResponseWithOnlyApps() {
	expectedCatalog := mockData(suite, 1, 0, 0, 0, 0)

	suite.testContext.SystemBroker.GET("/v2/catalog").WithHeader("X-Broker-API-Version", "2.15").Expect().
		Status(http.StatusOK).
		JSON().Equal(expectedCatalog)
}

func (suite *OSBCatalogTestSuite) TestResponseWithOnePackagePage() {
	expectedCatalog := mockData(suite, 1, 1, 1, 1, 1)

	suite.testContext.SystemBroker.GET("/v2/catalog").WithHeader("X-Broker-API-Version", "2.15").Expect().
		Status(http.StatusOK).
		JSON().Equal(expectedCatalog)
}

func (suite *OSBCatalogTestSuite) TestResponseWithTwoPackagePage() {
	expectedCatalog := mockData(suite, 1, 2, 1, 1, 1)

	suite.testContext.SystemBroker.GET("/v2/catalog").WithHeader("X-Broker-API-Version", "2.15").Expect().
		Status(http.StatusOK).
		JSON().Equal(expectedCatalog)
}

// func (suite *OSBCatalogTestSuite) TestResponseWithTwoPackagePage() {
// 	expectedCatalog := mockData(suite, 1, 2, 1, 1, 1)

// 	suite.testContext.SystemBroker.GET("/v2/catalog").WithHeader("X-Broker-API-Version", "2.15").Expect().
// 		Status(http.StatusOK).
// 		JSON().Equal(expectedCatalog)
// }

func mockData(suite *OSBCatalogTestSuite, apps, packages, apiDefs, evDefs, docs int) interface{} {
	var appsResult director.ApplicationsOutput

	if apps == 0 {
		appsResponseString := toDirectorResponse(suite.T(), appsEmptyResponse)
		err := suite.testContext.ConfigureResponse(suite.configURL, "query", "applications", appsResponseString)
		assert.NoError(suite.T(), err)
	}

	for i := 0; i < apps; i++ {
		appObj := genMockApp(i != (apps - 1))
		appsResponseString := toDirectorResponse(suite.T(), appObj)
		err := suite.testContext.ConfigureResponse(suite.configURL, "query", "applications", appsResponseString)
		assert.NoError(suite.T(), err)
		appsResult = append(appsResult, appObj.Result.Data...)
	}

	for i := 0; i < apps; i++ {

		if packages == 0 {
			packageResponseString := toDirectorResponse(suite.T(), packagesEmptyResponse)
			err := suite.testContext.ConfigureResponse(suite.configURL, "query", "application", packageResponseString)
			assert.NoError(suite.T(), err)
		}

		var packagesResult director.PackagessOutput
		for j := 0; j < packages; j++ {
			packObj := genMockPackage(j != (packages - 1))
			packageResponseString := toDirectorResponse(suite.T(), packObj)
			err := suite.testContext.ConfigureResponse(suite.configURL, "query", "application", packageResponseString)
			assert.NoError(suite.T(), err)
			packagesResult = append(packagesResult, packObj.Result.Packages.Data...)
		}

		for j := 0; j < packages; j++ {
			var apiDefsResult director.ApiDefinitionsOutput
			var eventDefsResult director.EventDefinitionsOutput
			var docsResult director.DocumentsOutput

			if apiDefs == 0 {
				apiDefsResponseString := toDirectorResponse(suite.T(), apiDefEmptyResponse)
				err := suite.testContext.ConfigureResponse(suite.configURL, "query", "application", apiDefsResponseString)
				assert.NoError(suite.T(), err)
			}

			for k := 0; k < apiDefs; k++ {
				apiDefObj := genMockApiDef(k != (apiDefs - 1))
				apiDefsResponse := toDirectorResponse(suite.T(), apiDefObj)
				err := suite.testContext.ConfigureResponse(suite.configURL, "query", "application", apiDefsResponse)
				assert.NoError(suite.T(), err)
				apiDefsResult = append(apiDefsResult, apiDefObj.Result.Package.ApiDefinitions.Data...)
			}

			if evDefs == 0 {
				eventDefsResponseString := toDirectorResponse(suite.T(), eventDefEmptyResponse)
				err := suite.testContext.ConfigureResponse(suite.configURL, "query", "application", eventDefsResponseString)
				assert.NoError(suite.T(), err)
			}

			for m := 0; m < evDefs; m++ {
				eventDefObj := genMockEventDef(m != (evDefs - 1))
				eventDefResponse := toDirectorResponse(suite.T(), eventDefObj)
				err := suite.testContext.ConfigureResponse(suite.configURL, "query", "application", eventDefResponse)
				assert.NoError(suite.T(), err)
				eventDefsResult = append(eventDefsResult, eventDefObj.Result.Package.EventDefinitions.Data...)
			}

			if docs == 0 {
				docsResponseString := toDirectorResponse(suite.T(), docsEmptyResponse)
				err := suite.testContext.ConfigureResponse(suite.configURL, "query", "application", docsResponseString)
				assert.NoError(suite.T(), err)
			}

			for m := 0; m < docs; m++ {
				docObj := genMockDocs(m != (docs - 1))
				docsResponse := toDirectorResponse(suite.T(), docObj)
				err := suite.testContext.ConfigureResponse(suite.configURL, "query", "application", docsResponse)
				assert.NoError(suite.T(), err)
				docsResult = append(docsResult, docObj.Result.Package.Documents.Data...)
			}

			packagesResult[j].APIDefinitions.Data = apiDefsResult
			packagesResult[j].EventDefinitions.Data = eventDefsResult
			packagesResult[j].Documents.Data = docsResult
		}
		appsResult[i].Packages.Data = packagesResult
	}

	return toCatalog(suite.T(), appsResult)
}

// func (suite *OSBCatalogTestSuite) TestErrorWhileFetchingApplicaitons() {
// 	err := suite.testContext.ConfigureResponse(suite.configURL, "query", "applications", appsPageResponse1)
// 	assert.NoError(suite.T(), err)
// 	err = suite.testContext.ConfigureResponse(suite.configURL, "query", "applications", appsErrorResponse)
// 	assert.NoError(suite.T(), err)

// 	suite.testContext.SystemBroker.GET("/v2/catalog").WithHeader("X-Broker-API-Version", "2.15").Expect().
// 		Status(http.StatusInternalServerError).
// 		JSON().Object().Value("description").Equal("could not build catalog")
// }

func toDirectorResponse(t *testing.T, mockApp interface{}) string {
	fixture := map[string]interface{}{
		"data": mockApp,
	}

	appsEmptyResponseBytes, err := json.Marshal(fixture)
	assert.NoError(t, err)
	apps := string(appsEmptyResponseBytes)
	return apps
}

func toCatalog(t *testing.T, mockApp director.ApplicationsOutput) interface{} {
	converter := osb.Converter{}
	svcs := make([]domain.Service, 0)
	for _, app := range mockApp {
		s, err := converter.Convert(&app)
		assert.NoError(t, err)
		svcs = append(svcs, s...)
	}

	catalogObj := map[string]interface{}{
		"services": svcs,
	}

	return catalogObj
}

func deepCopy(t *testing.T, src interface{}, dest interface{}) {
	bytes, err := json.Marshal(src)
	assert.NoError(t, err)

	err = json.Unmarshal(bytes, &dest)
	assert.NoError(t, err)
}
