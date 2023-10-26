package service

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	InternalError   = status.Error(codes.Internal, "internal error")
	InvalidArgument = status.Error(codes.InvalidArgument, "invalid argument")
	NotFound        = status.Error(codes.NotFound, "not found")
	Unauthenticated = status.Errorf(codes.Unauthenticated, "Unauthenticated")
)
