package db

import (
	"time"
)

type MemoLike struct {
	Id        int        `gorm:"column:id;primary_key;NOT NULL" json:"id,omitempty"`
	MemoId    int        `gorm:"column:memoId;NOT NULL;index:idx_memo_like_memo_user,unique" json:"memoId,omitempty"`
	UserId    int32      `gorm:"column:userId;NOT NULL;index:idx_memo_like_memo_user,unique" json:"userId,omitempty"`
	CreatedAt *time.Time `gorm:"column:createdAt;default:CURRENT_TIMESTAMP;NOT NULL" json:"createdAt,omitempty"`
}

func (m *MemoLike) TableName() string {
	return "MemoLike"
}
