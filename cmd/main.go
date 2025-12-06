package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	lldapv1 "github.com/dj346/lldap-operator/internal"

	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func main() {
	// Optional flag: limit to a namespace
	namespace := flag.String("namespace", "", "Namespace to look in (empty = all namespaces)")
	flag.Parse()

	ctx := context.Background()

	// 1) Build a scheme and register core + your CRDs
	scheme := runtime.NewScheme()

	if err := clientgoscheme.AddToScheme(scheme); err != nil {
		log.Fatalf("adding core scheme failed: %v", err)
	}
	if err := lldapv1.AddToScheme(scheme); err != nil {
		log.Fatalf("adding LLDAP scheme failed: %v", err)
	}

	// 2) Get a Kubernetes REST config
	cfg, err := rest.InClusterConfig()
	if err != nil {
		// fallback for running locally with ~/.kube/config
		kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		cfg, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Fatalf("failed to build kubeconfig: %v", err)
		}
	}

	// 3) Create a controller-runtime client with that scheme
	c, err := client.New(cfg, client.Options{Scheme: scheme})
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	// Optional namespace scoping
	var opts []client.ListOption
	if *namespace != "" {
		opts = append(opts, client.InNamespace(*namespace))
	}

	// 4) List and print users
	if err := printUsers(ctx, c, opts...); err != nil {
		log.Printf("error listing users: %v", err)
	}

	// 5) List and print groups
	if err := printGroups(ctx, c, opts...); err != nil {
		log.Printf("error listing groups: %v", err)
	}
}

func printUsers(ctx context.Context, c client.Client, opts ...client.ListOption) error {
	var users lldapv1.LLDAPUserList
	if err := c.List(ctx, &users, opts...); err != nil {
		return err
	}

	fmt.Println("=== LLDAPUsers ===")
	if len(users.Items) == 0 {
		fmt.Println("(none)")
		return nil
	}

	for _, u := range users.Items {
		fmt.Printf(
			"User %s/%s: username=%s, primaryEmail=%s, groups=%v\n",
			u.Namespace,
			u.Name,
			u.Spec.Username,
			u.Spec.PrimaryEmail,
			u.Spec.Groups,
		)
	}
	return nil
}

func printGroups(ctx context.Context, c client.Client, opts ...client.ListOption) error {
	var groups lldapv1.LLDAPGroupList
	if err := c.List(ctx, &groups, opts...); err != nil {
		return err
	}

	fmt.Println("=== LLDAPGroups ===")
	if len(groups.Items) == 0 {
		fmt.Println("(none)")
		return nil
	}

	for _, g := range groups.Items {
		fmt.Printf(
			"Group %s/%s: name=%s, gid=%v, members=%v\n",
			g.Namespace,
			g.Name,
			g.Spec.Name,
			g.Spec.GIDNumber,
			g.Spec.Members,
		)
	}
	return nil
}
