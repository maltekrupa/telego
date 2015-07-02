package telego

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"testing"
)

var (
	mux        *http.ServeMux
	server     *httptest.Server
	client     *telego
	ENV_APIKEY = os.Getenv("TELEGO_APIKEY")
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	serverUrl, _ := url.Parse(server.URL)

	client = NewTelego(ENV_APIKEY)
	client.ChangeUrl(serverUrl.String())
}

func teardown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if want != r.Method {
		t.Errorf("Request method = %v, want %v", r.Method, want)
	}
}

func TestGetUpdates(t *testing.T) {
	setup()
	defer teardown()

	getData := `{
		"ok":true,
		"result":[
			{
				"update_id":805750848,
				"message": {
					"message_id":2,"from":{
						"id":123456,
						"first_name":"foo",
						"last_name":"bar",
						"username":"foo_bar"
					},
					"chat":{
						"id":123456,
						"first_name":"foo",
						"last_name":"bar",
						"username":"foo_bar"
					},
					"date":1435771984,
					"text":"testing the messages"
				}
			}, {
				"update_id":805750849,
				"message":{
					"message_id":3,
					"from":{
						"id":123456,
						"first_name":"foo",
						"last_name":"bar",
						"username":"foo_bar"
					},
					"chat":{
						"id":123456,
						"first_name":"foo",
						"last_name":"bar",
						"username":"foo_bar"
					},
					"date":1435771995,
					"text":"\/start"
				}
			}
		]
	}`

	mux.HandleFunc("/bot"+client.token+"/getUpdates",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			fmt.Fprint(w, getData)
		})

	updates, err := client.GetUpdates()
	if err != nil {
		t.Errorf("GetUpdates() returned error: %v", err)
	}

	want := ResponseUpdate{true,
		[]Update{
			Update{805750848, Message{2, User{123456, "foo", "bar", "foo_bar"}, 1435771984, User{123456, "foo", "bar", "foo_bar"}, "testing the messages"}},
			Update{805750849, Message{3, User{123456, "foo", "bar", "foo_bar"}, 1435771995, User{123456, "foo", "bar", "foo_bar"}, "/start"}},
		}}

	if !reflect.DeepEqual(updates, want) {
		t.Errorf("GetUpdates() returned %+v, want %+v",
			updates, want)
	}
}

func TestGetUpdatesOkFalse(t *testing.T) {
	setup()
	defer teardown()

	getData := `{
		"ok":false,
	}`

	mux.HandleFunc("/bot"+client.token+"/getUpdates",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			fmt.Fprint(w, getData)
		})

	updates, err := client.GetUpdates()
	if err != nil {
		t.Errorf("GetUpdates() returned error: %v", err)
	}

	want := ResponseUpdate{false, nil}

	if !reflect.DeepEqual(updates, want) {
		t.Errorf("GetUpdates() returned %+v, want %+v",
			updates, want)
	}
}

func TestGetMe(t *testing.T) {
	setup()
	defer teardown()

	getData := `{
		"ok":true,
		"result":{
			"id":987654321,
			"first_name":"Foo",
			"username":"FooBot"
		}
	}`

	mux.HandleFunc("/bot"+client.token+"/getMe",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			fmt.Fprint(w, getData)
		})

	updates, err := client.GetMe()
	if err != nil {
		t.Errorf("GetMe() returned error: %v", err)
	}

	want := ResponseMe{true, Me{987654321, "Foo", "FooBot"}}

	if !reflect.DeepEqual(updates, want) {
		t.Errorf("GetMe() returned %+v, want %+v",
			updates, want)
	}
}

func TestSendMessage(t *testing.T) {
	setup()
	defer teardown()

	getData := `{
		"ok":true,
		"result":{
			"message_id":231
			,"from":{
				"id":987654321,
				"first_name":"Foo",
				"username":"FooBot"
			},
			"chat":{
				"id":123456,
				"first_name":"foo",
				"last_name":"bar",
				"username":"foo_bar"
			},
			"date":1435836484,
			"text":"test"
		}
	}`

	mux.HandleFunc("/bot"+client.token+"/sendMessage",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			fmt.Fprint(w, getData)
		})

	updates, err := client.SendMessage(123456, "test")
	if err != nil {
		t.Errorf("GetMe() returned error: %v", err)
	}

	want := ResponseSendMessage{true,
		Message{231, User{987654321, "Foo", "", "FooBot"}, 1435836484, User{123456, "foo", "bar", "foo_bar"}, "test"}}

	if !reflect.DeepEqual(updates, want) {
		t.Errorf("GetMe() returned %+v, want %+v",
			updates, want)
	}
}
