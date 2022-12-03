package internal

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/michimani/aws-lambda-api-go/alago"
)

type Header struct {
	Key   string
	Value string
}

// CallAPI execute http request using alago.Client
// and returns StatusCode, ResponseHeader, ResponseBody as slice of bytes and error.
func CallAPI(ctx context.Context, c alago.AlagoClient, method, url string, body io.Reader, headers ...Header) (int, map[string][]string, []byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return 0, nil, nil, err
	}

	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	// add header
	for _, h := range headers {
		req.Header.Set(h.Key, h.Value)
	}

	res, err := c.HttpClient().Do(req)
	if err != nil {
		return 0, nil, nil, err
	}
	defer res.Body.Close()

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, res.Body); err != nil {
		return 0, nil, nil, err
	}

	return res.StatusCode, res.Header, buf.Bytes(), nil
}
