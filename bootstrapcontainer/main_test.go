package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

func TestMainOutSideCluster(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r == nil {
				fmt.Println("unable to load in-cluster configuration, KUBERNETES_SERVICE_HOST and KUBERNETES_SERVICE_PORT must be defined", r)
			}
		}()
		// This function should cause a panic
		main()
	}()
}

func TestListPods(t *testing.T) {
	var tests = []struct {
		description string
		namespace   string
		expected    []string
		objs        []runtime.Object
	}{
		{"no pods", "", nil, nil},
		{"all namespaces", "", []string{"poda", "podb"}, []runtime.Object{pod("correct-namespace", "poda"), pod("wrong-namespace", "podb")}},
		{"filter namespace", "correct-namespace", []string{"poda"}, []runtime.Object{pod("correct-namespace", "poda"), pod("wrong-namespace", "podb")}},
		{"wrong namespace", "correct-namespace", nil, []runtime.Object{pod("wrong-namespace", "podb")}},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			client := fake.NewSimpleClientset(test.objs...)
			actual, err := ListPods(client.CoreV1(), test.namespace)
			if err != nil {
				t.Errorf("Unexpected error: %s", err)
				return
			}
			if diff := cmp.Diff(actual, test.expected); diff != "" {
				t.Errorf("%T differ (-got, +want): %s", test.expected, diff)
				return
			}
		})
	}
}

func pod(namespace, name string) *v1.Pod {
	return &v1.Pod{ObjectMeta: metav1.ObjectMeta{Namespace: namespace, Name: name}}
}

func testingHTTPClient(handler http.Handler) (*http.Client, func()) {
	s := httptest.NewTLSServer(handler)

	cli := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, s.Listener.Addr().String())
			},
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	return cli, s.Close
}

func TestClientGetToken(t *testing.T) {
	var okResponse = `{
		"status": "OK"
	}`
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		assert.Equal(t, "23rfwe23", r.Header.Get("Token"))
		w.Write([]byte(okResponse))
	})
	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	api := API{
		httpClient:    &http.Client{},
		URL:           "https://api.weather.gov/",
		RequestParams: map[string]string{"Accept": "application/json", "Token": "23rfwe23"}}
	api.httpClient = httpClient

	resp, err := GetToken(&api)

	assert.Nil(t, err)
	assert.Equal(t, 21, len(resp))
}
