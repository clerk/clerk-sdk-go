package clerk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestProxyChecksService_Create(t *testing.T) {
	client, mux, _, teardown := setup("token")
	defer teardown()
	var params CreateProxyCheckParams
	payload := `{
	"proxy_url":"https://example.com/__clerk",
	"domain_id": "dmn_1mebQggrD3xO5JfuHk7clQ94ysA"
}`
	_ = json.Unmarshal([]byte(payload), &params)

	mux.HandleFunc("/proxy_checks", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "POST")
		testHeader(t, req, "Authorization", "Bearer token")
		fmt.Fprint(w, dummyProxyCheckJSON)
	})

	got, err := client.ProxyChecks().Create(params)
	if err != nil {
		t.Fatal(err)
	}

	var want ProxyCheck
	err = json.Unmarshal([]byte(dummyProxyCheckJSON), &want)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, &want) {
		t.Errorf("Response = %v, want %v", got, &want)
	}
}

const dummyProxyCheckJSON = `{
	"object": "proxy_check",
	"id": "proxychk_1mebQggrD3xO5JfuHk7clQ94ysA",
	"successful": true,
	"domain_id": "dmn_1mebQggrD3xO5JfuHk7clQ94ysA",
	"proxy_url": "https://example.com/__clerk",
	"last_run_at": 1610783813,
	"created_at": 1610783813,
	"updated_at": 1610783813
}`
