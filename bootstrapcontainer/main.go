package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
)

type API struct {
	httpClient    *http.Client
	URL, Token    string
	RequestParams map[string]string
}

func main() {

	isInitContainer := flag.Bool("init", false, "to decide init or side-car container")
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)

	flag.Parse()

	var NamespaceFile = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"

	for {
		// Start - Bearer Token
		config, err := rest.InClusterConfig()
		if err != nil {
			log.Panicln(err.Error())
		}
		log.Println("Bear Token is : ", config.BearerToken)
		// End - Bearer Token

		// Get Namespace
		data, err := ioutil.ReadFile(NamespaceFile)
		if err != nil {
			log.Panic(err.Error())
		}

		namespace := strings.TrimSpace(string(data))
		if len(namespace) < 1 {
			log.Panic("Content of ", NamespaceFile, " can not be null")
		}

		// Start - K8S Client Operations to get Pods
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			log.Panic(err.Error())
		}

		pods, err := ListPods(clientset.CoreV1(), namespace)
		if err != nil {
			log.Fatalf("listing pods: %s", err)
		}

		fmt.Println(strings.Join(pods, "\n"))
		// End - K8S Client Operations to get Pods

		// Start - Call rest api
		api := API{
			httpClient:    &http.Client{},
			URL:           "https://api.weather.gov/",
			RequestParams: map[string]string{"Accept": "application/json", "Token": config.BearerToken}}
		resp_body, err := GetToken(&api)
		if err != nil {
			log.Fatalf("Errored when sending request to the server : %s", err)
		}
		log.Println("Response Body is : ", string(resp_body))
		// End - Call rest api

		//check for init container flag - break the loop; if initcontianer flag is enabled
		if *isInitContainer {
			break
		}

		time.Sleep(10 * time.Second)
	}
}

// ListPods returns a list of Pods running in the provided namespace
func ListPods(client corev1.CoreV1Interface, namespace string) ([]string, error) {
	pl, err := client.Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "getting pods")
	}

	var pods []string
	for _, p := range pl.Items {
		pods = append(pods, p.GetName())
	}

	return pods, nil
}

func GetToken(api *API) ([]byte, error) {

	log.Println("URL : ", api.URL)
	log.Println("RequestParams : ", api.RequestParams)

	req, err := http.NewRequest("GET", api.URL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build request")
	}

	for key, value := range api.RequestParams {
		log.Println("Key:", key, "Value:", value)
		req.Header.Add(key, value)
	}
	resp, err := api.httpClient.Do(req)

	log.Println("Response Status : ", resp.StatusCode)

	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}

	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	// handling error and doing stuff with body that needs to be unit tested
	return resp_body, err
}
