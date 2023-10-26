package network

import (
	"context"
	"encoding/base64"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"quick_web_golang/lib"
	"quick_web_golang/log"
	pb "quick_web_golang/protos"
	"quick_web_golang/provider"
	"strings"
)

var (
	MethodDescriptor     map[string]*desc.MethodDescriptor
	HTTPMethodDescriptor map[string]*desc.MethodDescriptor
)

func encodeKeyValue(k, v string) (string, string) {
	k = strings.ToLower(k)
	if strings.HasSuffix(k, "-bin") {
		val := base64.StdEncoding.EncodeToString([]byte(v))
		v = string(val)
	}
	return k, v
}

// NiceMD is a convenience wrapper definiting extra functions on the metadata.
type NiceMD metadata.MD

// Get retrieves a single value from the metadata.
//
// It works analogously to http.Header.Get, returning the first value if there are many set. If the value is not set,
// an empty string is returned.
//
// The function is binary-key safe.
func (m NiceMD) Get(key string) string {
	k, _ := encodeKeyValue(key, "")
	vv, ok := m[k]
	if !ok {
		return ""
	}
	return vv[0]
}

// ExtractIncoming extracts an inbound metadata from the server-side context.
//
// This function always returns a NiceMD wrapper of the metadata.MD, in case the context doesn't have metadata it returns
// a new empty NiceMD.
func ExtractIncoming(ctx context.Context) NiceMD {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return NiceMD(metadata.Pairs())
	}
	return NiceMD(md)
}

func FromMD(ctx context.Context) (string, error) {
	val := ExtractIncoming(ctx).Get("sid")
	if val == "" {
		return "", status.Errorf(codes.Unauthenticated, "Unauthenticated")
	}

	return val, nil
}

// AuthFunc is the pluggable function that performs authentication.
//
// The passed in `Context` will contain the gRPC metadata.MD object (for header-based authentication) and
// the peer.Peer information that can contain transport-based credentials (e.g. `credentials.AuthInfo`).
//
// The returned context will be propagated to handlers, allowing user changes to `Context`. However,
// please make sure that the `Context` returned is a child `Context` of the one passed in.
//
// If error is returned, its `grpc.Code()` will be returned to the user as well as the verbatim message.
// Please make sure you use `codes.Unauthenticated` (lacking auth) and `codes.PermissionDenied`
// (authed, but lacking perms) appropriately.
type AuthFunc func(ctx context.Context, info string) (context.Context, error)

// ServiceAuthFuncOverride allows a given gRPC service implementation to override the global `AuthFunc`.
//
// If a service implements the AuthFuncOverride method, it takes precedence over the `AuthFunc` method,
// and will be called instead of AuthFunc for all method invocations within that service.
type ServiceAuthFuncOverride interface {
	AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error)
}

// UnaryServerInterceptor returns a new unary server interceptors that performs per-request auth.
func UnaryServerInterceptor(authFunc AuthFunc) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var newCtx context.Context
		var err error
		if overrideSrv, ok := info.Server.(ServiceAuthFuncOverride); ok {
			newCtx, err = overrideSrv.AuthFuncOverride(ctx, info.FullMethod)
		} else {
			newCtx, err = authFunc(ctx, info.FullMethod)
		}
		if err != nil {
			return nil, err
		}
		return handler(newCtx, req)
	}
}

// StreamServerInterceptor returns a new unary server interceptors that performs per-request auth.
func StreamServerInterceptor(authFunc AuthFunc) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		var newCtx context.Context
		var err error
		if overrideSrv, ok := srv.(ServiceAuthFuncOverride); ok {
			newCtx, err = overrideSrv.AuthFuncOverride(stream.Context(), info.FullMethod)
		} else {
			newCtx, err = authFunc(stream.Context(), info.FullMethod)
		}
		if err != nil {
			return err
		}
		wrapped := WrapServerStream(stream)
		wrapped.WrappedContext = newCtx
		return handler(srv, wrapped)
	}
}

// WrappedServerStream is a thin wrapper around grpc.ServerStream that allows modifying context.
type WrappedServerStream struct {
	grpc.ServerStream
	// WrappedContext is the wrapper's own Context. You can assign it.
	WrappedContext context.Context
}

// Context returns the wrapper's WrappedContext, overwriting the nested grpc.ServerStream.Context()
func (w *WrappedServerStream) Context() context.Context {
	return w.WrappedContext
}

// WrapServerStream returns a ServerStream that has the ability to overwrite context.
func WrapServerStream(stream grpc.ServerStream) *WrappedServerStream {
	if existing, ok := stream.(*WrappedServerStream); ok {
		return existing
	}
	return &WrappedServerStream{ServerStream: stream, WrappedContext: stream.Context()}
}

func Handler(ctx context.Context, method string) (context.Context, error) {
	methodDesc, ok := MethodDescriptor[method]
	if !ok {
		return ctx, nil
	}
	if !proto.HasExtension(methodDesc.GetMethodOptions(), pb.E_Auth) {
		return ctx, nil
	}

	sessionId := ExtractIncoming(ctx).Get(lib.SessionId)

	ctx = context.WithValue(ctx, lib.SessionId, sessionId)

	ctx, err := provider.SessionManager.Manager.Load(ctx, sessionId)
	if err != nil {
		_ = log.Error(err)
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return ctx, nil
}
