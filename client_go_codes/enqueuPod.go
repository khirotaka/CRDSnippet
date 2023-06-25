package main

import (
	"flag"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/workqueue"
	"log"
	"path/filepath"
	"time"
)

func main() {
	var defaultKubeConfigPath string

	if home := homedir.HomeDir(); home != "" {
		defaultKubeConfigPath = filepath.Join(home, ".kube", "config")
	}

	kubeconfig := flag.String("kubeconfig", defaultKubeConfigPath, "kubeconfig config file")
	flag.Parse()

	config, _ := clientcmd.BuildConfigFromFlags("", *kubeconfig)

	clientset, _ := kubernetes.NewForConfig(config)

	// Informerを生成
	// 30*time.Second で30秒に一回 UpdateFuncが実行され in-memory-cacheを更新する。
	informerFactory := informers.NewSharedInformerFactory(clientset, 30*time.Second)

	podInformer := informerFactory.Core().V1().Pods()
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	defer queue.ShutDown()
	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(old interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(old)
			if err != nil {
				runtime.HandleError(err)
				return
			}
			queue.Add(key)
			log.Println("Added: " + key)
		},
		UpdateFunc: func(old, new interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(old)

			if err != nil {
				runtime.HandleError(err)
				return
			}
			queue.Add(key)
			log.Println("Update: " + key)
		},
		DeleteFunc: func(old interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(old)

			if err != nil {
				runtime.HandleError(err)
				return
			}
			queue.Add(key)
			log.Println("Deleted: " + key)
		},
	})

	informerFactory.Start(wait.NeverStop)
	informerFactory.WaitForCacheSync(wait.NeverStop)

	podLister := podInformer.Lister()
	_, err := podLister.List(labels.Nothing())

	if err != nil {
		log.Fatal(err)
	}

	for {
		time.Sleep(1 * time.Second)
	}
}
