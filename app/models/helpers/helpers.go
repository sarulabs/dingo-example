package helpers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func JSONResponse(w http.ResponseWriter, status int, data interface{}) {
	resp, _ := json.Marshal(data)
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func ReadJSONBody(r *http.Request, data interface{}) {
	return
}

// ReadBody reads the body of the http.Request.
// The method can be called multiple times on the same request.
func ReadBody(r *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	r.Body.Close()

	// replace the body in the request so that it can be red again
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return body, nil
}

// BytesBody reads the body of the http.Request with ReadBody.
// But the errors are not reported.
func BytesBody(r *http.Request) []byte {
	body, _ := ReadBody(r)
	return body
}

// StringBody reads the body of the http.Request with ReadBody.
// But the response is a string and errors are not reported.
func StringBody(r *http.Request) string {
	body, _ := ReadBody(r)
	return string(body)
}

// JSONBody reads the body of the given request
// and tries to decode it as a json in data.
// It uses jsoniter.
func JSONBody(r *http.Request, data interface{}) error {
	body, err := ReadBody(r)
	if err != nil {
		return err
	}

	return jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(body, data)
}
