package handler

import (
	"time"

	"github.com/kingwrcy/moments/db"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
)

type MemoViewHandler struct {
	base BaseHandler
}

func NewMemoViewHandler(injector do.Injector) *MemoViewHandler {
	return &MemoViewHandler{
		base: do.MustInvoke[BaseHandler](injector),
	}
}

// AddViewReq 标记已读请求
type AddViewReq struct {
	MemoId int `json:"memoId"`
}

// AddView 标记已读（UPSERT）
func (h *MemoViewHandler) AddView(c echo.Context) error {
	ctx := c.(CustomContext)
	currentUser := ctx.CurrentUser()

	var req AddViewReq
	if err := c.Bind(&req); err != nil {
		return FailResp(c, ParamError)
	}

	if req.MemoId <= 0 {
		return FailResp(c, ParamError)
	}

	now := time.Now()
	view := db.MemoView{
		MemoId:   req.MemoId,
		UserId:   int(currentUser.Id),
		ViewedAt: &now,
	}

	// UPSERT: 重复查看只更新时间
	err := h.base.db.Exec(`
		INSERT INTO memo_view (memo_id, user_id, viewed_at)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE viewed_at = VALUES(viewed_at)
	`, view.MemoId, view.UserId, view.ViewedAt).Error

	if err != nil {
		h.base.log.Error().Msgf("标记已读失败: %v", err)
		return FailRespWithMsg(c, Fail, "标记已读失败")
	}

	return SuccessResp(c, struct{}{})
}

// GetViewersReq 获取已读用户列表请求
type GetViewersReq struct {
	MemoId int `json:"memoId"`
	Page   int `json:"page"`
	Size   int `json:"size"`
}

// ViewerInfo 已读用户信息
type ViewerInfo struct {
	UserId    int32     `json:"userId"`
	Username  string    `json:"username"`
	Nickname  string    `json:"nickname"`
	AvatarUrl string    `json:"avatarUrl"`
	ViewedAt  time.Time `json:"viewedAt"`
}

// GetViewersResp 已读用户列表响应
type GetViewersResp struct {
	List    []ViewerInfo `json:"list"`
	Total   int64        `json:"total"`
	HasNext bool         `json:"hasNext"`
}

// GetViewers 获取已读用户列表
func (h *MemoViewHandler) GetViewers(c echo.Context) error {
	var req GetViewersReq
	if err := c.Bind(&req); err != nil {
		return FailResp(c, ParamError)
	}

	if req.MemoId <= 0 {
		return FailResp(c, ParamError)
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 20
	}

	offset := (req.Page - 1) * req.Size

	var total int64
	h.base.db.Model(&db.MemoView{}).Where("memo_id = ?", req.MemoId).Count(&total)

	var viewers []ViewerInfo
	h.base.db.Table("memo_view").
		Select("memo_view.user_id, u.username, u.nickname, u.avatar_url as avatarUrl, memo_view.viewed_at").
		Joins("LEFT JOIN user u ON memo_view.user_id = u.id").
		Where("memo_view.memo_id = ?", req.MemoId).
		Order("memo_view.viewed_at DESC").
		Limit(req.Size).Offset(offset).
		Scan(&viewers)

	if viewers == nil {
		viewers = []ViewerInfo{}
	}

	return SuccessResp(c, GetViewersResp{
		List:    viewers,
		Total:   total,
		HasNext: int64(req.Page*req.Size) < total,
	})
}

// GetMyViewsReq 获取我的已读列表请求
type GetMyViewsReq struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

// MyViewInfo 我的已读记录
type MyViewInfo struct {
	MemoId   int       `json:"memoId"`
	ViewedAt time.Time `json:"viewedAt"`
}

// GetMyViewsResp 我的已读列表响应
type GetMyViewsResp struct {
	List    []MyViewInfo `json:"list"`
	Total   int64        `json:"total"`
	HasNext bool         `json:"hasNext"`
}

// GetMyViews 获取我的已读列表
func (h *MemoViewHandler) GetMyViews(c echo.Context) error {
	ctx := c.(CustomContext)
	currentUser := ctx.CurrentUser()

	var req GetMyViewsReq
	if err := c.Bind(&req); err != nil {
		return FailResp(c, ParamError)
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 20
	}

	offset := (req.Page - 1) * req.Size

	var total int64
	h.base.db.Model(&db.MemoView{}).Where("user_id = ?", currentUser.Id).Count(&total)

	var views []MyViewInfo
	h.base.db.Table("memo_view").
		Select("memo_id, viewed_at").
		Where("user_id = ?", currentUser.Id).
		Order("viewed_at DESC").
		Limit(req.Size).Offset(offset).
		Scan(&views)

	if views == nil {
		views = []MyViewInfo{}
	}

	return SuccessResp(c, GetMyViewsResp{
		List:    views,
		Total:   total,
		HasNext: int64(req.Page*req.Size) < total,
	})
}
