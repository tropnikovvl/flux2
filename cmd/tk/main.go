package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
)

var VERSION = "0.0.1"

var rootCmd = &cobra.Command{
	Use:           "tk",
	Short:         "Command line utility for assembling Kubernetes CD pipelines",
	Long:          `Command line utility for assembling Kubernetes CD pipelines.`,
	Version:       VERSION,
	SilenceUsage:  true,
	SilenceErrors: true,
}

var (
	kubeconfig string
	namespace  string
	timeout    time.Duration
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "", "gitops-system",
		"the namespace scope for this operation")
	rootCmd.PersistentFlags().DurationVarP(&timeout, "timeout", "", 5*time.Minute,
		"timeout for this operation")
}

func main() {
	log.SetFlags(0)
	generateDocs()
	kubeconfigFlag()
	if err := rootCmd.Execute(); err != nil {
		logFailure("%v", err)
		os.Exit(1)
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func kubernetesClient() (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func execCommand(command string) (string, error) {
	c := exec.Command("/bin/sh", "-c", command)
	output, err := c.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func logAction(format string, a ...interface{}) {
	fmt.Println(`✚`, fmt.Sprintf(format, a...))
}

func logSuccess(format string, a ...interface{}) {
	fmt.Println(`✔`, fmt.Sprintf(format, a...))
}

func logFailure(format string, a ...interface{}) {
	fmt.Println(`✗`, fmt.Sprintf(format, a...))
}

func generateDocs() {
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "docgen" {
		rootCmd.PersistentFlags().StringVarP(&kubeconfig, "kubeconfig", "", "~/.kube/config",
			"path to the kubeconfig file")
		err := doc.GenMarkdownTree(rootCmd, "./docs/cmd")
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}
}

func kubeconfigFlag() {
	if home := homeDir(); home != "" {
		rootCmd.PersistentFlags().StringVarP(&kubeconfig, "kubeconfig", "", filepath.Join(home, ".kube", "config"),
			"path to the kubeconfig file")
	} else {
		rootCmd.PersistentFlags().StringVarP(&kubeconfig, "kubeconfig", "", "",
			"absolute path to the kubeconfig file")
	}
}
