package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	//"github.com/gorilla/mux"
	"github.com/rancher/auth/identities"
	//"github.com/rancher/auth/tokens"
	//"github.com/sirupsen/logrus"
	//"github.com/urfave/cli"

	"github.com/rancher/types/config"

	"k8s.io/client-go/tools/clientcmd"
)

//var VERSION = "v0.0.0-dev"
/*
func main() {
	app := cli.NewApp()
	app.Name = "auth"
	app.Version = VERSION
	app.Usage = "You need help!"
	app.Action = run

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "cluster-config",
			Usage: "Kube config for accessing cluster",
		},
		cli.StringFlag{
			Name:  "cluster-name",
			Usage: "name of the cluster",
		},
		cli.StringFlag{
			Name:  "httpHost",
			Usage: "host:port to listen on",
		},
	}

	app.Run(os.Args)
}

func run(c *cli.Context) {

	mgmtCtx, err := setupClient(c.String("cluster-config"), c.String("cluster-config"), "")
	if err != nil {
		log.Fatalf("Failed to create ManagementContext: %v", err)
	}

	tokenHandler, err := tokens.NewTokenAPIHandler(nil, mgmtCtx)
	if err != nil {
		log.Fatalf("Failed to get tokenAndIdentity handler: %v", err)
	}

	identityHandler, err := identities.NewIdentityAPIHandler(nil, mgmtCtx)
	if err != nil {
		log.Fatalf("Failed to get NewIdentityAPIHandler handler: %v", err)
	}

	if c.GlobalBool("debug") {
		log.SetLevel(log.DebugLevel)
	}

	textFormatter := &log.TextFormatter{
		FullTimestamp: true,
	}
	log.SetFormatter(textFormatter)

	log.Info("Starting Rancher Auth proxy")

	httpHost := c.GlobalString("httpHost")

	router := mux.NewRouter()
	router.Handle("/v3/tokens", tokenHandler).Methods("GET", "POST", "PUT", "DELETE", "PATCH", "HEAD")
	router.Handle("/v3/identities", identityHandler).Methods("GET", "POST", "PUT", "DELETE", "PATCH", "HEAD")

	log.Infof("Starting http server listening on %v.", httpHost)
	log.Fatal(http.ListenAndServe(httpHost, router))

}

func setupClient(clusterManagerCfg string, clusterCfg string, clusterName string) (*config.ManagementContext, error) {
	clusterManagementKubeConfig, err := clientcmd.BuildConfigFromFlags("", clusterManagerCfg)
	if err != nil {
		return nil, err
	}
	workload, err := config.NewManagementContext(*clusterManagementKubeConfig)
	if err != nil {
		return nil, err
	}

	return workload, nil

}
*/

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

	handler, err := identities.NewIdentityAPIHandler(context.Background(), mgmtCtx)
	if err != nil {
		return err
	}

	fmt.Println("Listening on 0.0.0.0:1234")
	return http.ListenAndServe("0.0.0.0:1234", handler)
}
