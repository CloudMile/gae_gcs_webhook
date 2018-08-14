# GCS Webhook and Invalidated CDN Cache by GAE

## Abstract
When GCS bucket's objects update/create/delete, the GAE will be got webhook and invalidate CDN cache

## Setup GAE

Create a new project
1. Choose App Engine.
2. Choose language for development, we using `Go`.
3. Choose deploy location.
4. You can skip the tutorial.
5. Refresh the page, and you can see the GAE main page.

## Setup Service Account
You have two GCP project, one is for CDN, the other is for GAE to invalidate CDN cache.
1. Go to CDN project, create a service account for do GCS webhook, the role is `Storage Admin`.
2. Go to CDN project, create a custom role for invalidate CDN cache, the permission is `compute.urlMaps.invalidateCache`.
3. Go to GAE project, copy the GAE service account email.
4. Go to CDN project, create IAM member with GAE service account email and the role is custom role for invalidate CDN cache.

## Domain Verification
Go to CDN project, `APIs & Service` -> `Credentials` to verificate the GAE url.
Chose the `HTML file upload` and download the HTML file.


## Deploy GAE
Get source code.
```
$ git clone git@github.com:CloudMile/gae_gcs_webhook.git
```

Edit the `app.yaml`
```
env_variables:
  PROJECT_ID: <YOUR_CDN_PROJECT_ID>
  URL_MAP: <YOUR_CDN_URL_MAP>
  DOMAIN_VERIFICATION: <GOOGLE_VERIFICATION>
```

Edit the env variables,
- PROJECT_ID: your CDN project id
- URL_MAP: your CDN url map name
- DOMAIN_VERIFICATION: the `HTML file name` from the `Domain Verification`

and aslo copy the HTML file into this folder.

Deploy
```
$ gcloud config set project <YOUR_GAE_PROJECT_ID>
$ gcloud app deploy app.yaml queue.yaml
```

## Setup GCS Watchbucket
```
$ gsutil acl ch -u <YOUR_GAE_SERVICE_ACCOUNT_EMAIL>:OWNER gs://<YOUR_CDN_SOURCE_BUCKET>
```

change gcloud auth to service account which you create before, the role is `Storage Admin`.
```
$ gcloud auth activate-service-account <GCS_WEBHOOK_SERVICE_ACCOUNT_EMAIL> --key-file <YOUR_KEY_PATH>
```

setup webhook
```
$ gsutil notification watchbucket https://<YOUR_GAE_PROJECT_ID>.appspot.com/ gs://<YOUR_CDN_SOURCE_BUCKET>
```
and remember the `CHANNEL_ID` and `RESOURCE_ID`, it will be useful when remove channel.

## Test
```
$ gsutil -m rsync ./<CDN_FILE_PATH>/ gs://<YOUR_CDN_SOURCE_BUCKET>/<CDN_FILE_PATH>/
```

go to CDN project, and you will see the invalidate CDN cache work.


## Remove Channel
```
$ gcloud auth activate-service-account <GCS_WEBHOOK_SERVICE_ACCOUNT_EMAIL> --key-file <YOUR_KEY_PATH>
$ gsutil notification stopchannel <CHANNEL_ID> <RESOURCE_ID>
```

## Refs
- https://cloud.google.com/storage/docs/object-change-notification
- https://cloud.google.com/storage/docs/gsutil/commands/notification
- https://cloud.google.com/compute/docs/reference/rest/v1/urlMaps/invalidateCache
