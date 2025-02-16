package response

import (
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
			want: "{\n\t \"anwser\": \"error\",\n\t \"can't read\": \"db error\"\n}",
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
			want: "{\n\t \"anwser\": \"error\",\n\t \"can't read\": \"db error\",\n\t \"can't parse\": \"parse error\",\n\t \"internal error\": \"internal error\"\n}",
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
			want: "{\n\t \"anwser\": \"error\",\n\t \"empty error\": \"\"\n}",
		},
		{
			name: "error with special characters",
			fields: fields{
				Anwer: "error",
				mapErrorCause: map[string]string{
					"db error with \n newline and \" quote": "can't read",
				},
			},
			want: "{\n\t \"anwser\": \"error\",\n\t \"can't read\": \"db error with \n newline and \" quote\"\n}",
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
				expected1 := "{\n\t \"anwser\": \"error\",\n\t \"can't read\": \"db error\",\n\t \"can't parse\": \"parse error\",\n\t \"internal error\": \"internal error\"\n}"
				expected2 := "{\n\t \"anwser\": \"error\",\n\t \"can't read\": \"db error\",\n\t \"internal error\": \"internal error\",\n\t \"can't parse\": \"parse error\"\n}"
				expected3 := "{\n\t \"anwser\": \"error\",\n\t \"can't parse\": \"parse error\",\n\t \"can't read\": \"db error\",\n\t \"internal error\": \"internal error\"\n}"
				expected4 := "{\n\t \"anwser\": \"error\",\n\t \"can't parse\": \"parse error\",\n\t \"internal error\": \"internal error\",\n\t \"can't read\": \"db error\"\n}"
				expected5 := "{\n\t \"anwser\": \"error\",\n\t \"internal error\": \"internal error\",\n\t \"can't parse\": \"parse error\",\n\t \"can't read\": \"db error\"\n}"
				expected6 := "{\n\t \"anwser\": \"error\",\n\t \"internal error\": \"internal error\",\n\t \"can't read\": \"db error\",\n\t \"can't parse\": \"parse error\"\n}"

				if got != expected1 && got != expected2 && got != expected3 && got != expected4 && got != expected5 && got != expected6 {
					t.Errorf("ErrorResp.JsonString() = %s, want %s or %s or %s or %s or %s or %s", got, expected1, expected2, expected3, expected4, expected5, expected6)
				}
			} else {
				if got != tt.want {
					t.Errorf("ErrorResp.JsonString() = \n%s\n,\n want \n%s", got, tt.want)
				}
			}

		})
	}
}
