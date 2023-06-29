package render

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func DecodeJSON(req *http.Request, data any) error {
	if !strings.Contains(req.Header.Get("Content-Type"), "application/json") {
		return fmt.Errorf("Content-Type header is not application/json")
	}
	dec := json.NewDecoder(req.Body)
	//dec numbers as numbers instead of floats
	dec.UseNumber()
	if err := dec.Decode(data); err != nil {
		return err
	}
	defer req.Body.Close()

	return nil
}
