package network

import (
	"fmt"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"net"
	"quick_web_golang/config"
	"quick_web_golang/log"
	pb "quick_web_golang/protos"
	"quick_web_golang/service"
	"runtime/debug"
	"time"
)

type Rpc struct {
	Server *grpc.Server
}

func (rpc *Rpc) New() *Rpc {
	handler := Handler
	recoveryOpts := []grpcrecovery.Option{
		grpcrecovery.WithRecoveryHandler(func(p interface{}) (err error) {
			_ = log.Errorf("%v %v", p, string(debug.Stack()))
			return status.Errorf(codes.Unknown, "panic triggered")
		}),
	}
	opts := []grpc.ServerOption{
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				Time:              10 * time.Minute,
				MaxConnectionIdle: 5 * time.Minute,
			},
		),
		grpc.MaxRecvMsgSize(1024 * 1024 * 8),
		grpc.ConnectionTimeout(2 * time.Second),
		grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
			UnaryServerInterceptor(handler),
			grpcrecovery.UnaryServerInterceptor(recoveryOpts...),
		)),
		grpc.StreamInterceptor(grpcmiddleware.ChainStreamServer(
			StreamServerInterceptor(handler),
			grpcrecovery.StreamServerInterceptor(recoveryOpts...),
		)),
	}
	rpc.Server = grpc.NewServer(opts...)
	return rpc
}

func (rpc *Rpc) Start() {
	userServer := service.UserService{}
	pb.RegisterUserServiceServer(rpc.Server, &userServer)
	MethodDescriptor = map[string]*desc.MethodDescriptor{}
	HTTPMethodDescriptor = map[string]*desc.MethodDescriptor{}
	sds, _ := grpcreflect.LoadServiceDescriptors(rpc.Server)
	for _, sd := range sds {
		for _, md := range sd.GetMethods() {
			methodName := fmt.Sprintf("/%s/%s", sd.GetFullyQualifiedName(), md.GetName())
			MethodDescriptor[methodName] = md

			b := proto.GetExtension(md.GetMethodOptions(), annotations.E_Http)
			v, ok := b.(*annotations.HttpRule)
			if !ok {
				continue
			}
			httpURL := ""
			switch v.Pattern.(type) {
			case *annotations.HttpRule_Get:
				httpURL = v.Pattern.(*annotations.HttpRule_Get).Get
				break
			case *annotations.HttpRule_Post:
				httpURL = v.Pattern.(*annotations.HttpRule_Post).Post
				break
			case *annotations.HttpRule_Put:
				httpURL = v.Pattern.(*annotations.HttpRule_Put).Put
				break
			case *annotations.HttpRule_Delete:
				httpURL = v.Pattern.(*annotations.HttpRule_Delete).Delete
				break
			case *annotations.HttpRule_Patch:
				httpURL = v.Pattern.(*annotations.HttpRule_Patch).Patch
				break
			case *annotations.HttpRule_Custom:
			default:
			}
			if len(httpURL) == 0 {
				continue
			}
			HTTPMethodDescriptor[httpURL] = md
		}
	}

	reflection.Register(rpc.Server)

	go func() {
		_ = log.Infof("grpc listening on %s", config.Get(config.GrpcAddress))
		listener, err := net.Listen("tcp", config.Get(config.GrpcAddress))
		if err != nil {
			log.Fatal(err)
		}

		if err := rpc.Server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

}

func (rpc *Rpc) Close() {
	rpc.Server.GracefulStop()
}
