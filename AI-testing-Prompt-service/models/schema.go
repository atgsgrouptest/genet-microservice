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

type CompletedRequest struct {
	RequestId         primitive.ObjectID `json:"requestId,omitempty" bson:"requestId,omitempty"`
	RequestMaterial   []byte             `json:"requestMaterial,omitempty" bson:"requestMaterial,omitempty"`
	CompanyID         string             `json:"companyId" bson:"companyId"`
	PositiveResponses []PositiveResponse `json:"positiveResponses,omitempty" bson:"positiveResponses,omitempty"`
	NegativeResponses []NegativeResponse `json:"negativeResponses,omitempty" bson:"negativeResponses,omitempty"`
}

type PositiveResponse struct {
	RequestId        primitive.ObjectID `json:"requestId" bson:"requestId"`
	RequestNo        int                `json:"requestNo" bson:"requestNo"`
	PositiveCase     string             `json:"positiveCase" bson:"positiveCase"`
	PositiveVideoUrl string             `json:"positiveVideoUrl" bson:"positiveVideoUrl"`
	ExtractedContent string             `json:"extractedcontent" bson:"extractedcontent"`
	Success          bool               `json:"success" bson:"success"`
}

type NegativeResponse struct {
	RequestId        primitive.ObjectID `json:"requestId" bson:"requestId"`
	RequestNo        int                `json:"requestNo" bson:"requestNo"`
	NegativeCase     string             `json:"negativeCase" bson:"negativeCase"`
	NegativeVideoUrl string             `json:"negativeVideoUrl" bson:"negativeVideoUrl"`
	ExtractedContent string             `json:"extractedcontent" bson:"extractedcontent"`
	Success          bool               `json:"success" bson:"success"`
}
