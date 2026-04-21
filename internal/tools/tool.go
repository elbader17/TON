package tools

import (
	"context"

	"github.com/ton/framework/pkg/errors"
)

type Tool interface {
	Name() string
	Description() string
	Execute(ctx context.Context, req interface{}) (interface{}, *errors.TONError)
}