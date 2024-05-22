package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var clientset *kubernetes.Clientset

//go:embed ExamplePage.tsx
var indexHTML []byte

func main() {
	// Load Kubernetes configuration from file
	// Get the value of the KUBECONFIG environment variable
	// Get the kubeconfig file path
	kubeConfigFile := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	if envKubeconfig := os.Getenv("KUBECONFIG"); envKubeconfig != "" {
		kubeConfigFile = envKubeconfig
	}
	fmt.Printf("Using kubeconfig file: %s\n", kubeConfigFile)
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigFile)
	if err != nil {
		panic(err.Error())
	}

	// Create Kubernetes client
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Create a new HTTP serve mux
	router := mux.NewRouter()

	// Define route to list pods
	router.HandleFunc("/api/pods", listPods).Methods("GET")

	// Define route to get pod logs
	router.HandleFunc("/api/logs/{podName}/{containerName}", getPodLogs).Methods("GET")

	// Enable CORS
	// Enable CORS
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:9000"}),
		handlers.AllowedMethods([]string{"GET", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)

	// Wrap your HTTP handler with the CORS middleware
	http.Handle("/", corsHandler(router))

	// Start the server
	if err := http.ListenAndServe(":9002", nil); err != nil {
		// Handle error
		panic(err.Error())
	}
}

// Handler function to serve ProxyTestPage.ts

func listPods(w http.ResponseWriter, r *http.Request) {
	fmt.Print("called get logs")

	pods, err := clientset.CoreV1().Pods("cnf-certsuite-operator").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type PodInfo struct {
		Name       string   `json:"name"`
		Containers []string `json:"containers"`
	}

	var podInfos []PodInfo
	for _, pod := range pods.Items {
		var containerNames []string
		for _, container := range pod.Spec.Containers {
			containerNames = append(containerNames, container.Name)
		}
		podInfos = append(podInfos, PodInfo{Name: pod.Name, Containers: containerNames})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(podInfos); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getPodLogs(w http.ResponseWriter, r *http.Request) {
	fmt.Print("called get logs")
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/logs/"), "/")
	if len(parts) != 2 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	podName := parts[0]
	containerName := parts[1]

	podLogOpts := corev1.PodLogOptions{Container: containerName}
	req := clientset.CoreV1().Pods("cnf-certsuite-operator").GetLogs(podName, &podLogOpts)
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer podLogs.Close()

	_, err = io.Copy(w, podLogs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
