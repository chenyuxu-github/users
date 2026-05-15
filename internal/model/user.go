package model

import "time"

// User 是用户表的数据模型，同时用于接口 JSON 返回。
type User struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"` // 用户 ID，主键自增。
	Name      string    `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Email     string    `gorm:"column:email;type:varchar(255);not null" json:"email"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}
