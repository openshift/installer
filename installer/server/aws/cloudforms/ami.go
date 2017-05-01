package cloudforms

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func getAMI(region, channel string) (string, error) {

	regions, err := getAMIData(channel)

	if err != nil {
		return "", fmt.Errorf("error getting ami data for channel %s: %v", channel, err)
	}

	amis, ok := regions[region]
	if !ok {
		return "", fmt.Errorf("could not find region %s for channel %s", region, channel)
	}

	if ami, ok := amis["hvm"]; ok {
		return ami, nil
	}

	return "", fmt.Errorf("could not find hvm image for region %s, channel %s", region, channel)
}

func getAMIData(channel string) (map[string]map[string]string, error) {
	r, err := http.Get(fmt.Sprintf("https://coreos.com/dist/aws/aws-%s.json", channel))
	if err != nil {
		return nil, fmt.Errorf("failed to get AMI data: %s: %v", channel, err)
	}

	if r.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get AMI data: %s: invalid status code: %d", channel, r.StatusCode)
	}

	output := map[string]map[string]string{}

	err = json.NewDecoder(r.Body).Decode(&output)
	if err != nil {
		return nil, fmt.Errorf("failed to parse AMI data: %s: %v", channel, err)
	}
	r.Body.Close()

	return output, nil
}
