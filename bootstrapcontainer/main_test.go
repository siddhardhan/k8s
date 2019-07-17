package main

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

func TestOutSideCluster(t *testing.T) {
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
