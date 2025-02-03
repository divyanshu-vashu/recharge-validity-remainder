package models

import "time"

type Sim struct {
    ID               uint      `json:"id" gorm:"primaryKey"`
    Name             string    `json:"name" gorm:"not null"`
    Number           string    `json:"number" gorm:"not null"`
    LastRechargeDate time.Time `json:"lastRechargeDate" gorm:"not null"`
    RechargeValidity time.Time `json:"rechargeValidity" gorm:"not null"`
    IncomingValidity time.Time `json:"incomingValidity" gorm:"not null"`
    SimExpiry        time.Time `json:"simExpiry" gorm:"not null"`
}