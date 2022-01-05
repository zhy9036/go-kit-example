package transport

import (
	"add-svc/endpoint"
	pb "add-svc/proto_gen"
	"context"
	"errors"

	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"
	"github.com/go-kit/kit/transport"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"
)

type grpcServer struct {
	pb.UnimplementedAddServer
	add grpctransport.Handler
}

func NewGRPCServer(endpoints endpoint.Set, logger log.Logger, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer) pb.AddServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}
	if zipkinTracer != nil {
		options = append(options, zipkin.GRPCServerTrace(zipkinTracer))
	}
	return &grpcServer{
		add: grpctransport.NewServer(
			endpoints.AddEndpoint,
			decodeGRPCAddRequest,
			encodeGRPCAddResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(otTracer, "Sum", logger)))...,
		),
	}
}

func (s *grpcServer) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	_, res, err := s.add.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.(*pb.AddResponse), nil
}

func decodeGRPCAddRequest(_ context.Context, req interface{}) (interface{}, error) {
	request := req.(*pb.AddRequest)
	return endpoint.AddRequest{A: request.A, B: request.B}, nil
}

func encodeGRPCAddRequest(_ context.Context, req interface{}) (interface{}, error) {
	request := req.(endpoint.AddRequest)
	return &pb.AddRequest{A: request.A, B: request.B}, nil
}

func decodeGRPCAddResponse(_ context.Context, res interface{}) (interface{}, error) {
	response := res.(*pb.AddResponse)
	return endpoint.AddResponse{V: response.V, Err: str2err(response.Err)}, nil
}

func encodeGRPCAddResponse(_ context.Context, res interface{}) (interface{}, error) {
	response := res.(endpoint.AddResponse)
	return &pb.AddResponse{V: response.V, Err: err2str(response.Err)}, nil
}

func str2err(s string) error {
	if s == "" {
		return nil
	}
	return errors.New(s)
}

func err2str(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
