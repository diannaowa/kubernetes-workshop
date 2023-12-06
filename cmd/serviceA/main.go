package main

import (
	"bytes"

	"github.com/gin-gonic/gin"
	"k8s.io/klog"

	"kubernetes-workshop/pkg/handlers"
)

func main() {
	h := &handlers.KubernetesWorkshop{
		ServiceName:     "serviceA",
		MemoryBlackHole: bytes.Buffer{},
	}
	r := gin.Default()
	r.GET("/info", h.Info)
	r.GET("/mem", h.Mem)
	r.GET("/service", h.GetInfoFromService)
	if err := r.Run(); err != nil {
		klog.Fatal(err)
	}
}
