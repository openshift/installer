package transfer

import (
	"fmt"
	"strings"
)

const userResourceIDSeparator = "/"

func UserCreateResourceID(serverID, userName string) string {
	parts := []string{serverID, userName}
	id := strings.Join(parts, userResourceIDSeparator)

	return id
}

func UserParseResourceID(id string) (string, string, error) {
	parts := strings.Split(id, userResourceIDSeparator)

	if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
		return parts[0], parts[1], nil
	}

	return "", "", fmt.Errorf("unexpected format for ID (%[1]s), expected SERVERID%[2]sUSERNAME", id, userResourceIDSeparator)
}

const accessResourceIDSeparator = "/"

func AccessCreateResourceID(serverID, externalID string) string {
	parts := []string{serverID, externalID}
	id := strings.Join(parts, accessResourceIDSeparator)

	return id
}

func AccessParseResourceID(id string) (string, string, error) {
	parts := strings.Split(id, accessResourceIDSeparator)

	if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
		return parts[0], parts[1], nil
	}

	return "", "", fmt.Errorf("unexpected format for ID (%[1]s), expected SERVERID%[2]sEXTERNALID", id, accessResourceIDSeparator)
}
