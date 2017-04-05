package util

import (
	"common/logging"
	"net/http"
	"encoding/json"
	"strconv"
	"io"
	"io/ioutil"
	"fmt"
	"errors"
)

func ParseJsonBody(body io.Reader, v interface{}) error {
	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, v)
	if err != nil {
		return errors.New(fmt.Sprintf("error: %s, input: %s", err.Error(), string(bytes)))
	}
	return nil
}

func ResponseJson(w http.ResponseWriter, v interface{}) {
	content, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Length", strconv.Itoa(len(content)))
	w.Header().Set("Content-Type", "application/json")
	w.Write(content)
	
	logging.Debug("Response %s", content)
}