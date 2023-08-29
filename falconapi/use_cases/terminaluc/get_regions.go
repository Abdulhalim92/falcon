package terminaluc

import (
	"context"
	"falconapi/domain/entities"
)

type getRegionsUseCase struct {
	dataStore TerminalDataStores
}

func NewGetRegionsUseCase(ds TerminalDataStores) *getRegionsUseCase {
	return &getRegionsUseCase{
		dataStore: ds,
	}
}

func (uc *getRegionsUseCase) GetRegions(ctx context.Context) ([]entities.TRegion, error) {
	return uc.dataStore.GerRegions(ctx)
}
