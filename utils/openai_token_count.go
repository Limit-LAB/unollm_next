package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/pkoukk/tiktoken-go"
	"github.com/sashabaranov/go-openai"
)

var tkm *tiktoken.Tiktoken

func init() {
	var err error
	tkm, err = loadCl00k()
	if err != nil {
		err = fmt.Errorf("encoding for model failed: %v", err)
		log.Println(err)
	}
}

func downloadOrCache(url string, path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	http.Get(url)
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

// GPT3 & GPT4 use CL100K encoding
func loadCl00k() (*tiktoken.Tiktoken, error) {
	bpeLoader := tiktoken.NewDefaultBpeLoader()
	path := "./cl100k_base.tiktoken"
	err := downloadOrCache("https://openaipublic.blob.core.windows.net/encodings/cl100k_base.tiktoken", path)
	ranks, err := bpeLoader.LoadTiktokenBpe(path)
	if err != nil {
		return nil, err
	}
	special_tokens := map[string]int{
		tiktoken.ENDOFTEXT:   100257,
		tiktoken.FIM_PREFIX:  100258,
		tiktoken.FIM_MIDDLE:  100259,
		tiktoken.FIM_SUFFIX:  100260,
		tiktoken.ENDOFPROMPT: 100276,
	}
	enc := &tiktoken.Encoding{
		Name:           tiktoken.MODEL_CL100K_BASE,
		PatStr:         `(?i:'s|'t|'re|'ve|'m|'ll|'d)|[^\r\n\p{L}\p{N}]?\p{L}+|\p{N}{1,3}| ?[^\s\p{L}\p{N}]+[\r\n]*|\s*[\r\n]+|\s+(?!\S)|\s+`,
		MergeableRanks: ranks,
		SpecialTokens:  special_tokens,
	}
	pbe, err := tiktoken.NewCoreBPE(enc.MergeableRanks, enc.SpecialTokens, enc.PatStr)
	if err != nil {
		return nil, err
	}
	specialTokensSet := map[string]any{}
	for k := range enc.SpecialTokens {
		specialTokensSet[k] = true
	}
	return tiktoken.NewTiktoken(pbe, enc, specialTokensSet), nil
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
