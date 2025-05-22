package credentials_provider_file

import (
	"encoding/json"
	"os"
)

type FileUserInfo struct {
	RealmID string `json:"realm_id"`
	UserID  string `json:"user_id"`
	Role    string `json:"role"`
	Login   string `json:"login"`
	Pwdhash string `json:"pwdhash"`
}

func ReadUsers(path string) ([]*FileUserInfo, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var result []*FileUserInfo
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
