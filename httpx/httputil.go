package httpx

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func DoPost(ctx context.Context, body interface{}, url string, headers map[string][]string, retryTimes int, retryDelay time.Duration) ([]byte, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, errors.Wrap(err, "marshal interface")
	}
	request, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, errors.Wrap(err, "new request with context")
	}
	request.Header = headers
	client := http.DefaultClient
	var resp []byte
	// if meet error, retry times that you set
	for k := 0; k < retryTimes; k++ {
		resp, err = doRequest(client, request)
		if err != nil {
			// sleep retry delay
			time.Sleep(retryDelay)
			continue
		}
		break
	}
	return resp, err
}

func doRequest(client *http.Client, request *http.Request) ([]byte, error) {
	res, err := client.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "do request")
	}
	if res == nil {
		return nil, errors.New("http response is nil")
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read response body")
	}
	return data, nil
}
