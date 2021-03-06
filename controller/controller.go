package controller

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/CloudMile/gae_gcs_webhook/model"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/taskqueue"
)

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	switch status {
	case http.StatusNotFound:
		fmt.Fprint(w, "404 Not Found")
	case http.StatusMethodNotAllowed:
		fmt.Fprint(w, "405 Method Not Allow")
	case http.StatusInternalServerError:
		fmt.Fprint(w, "500 Internal Server Error")
	default:
		fmt.Fprint(w, "Bad Request")
	}
}

// GooHandle is for domain verification
func GooHandle(w http.ResponseWriter, r *http.Request) {
	gcsObj := model.GCSObj{}
	domainVerification := os.Getenv("DOMAIN_VERIFICATION")
	t, _ := template.ParseFiles(domainVerification)
	t.Execute(w, gcsObj)
}

// HomeHandle is call a queue job for invalidate-cache
func HomeHandle(w http.ResponseWriter, r *http.Request) {
	body, isError := checkRequest(w, r)
	if isError {
		return
	}

	ctx := appengine.NewContext(r)
	var gcsObj model.GCSObj
	json.Unmarshal(body, &gcsObj)

	urlValues := url.Values{
		"Bucket":     {gcsObj.Bucket},
		"ObjectName": {gcsObj.ObjectName},
		"Md5Hash":    {gcsObj.Md5Hash},
	}
	log.Infof(ctx, "urlValues: %#v", urlValues)
	log.Infof(ctx, "channel id: %v", r.Header.Get(`X-Goog-Channel-Id`))
	log.Infof(ctx, "resource id: %v", r.Header.Get(`X-Goog-Resource-Id`))

	t := taskqueue.NewPOSTTask("/queue", urlValues)
	if _, err := taskqueue.Add(ctx, t, "invalidate-cache"); err != nil {
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}
}

// QueueHandle is doing invalidate-cache
func QueueHandle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/queue" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	ctx := appengine.NewContext(r)
	cs := model.ComputeService{Ctx: ctx}
	cs.Get()

	if cs.Error != nil {
		errorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	gcsObj := model.GCSObj{
		Ctx:            ctx,
		ComputeService: cs.ComputeService,
		Bucket:         r.FormValue(`Bucket`),
		ObjectName:     r.FormValue(`ObjectName`),
		Md5Hash:        r.FormValue(`Md5Hash`),
	}
	log.Infof(ctx, "body is %#v:", gcsObj)
	gcsObj.InvalidateCache()
}

func checkRequest(w http.ResponseWriter, r *http.Request) (body []byte, isError bool) {
	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusNotFound)
		isError = true
		return
	}

	if r.Method != http.MethodPost {
		errorHandler(w, r, http.StatusMethodNotAllowed)
		isError = true
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError)
		isError = true
		return
	}
	isError = false
	return
}
