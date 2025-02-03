package models

import "time"

type Sim struct {
    ID               uint      `json:"id" gorm:"primaryKey"`
    Name             string    `json:"name"`
    Number           string    `json:"number"`
    LastRechargeDate time.Time `json:"last_recharge_date"`
    RechargeValidity time.Time `json:"recharge_validity"`
    IncomingValidity time.Time `json:"incoming_validity"`
    SimExpiry        time.Time `json:"sim_expiry"`
}