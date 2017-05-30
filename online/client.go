package online

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	serverEndPoint = "https://api.online.net/api/v1/server"
)

type Client interface {
	Server(id int) (*Server, error)
	SetServer(s *Server) error
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

func (c *client) doPUT(target string, values map[string]string) (*http.Response, error) {
	form := url.Values{}
	for k, v := range values {
		form.Add(k, v)
	}

	req, err := http.NewRequest("PUT", target, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return c.Do(req)
}

func (c *client) SetServer(s *Server) error {
	target := fmt.Sprintf("%s/%d", serverEndPoint, s.ID)
	resp, err := c.doPUT(target, map[string]string{
		"hostname": s.Hostname,
	})

	fmt.Println(resp)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	js, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(js))
	return nil
}

func (c *client) Server(id int) (*Server, error) {
	target := fmt.Sprintf("%s/%d", serverEndPoint, id)
	fmt.Println(target)

	req, err := http.NewRequest("GET", target, nil)
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

	fmt.Println(string(js))

	s := &Server{}
	return s, json.Unmarshal(js, s)
}
