package model

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Brand struct {
	Id      bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name    string        `bson:"brand_name" json:"brand_name" binding:"required"`
	Devices []Device      `bson:"devices,omitempty" json:"devices,omitempty"`
}
