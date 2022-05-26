package django_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mini-kubernetes/tools/httpget"
	"testing"
)

type Person struct {
	UserType int
}

func Test(t *testing.T) {
	request2 := Person{
		UserType: 2,
	}
	uri := fmt.Sprintf("http://127.0.0.1:37889")
	body2, _ := json.Marshal(request2)
	_, _ = httpget.Get(uri).
		ContentType("application/json").
		Body(bytes.NewReader([]byte(body2))).
		Execute()
}
