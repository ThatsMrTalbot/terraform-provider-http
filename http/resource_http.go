package http

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceHTTP() *schema.Resource {
	return &schema.Resource{
		Read:   resourceHTTPRead,
		Create: resourceHTTPCreate,
		Update: resourceHTTPCreate,
		Delete: resourceHTTPDelete,

		Schema: map[string]*schema.Schema{
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"http_user": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HTTP_USER", ""),
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"http_pass": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HTTP_PASS", ""),
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"request_headers": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"body": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceHTTPRead(d *schema.ResourceData, meta interface{}) error {
	url := d.Get("url").(string)
	headers := d.Get("request_headers").(map[string]interface{})

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("Error creating request: %s", err)
	}

	for name, value := range headers {
		req.Header.Set(name, value.(string))
	}

	if _, ok := d.GetOk("http_user"); ok {
		req.SetBasicAuth(d.Get("http_user").(string), d.Get("http_pass").(string))
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error during making a request: %s", url)
	}

	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		d.SetId("")
		return nil
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP request error. Response code: %d", resp.StatusCode)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error while reading response body. %s", err)
	}

	d.Set("body", string(bytes))
	d.SetId(fmt.Sprintf("%x", md5.Sum(bytes)))

	return nil
}

func resourceHTTPCreate(d *schema.ResourceData, meta interface{}) error {
	url := d.Get("url").(string)
	body := d.Get("body").(string)
	headers := d.Get("request_headers").(map[string]interface{})

	client := &http.Client{}

	req, err := http.NewRequest("PUT", url, bytes.NewReader([]byte(body)))
	if err != nil {
		return fmt.Errorf("Error creating request: %s", err)
	}

	for name, value := range headers {
		req.Header.Set(name, value.(string))
	}

	if _, ok := d.GetOk("http_user"); ok {
		req.SetBasicAuth(d.Get("http_user").(string), d.Get("http_pass").(string))
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error during making a request: %s", url)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP request error. Response code: %d", resp.StatusCode)
	}

	d.SetId(fmt.Sprintf("%x", md5.Sum([]byte(body))))

	return nil
}

func resourceHTTPDelete(d *schema.ResourceData, meta interface{}) error {
	url := d.Get("url").(string)
	headers := d.Get("request_headers").(map[string]interface{})

	client := &http.Client{}

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("Error creating request: %s", err)
	}

	for name, value := range headers {
		req.Header.Set(name, value.(string))
	}

	if _, ok := d.GetOk("http_user"); ok {
		req.SetBasicAuth(d.Get("http_user").(string), d.Get("http_pass").(string))
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error during making a request: %s", url)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP request error. Response code: %d", resp.StatusCode)
	}

	d.SetId("")

	return nil
}
