package utils

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

func GLMBlockingRequest(body map[string]interface{}, modelName string, token string) (map[string]interface{}, error) {
	expire := time.Duration(10000) * time.Second
	token, err := CreateJWTToken(token, expire)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	reqBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://open.bigmodel.cn/api/paas/v3/model-api/"+modelName+"/invoke", strings.NewReader(string(reqBody)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
