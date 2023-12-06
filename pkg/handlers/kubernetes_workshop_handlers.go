package handlers

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"

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
	ServiceName     string
	MemoryBlackHole []byte
}

func (k *KubernetesWorkshop) Info(c *gin.Context) {
	entity := k.generateServiceInfo()
	c.JSON(http.StatusOK, entity)
}

func (k *KubernetesWorkshop) Mem(c *gin.Context) {
	// 10M
	for i := 0; i < 10*128*1024; i++ {
		v := uint64(100)
		k.MemoryBlackHole = binary.BigEndian.AppendUint64(k.MemoryBlackHole, v)
	}
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	//
	c.JSON(http.StatusOK, gin.H{"Alloc": float64(rtm.Alloc) / 1024. / 1024.})
	return
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
		klog.Error(err)
		return
	}
	if err := json.Unmarshal(resp.Body(), entityList); err != nil {
		klog.Error(err)
		return
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
