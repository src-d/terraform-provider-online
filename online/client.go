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
	rpnv2EndPoint  = "https://api.online.net/api/v1/rpn/v2"
)

type Client interface {
	Server(id int) (*Server, error)
	SetServer(s *Server) error

	ListRPNv2() ([]*RPNv2, error)
	RPNv2(id int) (*RPNv2, error)
	RPNv2ByName(name string) (*RPNv2, error)
	SetRPNv2(r *RPNv2) error
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

func (c *client) doGET(target string) ([]byte, error) {
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

	return js, nil
}

func (c *client) doPUT(target string, values map[string]string) error {
	return c.doBooleanRequest("PUT", target, values)
}

func (c *client) doPOST(target string, values map[string]string) error {
	return c.doBooleanRequest("POST", target, values)
}

func (c *client) doBooleanRequest(method, target string, values map[string]string) error {
	form := url.Values{}
	for k, v := range values {
		form.Add(k, v)
	}

	req, err := http.NewRequest(method, target, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	return c.handleBooleanResponse(resp)
}

func (c *client) handleBooleanResponse(r *http.Response) error {
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	bool := string(b)
	switch bool {
	case "true":
		return nil
	case "false":
		return fmt.Errorf("requested operation has failed")
	default:
		e := &ErrorResponse{}
		err := json.Unmarshal(b, e)
		if err != nil || e.Message == "" {
			return fmt.Errorf("unexpected boolean answer from server: %s", bool)
		}

		return e.Error()
	}
}

func (c *client) SetServer(s *Server) error {
	target := fmt.Sprintf("%s/%d", serverEndPoint, s.ID)
	err := c.doPUT(target, map[string]string{
		"hostname": s.Hostname,
	})

	if err != nil {
		return err
	}

	public := s.InterfaceByType(Public)
	if public == nil {
		return nil
	}

	return c.doSetServerIP(public)
}

func (c *client) doSetServerIP(i *Interface) error {
	target := fmt.Sprintf("%s/ip/edit", serverEndPoint)
	err := c.doPOST(target, map[string]string{
		"address": i.Address,
		"reverse": i.Reverse,
	})

	fmt.Println(target, err)
	return err
}

func (c *client) Server(id int) (*Server, error) {
	target := fmt.Sprintf("%s/%d", serverEndPoint, id)
	js, err := c.doGET(target)
	if err != nil {
		return nil, err
	}

	s := &Server{}
	return s, json.Unmarshal(js, s)
}

func (c *client) ListRPNv2() ([]*RPNv2, error) {
	js, err := c.doGET(rpnv2EndPoint)
	if err != nil {
		return nil, err
	}

	var list []*RPNv2
	return list, json.Unmarshal(js, &list)
}

func (c *client) RPNv2(id int) (*RPNv2, error) {
	target := fmt.Sprintf("%s/%d", rpnv2EndPoint, id)
	js, err := c.doGET(target)
	if err != nil {
		return nil, err
	}

	r := &RPNv2{}
	return r, json.Unmarshal(js, r)
}

func (c *client) RPNv2ByName(name string) (*RPNv2, error) {
	list, err := c.ListRPNv2()
	if err != nil {
		return nil, err
	}

	for _, rpn := range list {
		if rpn.Name == name {
			return rpn, nil
		}
	}

	return nil, nil
}

func (c *client) SetRPNv2(r *RPNv2) error {
	if r.ID == 0 {
		return c.doCreateRPNv2(r)
	}

	return nil
}

func (c *client) doCreateRPNv2(r *RPNv2) error {
	var ids []int
	for _, m := range r.Members {
		ids = append(ids, m.ID)
	}

	idsJSON, _ := json.Marshal(ids)

	fmt.Println(string(idsJSON))
	return c.doPOST(rpnv2EndPoint, map[string]string{
		"type":        string(r.Type),
		"description": r.Name,
		"server_ids":  string(idsJSON),
	})
}

type ErrorResponse struct {
	Code    int
	Message string `json:"error"`
}

func (e *ErrorResponse) Error() error {
	return fmt.Errorf("%s (code: %d)", e.Message, e.Code)
}
