package model

import (
	"context"
	"os"
	"time"

	compute "google.golang.org/api/compute/v1"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
)

// GCSObj is request body
type GCSObj struct {
	Ctx            context.Context
	ComputeService *compute.Service
	Bucket         string `json:"bucket"`
	ObjectName     string `json:"name"`
	Md5Hash        string `json:"md5Hash"`
}

// InvalidateCache will invalidate CDN cache
func (gcsObj *GCSObj) InvalidateCache() {
	project := os.Getenv("PROJECT_ID")
	urlMap := os.Getenv("URL_MAP")

	canInvalidateCache := checkFromGAEMemcache(gcsObj)

	if canInvalidateCache {
		rb := &compute.CacheInvalidationRule{
			Path: `/` + gcsObj.ObjectName,
		}
		resp, err := gcsObj.ComputeService.UrlMaps.InvalidateCache(project, urlMap, rb).Context(gcsObj.Ctx).Do()
		if err != nil {
			log.Errorf(gcsObj.Ctx, "InvalidateCache error %v", err)
		}
		log.Infof(gcsObj.Ctx, "InvalidateCache %#v\n", resp)
	}
}

func checkFromGAEMemcache(gcsObj *GCSObj) (canInvalidateCache bool) {
	memcacheKey := `MD5_` + gcsObj.Md5Hash + `_NAME_` + gcsObj.ObjectName
	canInvalidateCache = false
	_, err := memcache.Get(gcsObj.Ctx, memcacheKey)
	if err == memcache.ErrCacheMiss {
		item := &memcache.Item{
			Key:        memcacheKey,
			Value:      []byte(`true`),
			Expiration: time.Duration(300) * time.Second,
		}

		if err := memcache.Set(gcsObj.Ctx, item); err != nil {
			log.Errorf(gcsObj.Ctx, "error adding item: %v", err)
		}
		canInvalidateCache = true
	}
	return
}
