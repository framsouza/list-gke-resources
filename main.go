package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/framsouza/list-gke-resources/pkg/kubectl"
	"golang.org/x/oauth2/google"
	container "google.golang.org/api/container/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	projectID = flag.String("project", "", "Enter the Project ID")
	zone      = flag.String("zone", "-", "Enter the Compute zone")
	resource  = flag.String("resource", "", "Enter the resource type")
	ns        string
)

func main() {

	flag.Parse()
	if *resource == "" {
		fmt.Fprintln(os.Stderr, "Resource not found")
		flag.Usage()
		os.Exit(2)
	}

	if *projectID == "" {
		fmt.Fprintln(os.Stderr, "Missing project")
		flag.Usage()
		os.Exit(2)
	}
	if *zone == "" {
		fmt.Fprintln(os.Stderr, "Missing zone")
		flag.Usage()
		os.Exit(2)
	}

	ctx := context.Background()

	hc, err := google.DefaultClient(ctx, container.CloudPlatformScope)
	if err != nil {
		log.Fatalf("Could not get authenticated client: %v", err)
	}

	svc, err := container.New(hc)
	if err != nil {
		log.Fatalf("Could not initialize gke client: %v", err)
	}

	listResources(svc, *projectID, *zone, *resource)

}

func listResources(svc *container.Service, projectID, zone, rs string) string {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 10, 10, 0, ' ', 0)
	defer w.Flush()

	// retriving cluster list
	list, err := svc.Projects.Zones.Clusters.List(projectID, zone).Do()
	if err != nil {
		return rs
	}

	kubeConfig, err := kubectl.GetK8sClusterConfigs(context.TODO(), projectID)
	if err != nil {
		return rs
	}

	// gathering cluster name, pods name, volume path
	for _, v := range list.Clusters {

		cfg, err := clientcmd.NewNonInteractiveClientConfig(*kubeConfig, v.Name, &clientcmd.ConfigOverrides{CurrentContext: v.Name}, nil).ClientConfig()
		if err != nil {
			return rs
		}

		k8s, err := kubernetes.NewForConfig(cfg)
		if err != nil {
			return rs
		}

		var resources []string

		switch rs {
		case "policysecurity":
			policies, err := k8s.PolicyV1beta1().PodSecurityPolicies().List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				panic(err.Error())
			}

			for _, policy := range policies.Items {
				resources = append(resources, policy.Name)
				ns = policy.Namespace
			}

			fmt.Fprintf(w, "\n%s\t\t\t%s\t\t\n", "GKE NAME", "RESOURCE NAME")

			for index := range resources {
				fmt.Fprintf(w, "%s\t\t\t%s\t\t\n", v.Name, resources[index])
			}

		case "pods":
			pods, err := k8s.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				panic(err.Error())
			}
			for _, pod := range pods.Items {
				resources = append(resources, pod.Name)
				ns = pod.Namespace
			}

			fmt.Fprintf(w, "\n%s\t\t\t%s\t\t%s\t\t\n", "GKE NAME", "RESOURCE NAME", "NAMESPACE")

			for index := range resources {
				fmt.Fprintf(w, "%s\t\t\t%s\t\t%s\t\t\n", v.Name, resources[index], ns)
			}

		case "deployments":
			deploy, _ := k8s.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
			for _, deploy := range deploy.Items {
				resources = append(resources, deploy.Name)
				ns = deploy.Namespace
			}

			fmt.Fprintf(w, "\n%s\t\t\t%s\t\t%s\t\t\n", "GKE NAME", "RESOURCE NAME", "NAMESPACE")

			for index := range resources {
				fmt.Fprintf(w, "%s\t\t\t%s\t\t%s\t\t\n", v.Name, resources[index], ns)
			}

		case "nodes":
			nodes, _ := k8s.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
			for _, node := range nodes.Items {
				resources = append(resources, node.Name)

			}
			fmt.Fprintf(w, "\n%s\t\t%s\t\t\n", "GKE NAME", "RESOURCE NAME")

			for index := range resources {
				fmt.Fprintf(w, "%s\t\t%s\t\t\n", v.Name, resources[index])
			}

		case "statefulsets":
			statefulsets, _ := k8s.AppsV1().StatefulSets("").List(context.TODO(), metav1.ListOptions{})
			for _, sts := range statefulsets.Items {
				resources = append(resources, sts.Name)
				ns = sts.Namespace
			}

			fmt.Fprintf(w, "\n%s\t\t\t%s\t\t%s\t\t\n", "GKE NAME", "RESOURCE NAME", "NAMESPACE")

			for index := range resources {
				fmt.Fprintf(w, "%s\t\t\t%s\t\t%s\t\t\n", v.Name, resources[index], ns)
			}

		case "version":
			fmt.Fprintf(w, "\n%s\t\t\t%s\t\t%s\t\t\n", "GKE NAME", "MASTER VERSION", "NODE VERSION")
			fmt.Fprintf(w, "%s\t\t\t", v.Name)
			fmt.Fprintf(w, "%s\t\t", v.CurrentMasterVersion)
			fmt.Fprintf(w, "%s\t\t\n", v.CurrentNodeVersion)

		}

	}

	return rs
}
