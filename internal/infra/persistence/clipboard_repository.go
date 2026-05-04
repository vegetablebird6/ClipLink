package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/xiaojiu/cliplink/internal/domain/model"
	"github.com/xiaojiu/cliplink/internal/domain/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// clipboardRepository 剪贴板仓库实现
type clipboardRepository struct {
	gdb *gorm.DB
}

type duplicateCandidate struct {
	ID          string
	ContentHash string
	CreatedAt   time.Time
}

// NewClipboardRepository 创建新的剪贴板仓库
func NewClipboardRepository(gdb *gorm.DB) repository.ClipboardRepository {
	return &clipboardRepository{gdb: gdb}
}

// Save 保存剪贴板项目
func (r *clipboardRepository) Save(ctx context.Context, item *model.ClipboardItem) error {
	return r.gdb.WithContext(ctx).Create(item).Error
}

// FindByID 通过ID查找剪贴板项目
func (r *clipboardRepository) FindByID(ctx context.Context, id, channelID string) (*model.ClipboardItem, error) {
	var item model.ClipboardItem
	result := r.gdb.WithContext(ctx).Where("id = ? AND channel_id = ?", id, channelID).First(&item)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, model.ErrClipboardNotFound
		}
		return nil, result.Error
	}
	return &item, nil
}

// FindLatest 获取最新的剪贴板项目
func (r *clipboardRepository) FindLatest(ctx context.Context, channelID string, limit int) ([]*model.ClipboardItem, error) {
	var items []*model.ClipboardItem
	query := r.gdb.WithContext(ctx).Model(&model.ClipboardItem{})
	if channelID != "" {
		query = query.Where("channel_id = ?", channelID)
	}
	err := query.Order("created_at DESC").Limit(limit).Find(&items).Error
	return items, err
}

// FindWithKeyset 分页获取剪贴板项目（keyset 游标分页）
func (r *clipboardRepository) FindWithKeyset(ctx context.Context, channelID string, afterCreatedAt *time.Time, afterID *string, size int) ([]*model.ClipboardItem, error) {
	var items []*model.ClipboardItem
	query := r.gdb.WithContext(ctx).Where("channel_id = ?", channelID)
	if afterCreatedAt != nil && afterID != nil {
		query = query.Where("(created_at < ?) OR (created_at = ? AND id < ?)",
			*afterCreatedAt, *afterCreatedAt, *afterID)
	}
	err := query.Order("created_at DESC, id DESC").Limit(size + 1).Find(&items).Error
	return items, err
}

// FindByType 按类型查找剪贴板项目（keyset 游标分页）
func (r *clipboardRepository) FindByType(ctx context.Context, contentType, channelID string, afterCreatedAt *time.Time, afterID *string, size int) ([]*model.ClipboardItem, error) {
	var items []*model.ClipboardItem
	query := r.gdb.WithContext(ctx).Where("channel_id = ? AND type = ?", channelID, contentType)
	if afterCreatedAt != nil && afterID != nil {
		query = query.Where("(created_at < ?) OR (created_at = ? AND id < ?)",
			*afterCreatedAt, *afterCreatedAt, *afterID)
	}
	err := query.Order("created_at DESC, id DESC").Limit(size + 1).Find(&items).Error
	return items, err
}

// FindByDeviceType 按设备类型查找剪贴板项目（keyset 游标分页）
func (r *clipboardRepository) FindByDeviceType(ctx context.Context, deviceType, channelID string, afterCreatedAt *time.Time, afterID *string, size int) ([]*model.ClipboardItem, error) {
	var items []*model.ClipboardItem
	query := r.gdb.WithContext(ctx).Where("channel_id = ? AND device_type = ?", channelID, deviceType)
	if afterCreatedAt != nil && afterID != nil {
		query = query.Where("(created_at < ?) OR (created_at = ? AND id < ?)",
			*afterCreatedAt, *afterCreatedAt, *afterID)
	}
	err := query.Order("created_at DESC, id DESC").Limit(size + 1).Find(&items).Error
	return items, err
}

// FindByTypeAndDeviceType 同时按内容类型和设备类型查找剪贴板项目（keyset 游标分页）
func (r *clipboardRepository) FindByTypeAndDeviceType(ctx context.Context, contentType, deviceType, channelID string, afterCreatedAt *time.Time, afterID *string, size int) ([]*model.ClipboardItem, error) {
	var items []*model.ClipboardItem
	query := r.gdb.WithContext(ctx).Where("channel_id = ?", channelID)
	if contentType != "" {
		query = query.Where("type = ?", contentType)
	}
	if deviceType != "" {
		query = query.Where("device_type = ?", deviceType)
	}
	if afterCreatedAt != nil && afterID != nil {
		query = query.Where("(created_at < ?) OR (created_at = ? AND id < ?)",
			*afterCreatedAt, *afterCreatedAt, *afterID)
	}
	err := query.Order("created_at DESC, id DESC").Limit(size + 1).Find(&items).Error
	return items, err
}

// FindFavorites 查找收藏的剪贴板项目（keyset 游标分页）
func (r *clipboardRepository) FindFavorites(ctx context.Context, channelID string, afterCreatedAt *time.Time, afterID *string, size int) ([]*model.ClipboardItem, error) {
	var items []*model.ClipboardItem
	query := r.gdb.WithContext(ctx).Where("channel_id = ? AND favorite = ?", channelID, true)
	if afterCreatedAt != nil && afterID != nil {
		query = query.Where("(updated_at < ?) OR (updated_at = ? AND id < ?)",
			*afterCreatedAt, *afterCreatedAt, *afterID)
	}
	err := query.Order("updated_at DESC, id DESC").Limit(size + 1).Find(&items).Error
	return items, err
}

// Update 更新剪贴板项目
func (r *clipboardRepository) Update(ctx context.Context, id, channelID string, updates map[string]interface{}) error {
	result := r.gdb.WithContext(ctx).Model(&model.ClipboardItem{}).
		Where("id = ? AND channel_id = ?", id, channelID).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return model.ErrClipboardNotFound
	}

	return nil
}

// Delete 删除剪贴板项目
func (r *clipboardRepository) Delete(ctx context.Context, id, channelID string) error {
	result := r.gdb.WithContext(ctx).Where("id = ? AND channel_id = ?", id, channelID).
		Delete(&model.ClipboardItem{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return model.ErrClipboardNotFound
	}

	return nil
}

// DeleteByContentHash 基于内容哈希删除同通道下的重复项，保留指定项目。
func (r *clipboardRepository) DeleteByContentHash(ctx context.Context, channelID, contentHash, keepID string) (int64, error) {
	if channelID == "" || contentHash == "" || keepID == "" {
		return 0, nil
	}

	result := r.gdb.WithContext(ctx).
		Where("channel_id = ? AND id <> ? AND content_hash = ?", channelID, keepID, contentHash).
		Delete(&model.ClipboardItem{})
	return result.RowsAffected, result.Error
}

// cleanupBatchSize 批量扫描时每批查询的记录数
const cleanupBatchSize = 1000

// CleanupDuplicateContents 清理同通道下已存在的重复内容，保留每组最新项目。
func (r *clipboardRepository) CleanupDuplicateContents(ctx context.Context, channelID string) (int64, error) {
	if channelID == "" {
		return 0, nil
	}

	seen := make(map[string]struct{})
	duplicateIDs := make([]string, 0)
	var cursorCreatedAt time.Time
	var cursorID string

	for {
		candidates, lastCreatedAt, lastID, err := r.fetchDuplicateBatch(ctx, channelID, cursorCreatedAt, cursorID)
		if err != nil {
			return 0, err
		}
		if len(candidates) == 0 {
			break
		}

		r.collectDuplicates(candidates, seen, &duplicateIDs)

		cursorCreatedAt = lastCreatedAt
		cursorID = lastID

		if len(candidates) < cleanupBatchSize {
			break
		}
	}

	if len(duplicateIDs) == 0 {
		return 0, nil
	}

	var totalDeleted int64
	for i := 0; i < len(duplicateIDs); i += cleanupBatchSize {
		end := i + cleanupBatchSize
		if end > len(duplicateIDs) {
			end = len(duplicateIDs)
		}
		result := r.gdb.WithContext(ctx).
			Where("id IN ?", duplicateIDs[i:end]).
			Delete(&model.ClipboardItem{})
		if result.Error != nil {
			return totalDeleted, result.Error
		}
		totalDeleted += result.RowsAffected
	}

	return totalDeleted, nil
}

// duplicateBatchResult 用于接收 keyset 分页查询的完整行（含 created_at 作为游标）
type duplicateBatchResult struct {
	ID          string
	ContentHash string
	CreatedAt   time.Time
}

// fetchDuplicateBatch 使用 keyset 游标分批读取候选记录（仅读取，不删除）。
func (r *clipboardRepository) fetchDuplicateBatch(ctx context.Context, channelID string, cursorCreatedAt time.Time, cursorID string) ([]duplicateCandidate, time.Time, string, error) {
	var rows []duplicateBatchResult

	query := r.gdb.WithContext(ctx).
		Model(&model.ClipboardItem{}).
		Select([]string{"id", "content_hash", "created_at"}).
		Where("channel_id = ? AND content_hash <> '' AND content_hash IS NOT NULL", channelID)

	if !cursorCreatedAt.IsZero() {
		query = query.Where("(created_at < ?) OR (created_at = ? AND id < ?)",
			cursorCreatedAt, cursorCreatedAt, cursorID)
	}

	if err := query.
		Order("created_at DESC, id DESC").
		Limit(cleanupBatchSize).
		Find(&rows).Error; err != nil {
		return nil, time.Time{}, "", err
	}

	candidates := make([]duplicateCandidate, len(rows))
	var lastCreatedAt time.Time
	var lastID string
	for i, row := range rows {
		candidates[i] = duplicateCandidate{
			ID:          row.ID,
			ContentHash: row.ContentHash,
		}
		lastCreatedAt = row.CreatedAt
		lastID = row.ID
	}

	return candidates, lastCreatedAt, lastID, nil
}

// collectDuplicates 从一批候选记录中识别重复项，将重复 ID 追加到结果切片。
func (r *clipboardRepository) collectDuplicates(candidates []duplicateCandidate, seen map[string]struct{}, duplicateIDs *[]string) {
	for _, c := range candidates {
		key := c.ContentHash
		if key == "" {
			continue
		}
		if _, exists := seen[key]; exists {
			*duplicateIDs = append(*duplicateIDs, c.ID)
			continue
		}
		seen[key] = struct{}{}
	}
}

// Count 统计剪贴板项目数量
func (r *clipboardRepository) Count(ctx context.Context, channelID string) (int64, error) {
	var count int64
	query := r.gdb.WithContext(ctx).Model(&model.ClipboardItem{})

	if channelID != "" {
		query = query.Where("channel_id = ?", channelID)
	}

	err := query.Count(&count).Error
	return count, err
}

// CountByType 按类型统计剪贴板项目数量
func (r *clipboardRepository) CountByType(ctx context.Context, contentType, channelID string) (int64, error) {
	var count int64
	query := r.gdb.WithContext(ctx).Model(&model.ClipboardItem{}).Where("type = ?", contentType)

	if channelID != "" {
		query = query.Where("channel_id = ?", channelID)
	}

	err := query.Count(&count).Error
	return count, err
}

// SearchByKeyword 按关键词搜索剪贴板项目（支持标题和内容搜索，offset 分页）
func (r *clipboardRepository) SearchByKeyword(ctx context.Context, keyword, channelID string, page, size int) ([]*model.ClipboardItem, int64, int, error) {
	offset := (page - 1) * size
	var items []*model.ClipboardItem
	var total int64

	searchPattern := "%" + keyword + "%"
	query := r.gdb.WithContext(ctx).Model(&model.ClipboardItem{}).Where(
		"(title LIKE ? OR content LIKE ?)",
		searchPattern, searchPattern,
	)

	if channelID != "" {
		query = query.Where("channel_id = ?", channelID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	totalPages := int(total / int64(size))
	if total%int64(size) > 0 {
		totalPages++
	}

	orderClause := clause.OrderBy{
		Expression: clause.Expr{
			SQL:                "CASE WHEN title LIKE ? THEN 0 ELSE 1 END, created_at DESC",
			Vars:               []interface{}{searchPattern},
			WithoutParentheses: true,
		},
	}
	if err := query.Clauses(orderClause).Offset(offset).Limit(size).Find(&items).Error; err != nil {
		return nil, 0, 0, err
	}

	return items, total, totalPages, nil
}
