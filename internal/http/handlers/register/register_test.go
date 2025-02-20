package register

import (
	"context"
	"net/http"
	"testing"

	"github.com/SapolovichSV/durak/auth/internal/config"
	"github.com/SapolovichSV/durak/auth/internal/http/handlers/register/mocks"

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
		{
			name: "defaultTestCase",
			fields: fields{
				logger.New(config.Config{LogLevel: -4}),
				mocks.
			},
		},
	}
}
