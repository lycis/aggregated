package http

import (
	"net/http"
	"strings"
	"fmt"
)

type ConnectionData struct {
	Url string
	Headers []string
	Method string
}

func (c ConnectionData) eval() *http.Response {
	client := http.Client{}
	req, err := http.NewRequest(c.Method, c.Url, nil)
	if err != nil {
		panic(err)
	}
	
	for _, h := range c.Headers {
		header := strings.Split(h, ":")
		if len(header) < 2 {
			panic(fmt.Errorf("invalid header: %s", header))
		}
		req.Header.Set(header[0], header[1])
	}
	
	response, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	
	return response
}
