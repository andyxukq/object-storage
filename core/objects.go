package core

import (
	"encoding/json"
	"net/http"

	//"labix.org/v2/mgo/bson"
)

type Response struct {
	Code    int     `json:"code"`
	Message string  `json:"message"`
	Data    Objects `json:"data,omitempty"`
}

func (self *Response) Set(code int, msg string) {
	self.Code = code
	self.Message = msg
}

func (self *Response) Write(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(self.Code)

	body, err := json.Marshal(self)
	if err != nil {
		return err
	}

	w.Write(body)
	return nil
}

type Objects struct {
	Files []File `json:"files"`
}

type File struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}
