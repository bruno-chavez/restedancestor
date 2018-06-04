// Package lib contains various functions that may be useful for more than one part of the API.
package lib

import "encoding/json"

/*
//BadRequest returns a ready to use badrequest text, not used for now
func BadRequest(requestText string) []byte{
	badRequest := map[string]string{"code": "400", "message": requestText}
	badRequestJson, _ := json.MarshalIndent(badRequest, "", "")
	return badRequestJson
}
*/
func NotFound(notFoundWord string) []byte {
	notFound := map[string]string{"code": "404", "message": "'" + notFoundWord + "' was not found in the database"}
	notFoundJSON, _ := json.MarshalIndent(notFound, "", "")

	return notFoundJSON
}
