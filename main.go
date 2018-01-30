package main

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"

	"k8s.io/client-go/tools/clientcmd"

	"github.com/rancher/auth/server"
	"github.com/rancher/types/config"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	logrus.SetLevel(logrus.DebugLevel)
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	if err != nil {
		return err
	}

	mgmtCtx, err := config.NewManagementContext(*kubeConfig)
	if err != nil {
		return err
	}

	handler, err := server.NewTokenAPIHandler(context.Background(), mgmtCtx)
	if err != nil {
		return err
	}

	fmt.Println("Listening on 0.0.0.0:1234")
	return http.ListenAndServe("0.0.0.0:1234", handler)
}
