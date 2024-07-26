package models

import "gorm.io/gorm"

type IntentsUsers struct {
	gorm.Model
	SellerID string `json:"seller_id" gorm:"not null"`
	BuyerID  string `json:"buyer_id" gorm:"not null"`
}

type CarIntent struct {
	gorm.Model
	PayoutIDs    string        `json:"payout" gorm:"not null"`
	IntentUsers  *IntentsUsers `json:"intent_users" gorm:"foreignKey:IntentUsersID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	IntentUsersID uint          `json:"intent_users_id"`
	Status       string        `json:"status" gorm:"not null"`
	Currency     string        `json:"currency" gorm:"not null"`
	Amount       int64         `json:"amount" gorm:"not null"`
}
