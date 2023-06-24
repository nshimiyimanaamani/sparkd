package render

import (
	"encoding/json"
	"net/http"
	"reflect"
)

func JSON(w http.ResponseWriter, v any, status int) {
	if rv := reflect.ValueOf(v); rv.Type().Kind() == reflect.Slice && rv.IsNil() {
		v = []struct{}{}
	}
	var result []byte
	var err error
	result, err = json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(result)
}
