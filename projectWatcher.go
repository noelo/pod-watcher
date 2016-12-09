package main

import (
	// "client-go/kubernetes"
	"log"
	"time"

	"github.com/openshift/origin/pkg/client/cache"
	"github.com/openshift/origin/pkg/cmd/util/clientcmd"

	"github.com/spf13/pflag"
	"k8s.io/client-go/kubernetes"
	kapi "k8s.io/kubernetes/pkg/api"
	kcache "k8s.io/kubernetes/pkg/client/cache"
	// kclient "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/watch"
)

func main() {
	// var kubeClient kclient.Interface

	config, err := clientcmd.DefaultClientConfig(pflag.NewFlagSet("empty", pflag.ContinueOnError)).ClientConfig()
	if err != nil {
		log.Printf("Error creating cluster config: %s", err)
	}

	// kubeClient, err = kclient.New(config)

	kubeClient, err = kubernetes.NewForConfig(config)
	podQueue := cache.NewEventQueue(kcache.MetaNamespaceKeyFunc)

	podLW := &kcache.ListWatch{
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return kubeClient.Pods(kapi.NamespaceAll).List(options)
		},
		WatchFunc: func(options kapi.ListOptions) (watch.Interface, error) {
			return kubeClient.Pods(kapi.NamespaceAll).Watch(options)
		},
	}
	kcache.NewReflector(podLW, &kapi.Pod{}, podQueue, 0).Run()

	go func() {
		for {
			event, pod, err := podQueue.Pop()
			err = handlePod(event, pod.(*kapi.Pod), kubeClient)
			if err != nil {
				log.Fatalf("Error capturing pod event: %s", err)
			}
		}
	}()
}

func handlePod(eventType watch.EventType, pod *kapi.Pod, kubeClient kclient.Interface) {
	switch eventType {
	case watch.Added:
		log.Printf("Pod %s added", pod.Name)
		if pod.Namespace == "namespaceWeWantToRestrict" {
			hour := time.Now().Hour()
			if hour >= 5 && hour <= 10 {
				err := kubeClient.Pods(pod.Namespace).Delete(pod.Name, &kapi.DeleteOptions{})
				if err != nil {
					log.Printf("Error deleting pod %s: %s", pod.Name, err)
				}
			}
		}
	case watch.Modified:
		log.Printf("Pod %s modified", pod.Name)
	case watch.Deleted:
		log.Printf("Pod %s deleted", pod.Name)
	}
}
