package incapsula

import (
	"strings"
	"testing"
)

// Bodies copied from siem-config-service's ApiErrorResponse examples.
func TestFormatSiemApiError(t *testing.T) {
	testCases := []struct {
		name       string
		statusCode int
		body       string
		expected   string
	}{
		{
			name:       "test connection failure surfaces code and detail",
			statusCode: 400,
			body:       `{"errors":[{"status":400,"id":"2da889aa00000000dea6da80b802440c","code":"1000","source":{"pointer":"/siem-config-service/v3/connections"},"title":"Bad Request","detail":"The destination rejected the provided credentials. Check the connection credentials."}]}`,
			expected:   "received failure response for operation: create_siem_connection on SIEM connection\nstatus code: 400\ncode: 1000\ndetail: The destination rejected the provided credentials. Check the connection credentials.\ntrace id: 2da889aa00000000dea6da80b802440c",
		},
		{
			name:       "not found has detail but no code",
			statusCode: 404,
			body:       `{"errors":[{"status":404,"id":"2da889aa00000000dea6da80b802440c","source":{"pointer":"/siem-config-service/v3/"},"title":"Not Found","detail":"The requested resource was not found."}]}`,
			expected:   "received failure response for operation: create_siem_connection on SIEM connection\nstatus code: 404\ndetail: The requested resource was not found.\ntrace id: 2da889aa00000000dea6da80b802440c",
		},
		{
			name:       "server error without detail falls back to title",
			statusCode: 500,
			body:       `{"errors":[{"status":500,"id":"2da889aa00000000dea6da80b802440c","code":"1002","source":{"pointer":"/siem-config-service/v3/"},"title":"Internal Server Error"}]}`,
			expected:   "received failure response for operation: create_siem_connection on SIEM connection\nstatus code: 500\ncode: 1002\ntitle: Internal Server Error\ntrace id: 2da889aa00000000dea6da80b802440c",
		},
		{
			name:       "non standard body is returned raw",
			statusCode: 502,
			body:       "<html>Bad Gateway</html>",
			expected:   "received failure response for operation: create_siem_connection on SIEM connection\nstatus code: 502\nbody: <html>Bad Gateway</html>",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := formatSiemApiError(CreateSiemConnection, "SIEM connection", testCase.statusCode, testCase.body)

			if err.Error() != testCase.expected {
				t.Errorf("Expected:\n%s\ngot:\n%s", testCase.expected, err)
			}
		})
	}
}

func TestFormatSiemApiErrorDoesNotRepeatRawBody(t *testing.T) {
	body := `{"errors":[{"status":400,"id":"abc","code":"1006","title":"Bad Request","detail":"The connection is assigned to a log configuration."}]}`

	err := formatSiemApiError(DeleteSiemConnection, "SIEM connection", 400, body)

	if strings.Contains(err.Error(), `"errors"`) {
		t.Errorf("Should not include the raw JSON body once parsed, got: %s", err)
	}
}
