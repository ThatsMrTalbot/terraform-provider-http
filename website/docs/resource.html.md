---
layout: "http"
page_title: "HTTP Resource"
sidebar_current: "docs-http-resource"
description: |-
  Stores content on a HTTP server using PUT requests
---

# `http` Resource

The `http` data source makes an HTTP PUT request to set data on a given url. It
uses a HTTP GET request to check the existence of the resource and a HTTP DELETE 
request when the resource is deleted.

The given URL may be either an `http` or `https` URL.

~> **Important** Although `https` URLs can be used, there is currently no
mechanism to authenticate the remote server except for general verification of
the server certificate's chain of trust. Data retrieved from servers not under
your control should be treated as untrustworthy.

## Example Usage

```hcl
resource "http" "example" {
  url  = "https://example.com/resource.txt"
  data = "${file("${path.module}/resource.txt")}"

  # Optional request headers
  request_headers {
    "Accept" = "application/json"
  }

  # Optional basic auth
  http_user = "username"
  http_pass = "password"
}
```

## Argument Reference

The following arguments are supported:

* `url` - (Required) The URL to request data from. 

* `data` - (Required) The data to store on the URL. 

* `http_user` - (Optional) The HTTP username to send with basic auth. Can also 
  be set with the HTTP_USER environment variable.

* `http_pass` - (Optional) The HTTP password to send with basic auth. Can also 
  be set with the HTTP_PASS environment variable.

* `request_headers` - (Optional) A map of strings representing additional HTTP
  headers to include in the request.
  