package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Brand struct {
	Id        bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string        `bson:"brand_name" json:"brand_name" binding:"required"`
	Devices   []Device      `bson:"devices,omitempty" json:"devices,omitempty"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
}
