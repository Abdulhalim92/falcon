package entities

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	Id        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	Name      string
	Price     float32
}

type ErrorModel struct {
	Err     error
	Message string
	Code    int
}

type TerminalStatus struct {
	EndpointID        int64      `gorm:"column:endpoint_id"`
	EndpointNum       int64      `gorm:"column:endpoint_num"`
	Phone             string     `gorm:"column:phone"`
	Address           string     `gorm:"column:address"`
	LastPayment       *time.Time `gorm:"column:last_created_payment"`
	LastPing          *time.Time `gorm:"column:lastping"`
	Region_id         int64      `gorm:"column:region_id"`
	RegionName        string     `gorm:"column:"`
	Status            string     `gorm:"-"`
	LastPaymentDetail string     `gorm:"-"`
	DetailStatus      string     `gorm:"column:status"`
	EndpointDisabled  bool       `gorm:"column:endpoint_disabled"`
}

type TRegion struct {
	ID   int64  `gorm:"primary_key" json:"id"`
	Name string `gorm:"column:name" json:"name"`
}

func (tr *TRegion) TableName() string {
	return "tregion"
}
