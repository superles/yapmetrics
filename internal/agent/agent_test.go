package agent

import (
	"context"
	"github.com/superles/yapmetrics/internal/agent/client"
	"github.com/superles/yapmetrics/internal/agent/config"
	types "github.com/superles/yapmetrics/internal/metric"
	"go.uber.org/zap"
	"reflect"
	"testing"
	"time"
)

func TestAgent_Run(t *testing.T) {
	type fields struct {
		storage metricProvider
		config  *config.Config
		client  client.Client
		logger  *zap.SugaredLogger
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Agent{
				storage: tt.fields.storage,
				config:  tt.fields.config,
				client:  tt.fields.client,
				logger:  tt.fields.logger,
			}
			if err := a.Run(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAgent_poolTickPsutil(t *testing.T) {
	type fields struct {
		storage metricProvider
		config  *config.Config
		client  client.Client
		logger  *zap.SugaredLogger
	}
	type args struct {
		ctx          context.Context
		pollInterval time.Duration
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
			a := &Agent{
				storage: tt.fields.storage,
				config:  tt.fields.config,
				client:  tt.fields.client,
				logger:  tt.fields.logger,
			}
			a.poolTickPsutil(tt.args.ctx, tt.args.pollInterval)
		})
	}
}

func TestAgent_poolTickRuntime(t *testing.T) {
	type fields struct {
		storage metricProvider
		config  *config.Config
		client  client.Client
		logger  *zap.SugaredLogger
	}
	type args struct {
		ctx          context.Context
		pollInterval time.Duration
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
			a := &Agent{
				storage: tt.fields.storage,
				config:  tt.fields.config,
				client:  tt.fields.client,
				logger:  tt.fields.logger,
			}
			a.poolTickRuntime(tt.args.ctx, tt.args.pollInterval)
		})
	}
}

func TestAgent_send(t *testing.T) {
	type fields struct {
		storage metricProvider
		config  *config.Config
		client  client.Client
		logger  *zap.SugaredLogger
	}
	type args struct {
		url         string
		contentType string
		body        []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Agent{
				storage: tt.fields.storage,
				config:  tt.fields.config,
				client:  tt.fields.client,
				logger:  tt.fields.logger,
			}
			if err := a.send(tt.args.url, tt.args.contentType, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAgent_sendAll(t *testing.T) {
	type fields struct {
		storage metricProvider
		config  *config.Config
		client  client.Client
		logger  *zap.SugaredLogger
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Agent{
				storage: tt.fields.storage,
				config:  tt.fields.config,
				client:  tt.fields.client,
				logger:  tt.fields.logger,
			}
			if err := a.sendAll(); (err != nil) != tt.wantErr {
				t.Errorf("sendAll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAgent_sendAllJSON(t *testing.T) {
	type fields struct {
		storage metricProvider
		config  *config.Config
		client  client.Client
		logger  *zap.SugaredLogger
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Agent{
				storage: tt.fields.storage,
				config:  tt.fields.config,
				client:  tt.fields.client,
				logger:  tt.fields.logger,
			}
			if err := a.sendAllJSON(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("sendAllJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAgent_sendJSON(t *testing.T) {
	type fields struct {
		storage metricProvider
		config  *config.Config
		client  client.Client
		logger  *zap.SugaredLogger
	}
	type args struct {
		data *types.Metric
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Agent{
				storage: tt.fields.storage,
				config:  tt.fields.config,
				client:  tt.fields.client,
				logger:  tt.fields.logger,
			}
			if err := a.sendJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("sendJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAgent_sendPlain(t *testing.T) {
	type fields struct {
		storage metricProvider
		config  *config.Config
		client  client.Client
		logger  *zap.SugaredLogger
	}
	type args struct {
		data *types.Metric
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Agent{
				storage: tt.fields.storage,
				config:  tt.fields.config,
				client:  tt.fields.client,
				logger:  tt.fields.logger,
			}
			if err := a.sendPlain(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("sendPlain() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAgent_sendTicker(t *testing.T) {
	type fields struct {
		storage metricProvider
		config  *config.Config
		client  client.Client
		logger  *zap.SugaredLogger
	}
	type args struct {
		ctx            context.Context
		reportInterval time.Duration
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
			a := &Agent{
				storage: tt.fields.storage,
				config:  tt.fields.config,
				client:  tt.fields.client,
				logger:  tt.fields.logger,
			}
			a.sendTicker(tt.args.ctx, tt.args.reportInterval)
		})
	}
}

func TestAgent_sendWithRetry(t *testing.T) {
	type fields struct {
		storage metricProvider
		config  *config.Config
		client  client.Client
		logger  *zap.SugaredLogger
	}
	type args struct {
		url         string
		contentType string
		body        []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Agent{
				storage: tt.fields.storage,
				config:  tt.fields.config,
				client:  tt.fields.client,
				logger:  tt.fields.logger,
			}
			if err := a.sendWithRetry(tt.args.url, tt.args.contentType, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("sendWithRetry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		s   metricProvider
		cfg *config.Config
	}
	tests := []struct {
		name string
		args args
		want *Agent
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.s, tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
