package grpc

import (
	"context"
	"errors"
	"time"

	pb "expo/gen/pb/part/v1"
	"expo/internal/model"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ pb.PartServiceServer = (*partService)(nil)

type partService struct {
	pb.UnimplementedPartServiceServer

	parts PartRepository
}

func NewExpoServiceServer(
	repo PartRepository,
) pb.PartServiceServer {
	return &partService{
		parts: repo,
	}
}

func (s *partService) GetPart(ctx context.Context, request *pb.GetPartRequest) (*pb.GetPartResponse, error) {
	// 1. Валидация входного ID
	// Мы проверяем, является ли строка валидным UUID, прежде чем идти в базу
	id, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid uuid format: %v", err)
	}

	// 2. Вызов репозитория
	// Используем наш тип PartID (слайс байт)
	part, err := s.parts.Get(ctx, id)
	if err != nil {
		// Проверяем, это ошибка "не найдено" или системный сбой
		if errors.Is(err, model.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "part with id %s not found", request.GetId())
		}
		return nil, status.Errorf(codes.Internal, "failed to fetch part: %v", err)
	}

	// 3. Маппинг доменной модели в ответ Protobuf
	return &pb.GetPartResponse{
		Part: mapPartToPb(part),
	}, nil
}

func (s *partService) CreatePart(ctx context.Context, request *pb.CreatePartRequest) (*pb.CreatePartResponse, error) {
	partID, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid uuid_v7 format: %v", err)
	}
	if request.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "part name is required")
	}

	partName := request.GetName()
	newPart := model.NewPart(partID, partName)

	if err = s.parts.Create(ctx, newPart); err != nil {
		if errors.Is(err, model.ErrAlreadyExists) {
			return nil, status.Errorf(codes.AlreadyExists, "part with id %s already exists", request.GetId())
		}
		return nil, status.Errorf(codes.Internal, "failed to create part: %v", err)
	}

	return &pb.CreatePartResponse{
		Part: mapPartToPb(newPart),
	}, nil
}

func (s *partService) ListParts(ctx context.Context, _ *pb.ListPartsRequest) (*pb.ListPartsResponse, error) {
	parts, err := s.parts.List(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch parts: %v", err)
	}

	return &pb.ListPartsResponse{
		Parts: mapPartsToPbs(parts),
		// NextPageToken можно будет реализовать позже при добавлении пагинации в SQL
	}, nil
}

func (s *partService) UpdatePart(ctx context.Context, request *pb.UpdatePartRequest) (*pb.UpdatePartResponse, error) {
	id, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid uuid_v7 format: %v", err)
	}

	if request.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "part name is required")
	}

	current, err := s.parts.Get(ctx, id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "not found")
	}

	name := request.GetName()
	current.Update(name)
	oldVersion := request.GetOldVersion()

	err = s.parts.Update(ctx, current, oldVersion)
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

	return &pb.UpdatePartResponse{Part: mapPartToPb(current)}, nil
}

func (s *partService) DeletePart(ctx context.Context, request *pb.DeletePartRequest) (*pb.DeletePartResponse, error) {
	partID, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid uuid format: %v", err)
	}

	oldVersion := request.GetOldVersion()
	deletedAt := time.Now().UTC()

	err = s.parts.Delete(ctx, partID, oldVersion, deletedAt)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrNotFound):
			return nil, status.Error(codes.NotFound, "part not found")
		case errors.Is(err, model.ErrOptimisticLockConflict):
			return nil, status.Error(codes.Aborted, "optimistic lock failed: oldVersion mismatch")
		default:
			return nil, status.Errorf(codes.Internal, "failed to delete part: %v", err)
		}
	}

	return &pb.DeletePartResponse{}, nil
}
