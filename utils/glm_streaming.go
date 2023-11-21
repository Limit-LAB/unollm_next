package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"limit.dev/unollm/model/zhipu"
)

func GLMStreamingRequest(body zhipu.ChatCompletionRequest, modelName string, token string) (chan string, chan zhipu.ChatCompletionStreamResponse, error) {
	expire := time.Duration(10000) * time.Second
	token, err := CreateJWTToken(token, expire)
	if err != nil {
		return nil, nil, err
	}

	client := &http.Client{}
	reqBody, err := json.Marshal(body)
	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequest("POST", "https://open.bigmodel.cn/api/paas/v3/model-api/"+modelName+"/sse-invoke", strings.NewReader(string(reqBody)))
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	reader := NewEventStreamReader(resp.Body, 4096)

	llmCh := make(chan string)
	resultCh := make(chan zhipu.ChatCompletionStreamResponse, 1)

	go func() {
		for reader.scanner.Scan() {
			kv := strings.Split(reader.scanner.Text(), "\n")
			switch kv[0] {
			case "event:add":
				llmCh <- kv[2][5:]
			case "event:finish":
				var usage zhipu.ChatCompletionStreamResponse
				json.NewDecoder(strings.NewReader(kv[3][5:])).Decode(&usage)
				resultCh <- usage
			}
		}
		defer resp.Body.Close()
	}()

	return llmCh, resultCh, nil
}

// EventStreamReader scans an io.Reader looking for EventStream messages.
type EventStreamReader struct {
	scanner *bufio.Scanner
}

// Returns the minimum non-negative value out of the two values. If both
// are negative, a negative value is returned.
func minPosInt(a, b int) int {
	if a < 0 {
		return b
	}
	if b < 0 {
		return a
	}
	if a > b {
		return b
	}
	return a
}

// Returns a tuple containing the index of a double newline, and the number of bytes
// represented by that sequence. If no double newline is present, the first value
// will be negative.
func containsDoubleNewline(data []byte) (int, int) {
	// Search for each potentially valid sequence of newline characters
	crcr := bytes.Index(data, []byte("\r\r"))
	lflf := bytes.Index(data, []byte("\n\n"))
	crlflf := bytes.Index(data, []byte("\r\n\n"))
	lfcrlf := bytes.Index(data, []byte("\n\r\n"))
	crlfcrlf := bytes.Index(data, []byte("\r\n\r\n"))
	// Find the earliest position of a double newline combination
	minPos := minPosInt(crcr, minPosInt(lflf, minPosInt(crlflf, minPosInt(lfcrlf, crlfcrlf))))
	// Detemine the length of the sequence
	nlen := 2
	if minPos == crlfcrlf {
		nlen = 4
	} else if minPos == crlflf || minPos == lfcrlf {
		nlen = 3
	}
	return minPos, nlen
}

func NewEventStreamReader(eventStream io.Reader, maxBufferSize int) *EventStreamReader {
	scanner := bufio.NewScanner(eventStream)

	initBufferSize := minPosInt(4096, maxBufferSize)
	scanner.Buffer(make([]byte, initBufferSize), maxBufferSize)

	split := func(data []byte, atEOF bool) (int, []byte, error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		// We have a full event payload to parse.
		if i, nlen := containsDoubleNewline(data); i >= 0 {
			return i + nlen, data[0:i], nil
		}
		// If we're at EOF, we have all of the data.
		if atEOF {
			return len(data), data, nil
		}
		// Request more data.
		return 0, nil, nil
	}
	// Set the split function for the scanning operation.
	scanner.Split(split)

	return &EventStreamReader{
		scanner: scanner,
	}
}
