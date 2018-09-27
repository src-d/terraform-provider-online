package online

import (
	"encoding/json"
	"fmt"
)

func (c *client) GetRescueImages(serverID int) ([]string, error) {
	target := fmt.Sprintf("%s/rescue_images/%d", serverEndPoint, serverID)
	body, err := c.doGET(target)
	if err != nil {
		return nil, err
	}

	images := []string{}
	err = json.Unmarshal(body, &images)
	if err != nil {
		return nil, err
	}

	return images, nil
}
