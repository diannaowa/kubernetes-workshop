package main

import (
	"github.com/gin-gonic/gin"
	"k8s.io/klog"

	"kubernetes-workshop/pkg/handlers"
)

func main() {
	h := handlers.KubernetesWorkshop{
		ServiceName: "serviceA",
	}
	r := gin.Default()
	r.GET("/info", h.Info)
	r.GET("/service", h.GetInfoFromService)
	if err := r.Run(); err != nil {
		klog.Fatal(err)
	}
}
