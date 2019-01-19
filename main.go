package main

import (
	"net/http"
	"os"

	"github.com/CloudMile/gae_gcs_webhook/controller"
	"google.golang.org/appengine"
)

func main() {
	http.HandleFunc("/", controller.HomeHandle)
	http.HandleFunc("/queue", controller.QueueHandle)
	domainVerification := os.Getenv("DOMAIN_VERIFICATION")

	http.HandleFunc(`/`+domainVerification, controller.GooHandle)
	appengine.Main()
}
