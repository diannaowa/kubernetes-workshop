package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"k8s.io/klog"
)

type KubernetesWorkshop struct {
	ServiceTag string
}

func (k *KubernetesWorkshop) Info(c *gin.Context) {
	hostname, err := os.Hostname()
	version := os.Getenv("VERSION")
	if err != nil {
		klog.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{
		k.ServiceTag: fmt.Sprintf("hostname=%s,version=%s", hostname, version),
	})
}
