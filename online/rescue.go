package online

import (
	"encoding/json"
	"fmt"
)

// RescueCredentials contain the login details for a server that booted into rescue mode
type RescueCredentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Protocol string `json:"protocol"`
	IP       string `json:"ip"`
}

func (c *client) BootRescueMode(serverID int, image string) (*RescueCredentials, error) {
	target := fmt.Sprintf("%s/boot/rescue/%d", serverEndPoint, serverID)
	body, err := c.doPOST(target, map[string]string{
		"image": image,
	})
	if err != nil {
		return nil, err
	}

	credentials := RescueCredentials{}
	err = json.Unmarshal(body, &credentials)
	if err != nil {
		return nil, err
	}

	return &credentials, nil
}

func (c *client) BootNormalMode(serverID int) error {
	target := fmt.Sprintf("%s/boot/normal/%d", serverEndPoint, serverID)
	_, err := c.doPOST(target, map[string]string{})
	if err != nil {
		return err
	}

	return nil
}
