package online

import (
	"fmt"
)

func (c *client) EditFailoverIP(source, destination string) error {
	target := fmt.Sprintf("%s/failover/edit", serverEndPoint)
	_, err := c.doPOST(target, map[string]string{
		"source":      source,
		"destination": destination,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *client) GenerateMACFailoverIP(address, macType string) (string, error) {
	target := fmt.Sprintf("%s/failover/generateMac", serverEndPoint)
	body, err := c.doPOST(target, map[string]string{
		"address": address,
		"type":    macType,
	})
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (c *client) DeleteMACFailoverIP(address string) error {
	target := fmt.Sprintf("%s/failover/deleteMac", serverEndPoint)
	_, err := c.doPOST(target, map[string]string{
		"address": address,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *client) SetReverseFailoverIP(address, reverse string) error {
	target := fmt.Sprintf("%s/ip/edit", serverEndPoint)
	_, err := c.doPOST(target, map[string]string{
		"address": address,
		"reverse": reverse,
	})
	if err != nil {
		return err
	}

	return nil
}
