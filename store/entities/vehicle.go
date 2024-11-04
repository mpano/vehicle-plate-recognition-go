package entities

import (
	"gorm.io/gorm"
	"time"

	"github.com/google/uuid"
)

type Vehicle struct {
	Id        string    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	Plate     string    `gorm:"type:varchar(255);not null;column:plate" json:"plate"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}

func (*Vehicle) TableName() string {
	return "vehicle"
}

func (vehicle *Vehicle) BeforeCreate(tx *gorm.DB) error {
	if vehicle.Id == "" {
		id := uuid.New()
		vehicle.Id = id.String()
	}
	return nil
}
