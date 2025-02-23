package register

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/SapolovichSV/durak/auth/internal/config"
	"github.com/SapolovichSV/durak/auth/internal/entities/response"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	register "github.com/SapolovichSV/durak/auth/internal/http/handlers/register/mocks"
	"github.com/SapolovichSV/durak/auth/internal/logger"
)

func TestHandler_Register(t *testing.T) {
	type fields struct {
		log  logger.Logger
		repo strge
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Handler{
				log:  tt.fields.log,
				repo: tt.fields.repo,
			}
			c.Register(tt.args.w, tt.args.r)
		})
	}
}
func TestHandler_Register_OkCases(t *testing.T) {
	var methodPost = "POST"
	var handlerPath = "/auth/register"
	type fields struct {
		log  logger.Logger
		repo strge
	}
	type args struct {
		w    *httptest.ResponseRecorder
		body io.Reader
	}
	type want struct {
		code int
		body string
	}
	defaultFields := func() fields {
		return fields{
			logger.New(config.Config{LogLevel: -4}),

			func() strge {
				mockStor := register.NewMockstrge(t)
				mockStor.EXPECT().AddUser(
					mock.Anything, mock.Anything, mock.Anything, mock.Anything,
				).Return(nil)
				return mockStor
			}(),
		}
	}

	defaultWant := func() want {
		return want{
			code: 201,
			body: `{
    "anwser": "created"
}`,
		}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name:   "defaultTestCase",
			fields: defaultFields(),
			args: args{
				w:    httptest.NewRecorder(),
				body: strings.NewReader(`{"email":"ballpla45@gmail.com","username":"Stas","password":"12347112"}`),
			},
			want: defaultWant(),
		},
		{
			name:   "defaultTestCase2",
			fields: defaultFields(),
			args: args{
				w:    httptest.NewRecorder(),
				body: strings.NewReader(`{"email":"ballplayer45@gmail.com","username":"Stasyan125","password":"antiHype"}`),
			},
			want: defaultWant(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(methodPost, handlerPath, tt.args.body)

			c := Handler{
				tt.fields.log,
				tt.fields.repo,
			}
			c.Register(tt.args.w, r)

			res := tt.args.w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.want.code {
				t.Fatalf("invalid code want =  %d,have %d", tt.want.code, res.StatusCode)
			}

			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("error at reading response body %s", err.Error())
			}

			require.JSONEq(t, tt.want.body, string(data))
		})
	}

}
func TestHandler_Register_ValidationErrorCases(t *testing.T) {
	var methodPost = "POST"
	var handlerPath = "/auth/register"
	type fields struct {
		log  logger.Logger
		repo strge
	}
	defaultFields := func() fields {
		return fields{
			logger.New(config.Config{LogLevel: -4}),

			func() strge {
				mockStor := register.NewMockstrge(t)
				return mockStor
			}(),
		}
	}

	tests := []struct {
		name     string
		fields   fields
		wantCode int
		body     io.Reader
	}{
		{
			name:     "Empty body fields",
			fields:   defaultFields(),
			wantCode: 400,
			body: strings.NewReader(
				`{"email":"","username":"","password":""}`,
			),
		},
		{
			name:     "Incorrect email",
			fields:   defaultFields(),
			wantCode: 400,
			body: strings.NewReader(
				`{"email":"asdadskjh@@@mail.com","username":"someGoodName","password":"123123"}`,
			),
		},
		{
			name:     "Only empty password",
			fields:   defaultFields(),
			wantCode: 400,
			body: strings.NewReader(
				`{"email":"asdadskjh@gmail.com","username":"someGoodName","password":""}`,
			),
		}, {
			name:     "Only empty email",
			fields:   defaultFields(),
			wantCode: 400,
			body: strings.NewReader(
				`{"email":"","username":"someGoodName","password":"asads"}`,
			),
		}, {
			name:     "Only empty username",
			fields:   defaultFields(),
			wantCode: 400,
			body: strings.NewReader(
				`{"email":"ballplayerr@gmail.com","username":"","password":"asads"}`,
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(methodPost, handlerPath, tt.body)
			c := Handler{
				log:  tt.fields.log,
				repo: tt.fields.repo,
			}
			c.Register(w, r)
			res := w.Result()
			if res.StatusCode != tt.wantCode {
				resBody, _ := io.ReadAll(res.Body)
				t.Logf("resp body: %s", string(resBody))
				t.Fatalf("want code : %d,have code : %d", tt.wantCode, res.StatusCode)
			}
		})
	}
}
func TestHandler_Register_RepoErrorCases(t *testing.T) {
	var methodPost = "POST"
	var handlerPath = "/auth/registration"
	type fields struct {
		log  logger.Logger
		repo strge
	}
	type args struct {
		emailUserPass []string
		err           error
	}
	defaultFields := func(email, username, password string, returningError error) fields {
		return fields{
			logger.New(config.Config{LogLevel: -4}),

			func() strge {
				mockStor := register.NewMockstrge(t)
				mockStor.EXPECT().AddUser(
					mock.Anything, email, username, password,
				).Return(returningError)
				return mockStor
			}(),
		}
	}
	tests := []struct {
		name     string
		args     args
		fields   fields
		wantCode int
		wantBody string
		body     io.Reader
	}{
		{
			name: "AlreadyExistsError",
			args: args{
				emailUserPass: []string{"already@gmail.com", "staspiska", "pass123"},
				err:           errors.New("such user already exist"),
			},
			wantCode: 500,
			wantBody: response.NewErrorResp(
				map[string]error{"can't add user to repo": errors.New("such user already exist")},
			).JsonString(),
			body: strings.NewReader(
				`{"email":"already@gmail.com","username":"staspiska","password":"pass123"}`,
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = defaultFields(tt.args.emailUserPass[0], tt.args.emailUserPass[1], tt.args.emailUserPass[2], tt.args.err)
			w := httptest.NewRecorder()
			r := httptest.NewRequest(methodPost, handlerPath, tt.body)

			c := Handler{
				tt.fields.log,
				tt.fields.repo,
			}

			c.Register(w, r)
			res := w.Result()

			if res.StatusCode != tt.wantCode {
				t.Fatalf("want code: %d have code: %d", tt.wantCode, res.StatusCode)
			}
			data, err := io.ReadAll(res.Body)
			defer res.Body.Close()
			if err != nil {
				t.Fatalf("unexpected error at reading response data: %s", err.Error())
			}
			require.JSONEq(t, tt.wantBody, string(data))
		})
	}
}
func TestHandler_Register_UnmarshalErrorCases(t *testing.T) {
	var methodPost = "POST"
	var handlerPath = "/auth/register"
	type fields struct {
		log  logger.Logger
		repo strge
	}
	defaultFields := func() fields {
		return fields{
			logger.New(config.Config{LogLevel: -4}),
			func() strge {
				mockStor := register.NewMockstrge(t)
				// Для случая ошибки разбора JSON репозиторий не вызывается
				return mockStor
			}(),
		}
	}

	tests := []struct {
		name        string
		fields      fields
		body        io.Reader
		wantCode    int
		errContains string
	}{
		{
			name:   "Invalid JSON input",
			fields: defaultFields(),
			// Передаём некорректный JSON (например, отсутствует закрывающая фигурная скобка)
			body:        strings.NewReader(`{"email":"test@mail.com", "username": "test"`),
			wantCode:    http.StatusBadRequest,
			errContains: "can't parse json data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(methodPost, handlerPath, tt.body)

			c := Handler{
				log:  tt.fields.log,
				repo: tt.fields.repo,
			}
			c.Register(w, r)
			res := w.Result()
			defer res.Body.Close()

			require.Equal(t, tt.wantCode, res.StatusCode, "unexpected status code")
			data, err := io.ReadAll(res.Body)
			require.NoError(t, err, "failed to read response body")

			// Парсим ответ и проверяем, что он содержит требуемое сообщение об ошибке
			var parsed map[string]string
			err = json.Unmarshal(data, &parsed)
			require.NoError(t, err, "response is not valid JSON")
			require.Equal(t, "error", parsed["anwser"], "expected 'anwser' to be 'error'")
			found := false
			for key := range parsed {
				if key != "anwser" && strings.Contains(key, tt.errContains) {
					found = true
					break
				}
			}
			require.True(t, found, "expected error key to contain %q, got: %v", tt.errContains, parsed)
		})
	}
}
