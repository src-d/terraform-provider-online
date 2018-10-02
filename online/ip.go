package online

import "fmt"

func (c *client) EditReverseHostname(ip, revere string) error {
	target := fmt.Sprintf("%s/ip/edit", serverEndPoint)
	_, err := c.doPOST(target, map[string]string{
		"address": ip,
		"reverse": revere,
	})
	if err != nil {
		return err
	}
	return nil
}
