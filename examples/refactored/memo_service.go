// Package service 提供业务逻辑层实现
// 这是重构后的 memo 服务层示例，展示了分层架构和依赖注入的最佳实践
package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kingwrcy/moments/db"
	"github.com/kingwrcy/moments/vo"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// 定义领域错误
var (
	ErrMemoNotFound     = errors.New("memo not found")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrInvalidContent   = errors.New("invalid content")
	ErrContentTooLong   = errors.New("content too long")
	MaxContentLength    = 10000
	MaxSearchLength     = 100
)

// MemoRepository 定义 memo 数据访问接口
type MemoRepository interface {
	GetByID(ctx context.Context, id int64) (*db.Memo, error)
	List(ctx context.Context, filter ListMemoFilter) ([]*db.Memo, int64, error)
	Create(ctx context.Context, memo *db.Memo) error
	Update(ctx context.Context, memo *db.Memo) error
	Delete(ctx context.Context, id int64) error
	Like(ctx context.Context, id int64) error
	SetPinned(ctx context.Context, id int64, pinned bool) error
}

// ListMemoFilter 定义列表查询过滤器
type ListMemoFilter struct {
	Page            int
	Size            int
	UserID          *int64
	Username        string
	Tag             string
	ContentContains string
	StartTime       *time.Time
	EndTime         *time.Time
	ShowType        *int
	CurrentUserID   *int64
}

// memoRepository 实现 MemoRepository 接口
type memoRepository struct {
	db  *gorm.DB
	log zerolog.Logger
}

// NewMemoRepository 创建 memo 仓库实例
func NewMemoRepository(db *gorm.DB, log zerolog.Logger) MemoRepository {
	return &memoRepository{
		db:  db,
		log: log.With().Str("component", "memo_repository").Logger(),
	}
}

func (r *memoRepository) GetByID(ctx context.Context, id int64) (*db.Memo, error) {
	var memo db.Memo
	
	err := r.db.WithContext(ctx).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("username", "nickname", "slogan", "id", "avatar_url", "cover_url")
		}).
		First(&memo, id).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMemoNotFound
		}
		r.log.Error().Err(err).Int64("memo_id", id).Msg("failed to get memo by id")
		return nil, fmt.Errorf("database error: %w", err)
	}
	
	return &memo, nil
}

func (r *memoRepository) List(ctx context.Context, filter ListMemoFilter) ([]*db.Memo, int64, error) {
	var (
		memos []*db.Memo
		total int64
	)
	
	// 构建基础查询
	tx := r.db.WithContext(ctx).Model(&db.Memo{}).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("username", "nickname", "slogan", "id", "avatar_url", "cover_url")
		})
	
	// 应用过滤条件
	tx = r.applyFilters(tx, filter)
	
	// 查询总数
	if err := tx.Count(&total).Error; err != nil {
		r.log.Error().Err(err).Msg("failed to count memos")
		return nil, 0, fmt.Errorf("database error: %w", err)
	}
	
	// 分页查询
	offset := (filter.Page - 1) * filter.Size
	if err := tx.Order("pinned DESC, created_at DESC").
		Limit(filter.Size).
		Offset(offset).
		Find(&memos).Error; err != nil {
		r.log.Error().Err(err).Msg("failed to list memos")
		return nil, 0, fmt.Errorf("database error: %w", err)
	}
	
	return memos, total, nil
}

func (r *memoRepository) applyFilters(tx *gorm.DB, filter ListMemoFilter) *gorm.DB {
	// 时间范围过滤
	if filter.StartTime != nil {
		tx = tx.Where("created_at >= ?", *filter.StartTime)
	}
	if filter.EndTime != nil {
		tx = tx.Where("created_at <= ?", *filter.EndTime)
	}
	
	// 内容搜索 - 使用参数化查询防止 SQL 注入
	if filter.ContentContains != "" {
		// 限制搜索长度
		if len(filter.ContentContains) > MaxSearchLength {
			filter.ContentContains = filter.ContentContains[:MaxSearchLength]
		}
		tx = tx.Where("content LIKE ?", "%"+filter.ContentContains+"%")
	}
	
	// 显示类型过滤
	if filter.ShowType != nil && *filter.ShowType >= 0 {
		tx = tx.Where("show_type = ?", *filter.ShowType)
	}
	
	// 权限过滤
	if filter.CurrentUserID == nil {
		// 未登录用户只能看到公开内容
		tx = tx.Where("show_type = 1").
			Where("created_at <= ?", time.Now())
	} else {
		// 登录用户可以看到自己的所有内容和别人的公开内容
		tx = tx.Where("user_id = ? OR (user_id <> ? AND show_type = 1)", 
			*filter.CurrentUserID, *filter.CurrentUserID).
			Where("user_id = ? OR created_at <= ?", 
				*filter.CurrentUserID, time.Now())
	}
	
	// 标签过滤
	if filter.Tag != "" {
		tx = tx.Where("tags LIKE ?", "%"+filter.Tag+",%")
	}
	
	// 用户过滤
	if filter.UserID != nil {
		tx = tx.Where("user_id = ?", *filter.UserID)
	}
	
	return tx
}

func (r *memoRepository) Create(ctx context.Context, memo *db.Memo) error {
	if err := r.db.WithContext(ctx).Create(memo).Error; err != nil {
		r.log.Error().Err(err).Interface("memo", memo).Msg("failed to create memo")
		return fmt.Errorf("database error: %w", err)
	}
	return nil
}

func (r *memoRepository) Update(ctx context.Context, memo *db.Memo) error {
	if err := r.db.WithContext(ctx).Save(memo).Error; err != nil {
		r.log.Error().Err(err).Interface("memo", memo).Msg("failed to update memo")
		return fmt.Errorf("database error: %w", err)
	}
	return nil
}

func (r *memoRepository) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&db.Memo{}, id)
	if result.Error != nil {
		r.log.Error().Err(result.Error).Int64("memo_id", id).Msg("failed to delete memo")
		return fmt.Errorf("database error: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrMemoNotFound
	}
	return nil
}

func (r *memoRepository) Like(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Model(&db.Memo{}).
		Where("id = ?", id).
		UpdateColumn("fav_count", gorm.Expr("fav_count + 1"))
	
	if result.Error != nil {
		r.log.Error().Err(result.Error).Int64("memo_id", id).Msg("failed to like memo")
		return fmt.Errorf("database error: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrMemoNotFound
	}
	return nil
}

func (r *memoRepository) SetPinned(ctx context.Context, id int64, pinned bool) error {
	result := r.db.WithContext(ctx).Model(&db.Memo{}).
		Where("id = ?", id).
		Update("pinned", pinned)
	
	if result.Error != nil {
		r.log.Error().Err(result.Error).Int64("memo_id", id).Bool("pinned", pinned).Msg("failed to set pinned")
		return fmt.Errorf("database error: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrMemoNotFound
	}
	return nil
}

// MemoService 定义 memo 业务逻辑接口
type MemoService interface {
	GetByID(ctx context.Context, id int64, currentUserID *int64) (*db.Memo, error)
	List(ctx context.Context, filter ListMemoFilter) ([]*db.Memo, int64, bool, error)
	Create(ctx context.Context, req *vo.SaveMemoReq, userID int64) (*db.Memo, error)
	Update(ctx context.Context, req *vo.SaveMemoReq, userID int64) (*db.Memo, error)
	Delete(ctx context.Context, id int64, userID int64, isAdmin bool) error
	Like(ctx context.Context, id int64) error
	SetPinned(ctx context.Context, id int64, userID int64, isAdmin bool) error
}

// memoService 实现 MemoService 接口
type memoService struct {
	repo MemoRepository
	log  zerolog.Logger
}

// NewMemoService 创建 memo 服务实例
func NewMemoService(repo MemoRepository, log zerolog.Logger) MemoService {
	return &memoService{
		repo: repo,
		log:  log.With().Str("component", "memo_service").Logger(),
	}
}

func (s *memoService) GetByID(ctx context.Context, id int64, currentUserID *int64) (*db.Memo, error) {
	memo, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	// 权限检查
	if memo.ShowType != nil && *memo.ShowType != 1 {
		if currentUserID == nil || *currentUserID != memo.UserID {
			return nil, ErrUnauthorized
		}
	}
	
	return memo, nil
}

func (s *memoService) List(ctx context.Context, filter ListMemoFilter) ([]*db.Memo, int64, bool, error) {
	// 参数校验和默认值
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Size <= 0 {
		filter.Size = 10
	}
	if filter.Size > 100 {
		filter.Size = 100 // 限制最大分页大小
	}
	
	memos, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, 0, false, err
	}
	
	hasNext := int64(filter.Page*filter.Size) < total
	return memos, total, hasNext, nil
}

func (s *memoService) Create(ctx context.Context, req *vo.SaveMemoReq, userID int64) (*db.Memo, error) {
	// 业务校验
	if err := s.validateMemo(req); err != nil {
		return nil, err
	}
	
	now := time.Now()
	memo := &db.Memo{
		UserID:    userID,
		Content:   req.Content,
		CreatedAt: &now,
		UpdatedAt: &now,
		FavCount:  0,
	}
	
	// 处理标签
	if len(req.Tags) > 0 {
		tags := ""
		for _, tag := range req.Tags {
			tags += tag + ","
		}
		memo.Tags = &tags
	}
	
	// 处理图片
	if len(req.Imgs) > 0 {
		memo.Imgs = req.Imgs[0]
		for i := 1; i < len(req.Imgs); i++ {
			memo.Imgs += "," + req.Imgs[i]
		}
	}
	
	// 处理扩展字段
	if req.Ext != nil {
		extJSON, err := json.Marshal(req.Ext)
		if err != nil {
			s.log.Error().Err(err).Interface("ext", req.Ext).Msg("failed to marshal ext")
			return nil, fmt.Errorf("invalid ext data: %w", err)
		}
		extStr := string(extJSON)
		memo.Ext = &extStr
	}
	
	if err := s.repo.Create(ctx, memo); err != nil {
		return nil, err
	}
	
	s.log.Info().Int64("memo_id", memo.ID).Int64("user_id", userID).Msg("memo created")
	return memo, nil
}

func (s *memoService) Update(ctx context.Context, req *vo.SaveMemoReq, userID int64) (*db.Memo, error) {
	if req.ID <= 0 {
		return nil, ErrInvalidContent
	}
	
	// 获取现有 memo
	memo, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	
	// 权限检查
	if memo.UserID != userID {
		return nil, ErrUnauthorized
	}
	
	// 业务校验
	if err := s.validateMemo(req); err != nil {
		return nil, err
	}
	
	// 更新字段
	now := time.Now()
	memo.Content = req.Content
	memo.UpdatedAt = &now
	
	// 更新标签
	if len(req.Tags) > 0 {
		tags := ""
		for _, tag := range req.Tags {
			tags += tag + ","
		}
		memo.Tags = &tags
	} else {
		memo.Tags = nil
	}
	
	if err := s.repo.Update(ctx, memo); err != nil {
		return nil, err
	}
	
	s.log.Info().Int64("memo_id", memo.ID).Int64("user_id", userID).Msg("memo updated")
	return memo, nil
}

func (s *memoService) Delete(ctx context.Context, id int64, userID int64, isAdmin bool) error {
	memo, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	// 权限检查：只有作者或管理员可以删除
	if memo.UserID != userID && !isAdmin {
		return ErrUnauthorized
	}
	
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	
	s.log.Info().Int64("memo_id", id).Int64("user_id", userID).Bool("is_admin", isAdmin).Msg("memo deleted")
	return nil
}

func (s *memoService) Like(ctx context.Context, id int64) error {
	return s.repo.Like(ctx, id)
}

func (s *memoService) SetPinned(ctx context.Context, id int64, userID int64, isAdmin bool) error {
	if !isAdmin {
		return ErrUnauthorized
	}
	
	// 先取消所有置顶
	// 注意：这里需要事务处理，简化示例中省略
	
	memo, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	pinned := memo.Pinned == nil || !*memo.Pinned
	return s.repo.SetPinned(ctx, id, pinned)
}

func (s *memoService) validateMemo(req *vo.SaveMemoReq) error {
	content := req.Content
	if content == "" {
		return ErrInvalidContent
	}
	
	if len(content) > MaxContentLength {
		return ErrContentTooLong
	}
	
	return nil
}
