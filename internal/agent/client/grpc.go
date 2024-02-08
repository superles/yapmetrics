package client

import (
	"context"
	"fmt"
	pb "github.com/superles/yapmetrics/internal/grpc/proto"
	"github.com/superles/yapmetrics/internal/metric"
	"github.com/superles/yapmetrics/internal/utils/encoder"
	"github.com/superles/yapmetrics/internal/utils/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type GrpcClientParams struct {
	Key     string
	RealIP  string
	Encoder *encoder.Encoder
}

type GrpcClient struct {
	params GrpcClientParams
}

func NewGrpcClient(params GrpcClientParams) Client {
	return &GrpcClient{params}
}

func (c GrpcClient) Send(ctx context.Context, endpoint string, metrics []metric.Metric) error {
	// Замените "localhost:50051" на адрес вашего gRPC сервера
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	defer func(conn *grpc.ClientConn) {
		if err = conn.Close(); err != nil {
			logger.Log.Error(err)
		}
	}(conn)

	// Создаем клиент
	serviceClient := pb.NewServerServiceClient(conn)

	data := make([]*pb.Metric, len(metrics))

	for i := 0; i < len(metrics); i++ {
		item := metrics[i]
		valueType := pb.Metric_GAUGE
		if item.Type == metric.CounterMetricType {
			valueType = pb.Metric_COUNTER
		}
		data[i] = &pb.Metric{
			Id:    item.Name,
			Type:  valueType,
			Value: item.Value,
		}
	}

	// добавляем метадаты
	if len(c.params.RealIP) != 0 {
		ctx = metadata.AppendToOutgoingContext(ctx, "X-RealIP", c.params.RealIP)
	}

	// Вызываем gRPC метод с передачей JSON данных и метаданных
	response, err := serviceClient.Updates(ctx, &pb.UpdateMetricsRequest{Metrics: data})

	if err != nil {
		return fmt.Errorf("failed to call gRPC method: %w", err)
	}

	if len(response.GetError()) > 0 {
		return fmt.Errorf("server response grpc error: %s", response.GetError())
	}

	return nil
}
