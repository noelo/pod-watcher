package main

import (
    "flag"
    "fmt"
    //"time"

    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/pkg/api/v1"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/kubernetes/test/e2e_node/services"
)

var (
    kubeconfig = flag.String("kubeconfig", "/Users/admin/.kube/config", "absolute path to the kubeconfig file")
)

func main() {
    flag.Parse()
    // uses the current context in kubeconfig
    config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
    if err != nil {
        panic(err.Error())
    }
    // creates the clientset
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        panic(err.Error())
    }

    nspsaces, err:= clientset.Core().Namespaces().List(v1.ListOptions{})
    for _,nspace := range nspsaces.Items {
        fmt.Printf("I have access to %s namespace in the cluster\n", nspace.ObjectMeta.Name)
        pods, err := clientset.Core().Pods(nspace.ObjectMeta.Name).List(v1.ListOptions{})
        if err != nil {
            panic(err.Error())
        }
        for _,pod :=  range pods.Items{
            fmt.Printf("\t Pod %s on %s \n", pod.ObjectMeta.Name,pod.Status.HostIP)
            fmt.Printf("\t\t Annotations %s \n", pod.ObjectMeta.Annotations["kubernetes.io/created-by"])

        }


        services,err := clientset.Core().Services(nspace.ObjectMeta.Name).List(v1.ListOptions{})
        if err != nil {
            panic(err.Error())
        }
        for _,service :=  range services.Items{
            fmt.Printf("\t Service %s @ %s type %s \n", service.ObjectMeta.Name,service.Spec.ClusterIP,service.Spec.Type)
        }

    }
}
