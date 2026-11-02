package incapsula

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestSiemResourceReadClearsIdWhenResourceIsNotFound(t *testing.T) {
	testCases := []struct {
		name     string
		resource *schema.Resource
	}{
		{name: "s3 connection", resource: resourceSiemConnection()},
		{name: "splunk connection", resource: resourceSiemSplunkConnection()},
		{name: "sftp connection", resource: resourceSiemSftpConnection()},
		{name: "generic http connection", resource: resourceSiemGenericHttpConnection()},
		{name: "log configuration", resource: resourceSiemLogConfiguration()},
	}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(`{"errors":[{"status":404,"id":"abc","title":"Not Found","detail":"The requested resource was not found."}]}`))
	}))
	defer server.Close()
	client := &Client{config: &Config{BaseURLAPI: server.URL}, httpClient: &http.Client{}}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			d := testCase.resource.TestResourceData()
			d.SetId(RandomLowLetterAndNumberString(24))

			err := testCase.resource.Read(d, client)

			if err != nil {
				t.Errorf("Should not have received an error, got: %s", err)
			}
			if d.Id() != "" {
				t.Errorf("Should have cleared the id, got: %s", d.Id())
			}
		})
	}
}
