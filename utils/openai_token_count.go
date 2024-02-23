package utils

import (
	"fmt"
	"log"

	"github.com/pkoukk/tiktoken-go"
	"github.com/sashabaranov/go-openai"
)

var tkm *tiktoken.Tiktoken

func init() {
	var err error

	// FIXME: runtime downloading tiktoken model, switch to offline tiktoken_loader instead
	// But its has Cache.
	// Offline token : https://openaipublic.blob.core.windows.net/encodings/cl100k_base.tiktoken
	tkm, err = tiktoken.EncodingForModel(tkm_model)
	if err != nil {
		err = fmt.Errorf("encoding for model failed: %v", err)
		log.Println(err)
	}
}

const (
	tkm_tokenPerMessage = 3
	tkm_tokenPerName    = 1
	tkm_model           = "gpt-3.5-turbo-0613"
)

// OpenAI Cookbook: https://github.com/openai/openai-cookbook/blob/main/examples/How_to_count_tokens_with_tiktoken.ipynb
func GetOpenAITokenCount(messages []openai.ChatCompletionMessage) (numTokens int) {
	if tkm == nil {
		return -1
	}

	// ignore old models
	// switch model {
	// case "gpt-3.5-turbo-0613",
	// 	"gpt-3.5-turbo-16k-0613",
	// 	"gpt-4-0314",
	// 	"gpt-4-32k-0314",
	// 	"gpt-4-0613",
	// 	"gpt-4-32k-0613":
	// 	tokensPerMessage = 3
	// 	tokensPerName = 1
	// case "gpt-3.5-turbo-0301":
	// 	tokensPerMessage = 4 // every message follows <|start|>{role/name}\n{content}<|end|>\n
	// 	tokensPerName = -1   // if there's a name, the role is omitted
	// default:
	// 	if strings.Contains(model, "gpt-3.5-turbo") {
	// 		log.Println("warning: gpt-3.5-turbo may update over time. Returning num tokens assuming gpt-3.5-turbo-0613.")
	// 		return GetOpenAITokenCount(messages, "gpt-3.5-turbo-0613")
	// 	} else if strings.Contains(model, "gpt-4") {
	// 		log.Println("warning: gpt-4 may update over time. Returning num tokens assuming gpt-4-0613.")
	// 		return GetOpenAITokenCount(messages, "gpt-4-0613")
	// 	} else {
	// 		err = fmt.Errorf("num_tokens_from_messages() is not implemented for model %s. See https://github.com/openai/openai-python/blob/main/chatml.md for information on how messages are converted to tokens.", model)
	// 		log.Println(err)
	// 		return
	// 	}
	// }

	for _, message := range messages {
		numTokens += tkm_tokenPerMessage
		numTokens += len(tkm.Encode(message.Content, nil, nil))
		numTokens += len(tkm.Encode(message.Role, nil, nil))
		numTokens += len(tkm.Encode(message.Name, nil, nil))
		if message.Name != "" {
			numTokens += tkm_tokenPerName
		}
	}
	numTokens += 3 // every reply is primed with <|start|>assistant<|message|>
	return numTokens
}
