package handlers

import (
	"context"
	"falconapi/domain/entities"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type GetTerminalsStatusesUseCase interface {
	GetTerminalsStatuses(ctx context.Context) ([]entities.TerminalStatus, error)
}

type GetRegionsUseCase interface {
	GetRegions(ctx context.Context) ([]entities.TRegion, error)
}

type GetTerminalInfoUseCase interface {
	GetTerminalsInfo(ctx context.Context) ([]entities.TerminalStatus, error)
}

// @Summary Метод получения статусов терминалов
// @Description Получение статусов терминалов
// @Tags Terminals
// @Accept json
// @Produce json
// @OperationId getTerminalsStatuses
// @Success 200 {object} map[string][]entities.TerminalStatus
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /terminalstatuses [get]
func GetTerminalsStatusesHandler(useCase GetTerminalsStatusesUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		terminalsStatuses, err := useCase.GetTerminalsStatuses(ctx)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, "something went wrong")
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": terminalsStatuses})
	}
}

// @Summary Метод получения инфо - статусов терминалов
// @Description Получение инфо - статусов терминалов
// @Tags Terminals
// @Accept json
// @Produce json
// @OperationId getTerminalsInfo
// @Success 200 {object} map[string][]entities.TerminalStatus
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /terminalsinfo [get]
func GetTerminalsInfoHandler(useCase GetTerminalInfoUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		terminalsStatuses, err := useCase.GetTerminalsInfo(ctx)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, "something went wrong")
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": terminalsStatuses})
	}
}

// @Summary Метод получения списка регионов
// @Description Получение списка регионов
// @Tags Terminals
// @Accept json
// @Produce json
// @OperationId getRegions
// @Success 200 {object} map[string][]entities.TRegion
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /region [get]
func GetRegionsHandler(useCase GetRegionsUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		regions, err := useCase.GetRegions(ctx)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, "something went wrong")
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": regions})
	}
}
