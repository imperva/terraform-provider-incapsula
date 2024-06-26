package incapsula

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetSslInstructions(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_ssl_instructions_test.TestGetSslInstructions")

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("%s", endpointSSLInstructions+"?extSiteId=123") {
			t.Errorf("Should have hit %s endpoint. Got: %s", endpointSSLInstructions, req.URL.String())
		}
		rw.WriteHeader(200)
		rw.Write([]byte("{\n  \"data\": [\n    {\n      \"domain\": \"example.com\",\n      \"validationMethod\": \"DNS\",\n      \"recordType\": \"TXT\",\n      \"verificationCode\": \"123456\",\n      \"verificationCodeExpirationDate\": 1714979200,\n      \"lastNotificationDate\": 1714892800,\n      \"relatedSansDetails\": [\n        {\n          \"sanId\": 1,\n          \"sanValue\": \"sub.example.com\",\n          \"domainIds\": [101, 102]\n        }\n      ]\n    },\n    {\n      \"domain\": \"anotherdomain.com\",\n      \"validationMethod\": \"Email\",\n      \"recordType\": \"CNAME\",\n      \"verificationCode\": \"654321\",\n      \"verificationCodeExpirationDate\": 1715979200,\n      \"lastNotificationDate\": 1715892800,\n      \"relatedSansDetails\": [\n        {\n          \"sanId\": 3,\n          \"sanValue\": \"mail.anotherdomain.com\",\n          \"domainIds\": [104, 105]\n        }\n      ]\n    }\n  ]\n}"))
	}))

	defer server.Close()
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	siteId := 123
	sslInstructions, _ := client.GetSiteSSLInstructions(siteId)

	checkSslInstructionsResponse(t, "Domain", "example.com", sslInstructions.Data[0].Domain)
	checkSslInstructionsResponse(t, "ValidationMethod", "DNS", sslInstructions.Data[0].ValidationMethod)
	checkSslInstructionsResponse(t, "RecordType", "TXT", sslInstructions.Data[0].RecordType)
	checkSslInstructionsResponse(t, "VerificationCode", "123456", sslInstructions.Data[0].VerificationCode)
	checkSslInstructionsResponse(t, "VerificationCodeExpirationDate", "1714979200", fmt.Sprintf("%d", sslInstructions.Data[0].VerificationCodeExpirationDate))
	checkSslInstructionsResponse(t, "LastNotificationDate", "1714892800", fmt.Sprintf("%d", sslInstructions.Data[0].LastNotificationDate))
	checkSslInstructionsResponse(t, "RelatedSansDetails", "sub.example.com", sslInstructions.Data[0].RelatedSansDetails[0].SanValue)
	checkSslInstructionsResponse(t, "RelatedSansDetails", "mail.anotherdomain.com", sslInstructions.Data[1].RelatedSansDetails[0].SanValue)
	checkSslInstructionsResponse(t, "RelatedSansDetails", "101", fmt.Sprintf("%d", sslInstructions.Data[0].RelatedSansDetails[0].DomainIds[0]))
	checkSslInstructionsResponse(t, "RelatedSansDetails", "102", fmt.Sprintf("%d", sslInstructions.Data[0].RelatedSansDetails[0].DomainIds[1]))
	checkSslInstructionsResponse(t, "RelatedSansDetails", "104", fmt.Sprintf("%d", sslInstructions.Data[1].RelatedSansDetails[0].DomainIds[0]))
}

func checkSslInstructionsResponse(t *testing.T, fieldName string, expected string, result string) {
	if expected != result {
		t.Errorf(" %s Should have  value %s.Got: %s", fieldName, expected, result)
	}
}
