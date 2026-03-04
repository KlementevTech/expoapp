package grpc

import (
	"time"

	pb "expo/gen/pb/part/v1"
	"expo/internal/model"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func mapPartToPb(p *model.Part) *pb.Part {
	if p == nil {
		return nil
	}

	return &pb.Part{
		Id:        p.ID.String(),
		Name:      p.Name,
		Version:   int64(p.Version),
		CreatedAt: timestamppb.New(p.CreatedAt),
		UpdatedAt: mapTimePtrToPb(p.UpdatedAt),
		DeletedAt: mapTimePtrToPb(p.DeletedAt),
	}
}

func mapPartsToPbs(parts []*model.Part) []*pb.Part {
	if len(parts) == 0 {
		return nil
	}

	pbParts := make([]*pb.Part, 0, len(parts))
	for _, p := range parts {
		pbParts = append(pbParts, mapPartToPb(p))
	}
	return pbParts
}

func mapTimePtrToPb(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}
