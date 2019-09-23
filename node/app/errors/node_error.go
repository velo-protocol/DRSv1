package nerrors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NodeError interface {
	GRPCError() error
	Error() string
}

type ErrInvalidArgument struct {
	Message string
}

func (e ErrInvalidArgument) Error() string {
	return e.Message
}

func (e ErrInvalidArgument) GRPCError() error {
	return status.Errorf(codes.InvalidArgument, e.Message)
}

type ErrNotFound struct {
	Message string
}

func (e ErrNotFound) Error() string {
	return e.Message
}

func (e ErrNotFound) GRPCError() error {
	return status.Errorf(codes.NotFound, e.Message)
}

type ErrAlreadyExists struct {
	Message string
}

func (e ErrAlreadyExists) Error() string {
	return e.Message
}

func (e ErrAlreadyExists) GRPCError() error {
	return status.Errorf(codes.AlreadyExists, e.Message)
}

type ErrPermissionDenied struct {
	Message string
}

func (e ErrPermissionDenied) Error() string {
	return e.Message
}

func (e ErrPermissionDenied) GRPCError() error {
	return status.Errorf(codes.PermissionDenied, e.Message)
}

type ErrUnAuthenticated struct {
	Message string
}

func (e ErrUnAuthenticated) Error() string {
	return e.Message
}

func (e ErrUnAuthenticated) GRPCError() error {
	return status.Errorf(codes.Unauthenticated, e.Message)
}

type ErrInternal struct {
	Message string
}

func (e ErrInternal) Error() string {
	return e.Message
}

func (e ErrInternal) GRPCError() error {
	return status.Errorf(codes.Internal, e.Message)
}
