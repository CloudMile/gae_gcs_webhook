package app

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

// ComputeService is GCP compute setvice
type ComputeService struct {
	Ctx            context.Context
	ComputeService *compute.Service
	isError        bool
}

// Get for setup GCP auth client service
func (cs *ComputeService) Get() {
	cs.isError = false
	client := &http.Client{
		Transport: &oauth2.Transport{
			Source: google.AppEngineTokenSource(cs.Ctx, compute.CloudPlatformScope),
			Base: &urlfetch.Transport{
				Context: cs.Ctx,
			},
		},
	}
	computeService, err := compute.New(client)

	if err != nil {
		log.Errorf(cs.Ctx, "compute error: %s", err)
		cs.isError = true
		return
	}
	cs.ComputeService = computeService
}
