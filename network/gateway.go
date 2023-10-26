package network

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"net/http"
	"quick_web_golang/config"
	"quick_web_golang/log"
	pb "quick_web_golang/protos"
	"quick_web_golang/provider"
	"strings"
	"time"
)

type Gateway struct {
	Server *http.Server
}

func (gateway *Gateway) New() *Gateway {
	gateway.Server = &http.Server{
		Addr:         config.Get(config.GatewayAddress),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Second * 60,
		Handler:      cors(provider.SessionManager.Manager.LoadAndSave(mux.NewRouter())),
	}
	return gateway
}

func (gateway *Gateway) Start() {
	var m runtime.ProtoMarshaller
	// use our hook to modify the response after the gRPC call comes back
	gwmux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: true}),
		runtime.WithMarshalerOption("application/protobuf", &m),
		runtime.WithMetadata(gatewayMetadataAnnotator))
	opts := []grpc.DialOption{grpc.WithInsecure()}

	if err := pb.RegisterUserServiceHandlerFromEndpoint(context.Background(), gwmux, config.Get(config.GrpcAddress), opts); err != nil {
		log.Fatalf("Error register service %v", err)
	}

	go func() {
		_ = log.Infof("gateway listening on %s", config.Get(config.GatewayAddress))
		if err := gateway.Server.ListenAndServe(); err != nil {
			_ = log.Error(err)
		}
	}()
}

func (gateway *Gateway) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_ = gateway.Server.Shutdown(ctx)
}

// look up session and pass userId in to context if it exists
func gatewayMetadataAnnotator(ctx context.Context, r *http.Request) metadata.MD {
	lang := r.Header.Get("lang")
	timestamp := r.Header.Get("timestamp")
	deviceId := r.Header.Get("deviceId")
	md := metadata.Pairs(
		"lang", lang,
		"timestamp", timestamp,
		"deviceId", deviceId,
	)

	uid, ok := provider.SessionManager.Manager.Get(ctx, "uid").(string)
	if !ok {
		return md
	}
	sessionId, _, err := provider.SessionManager.Manager.Commit(r.Context())
	if err != nil {
		return md
	}
	platform, ok := provider.SessionManager.Manager.Get(ctx, "platform").(string)
	if !ok {
		platform = ""
	}

	return metadata.Pairs(
		"uid", uid,
		"platform", platform,
		"lang", lang,
		"timestamp", timestamp,
		"sessionId", sessionId,
		"deviceId", deviceId,
	)
}

func cors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

// preflightHandler adds the necessary headers in order to serve
// CORS from any origin using the methods "GET", "HEAD", "POST", "PUT", "DELETE"
// We insist, don't do this without consideration in production systems.
func preflightHandler(w http.ResponseWriter, r *http.Request) {
	headers := []string{"Content-Type", "Accept", "Authorization", "lang", "Lang"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
}
