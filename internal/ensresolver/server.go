package ensresolver

import (
	"context"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/goverland-labs/goverland-core-storage/protocol/storagepb"
)

type Server struct {
	storagepb.UnimplementedEnsServer

	sp *Service
}

func NewServer(sp *Service) *Server {
	return &Server{
		sp: sp,
	}
}

func (s *Server) GetEnsByAddresses(_ context.Context, req *storagepb.EnsByAddressesRequest) (*storagepb.EnsByAddressesResponse, error) {
	list, err := s.sp.GetByAddresses(req.GetAddresses())
	if err != nil {
		log.Error().Err(err).Msgf("get ens names by addresses: %v", req.GetAddresses())
		return nil, status.Error(codes.Internal, "internal error")
	}

	res := make([]*storagepb.EnsName, len(list))

	for i, info := range list {
		res[i] = convertEnsNameToAPI(&info)
	}

	return &storagepb.EnsByAddressesResponse{
		EnsNames: res,
	}, nil
}

func (s *Server) GetAddressesByEnsNames(_ context.Context, req *storagepb.AddressesByEnsNamesRequest) (*storagepb.AddressesByEnsNamesResponse, error) {
	names := req.GetNames()
	if len(names) == 0 {
		return &storagepb.AddressesByEnsNamesResponse{
			EnsNames: []*storagepb.EnsName{},
		}, nil
	}

	list, err := s.sp.GetByNames(names)
	if err != nil {
		log.Error().Err(err).Msgf("get addresses by ens names: %v", names)
		return nil, status.Error(codes.Internal, "internal error")
	}

	res := make([]*storagepb.EnsName, len(list))
	for i, info := range list {
		res[i] = convertEnsNameToAPI(&info)
	}

	return &storagepb.AddressesByEnsNamesResponse{
		EnsNames: res,
	}, nil
}

func convertEnsNameToAPI(info *EnsName) *storagepb.EnsName {
	return &storagepb.EnsName{
		Address: info.Address,
		Name:    info.Name,
	}
}
