package util

import (
	"encoding/json"
	"net/http"
)

func ReplyJson(w http.ResponseWriter, httpCode int, resp interface{}) {
	w.WriteHeader(httpCode)
	w.Header().Set("Content-Type", "application/json")
	respJson, _ := json.Marshal(resp)
	w.Write(respJson)
}
