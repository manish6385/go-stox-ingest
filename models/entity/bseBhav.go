package models

import "gorm.io/datatypes"

type BSE_BHAV struct {
	SC_CODE    uint32         `gorm:"primaryKey:sc_code" json:"sc_code"`
	SC_NAME    string         `gorm:"sc_name" json:"sc_name"`
	SC_GROUP   string         `gorm:"sc_group" json:"sc_group"`
	SC_TYPE    string         `gorm:"sc_type" json:"sc_type"`
	OPEN       float64        `gorm:"open" json:"open"`
	HIGH       float64        `gorm:"high" json:"high"`
	LOW        float64        `gorm:"low" json:"low"`
	CLOSE      float64        `gorm:"close" json:"close"`
	LAST       float64        `gorm:"last" json:"last"`
	PREVCLOSE  float64        `gorm:"prev_close" json:"prev_close"`
	NO_TRADES  uint64         `gorm:"no_trades" json:"no_trades"`
	NET_TURNOV uint64         `gorm:"net_turnov" json:"net_turnov"`
	TDCLOINDI  string         `gorm:"tdcloindi" json:"tdcloindi"`
	CreatedAt  datatypes.Date `gorm:"primaryKey:created_at;autoCreateTime" json:"created_at"`
}
