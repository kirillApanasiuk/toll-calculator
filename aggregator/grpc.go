package main

import (
	"context"
	"github.com/kirillApanasiuk/toll-calculator/types"
)

type GRPCAggregatorServer struct {
	types.UnimplementedAggregatorServer
	svc Aggregator
}

func NewGRPCAggregatorServer(svc Aggregator) *GRPCAggregatorServer {
	return &GRPCAggregatorServer{
		svc: svc,
	}
}

// transport Layer
// JSON -> types.Distance-> all done(same type)
// GRPC -> types.Aggregate->Request
// business Layer -> business Layer type (main type everyone needs to convert to)
func (s *GRPCAggregatorServer) Aggregate(ctx context.Context, req *types.AggregateRequest) (*types.None, error) {
	distance := types.Distance{
		Value: req.Value,
		OBUID: int(req.ObuId),
		Unix:  req.Unix,
	}
	return &types.None{}, s.svc.AggregateDistance(distance)
}
