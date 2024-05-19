package models

import "time"

type Training struct {
	ID           string    `json:"id"`
	Name         string    `json:"name" validation:"required"`
	Commentary   string    `json:"commentary"`
	Training_Day string    `json:"training_day" validation:"required"`
	Duration     int       `json:"duration" validation:"required"`
	Time         time.Time `json:"time" validation:"required"`
}
