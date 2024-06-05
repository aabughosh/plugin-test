package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"

	"io"
	"net/http"
	"os"

	"strings"

	"github.com/gorilla/mux"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var clientset *kubernetes.Clientset

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	route := mux.NewRouter()
	staticPathArg := "opt/app-root/web/dist"

	route.HandleFunc("/example", healthHandler)
	route.HandleFunc("/plugin-manifest.json", manifesthHandler)
	route.PathPrefix("/api/pods").HandlerFunc(listPods)
	route.PathPrefix("/api/logs/{podName}/{containerName}").HandlerFunc(getPodLogs)
	route.PathPrefix("/").Handler(filesHandler(http.Dir(staticPathArg)))

	// Start the server
	fmt.Print("Starting server on :9443\n")
	if err := http.ListenAndServeTLS(":9443", "/var/cert/tls.crt", "/var/cert/tls.key", route); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		panic(err.Error())
	}
}
func filesHandler(root http.FileSystem) http.Handler {
	fileServer := http.FileServer(root)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "text/javascript")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Expires", "0")

		fileServer.ServeHTTP(w, r)
	})
}
func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print("health check worked!\n")
	w.Write([]byte("health check worked!\n"))
}

func manifesthHandler(w http.ResponseWriter, r *http.Request) {
	manifestData, err := os.ReadFile("/opt/app-root/web/dist/plugin-manifest.json")
	if err != nil {
		fmt.Errorf("cannot read base manifest file: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// w.Header().Set("Content-Type", "application/json")
	// w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	// w.Header().Set("Expires", "0")

	w.Write(manifestData)
}

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
