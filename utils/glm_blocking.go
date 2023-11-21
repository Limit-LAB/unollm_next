package utils

import (
	"encoding/json"
	"limit.dev/unollm/model/zhipu"
	"net/http"
	"strings"
	"time"
)

func GLMBlockingRequest(body zhipu.ChatCompletionRequest, modelName string, token string) (result zhipu.ChatCompletionResponse, err error) {
	expire := time.Duration(10000) * time.Second
	token, err = CreateJWTToken(token, expire)
	if err != nil {
		return
	}

	client := &http.Client{}
	reqBody, err := json.Marshal(body)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://open.bigmodel.cn/api/paas/v3/model-api/"+modelName+"/invoke", strings.NewReader(string(reqBody)))
	if err != nil {
		return
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&result)
	
	return
}
