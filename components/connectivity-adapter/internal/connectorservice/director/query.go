package director

import "fmt"

func viewerQuery() string {
	return `query{
 		result: viewer
        {
 			 id
		}
   }`
}

func applicationQuery(appID string) string {
	return fmt.Sprintf(`query{
 		result: application(id: "%s")
        {
			 id
 			 name
			 labels
			 eventingConfiguration {
  				defaultURL
			}
		}
   }`, appID)
}
