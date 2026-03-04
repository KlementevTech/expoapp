package handler

import (
	"net/http"

	"expo/gen/pb/part/v1/partv1connect"

	"connectrpc.com/grpcreflect"
)

func NewHandler(parts PartRepository, reflection bool) http.Handler {
	mux := http.NewServeMux()

	path, handler := partv1connect.NewPartServiceHandler(newPartServiceHandler(parts))

	mux.Handle(path, handler)

	if reflection {
		reflector := grpcreflect.NewStaticReflector(
			partv1connect.PartServiceName,
		)

		mux.Handle(grpcreflect.NewHandlerV1(reflector))
		mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))
	}

	return mux
}
