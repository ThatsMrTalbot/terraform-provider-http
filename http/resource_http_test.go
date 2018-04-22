package http

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const testResourceHTTPConfig_basic = `
resource "http" "http_test" {
  url  = "%s/%s"
  body = "%s"
}

data "http" "http_test" {
  url  = "${http.http_test.url}"
}

output "body" {
  value = "${data.http.http_test.body}"
}
`

func TestResourceSource_http200(t *testing.T) {
	testHttpMock := setUpMockHttpServer()

	defer testHttpMock.server.Close()

	resource.UnitTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(testResourceHTTPConfig_basic, testHttpMock.server.URL, "resource.txt", "1.0.0"),
				Check: func(s *terraform.State) error {
					_, ok := s.RootModule().Resources["data.http.http_test"]
					if !ok {
						return fmt.Errorf("missing data resource")
					}

					outputs := s.RootModule().Outputs

					if outputs["body"].Value != "1.0.0" {
						return fmt.Errorf(
							`'body' output is %s; want '1.0.0'`,
							outputs["body"].Value,
						)
					}

					return nil
				},
			},
		},
	})
}

const testResourceHTTPBasicConfig_basic = `
resource "http" "http_test" {
  url       = "%s/%s"
  body      = "%s"
  http_user = "%s"
  http_pass = "%s"
}

data "http" "http_test" {
  url       = "${http.http_test.url}"
  http_user = "${http.http_test.http_user}"
  http_pass = "${http.http_test.http_pass}"
}

output "body" {
  value = "${data.http.http_test.body}"
}
`

func TestResourceSource_basic200(t *testing.T) {
	testHttpMock := setUpMockHttpServer()

	defer testHttpMock.server.Close()

	resource.UnitTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(testResourceHTTPBasicConfig_basic, testHttpMock.server.URL, "basic/resource.txt", "1.0.0", "user", "pass"),
				Check: func(s *terraform.State) error {
					_, ok := s.RootModule().Resources["data.http.http_test"]
					if !ok {
						return fmt.Errorf("missing data resource")
					}

					outputs := s.RootModule().Outputs

					if outputs["body"].Value != "1.0.0" {
						return fmt.Errorf(
							`'body' output is %s; want '1.0.0'`,
							outputs["body"].Value,
						)
					}

					return nil
				},
			},
		},
	})
}

const testResourceHTTPConfig_withHeaders = `
resource "http" "http_test" {
  url  = "%s/%s"
  body = "%s"
	
  request_headers = {
    "Authorization" = "Zm9vOmJhcg=="
  }
}

data "http" "http_test" {
  url = "${http.http_test.url}"
	
  request_headers = {
    "Authorization" = "Zm9vOmJhcg=="
  }
}

output "body" {
  value = "${data.http.http_test.body}"
}
`

func TestResourceSource_withHeaders200(t *testing.T) {
	testHttpMock := setUpMockHttpServer()

	defer testHttpMock.server.Close()

	resource.UnitTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(testResourceHTTPConfig_withHeaders, testHttpMock.server.URL, "restricted/resource.txt", "1.0.0"),
				Check: func(s *terraform.State) error {
					_, ok := s.RootModule().Resources["data.http.http_test"]
					if !ok {
						return fmt.Errorf("missing data resource")
					}

					outputs := s.RootModule().Outputs

					if outputs["body"].Value != "1.0.0" {
						return fmt.Errorf(
							`'body' output is %s; want '1.0.0'`,
							outputs["body"].Value,
						)
					}

					return nil
				},
			},
		},
	})
}
