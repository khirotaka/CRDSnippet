package main

import (
	"context"
	"flag"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

func main() {
	var default_kube_config_path string

	if home := homedir.HomeDir(); home != "" {
		default_kube_config_path = filepath.Join(home, ".kube", "config")
	}

	kubeconfig := flag.String("kubeconfig", default_kube_config_path, "kubeconfig config file")
	flag.Parse()

	config, _ := clientcmd.BuildConfigFromFlags("", *kubeconfig)

	clientset, _ := kubernetes.NewForConfig(config)

	pods, _ := clientset.CoreV1().Pods("kube-system").List(context.TODO(), metav1.ListOptions{})

	for i, pod := range pods.Items {
		fmt.Printf("[Pod name %d] %s\n", i, pod.GetName())
	}
}
