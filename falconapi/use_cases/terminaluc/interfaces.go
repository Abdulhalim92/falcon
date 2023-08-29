package terminaluc

import (
	"context"
	"falconapi/domain/entities"
)

type TerminalDataStores interface {
	GetAll(ctx context.Context) ([]entities.TerminalStatus, error)
	GerRegions(ctx context.Context) ([]entities.TRegion, error)
	GetAllWithRegionsNames(ctx context.Context) ([]entities.TerminalStatus, error)
}
