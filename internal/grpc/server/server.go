package server

import (
	"context"
	"fmt"
	types "github.com/superles/yapmetrics/internal/metric"
	"github.com/superles/yapmetrics/internal/server"
	"github.com/superles/yapmetrics/internal/server/config"
	"github.com/superles/yapmetrics/internal/utils/logger"
	"github.com/superles/yapmetrics/internal/utils/network"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"time"

	pb "github.com/superles/yapmetrics/internal/grpc/proto"
	"google.golang.org/grpc"
)

type metricProvider interface {
	GetAll(ctx context.Context) (map[string]types.Metric, error)
	Get(ctx context.Context, name string) (types.Metric, error)
	Set(ctx context.Context, data types.Metric) error
	SetAll(ctx context.Context, data []types.Metric) error
	SetFloat(ctx context.Context, Name string, Value float64) error
	IncCounter(ctx context.Context, Name string, Value int64) error
	Ping(ctx context.Context) error
	Dump(ctx context.Context, path string) error
	Restore(ctx context.Context, path string) error
}

type GRPCServer struct {
	pb.UnimplementedServerServiceServer
	storage metricProvider
	config  *config.Config
}

func NewGrpcServer(storage metricProvider, cfg *config.Config) server.IServer {
	return &GRPCServer{storage: storage, config: cfg}
}

// Updates реализует метод gRPC Updates из вашего .proto файла
func (s *GRPCServer) Updates(ctx context.Context, req *pb.UpdateMetricsRequest) (*pb.UpdateMetricsResponse, error) {

	if len(s.config.TrustedSubnet) > 0 {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			values := md.Get("X-RealIP")
			if len(values) == 0 {
				return nil, status.Errorf(codes.Unauthenticated, "address not exist in metadata")
			}
			realIp := values[0]
			inNetwork, err := network.IsAddressInNetwork(realIp, s.config.TrustedSubnet)
			if err != nil {
				return nil, err
			}
			if !inNetwork {
				return nil, status.Errorf(codes.Unauthenticated, "address not exist in trusted network")
			}
		}
	}

	var err error

	data := make([]types.Metric, len(req.Metrics))

	idx := 0
	for _, metric := range req.Metrics {
		var metricType int
		switch metric.GetType() {
		case pb.Metric_GAUGE:
			metricType = types.GaugeMetricType
		case pb.Metric_COUNTER:
			metricType = types.CounterMetricType
		default:
			return nil, fmt.Errorf("undefined grpc metric type %d", metric.Type)
		}
		item := types.Metric{
			Name:  metric.GetId(),
			Type:  metricType,
			Value: metric.GetValue(),
		}
		if len(item.Name) == 0 {
			return nil, fmt.Errorf("undefined grpc metric name %s", metric.GetId())
		}
		data[idx] = item
		idx++
	}

	err = s.storage.SetAll(ctx, data)

	if err != nil {
		return nil, err
	}

	// Возвращаем успешный ответ (может быть расширено в зависимости от вашей логики)
	return &pb.UpdateMetricsResponse{Error: ""}, nil
}

func (s *GRPCServer) load(ctx context.Context) error {
	if len(s.config.FileStoragePath) != 0 {
		return s.storage.Restore(ctx, s.config.FileStoragePath)
	}
	return nil
}

func (s *GRPCServer) dump(ctx context.Context) error {
	if len(s.config.FileStoragePath) != 0 {
		return s.storage.Dump(ctx, s.config.FileStoragePath)
	}
	return nil
}

func (s *GRPCServer) startDumpWatcher(ctx context.Context) {
	if s.config.StoreInterval > 0 {
		ticker := time.NewTicker(time.Second * time.Duration(s.config.StoreInterval))
		go func() {
			for t := range ticker.C {
				logger.Log.Debug(fmt.Sprintf("Tick at: %v\n", t.UTC()))
				if err := s.dump(ctx); err != nil {
					logger.Log.Fatal(err.Error())
				}
			}
		}()
	}
}

func (s *GRPCServer) Run(ctx context.Context) error {

	if s.config.Restore && s.config.DatabaseDsn == "" {
		if err := s.load(ctx); err != nil {
			logger.Log.Error(err.Error())
			return err
		} else {
			logger.Log.Debug("бд загружена успешно")
		}
	}

	// Запускаем gRPC сервер на порту 50051
	lis, err := net.Listen("tcp", s.config.GrpcAddress)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Создаем gRPC сервер
	grpcSrv := grpc.NewServer()

	// Регистрируем наш сервер в gRPC
	pb.RegisterServerServiceServer(grpcSrv, s)

	go func() {
		// Запускаем сервер
		if err = grpcSrv.Serve(lis); err != nil {
			logger.Log.Error(fmt.Sprintf("не могу запустить сервер: %s", err))
		}
	}()

	s.startDumpWatcher(ctx)

	logger.Log.Info("Server Started")
	<-ctx.Done()
	logger.Log.Info("Server Stopped")

	grpcSrv.GracefulStop()

	logger.Log.Info("Server Exited Properly")

	if err := s.dump(ctx); err != nil {
		return err
	}
	return nil
}
