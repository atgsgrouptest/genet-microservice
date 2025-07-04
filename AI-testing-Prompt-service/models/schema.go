package models

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
)


type User struct {
    CompanyID   primitive.ObjectID `json:"companyId" bson:"companyId"`
    CompanyName string             `json:"companyName" bson:"companyName"`
}

type Request struct {
    RequestID       primitive.ObjectID `json:"requestId,omitempty" bson:"requestId,omitempty"`
    CompanyID       string `json:"companyId" bson:"companyId"` // foreign key reference to User
    RequestMaterial []byte             `json:"requestMaterial,omitempty" bson:"requestMaterial,omitempty"` // store image as binary
    PositiveCases   []string           `json:"positiveCases" bson:"positiveCases"`
    NegativeCases   [][]string         `json:"negativeCases" bson:"negativeCases"`
}