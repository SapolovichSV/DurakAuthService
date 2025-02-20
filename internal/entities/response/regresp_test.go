package response

import (
	"strings"
	"testing"
)

func TestErrorResp_JsonString(t *testing.T) {
	type fields struct {
		Anwer         string
		mapErrorCause map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "single error",
			fields: fields{
				Anwer: "error",
				mapErrorCause: map[string]string{
					"db error": "can't read",
				},
			},
			want: "{\n\t \"anwser\": \"error\",\n\t \"db error\": \"can't read\"\n}",
		},
		{
			name: "multiple errors",
			fields: fields{
				Anwer: "error",
				mapErrorCause: map[string]string{
					"db error":       "can't read",
					"parse error":    "can't parse",
					"internal error": "internal error",
				},
			},
			// Not used for direct string comparison
			want: "",
		},
		{
			name: "no errors",
			fields: fields{
				Anwer:         "error",
				mapErrorCause: map[string]string{},
			},
			want: `{
     "anwser": "error"
}`,
		},
		{
			name: "empty error message",
			fields: fields{
				Anwer: "error",
				mapErrorCause: map[string]string{
					"": "empty error",
				},
			},
			want: "{\n\t \"anwser\": \"error\",\n\t \"\": \"empty error\"\n}",
		},
		{
			name: "error with special characters",
			fields: fields{
				Anwer: "error",
				mapErrorCause: map[string]string{
					"db error with \n newline and \" quote": "can't read",
				},
			},
			want: "{\n\t \"anwser\": \"error\",\n\t \"db error with \n newline and \" quote\": \"can't read\"\n}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := ErrorResp{
				Anwer:         tt.fields.Anwer,
				mapErrorCause: tt.fields.mapErrorCause,
			}
			got := r.JsonString()

			if tt.name == "multiple errors" {
				// For multiple errors, check that output has correct prefix and all expected key-value pairs
				expectedPrefix := "{\n\t \"anwser\": \"error\","
				if !strings.HasPrefix(got, expectedPrefix) {
					t.Errorf("Expected output to have prefix %q, got %q", expectedPrefix, got)
				}
				// Check for presence of each key-value entry
				pairs := []struct {
					key   string
					value string
				}{
					{"db error", "can't read"},
					{"parse error", "can't parse"},
					{"internal error", "internal error"},
				}
				for _, pair := range pairs {
					entry := "\"" + pair.key + "\": \"" + pair.value + "\""
					if !strings.Contains(got, entry) {
						t.Errorf("Expected output to contain %q, got %q", entry, got)
					}
				}
				// Check that output starts with { and ends with }
				if !strings.HasPrefix(got, "{") || !strings.HasSuffix(got, "}") {
					t.Errorf("Expected output to be enclosed in braces, got %q", got)
				}
			} else {
				if got != tt.want {
					t.Errorf("ErrorResp.JsonString() = \n%s\n,\n want \n%s", got, tt.want)
				}
			}
		})
	}
}
