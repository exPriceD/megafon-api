package ports

import (
	"context"
	"megafon-buisness-reports/internal/domain/entities"
)

type CallRepository interface {
	History(ctx context.Context, f entities.CallFilter) ([]entities.Call, error)
}
