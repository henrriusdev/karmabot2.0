package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Karma struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID    int64              `json:"userID,omitempty"`
	FirstName string             `json:"firstName,omitempty"`
	Karma     int                `json:"karma,omitempty"`
	LastGived time.Time          `json:"lastGived,omitempty"`
}
