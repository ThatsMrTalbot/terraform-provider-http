package http

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

type TestHttpMock struct {
	server *httptest.Server
}

var testProviders = map[string]terraform.ResourceProvider{
	"http": Provider(),
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func setUpMockHttpServer() *TestHttpMock {
	files := make(map[string][]byte)
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			w.Header().Set("Content-Type", "text/plain")
			switch r.URL.Path {
			case "/meta_200.txt":
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("1.0.0"))
			case "/restricted/meta_200.txt":
				if r.Header.Get("Authorization") == "Zm9vOmJhcg==" {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("1.0.0"))
				} else {
					w.WriteHeader(http.StatusForbidden)
				}
			case "/utf-8/meta_200.txt":
				w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("1.0.0"))
			case "/utf-16/meta_200.txt":
				w.Header().Set("Content-Type", "application/json; charset=UTF-16")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("\"1.0.0\""))
			case "/meta_404.txt":
				w.WriteHeader(http.StatusNotFound)
			default:
				// Handle auth for paths
				switch {
				case strings.HasPrefix(r.URL.Path, "/basic"):
					if user, pass, _ := r.BasicAuth(); user != "user" || pass != "pass" {
						w.WriteHeader(http.StatusForbidden)
						return
					}
				case strings.HasPrefix(r.URL.Path, "/restricted"):
					if r.Header.Get("Authorization") != "Zm9vOmJhcg==" {
						w.WriteHeader(http.StatusForbidden)
						return
					}
				}

				// Handle PUT/GET for resource testing
				switch r.Method {
				case "GET":
					if file, ok := files[r.URL.Path]; !ok {
						w.WriteHeader(http.StatusNotFound)
					} else {
						w.WriteHeader(http.StatusOK)
						w.Write(file)
					}
				case "PUT":
					files[r.URL.Path], _ = ioutil.ReadAll(r.Body)
				}

			}
		}),
	)

	return &TestHttpMock{
		server: server,
	}
}
