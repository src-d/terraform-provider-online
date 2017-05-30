package online

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	serverEndPoint = "https://api.online.net/api/v1/server/"
)

type Client interface {
	Server(id string) (*Server, error)

	Do(req *http.Request) (*http.Response, error)
}

func NewClient(token string) Client {
	return &client{token: token, c: &http.Client{}}
}

type client struct {
	token string
	c     *http.Client
}

func (c *client) Do(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))
	return c.c.Do(req)
}

func (c *client) Server(id string) (*Server, error) {
	req, err := http.NewRequest("GET", serverEndPoint+id, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {

		return nil, err
	}

	defer resp.Body.Close()
	js, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	s := &Server{}
	return s, json.Unmarshal(js, s)
}
