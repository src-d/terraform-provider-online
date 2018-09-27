package online

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type responseType int

const (
	serverEndPoint = "https://api.online.net/api/v1/server"
	rpnv2EndPoint  = "https://api.online.net/api/v1/rpn/v2"

	responseBoolean responseType = iota
	responseJSON
	responseString
)

type Client interface {
	Server(id int) (*Server, error)
	SetServer(s *Server) error

	BootRescueMode(serverID int, image string) (*RescueCredentials, error)
	BootNormalMode(serverID int) error

	GetRescueImages(serverID int) ([]string, error)

	ListRPNv2() ([]*RPNv2, error)
	RPNv2(id int) (*RPNv2, error)
	RPNv2ByName(name string) (*RPNv2, error)
	SetRPNv2(r *RPNv2, wait time.Duration) error
	DeleteRPNv2(id int, wait time.Duration) error
}

func NewClient(token string) Client {
	return &client{token: token, c: &http.Client{}}
}

type client struct {
	token string
	c     *http.Client

	// rpn changes are controlled by a mutex
	rpnWriteLock sync.Mutex
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

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	return c.handleResponse(res)
}

func (c *client) doPUT(target string, values map[string]string) ([]byte, error) {
	return c.doRequest("PUT", target, values)
}

func (c *client) doPOST(target string, values map[string]string) ([]byte, error) {
	return c.doRequest("POST", target, values)
}

func (c *client) doPATCH(target string, values map[string]string) ([]byte, error) {
	return c.doRequest("PATCH", target, values)
}

func (c *client) doDELETE(target string, values map[string]string) ([]byte, error) {
	return c.doRequest("DELETE", target, values)
}

func (c *client) doRequest(method, target string, values map[string]string) ([]byte, error) {
	form := url.Values{}
	for k, v := range values {
		form.Add(k, v)
	}

	req, err := http.NewRequest(method, target, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	return c.handleResponse(resp)
}

func (c *client) handleResponse(r *http.Response) ([]byte, error) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	if r.StatusCode >= 200 && r.StatusCode <= 209 {
		return body, nil
	}

	return nil, decodeErrorResponse(body)
}

func (c *client) SetServer(s *Server) error {

	target := fmt.Sprintf("%s/%d", serverEndPoint, s.ID)
	_, err := c.doPUT(target, map[string]string{
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
	_, err := c.doPOST(target, map[string]string{
		"address": i.Address,
		"reverse": i.Reverse,
	})

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

func (c *client) SetRPNv2(r *RPNv2, wait time.Duration) error {
	c.rpnWriteLock.Lock()
	defer c.rpnWriteLock.Unlock()

	var err error
	if r.ID == 0 {
		err = c.doCreateRPNv2(r, wait)
	} else {
		err = c.doUpdateRPNv2(r, wait)
	}

	if err != nil {
		return err
	}

	return c.doSyncVLAN(r, wait)
}

func (c *client) doCreateRPNv2(r *RPNv2, wait time.Duration) error {
	ids := []int{}
	for _, m := range r.Members {
		ids = append(ids, m.Linked.ID)
	}

	idsJSON, _ := json.Marshal(ids)
	js, err := c.doPOST(rpnv2EndPoint, map[string]string{
		"type":        string(r.Type),
		"description": r.Name,
		"server_ids":  string(idsJSON),
	})

	if err != nil {
		return err
	}

	if err := json.Unmarshal(js, r); err != nil {
		return err
	}

	return c.waitRPNv2(r.ID, wait)
}

func (c *client) doUpdateRPNv2(r *RPNv2, wait time.Duration) error {
	prev, err := c.RPNv2(r.ID)
	if err != nil {
		return err
	}

	if r.Type != prev.Type {
		return fmt.Errorf("rpn type can't changed after creation")
	}

	var toAdd []int
	for _, new := range r.Members {
		if prev.MemberByServerID(new.Linked.ID) != nil {
			continue
		}

		toAdd = append(toAdd, new.Linked.ID)
	}

	if err := c.doAddMembers(r, toAdd); err != nil {
		return err
	}

	var toDelete []int
	for _, old := range prev.Members {
		if r.MemberByServerID(old.Linked.ID) != nil {
			continue
		}

		toDelete = append(toDelete, old.Linked.ID)
	}

	if err := c.doRemoveMembers(r, toDelete); err != nil {
		return err
	}

	return c.waitRPNv2(r.ID, wait)
}

func (c *client) doAddMembers(r *RPNv2, serverIDs []int) error {
	if len(serverIDs) == 0 {
		return nil
	}

	target := fmt.Sprintf("%s/%d/addMember", rpnv2EndPoint, r.ID)
	idsJSON, _ := json.Marshal(serverIDs)
	_, err := c.doPOST(target, map[string]string{
		"server_ids": string(idsJSON),
	})

	return err
}

func (c *client) doRemoveMembers(r *RPNv2, serverIDs []int) error {
	if len(serverIDs) == 0 {
		return nil
	}

	target := fmt.Sprintf("%s/%d/removeMember", rpnv2EndPoint, r.ID)
	idsJSON, _ := json.Marshal(serverIDs)
	_, err := c.doDELETE(target, map[string]string{
		"server_ids": string(idsJSON),
	})

	return err
}

func (c *client) doSyncVLAN(r *RPNv2, wait time.Duration) error {
	prev, err := c.RPNv2(r.ID)
	if err != nil {
		return err
	}

	for _, new := range r.Members {
		old := prev.MemberByServerID(new.Linked.ID)
		if old == nil {
			continue
		}

		if new.VLAN != old.VLAN {
			new.ID = old.ID
			if err := c.doEditVlanMember(r.ID, new); err != nil {
				return err
			}
		}
	}

	return c.waitRPNv2(r.ID, wait)
}

func (c *client) doEditVlanMember(groupID int, m *Member) error {
	target := fmt.Sprintf("%s/%d/editVlanMember/%d", rpnv2EndPoint, groupID, m.ID)
	_, err := c.doPATCH(target, map[string]string{
		"vlan_number": strconv.Itoa(m.VLAN),
		"reset_vlan":  "false",
	})

	return err
}

func (c *client) waitRPNv2(id int, wait time.Duration) error {
	until := time.Now().Add(wait)

	for now := range time.Tick(time.Second) {
		rpn, err := c.RPNv2(id)
		if err != nil {
			return err
		}

		membersUpdating := false
		for _, m := range rpn.Members {
			if m.Status != "ACTIVE" {
				membersUpdating = true
			}
		}

		if !membersUpdating && rpn.Status == "ACTIVE" {
			return nil
		}

		if now.After(until) {
			return fmt.Errorf("timeout waiting for RPNv2 changes")
		}
	}

	return nil
}

func (c *client) DeleteRPNv2(id int, wait time.Duration) error {
	c.rpnWriteLock.Lock()
	defer c.rpnWriteLock.Unlock()

	target := fmt.Sprintf("%s/%d", rpnv2EndPoint, id)
	_, err := c.doDELETE(target, nil)
	if err != nil {
		return err
	}

	err = c.waitRPNv2(id, wait)
	if err == nil {
		return nil
	}

	if er, ok := err.(*ErrorResponse); ok {
		// 7 is the code of RPNv2 not found error
		if er.Code == 7 {
			return nil
		}
	}

	return err

}

type ErrorResponse struct {
	Code    int
	Message string `json:"error"`
}

func decodeErrorResponse(b []byte) error {
	e := &ErrorResponse{}

	values := map[string]interface{}{}
	err := json.Unmarshal(b, &values)
	if err != nil {
		goto Unexpected
	}

	if msg, ok := values["error"]; ok {
		if e.Message, ok = msg.(string); !ok {
			goto Unexpected
		}
	}

	if msg, ok := values["error_description"]; ok {
		if e.Message, ok = msg.(string); !ok {
			goto Unexpected
		}
	}

	if code, ok := values["code"]; ok {
		code, ok := code.(float64)
		if !ok {
			goto Unexpected
		}

		e.Code = int(code)
	}

	return e

Unexpected:
	return fmt.Errorf("unexpected answer from server: %s", b)
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("%s (code: %d)", e.Message, e.Code)
}
