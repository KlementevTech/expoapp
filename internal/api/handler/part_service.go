package handler

import (
	"context"
	"errors"
	"time"

	pb "expo/gen/pb/part/v1"
	"expo/gen/pb/part/v1/partv1connect"
	"expo/internal/model"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ partv1connect.PartServiceHandler = (*partServiceHandler)(nil)

type partServiceHandler struct {
	parts PartRepository
}

func newPartServiceHandler(parts PartRepository) partv1connect.PartServiceHandler {
	return &partServiceHandler{parts: parts}
}

func (p partServiceHandler) CreatePart(
	ctx context.Context,
	c *connect.Request[pb.CreatePartRequest],
) (*connect.Response[pb.CreatePartResponse], error) {
	partID, err := uuid.Parse(c.Msg.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid uuid_v7 format: %v", err)
	}

	partName := c.Msg.GetName()
	if partName == "" {
		return nil, status.Error(codes.InvalidArgument, "part name is required")
	}

	newPart := model.NewPart(partID, partName)

	if err = p.parts.Create(ctx, newPart); err != nil {
		if errors.Is(err, model.ErrAlreadyExists) {
			return nil, status.Errorf(codes.AlreadyExists, "part with id %s already exists", partID.String())
		}
		return nil, status.Errorf(codes.Internal, "failed to create part: %v", err)
	}

	return &connect.Response[pb.CreatePartResponse]{
		Msg: &pb.CreatePartResponse{
			Part: mapPartToPb(newPart),
		},
	}, nil
}

func (p partServiceHandler) GetPart(
	ctx context.Context,
	c *connect.Request[pb.GetPartRequest],
) (*connect.Response[pb.GetPartResponse], error) {
	partID, err := uuid.Parse(c.Msg.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid uuid_v7 format: %v", err)
	}

	part, err := p.parts.Get(ctx, partID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "part with id %s not found", partID.String())
		}
		return nil, status.Errorf(codes.Internal, "failed to fetch part: %v", err)
	}

	return &connect.Response[pb.GetPartResponse]{
		Msg: &pb.GetPartResponse{
			Part: mapPartToPb(part),
		},
	}, nil
}

func (p partServiceHandler) ListParts(
	ctx context.Context,
	_ *connect.Request[pb.ListPartsRequest],
) (*connect.Response[pb.ListPartsResponse], error) {
	parts, err := p.parts.List(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch parts: %v", err)
	}

	return &connect.Response[pb.ListPartsResponse]{
		Msg: &pb.ListPartsResponse{
			Parts: mapPartsToPbs(parts),
		},
	}, nil
}

func (p partServiceHandler) UpdatePart(
	ctx context.Context,
	c *connect.Request[pb.UpdatePartRequest],
) (*connect.Response[pb.UpdatePartResponse], error) {
	partID, err := uuid.Parse(c.Msg.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid uuid_v7 format: %v", err)
	}

	partName := c.Msg.GetName()
	if partName == "" {
		return nil, status.Error(codes.InvalidArgument, "part name is required")
	}

	current, err := p.parts.Get(ctx, partID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "not found")
	}

	current.Update(partName)

	oldVersion := c.Msg.GetOldVersion()

	err = p.parts.Update(ctx, current, oldVersion)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrNotFound):
			return nil, status.Error(codes.NotFound, "part not found")
		case errors.Is(err, model.ErrOptimisticLockConflict):
			return nil, status.Error(codes.Aborted, "optimistic lock failed: version mismatch")
		default:
			return nil, status.Errorf(codes.Internal, "failed to update part: %v", err)
		}
	}

	return &connect.Response[pb.UpdatePartResponse]{
		Msg: &pb.UpdatePartResponse{
			Part: mapPartToPb(current),
		},
	}, nil
}

func (p partServiceHandler) DeletePart(
	ctx context.Context,
	c *connect.Request[pb.DeletePartRequest],
) (*connect.Response[pb.DeletePartResponse], error) {
	partID, err := uuid.Parse(c.Msg.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid uuid_v7 format: %v", err)
	}

	oldVersion := c.Msg.GetOldVersion()
	deletedAt := time.Now().UTC()

	err = p.parts.Delete(ctx, partID, oldVersion, deletedAt)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrNotFound):
			return nil, status.Error(codes.NotFound, "part not found")
		case errors.Is(err, model.ErrOptimisticLockConflict):
			return nil, status.Error(codes.Aborted, "optimistic lock failed: version mismatch")
		default:
			return nil, status.Errorf(codes.Internal, "failed to delete part: %v", err)
		}
	}

	return &connect.Response[pb.DeletePartResponse]{}, nil
}
