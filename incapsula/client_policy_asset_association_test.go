package incapsula

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientPolicyAssetAssociatedPositive(t *testing.T) {
	isAssociated, err := ClientPolicyAssetAssociatedBase(t, true)
	if err != nil {
		t.Errorf("unexpected error")
	}
	if !isAssociated {
		t.Errorf("expected policy to be assosiated")
	}
}

func TestClientPolicyAssetAssociatedNegative(t *testing.T) {
	isAssociated, err := ClientPolicyAssetAssociatedBase(t, false)
	if err != nil {
		t.Errorf("did not expected error but got one. error: %v", err)
	}
	if isAssociated {
		t.Errorf("expected policy to NOT be assosiated")
	}
}

func ClientPolicyAssetAssociatedBase(t *testing.T, shouldBeAssociated bool) (bool, error) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test Running test client_policy_asset_association.TestClientPolicyAssetAssociated, shouldBeAssociated set to %t", shouldBeAssociated)
	apiID := "foo"
	apiKey := "bar"
	assetID := "5432"
	policyID := "11"
	assetType := "WEBSITE"

	endpoint := fmt.Sprintf("/policies/v2/policies/%s/assets/%s/%s", policyID, assetType, assetID)
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		if shouldBeAssociated {
			rw.WriteHeader(200)
			rw.Write([]byte(`{"value":true,"isError":false}`))
		} else {
			rw.WriteHeader(404)
			rw.Write([]byte(`{"value":false,"isError":false}`))
		}
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	return client.isPolicyAssetAssociated(policyID, assetID, assetType, nil)

}
