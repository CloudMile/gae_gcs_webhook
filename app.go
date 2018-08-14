package app

import (
	"net/http"
	"os"

	"google.golang.org/appengine"
)

func init() {
	http.HandleFunc("/", homeHandle)
	http.HandleFunc("/queue", queueHandle)
	domainVerification := os.Getenv("DOMAIN_VERIFICATION")

	http.HandleFunc(`/`+domainVerification, gooHandle)
	appengine.Main()
}
