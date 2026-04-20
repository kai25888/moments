package db

import (
	"time"
)

// MemoView 动态已读记录
type MemoView struct {
	Id       int        `gorm:"column:id;primary_key;NOT NULL" json:"id,omitempty"`
	MemoId   int        `gorm:"column:memo_id;NOT NULL" json:"memoId,omitempty"`
	UserId   int        `gorm:"column:user_id;NOT NULL" json:"userId,omitempty"`
	ViewedAt *time.Time `gorm:"column:viewed_at;default:CURRENT_TIMESTAMP;NOT NULL" json:"viewedAt,omitempty"`
}

func (m *MemoView) TableName() string {
	return "memo_view"
}
