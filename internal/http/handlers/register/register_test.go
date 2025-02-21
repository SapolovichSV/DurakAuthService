package register

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/SapolovichSV/durak/auth/internal/config"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	register "github.com/SapolovichSV/durak/auth/internal/http/handlers/register/mocks"
	"github.com/SapolovichSV/durak/auth/internal/logger"
)

func TestHandler_Register(t *testing.T) {
	type fields struct {
		log  logger.Logger
		repo storage
		ctx  context.Context
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
				ctx:  tt.fields.ctx,
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
		repo storage
		ctx  context.Context
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

			func() storage {
				mockStor := register.NewMockstorage(t)
				mockStor.EXPECT().AddUser(
					mock.Anything, mock.Anything, mock.Anything, mock.Anything,
				).Return(nil)
				return mockStor
			}(),

			context.Background(),
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
				tt.fields.ctx,
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

}
