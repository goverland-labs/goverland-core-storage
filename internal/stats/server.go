package stats

import (
	"context"

	pb "github.com/goverland-labs/goverland-core-storage/protocol/storagepb"
)

type Server struct {
	pb.UnimplementedStatsServer

	service *Service
}

func NewServer(s *Service) *Server {
	return &Server{
		service: s,
	}
}

func (s *Server) GetTotals(_ context.Context, _ *pb.GetTotalsRequest) (*pb.GetTotalsResponse, error) {
	totals := s.service.GetTotals()

	return &pb.GetTotalsResponse{
		Dao: &pb.DaoStats{
			Total:         totals.Dao.Total,
			TotalVerified: totals.Dao.TotalVerified,
		},
		Proposals: &pb.ProposalsStats{
			Total: totals.Proposals.Total,
		},
	}, nil
}
