package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	// "time"

	// "k8s.io/apimachinery/pkg/api/errors"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
)

func main() {
	var kubeconfig *string

	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	var namespace *string
	var labels *string

    namespace = flag.String("namespace", "default", "(optional) namespace to filter pods")
    labels = flag.String("labels", "", "(optional) labels to filter nodes, comma-separated")

	flag.Parse()

    fmt.Printf("labels: %v\n", *labels)

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

    api := clientset.CoreV1()

	// for {
    fmt.Println(listNodes(api, *namespace, *labels))
		// time.Sleep(1 * time.Second)
	// }
	os.Exit(0)
}

func listNodes(api corev1.CoreV1Interface, namespace string, labels string) string {
    // pods, err := api.Pods("").List(metav1.ListOptions{})
    // if err != nil {
    //     panic(err.Error())
    // }

    nodes, err := api.Nodes().List(metav1.ListOptions{ LabelSelector: labels })
    if err != nil {
        panic(err.Error())
    }

    pods, err := api.Pods("").List(metav1.ListOptions{})
    if err != nil {
        panic(err.Error())
    }

    fmt.Printf("There are %d nodes in the cluster\n", len(nodes.Items))
    // fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

    for _, node := range nodes.Items {
        var env string
        if val, ok := node.ObjectMeta.Labels["env"]; ok {
            env = val
        } else {
            env = ""
        }

        var nodeName = node.ObjectMeta.Name

        fmt.Printf("\n  - %v  %v\n", env, node.ObjectMeta.Name)

        for _, pod := range pods.Items {
            if pod.Spec.NodeName == nodeName {
                fmt.Printf("    - %v %v\n", pod.ObjectMeta.Name, pod.TypeMeta.Kind)
            }
        }
    }

    return ""
    // pod := "example-xxxxx"
    // _, err = clientset.CoreV1().Pods(namespace).Get(pod, metav1.GetOptions{})
    // if errors.IsNotFound(err) {
    //     fmt.Printf("Pod %s in namespace %s not found\n", pod, namespace)
    // } else if statusError, isStatus := err.(*errors.StatusError); isStatus {
    //     fmt.Printf("Error getting pod %s in namespace %s: %v\n", pod, namespace, statusError.ErrStatus.Message)
    // } else if err != nil {
    //     panic(err.Error())
    // } else {
    //     fmt.Printf("Found pod %s in namespace %s\n", pod, namespace)
    // }
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
