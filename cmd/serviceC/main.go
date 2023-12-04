package main

import (
	"github.com/gin-gonic/gin"
	"k8s.io/klog"

	"kubernetes-workshop/pkg/handlers"
)

func main() {
	h := handlers.KubernetesWorkshop{
		ServiceTag: "serviceC",
	}
	r := gin.Default()
	r.GET("/info", h.Info)
	if err := r.Run(); err != nil {
		klog.Fatal(err)
	}
}
