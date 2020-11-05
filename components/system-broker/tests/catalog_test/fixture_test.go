package catalog_test

import (
	"github.com/google/uuid"
	schema "github.com/kyma-incubator/compass/components/director/pkg/graphql"
	"github.com/kyma-incubator/compass/components/system-broker/internal/director"
)

var appsEmptyResponse = director.ApplicationResponse{
	Result: struct {
		Data director.ApplicationsOutput `json:"data"`
		Page schema.PageInfo             `json:"pageInfo"`
	}{
		Data: director.ApplicationsOutput{},
		Page: schema.PageInfo{},
	},
}

// var appsMockResponse = director.ApplicationResponse{
// 	Result: struct {
// 		Data director.ApplicationsOutput `json:"data"`
// 		Page schema.PageInfo             `json:"pageInfo"`
// 	}{
// 		Data: director.ApplicationsOutput{
// 			schema.ApplicationExt{
// 				Application: schema.Application{
// 					Name: "test-app",
// 					ID:   "3e3cecce-74b3-4881-854e-58791021b522",
// 				},
// 			},
// 		},
// 		Page: schema.PageInfo{},
// 	},
// }

func genMockApp(withPage bool) director.ApplicationResponse {
	id := uuid.New().String()
	result := director.ApplicationResponse{
		Result: struct {
			Data director.ApplicationsOutput `json:"data"`
			Page schema.PageInfo             `json:"pageInfo"`
		}{
			Data: director.ApplicationsOutput{
				schema.ApplicationExt{
					Application: schema.Application{
						Name: "test-app-" + id,
						ID:   id,
					},
				},
			},
			Page: schema.PageInfo{},
		},
	}

	if withPage {
		result.Result.Page = pageInfoForMore
	}
	return result
}

func genMockPackage(withPage bool) director.PackagesResponse {
	id := uuid.New().String()
	result := director.PackagesResponse{
		Result: struct {
			Packages struct {
				Data director.PackagessOutput "json:\"data\""
				Page schema.PageInfo          "json:\"pageInfo\""
			} "json:\"packages\""
		}{
			Packages: struct {
				Data director.PackagessOutput "json:\"data\""
				Page schema.PageInfo          "json:\"pageInfo\""
			}{
				Data: director.PackagessOutput{
					&schema.PackageExt{
						Package: schema.Package{
							ID:   id,
							Name: "pack1-name-" + id,
						},
					},
				},
				Page: schema.PageInfo{},
			},
		},
	}

	if withPage {
		result.Result.Packages.Page = pageInfoForMore
	}
	return result
}

func genMockApiDef(withPage bool) director.ApiDefinitionsResponse {
	id := uuid.New().String()
	result := director.ApiDefinitionsResponse{
		Result: struct {
			Package struct {
				ApiDefinitions struct {
					Data director.ApiDefinitionsOutput "json:\"data\""
					Page schema.PageInfo               "json:\"pageInfo\""
				} "json:\"apiDefinitions\""
			} "json:\"package\""
		}{
			Package: struct {
				ApiDefinitions struct {
					Data director.ApiDefinitionsOutput "json:\"data\""
					Page schema.PageInfo               "json:\"pageInfo\""
				} "json:\"apiDefinitions\""
			}{
				ApiDefinitions: struct {
					Data director.ApiDefinitionsOutput "json:\"data\""
					Page schema.PageInfo               "json:\"pageInfo\""
				}{
					Data: director.ApiDefinitionsOutput{
						&schema.APIDefinitionExt{
							APIDefinition: schema.APIDefinition{
								ID:        id,
								Name:      "name-" + id,
								TargetURL: "target.url",
							},
						},
					},
					Page: schema.PageInfo{},
				},
			},
		},
	}

	if withPage {
		result.Result.Package.ApiDefinitions.Page = pageInfoForMore
	}
	return result
}

func genMockEventDef(withPage bool) director.EventDefinitionsResponse {
	id := uuid.New().String()
	result := director.EventDefinitionsResponse{
		Result: struct {
			Package struct {
				EventDefinitions struct {
					Data director.EventDefinitionsOutput "json:\"data\""
					Page schema.PageInfo                 "json:\"pageInfo\""
				} "json:\"eventDefinitions\""
			} "json:\"package\""
		}{
			Package: struct {
				EventDefinitions struct {
					Data director.EventDefinitionsOutput "json:\"data\""
					Page schema.PageInfo                 "json:\"pageInfo\""
				} "json:\"eventDefinitions\""
			}{
				EventDefinitions: struct {
					Data director.EventDefinitionsOutput "json:\"data\""
					Page schema.PageInfo                 "json:\"pageInfo\""
				}{
					Data: director.EventDefinitionsOutput{
						&schema.EventAPIDefinitionExt{
							EventDefinition: schema.EventDefinition{
								ID:   id,
								Name: "name-" + id,
							},
						},
					},
					Page: schema.PageInfo{},
				},
			},
		},
	}

	if withPage {
		result.Result.Package.EventDefinitions.Page = pageInfoForMore
	}
	return result
}

func genMockDocs(withPage bool) director.DocumentsResponse {
	id := uuid.New().String()
	result := director.DocumentsResponse{
		Result: struct {
			Package struct {
				Documents struct {
					Data director.DocumentsOutput "json:\"data\""
					Page schema.PageInfo          "json:\"pageInfo\""
				} "json:\"documents\""
			} "json:\"package\""
		}{
			Package: struct {
				Documents struct {
					Data director.DocumentsOutput "json:\"data\""
					Page schema.PageInfo          "json:\"pageInfo\""
				} "json:\"documents\""
			}{
				Documents: struct {
					Data director.DocumentsOutput "json:\"data\""
					Page schema.PageInfo          "json:\"pageInfo\""
				}{
					Data: director.DocumentsOutput{
						&schema.DocumentExt{
							Document: schema.Document{
								ID:          id,
								DisplayName: "display-name-" + id,
							},
						},
					},
					Page: schema.PageInfo{},
				},
			},
		},
	}

	if withPage {
		result.Result.Package.Documents.Page = pageInfoForMore
	}
	return result
}

var packagesEmptyResponse = director.PackagesResponse{
	Result: struct {
		Packages struct {
			Data director.PackagessOutput "json:\"data\""
			Page schema.PageInfo          "json:\"pageInfo\""
		} "json:\"packages\""
	}{
		Packages: struct {
			Data director.PackagessOutput "json:\"data\""
			Page schema.PageInfo          "json:\"pageInfo\""
		}{
			Data: director.PackagessOutput{},
			Page: schema.PageInfo{},
		},
	},
}

var apiDefEmptyResponse = director.ApiDefinitionsResponse{
	Result: struct {
		Package struct {
			ApiDefinitions struct {
				Data director.ApiDefinitionsOutput "json:\"data\""
				Page schema.PageInfo               "json:\"pageInfo\""
			} "json:\"apiDefinitions\""
		} "json:\"package\""
	}{
		Package: struct {
			ApiDefinitions struct {
				Data director.ApiDefinitionsOutput "json:\"data\""
				Page schema.PageInfo               "json:\"pageInfo\""
			} "json:\"apiDefinitions\""
		}{
			ApiDefinitions: struct {
				Data director.ApiDefinitionsOutput "json:\"data\""
				Page schema.PageInfo               "json:\"pageInfo\""
			}{
				Data: director.ApiDefinitionsOutput{},
				Page: schema.PageInfo{},
			},
		},
	},
}

var eventDefEmptyResponse = director.EventDefinitionsResponse{
	Result: struct {
		Package struct {
			EventDefinitions struct {
				Data director.EventDefinitionsOutput "json:\"data\""
				Page schema.PageInfo                 "json:\"pageInfo\""
			} "json:\"eventDefinitions\""
		} "json:\"package\""
	}{
		Package: struct {
			EventDefinitions struct {
				Data director.EventDefinitionsOutput "json:\"data\""
				Page schema.PageInfo                 "json:\"pageInfo\""
			} "json:\"eventDefinitions\""
		}{
			EventDefinitions: struct {
				Data director.EventDefinitionsOutput "json:\"data\""
				Page schema.PageInfo                 "json:\"pageInfo\""
			}{
				Data: director.EventDefinitionsOutput{},
				Page: schema.PageInfo{},
			},
		},
	},
}

var docsEmptyResponse = director.DocumentsResponse{
	Result: struct {
		Package struct {
			Documents struct {
				Data director.DocumentsOutput "json:\"data\""
				Page schema.PageInfo          "json:\"pageInfo\""
			} "json:\"documents\""
		} "json:\"package\""
	}{
		Package: struct {
			Documents struct {
				Data director.DocumentsOutput "json:\"data\""
				Page schema.PageInfo          "json:\"pageInfo\""
			} "json:\"documents\""
		}{
			Documents: struct {
				Data director.DocumentsOutput "json:\"data\""
				Page schema.PageInfo          "json:\"pageInfo\""
			}{
				Data: director.DocumentsOutput{},
				Page: schema.PageInfo{},
			},
		},
	},
}

var onePackageResponse = director.PackagesResponse{
	Result: struct {
		Packages struct {
			Data director.PackagessOutput "json:\"data\""
			Page schema.PageInfo          "json:\"pageInfo\""
		} "json:\"packages\""
	}{
		Packages: struct {
			Data director.PackagessOutput "json:\"data\""
			Page schema.PageInfo          "json:\"pageInfo\""
		}{
			Data: director.PackagessOutput{
				&schema.PackageExt{
					Package: schema.Package{
						ID:   "pack1-id",
						Name: "pack1-name",
					},
				},
			},
			Page: schema.PageInfo{},
		},
	},
}

var pageInfoForMore = schema.PageInfo{
	StartCursor: "",
	EndCursor:   "next",
	HasNextPage: true,
}

const (
	appsMockResponsed = `{
  "data": {
    "result": {
      "data": [
        {
          "id": "3e3cecce-74b3-4881-854e-58791021b522",
          "name": "test-app",
          "providerName": "test provider",
          "description": "a test application",
          "integrationSystemID": null,
          "labels": {
            "group": [
              "production",
              "experimental"
            ],
            "integrationSystemID": "",
            "name": "test-app",
            "scenarios": [
              "DEFAULT"
            ]
          },
          "status": {
            "condition": "INITIAL",
            "timestamp": "2020-10-21T18:23:59Z"
          },
          "webhooks": null,
          "healthCheckURL": "http://test-app.com/health",
          "packages": {
            "data": [],
            "pageInfo": {
              "startCursor": "",
              "endCursor": "",
              "hasNextPage": false
            },
            "totalCount": 0
          },
          "auths": null,
          "eventingConfiguration": {
            "defaultURL": ""
          }
        }
      ],
      "pageInfo": {
        "startCursor": "",
        "endCursor": "",
        "hasNextPage": false
      },
      "totalCount": 1
    }
  }
}`
	appsPageResponse1 = `{
  "data": {
    "result": {
      "data": [
        {
          "id": "4310851a-04a0-4217-a2b8-1766c6a3f0fe",
          "name": "test-app2",
          "providerName": "test provider",
          "description": "a test application",
          "integrationSystemID": null,
          "labels": {
            "group": [
              "production",
              "experimental"
            ],
            "integrationSystemID": "",
            "name": "test-app2",
            "scenarios": [
              "DEFAULT"
            ]
          },
          "status": {
            "condition": "INITIAL",
            "timestamp": "2020-10-22T12:17:57Z"
          },
          "webhooks": null,
          "healthCheckURL": "http://test-app.com/health",
          "packages": {
            "data": [],
            "pageInfo": {
              "startCursor": "",
              "endCursor": "",
              "hasNextPage": false
            },
            "totalCount": 0
          },
          "auths": null,
          "eventingConfiguration": {
            "defaultURL": ""
          }
        },
        {
          "id": "5cf79030-0433-4b32-8618-9844085ca7a6",
          "name": "test-app",
          "providerName": "test provider",
          "description": "a test application",
          "integrationSystemID": null,
          "labels": {
            "group": [
              "production",
              "experimental"
            ],
            "integrationSystemID": "",
            "name": "test-app",
            "scenarios": [
              "DEFAULT"
            ]
          },
          "status": {
            "condition": "INITIAL",
            "timestamp": "2020-10-22T12:17:50Z"
          },
          "webhooks": null,
          "healthCheckURL": "http://test-app.com/health",
          "packages": {
            "data": [],
            "pageInfo": {
              "startCursor": "",
              "endCursor": "",
              "hasNextPage": false
            },
            "totalCount": 0
          },
          "auths": null,
          "eventingConfiguration": {
            "defaultURL": ""
          }
        }
      ],
      "pageInfo": {
        "startCursor": "",
        "endCursor": "RHBLdEo0ajlqRHEy",
        "hasNextPage": true
      },
      "totalCount": 4
    }
  }
}`
	appsPageResponse2 = `{
  "data": {
    "result": {
      "data": [
        {
          "id": "75ab9628-24d1-4e39-bdae-9c5042e908f2",
          "name": "test-app3",
          "providerName": "test provider",
          "description": "a test application",
          "integrationSystemID": null,
          "labels": {
            "group": [
              "production",
              "experimental"
            ],
            "integrationSystemID": "",
            "name": "test-app3",
            "scenarios": [
              "DEFAULT"
            ]
          },
          "status": {
            "condition": "INITIAL",
            "timestamp": "2020-10-22T12:18:00Z"
          },
          "webhooks": null,
          "healthCheckURL": "http://test-app.com/health",
          "packages": {
            "data": [],
            "pageInfo": {
              "startCursor": "",
              "endCursor": "",
              "hasNextPage": false
            },
            "totalCount": 0
          },
          "auths": null,
          "eventingConfiguration": {
            "defaultURL": ""
          }
        },
        {
          "id": "f945951c-bcaf-46af-a017-b3e2b575bdbd",
          "name": "test-app1",
          "providerName": "test provider",
          "description": "a test application",
          "integrationSystemID": null,
          "labels": {
            "group": [
              "production",
              "experimental"
            ],
            "integrationSystemID": "",
            "name": "test-app1",
            "scenarios": [
              "DEFAULT"
            ]
          },
          "status": {
            "condition": "INITIAL",
            "timestamp": "2020-10-22T12:17:53Z"
          },
          "webhooks": null,
          "healthCheckURL": "http://test-app.com/health",
          "packages": {
            "data": [],
            "pageInfo": {
              "startCursor": "",
              "endCursor": "",
              "hasNextPage": false
            },
            "totalCount": 0
          },
          "auths": null,
          "eventingConfiguration": {
            "defaultURL": ""
          }
        }
      ],
      "pageInfo": {
        "startCursor": "RHBLdEo0ajlqRHEy",
        "endCursor": "",
        "hasNextPage": false
      },
      "totalCount": 4
    }
  }
}`
	appsExpectedCatalogPaging = `{"services":[{"id":"4310851a-04a0-4217-a2b8-1766c6a3f0fe","name":"test-app2","description":"a test application","bindable":true,"plan_updateable":false,"plans":null,"metadata":{"displayName":"test-app2","group":["production","experimental"],"integrationSystemID":"","name":"test-app2","providerDisplayName":"test provider","scenarios":["DEFAULT"]}},{"id":"5cf79030-0433-4b32-8618-9844085ca7a6","name":"test-app","description":"a test application","bindable":true,"plan_updateable":false,"plans":null,"metadata":{"displayName":"test-app","group":["production","experimental"],"integrationSystemID":"","name":"test-app","providerDisplayName":"test provider","scenarios":["DEFAULT"]}},{"id":"75ab9628-24d1-4e39-bdae-9c5042e908f2","name":"test-app3","description":"a test application","bindable":true,"plan_updateable":false,"plans":null,"metadata":{"displayName":"test-app3","group":["production","experimental"],"integrationSystemID":"","name":"test-app3","providerDisplayName":"test provider","scenarios":["DEFAULT"]}},{"id":"f945951c-bcaf-46af-a017-b3e2b575bdbd","name":"test-app1","description":"a test application","bindable":true,"plan_updateable":false,"plans":null,"metadata":{"displayName":"test-app1","group":["production","experimental"],"integrationSystemID":"","name":"test-app1","providerDisplayName":"test provider","scenarios":["DEFAULT"]}}]}` + "\n"
	appsErrorResponse         = `{
  "errors": [
    {
      "message": "Internal Server Error",
      "path": [
        "result"
      ],
      "extensions": {
        "error": "InternalError",
        "error_code": 10
      }
    }
  ],
  "data": null
}`
)
