package vcago

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

//AdminRequest represents model for admin requests.
type AdminRequest struct {
	URL string
}

func NewAdminRequest() *AdminRequest {
	return &AdminRequest{
		URL: Config.GetEnvString("ADMIN_URL", "n", "http://172.4.5.3"),
	}
}

func (i *AdminRequest) Get(path string) (r *Response, err error) {
	url := i.URL + path
	request := new(http.Request)
	request, err = http.NewRequest("GET", url, nil)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	response := new(http.Response)
	if response, err = client.Do(request); err != nil {
		return
	}
	defer response.Body.Close()
	var bodyBytes []byte
	if response.StatusCode != 200 {
		if bodyBytes, err = ioutil.ReadAll(response.Body); err != nil {
			return
		}
		body := new(Response)
		if err = json.Unmarshal(bodyBytes, body); err != nil {
			return
		}
		return
	}
	r = new(Response)
	if bodyBytes, err = ioutil.ReadAll(response.Body); err != nil {
		return
	}
	if err = json.Unmarshal(bodyBytes, r); err != nil {
		log.Print(err)
		return
	}
	return
}

func (i *AdminRequest) Post(path string, data interface{}) (r *Response, err error) {
	var jsonData []byte
	if jsonData, err = json.Marshal(data); err != nil {
		return
	}
	url := i.URL + path
	request := new(http.Request)
	if request, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonData)); err != nil {
		return
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	response := new(http.Response)
	response, err = client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	var bodyBytes []byte
	if response.StatusCode != 201 {
		if response.StatusCode == 404 {
			return nil, errors.New("request failed with status 404")
		}
		if bodyBytes, err = ioutil.ReadAll(response.Body); err != nil {
			return
		}
		body := new(interface{})
		if err = json.Unmarshal(bodyBytes, body); err != nil {
			return
		}
		return nil, errors.New(response.Status)
	}
	r = new(Response)
	if bodyBytes, err = ioutil.ReadAll(response.Body); err != nil {
		return
	}
	if err = json.Unmarshal(bodyBytes, r); err != nil {
		return
	}
	return

}
