package webhooks

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	//Ensuring the correct types
	_ "github.com/openshift/api/build/v1"
	buildv1 "github.com/openshift/client-go/build/clientset/versioned/typed/build/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

//ExposeKey is label key used to trigger webhook processing
const ExposeKey = "push-webhook"

func main() {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	buildV1Client, err := buildv1.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Getting builds")

	builds, err := buildV1Client.Builds("").List(metav1.ListOptions{})

	fmt.Println("Got builds")

	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d builds in the cluster\n", len(builds.Items))

	for _, build := range builds.Items {
		fmt.Println("Build = ", build)
	}

	fmt.Println("Getting buildconfigs")
	buildconfigs, err := buildV1Client.BuildConfigs("").List(metav1.ListOptions{})

	fmt.Println("Got buildconfigs")

	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d buildconfigs in the cluster\n", len(buildconfigs.Items))

	for _, buildconfig := range buildconfigs.Items {
		fmt.Println("Buildconfigs = ", buildconfig.GetName()+" in namespace "+buildconfig.GetNamespace())
		labels := buildconfig.GetLabels()

		dest, found := labels[ExposeKey]

		if found == true {
			switch dest {
			case "Github":
				b := githubWebhook{}

				b.Publish(buildconfig.Spec.Triggers[0], buildconfig.Spec.Source.Git)

				// for _, wh := range buildconfig.Spec.Triggers {
				// 	if wh.GitHubWebHook != nil {
				// 		fmt.Println("Github trigger found with secret " + wh.GitHubWebHook.Secret + " " + dest)
				// 		fmt.Println("Github URL " + buildconfig.Spec.Source.Git.URI)
				// 	}
				// }
			case "Gitlab":
				fmt.Println("Unable to process Gitlab webhook type")
			case "Bitbucket":
				fmt.Println("Unable to process Bitbucket webhook type")
			case "Generic":
				fmt.Println("Unable to process Generic webhook type")
			default:
				fmt.Println("Unable to determine webhook type")
			}
		}

		for key, val := range labels {
			fmt.Println("labels = " + key + ":" + val)
		}

		if found {
			fmt.Println("Exposing webhook")
		} else {
			fmt.Println("Not exposing webhook")
		}
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
