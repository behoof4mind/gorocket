// This package provides a RocketChat client. It tries to use the rest api whenever possible and
// ppd only for features that are exclusive to the web client.
//   	client := Client{Protocol: "http", Host: "127.0.0.1", Port: "3000"}
// You have to login to interact with a RocketChat server:
//      credentials := UserCredentials{Email: "user@mail.com", Password: "secret"}
//      client.Login(credentials)
package gorocket

import (
	"net/http"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"fmt"
)

type Client struct {
	Protocol string
	Host  string
	Port  string

	// Use this switch to see all network communication.
	Debug bool

	auth  *authInfo
}

type authInfo struct {
	token string
	id    string
}

// The base for the most of the json responses
type statusResponse struct {
	Status string `json:"status"`
	Message string `json:"message"`
}

func (c *Client) getUrl() string {
	return fmt.Sprintf("%v://%v:%v", c.Protocol, c.Host, c.Port)
}

func (c *Client) doRequest(request *http.Request, responseBody interface{}) error {

	if c.auth != nil {
		request.Header.Set("X-Auth-Token", c.auth.token)
		request.Header.Set("X-User-Id", c.auth.id)
	}

	if c.Debug {
		log.Println(request)
	}

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()
	bodyBytes, err := ioutil.ReadAll(response.Body)

	if c.Debug {
		log.Println(string(bodyBytes))
	}

	if response.StatusCode != http.StatusOK {
		return errors.New("Request error: " + response.Status)
	}

	if err != nil {
		return err
	}

	return json.Unmarshal(bodyBytes, responseBody)
}