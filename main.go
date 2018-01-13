package main

import (
	"net/http"
	"os"
	"context"
	"fmt"

	"k8s.io/client-go/tools/clientcmd"

	"github.com/rancher/types/config"
	"github.com/rancher/auth/authconfig"
)


func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	if err != nil {
		return err
	}

	mgmtCtx, err := config.NewManagementContext(*kubeConfig)
	if err != nil {
		return err
	}

	handler, err := authconfig.New(context.Background(), mgmtCtx)
	if err != nil {
		return err
	}

	fmt.Println("Listening on 0.0.0.0:1234")
	return http.ListenAndServe("0.0.0.0:1234", handler)
}

