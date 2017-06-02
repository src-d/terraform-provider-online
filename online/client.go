package online

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
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

func (c *client) doPATCH(target string, values map[string]string) error {
	return c.doBooleanRequest("PATCH", target, values)
}

func (c *client) doDELETE(target string, values map[string]string) error {
	return c.doBooleanRequest("DELETE", target, values)
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
	var err error
	if r.ID == 0 {
		err = c.doCreateRPNv2(r)
	} else {
		err = c.doUpdateRPNv2(r)
	}

	if err != nil {
		return err
	}

	return c.doSyncVLAN(r)
}

func (c *client) doCreateRPNv2(r *RPNv2) error {
	var ids []int
	for _, m := range r.Members {
		ids = append(ids, m.ID)
	}

	idsJSON, _ := json.Marshal(ids)
	return c.doPOST(rpnv2EndPoint, map[string]string{
		"type":        string(r.Type),
		"description": r.Name,
		"server_ids":  string(idsJSON),
	})
}

func (c *client) doUpdateRPNv2(r *RPNv2) error {
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

	return nil
}

func (c *client) doAddMembers(r *RPNv2, serverIDs []int) error {
	if len(serverIDs) == 0 {
		return nil
	}

	target := fmt.Sprintf("%s/%d/addMember", rpnv2EndPoint, r.ID)
	idsJSON, _ := json.Marshal(serverIDs)
	err := c.doPOST(target, map[string]string{
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
	err := c.doDELETE(target, map[string]string{
		"server_ids": string(idsJSON),
	})

	return err
}

func (c *client) doSyncVLAN(r *RPNv2) error {
	prev, err := c.RPNv2(r.ID)
	if err != nil {
		return err
	}

	for _, new := range r.Members {
		old := prev.MemberByServerID(new.Linked.ID)
		if old == nil {
			continue
		}

		if new.Vlan != old.Vlan {
			new.ID = old.ID
			if err := c.doEditVlanMember(r.ID, new); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *client) doEditVlanMember(groupID int, m *Member) error {
	target := fmt.Sprintf("%s/%d/editVlanMember/%d", rpnv2EndPoint, groupID, m.ID)
	fmt.Println(target)
	return c.doPATCH(target, map[string]string{
		"vlan_number": strconv.Itoa(m.Vlan),
		"reset_vlan":  "false",
	})
}

type ErrorResponse struct {
	Code    int
	Message string `json:"error"`
}

func (e *ErrorResponse) Error() error {
	return fmt.Errorf("%s (code: %d)", e.Message, e.Code)
}
