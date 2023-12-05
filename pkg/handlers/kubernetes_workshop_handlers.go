package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"k8s.io/klog"
)

const (
	UpstreamService = "UPSTREAM_SERVICE" //ENV KEY
	VERSION         = "VERSION"
)

type Entity struct {
	ServiceName string
	Version     string
	Hostname    string
}

type KubernetesWorkshop struct {
	ServiceName string
}

func (k *KubernetesWorkshop) Info(c *gin.Context) {
	entity := k.generateServiceInfo()
	c.JSON(http.StatusOK, entity)
}

func (k *KubernetesWorkshop) GetInfoFromService(c *gin.Context) {
	upStreamService := os.Getenv(UpstreamService)
	localServiceInfo := k.generateServiceInfo()
	if upStreamService == "" {
		c.JSON(http.StatusOK, []*Entity{localServiceInfo})
		return
	}
	client := resty.New()
	entityList := &[]*Entity{{}}
	resp, err := client.R().
		ForceContentType("application/json").
		Get(fmt.Sprintf("%s/service", upStreamService))
	if err != nil {
		klog.Fatal(err)
	}
	if err := json.Unmarshal(resp.Body(), entityList); err != nil {
		klog.Fatal(err)
	}
	*entityList = append(*entityList, localServiceInfo)
	c.JSON(http.StatusOK, entityList)
}

func (k *KubernetesWorkshop) generateServiceInfo() *Entity {
	hostname, err := os.Hostname()
	version := os.Getenv(VERSION)
	if err != nil {
		klog.Fatal(err)
	}
	entity := &Entity{
		ServiceName: k.ServiceName,
		Hostname:    hostname,
		Version:     version,
	}
	return entity
}
