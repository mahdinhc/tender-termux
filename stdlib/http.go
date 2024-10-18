package stdlib

import (
	"bytes"
	// "encoding/json"
	_"fmt"
	"io/ioutil"
	"net/http"
	"github.com/2dprototype/tender"
)

var httpModule = map[string]tender.Object{
	"get":     &tender.UserFunction{Name: "get", Value: httpGet},
	"post":    &tender.UserFunction{Name: "post", Value: httpPost},
	"put":     &tender.UserFunction{Name: "put", Value: httpPut},
	"delete":  &tender.UserFunction{Name: "delete", Value: httpDelete},
	"patch":   &tender.UserFunction{Name: "patch", Value: httpPatch},
	"options": &tender.UserFunction{Name: "options", Value: httpOptions},
	"head":    &tender.UserFunction{Name: "head", Value: httpHead},
	"trace":   &tender.UserFunction{Name: "trace", Value: httpTrace},
}

func httpRequest(method string, args ...tender.Object) (ret tender.Object, err error) {
	if len(args) < 1 || len(args) > 3 {
		return nil, tender.ErrWrongNumArguments
	}

	url, _ := tender.ToString(args[0])
	var body []byte
	if len(args) > 1 {
		body, _ = tender.ToByteSlice(args[1])
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return wrapError(err), nil
	}

	if len(args) == 3 {
		headersObj, ok := args[2].(*tender.Map)
		if !ok {
			return nil, tender.ErrInvalidArgumentType{
				Name:     "headers",
				Expected: "map",
				Found:    args[2].TypeName(),
			}
		}

		for key, value := range headersObj.Value {
			valueStr, _ := value.(*tender.String)
			req.Header.Set(key, valueStr.Value)
		}

	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return wrapError(err), nil
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return wrapError(err), nil
	}

	return &tender.Bytes{Value: respBody}, nil
}

func httpGet(args ...tender.Object) (ret tender.Object, err error) {
	return httpRequest("GET", args...)
}

func httpPost(args ...tender.Object) (ret tender.Object, err error) {
	return httpRequest("POST", args...)
}

func httpPut(args ...tender.Object) (ret tender.Object, err error) {
	return httpRequest("PUT", args...)
}

func httpDelete(args ...tender.Object) (ret tender.Object, err error) {
	return httpRequest("DELETE", args...)
}

func httpPatch(args ...tender.Object) (ret tender.Object, err error) {
	return httpRequest("PATCH", args...)
}

func httpOptions(args ...tender.Object) (ret tender.Object, err error) {
	return httpRequest("OPTIONS", args...)
}

func httpHead(args ...tender.Object) (ret tender.Object, err error) {
	return httpRequest("HEAD", args...)
}

func httpTrace(args ...tender.Object) (ret tender.Object, err error) {
	return httpRequest("TRACE", args...)
}
