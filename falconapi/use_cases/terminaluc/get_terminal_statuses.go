package terminaluc

import (
	"context"
	"falconapi/domain/entities"
	"fmt"
	"strings"
	"time"
)

type getTerminalsStatusesUseCase struct {
	datastore TerminalDataStores
}

func NewGetTerminalStatusesUseCase(ds TerminalDataStores) *getTerminalsStatusesUseCase {
	return &getTerminalsStatusesUseCase{
		datastore: ds,
	}
}

func (uc *getTerminalsStatusesUseCase) GetTerminalsStatuses(ctx context.Context) ([]entities.TerminalStatus, error) {
	terminalStatuses, err := uc.datastore.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	for idx := range terminalStatuses {
		terminalStatuses[idx].Status = entities.NonActive

		if terminalStatuses[idx].LastPayment != nil {
			diff := time.Now().Sub(*terminalStatuses[idx].LastPayment)
			_time := strings.Split(diff.String(), ".")[0]
			hours := diff.Hours()
			minutes := diff.Minutes()

			timeCaption := fmt.Sprintf("%ss", _time)

			if hours > 24 {
				timeCaption = fmt.Sprintf("%dd", int(hours/24))
			}

			terminalStatuses[idx].LastPaymentDetail = timeCaption

			if minutes <= 5 {
				terminalStatuses[idx].Status = entities.Active
			}
		}

		if terminalStatuses[idx].LastPing != nil && terminalStatuses[idx].Status == entities.NonActive {
			diff := time.Now().Sub(*terminalStatuses[idx].LastPing)
			minutes := diff.Minutes()

			if minutes <= 15 {
				terminalStatuses[idx].Status = entities.Awaiting
			}
		}
	}

	return terminalStatuses, nil
}
