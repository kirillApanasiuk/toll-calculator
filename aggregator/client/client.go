package client

import (
	"context"
	"github.com/kirillApanasiuk/toll-calculator/types"
)

type Client interface {
	Aggregate(ctx context.Context, request *types.AggregateRequest) error
}
