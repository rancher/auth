package main

import (
	"net/http"
	"os"

	"github.com/rancher/auth/tokens"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/rancher/types/config"

	"k8s.io/client-go/tools/clientcmd"
)

var VERSION = "v0.0.0-dev"

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

	handler, err := tokens.NewTokenAPIHandler(nil, mgmtCtx)
	if err != nil {
		log.Fatalf("Failed to get tokenAndIdentity handler: %v", err)
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
	server := &http.Server{
		Handler: handler,
		Addr:    httpHost,
	}
	log.Infof("Starting http server listening on %v.", httpHost)
	err = server.ListenAndServe()
	log.Infof("https server exited. Error: %v", err)

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
