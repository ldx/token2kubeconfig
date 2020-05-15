package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var (
	ServiceAccountDirectory = "/var/run/secrets/kubernetes.io/serviceaccount/"
)

func FromToken(userName, clusterName, serverURL, tokenFile, rootCAFile string) (*clientcmdapi.Config, error) {
	if clusterName == "" {
		clusterName = "default"
	}
	if serverURL == "" {
		host, port := os.Getenv("KUBERNETES_SERVICE_HOST"), os.Getenv("KUBERNETES_SERVICE_PORT")
		if len(host) == 0 || len(port) == 0 {
			return nil, fmt.Errorf("no Kubernetes master host or port set")
		}
		serverURL = fmt.Sprintf("https://%s:%s", host, port)
	}
	token, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		return nil, err
	}
	caCert, err := ioutil.ReadFile(rootCAFile)
	if err != nil {
		return nil, err
	}
	return CreateFromToken(serverURL, clusterName, userName, caCert, token), nil
}

func CreateFromToken(serverURL, clusterName, userName string, caCert, token []byte) *clientcmdapi.Config {
	contextName := fmt.Sprintf("%s@%s", userName, clusterName)
	config := &clientcmdapi.Config{
		Clusters: map[string]*clientcmdapi.Cluster{
			clusterName: {
				Server:                   serverURL,
				CertificateAuthorityData: caCert,
			},
		},
		Contexts: map[string]*clientcmdapi.Context{
			contextName: {
				Cluster:  clusterName,
				AuthInfo: userName,
			},
		},
		AuthInfos:      map[string]*clientcmdapi.AuthInfo{},
		CurrentContext: contextName,
	}
	config.AuthInfos[userName] = &clientcmdapi.AuthInfo{
		Token: string(token),
	}
	return config
}

func main() {
	serverURL := "" // Use environment variables.
	clusterName := "default"
	userName := "bootstrap"
	caCert := ServiceAccountDirectory + "ca.crt"
	token := ServiceAccountDirectory + "token"
	cfg, err := FromToken(userName, clusterName, serverURL, token, caCert)
	if err != nil {
		panic(err)
	}
	content, err := clientcmd.Write(*cfg)
	if err != nil {
		panic(err)
	}
	os.Stdout.Write(content)
}
