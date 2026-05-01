package persistence

import (
	"errors"
	"time"

	"github.com/xiaojiu/cliplink/internal/domain/model"
	"github.com/xiaojiu/cliplink/internal/domain/repository"
	"github.com/xiaojiu/cliplink/internal/infra/db"
	"gorm.io/gorm/clause"
)

// clipboardRepository 剪贴板仓库实现
type clipboardRepository struct{}

type duplicateCandidate struct {
	ID          string
	ContentHash string
	CreatedAt   time.Time
}

// NewClipboardRepository 创建新的剪贴板仓库
func NewClipboardRepository() repository.ClipboardRepository {
	return &clipboardRepository{}
}

// Save 保存剪贴板项目
func (r *clipboardRepository) Save(item *model.ClipboardItem) error {
	return db.GetDB().Create(item).Error
}

// FindByID 通过ID查找剪贴板项目
func (r *clipboardRepository) FindByID(id, channelID string) (*model.ClipboardItem, error) {
	var item model.ClipboardItem
	result := db.GetDB().Where("id = ? AND channel_id = ?", id, channelID).First(&item)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

// FindLatest 获取最新的剪贴板项目
func (r *clipboardRepository) FindLatest(channelID string, limit int) ([]*model.ClipboardItem, error) {
	var items []*model.ClipboardItem
	query := db.GetDB().Model(&model.ClipboardItem{})
	if channelID != "" {
		query = query.Where("channel_id = ?", channelID)
	}
	err := query.Order("created_at DESC").Limit(limit).Find(&items).Error
	return items, err
}

// FindWithPagination 分页获取剪贴板项目
func (r *clipboardRepository) FindWithPagination(channelID string, page, size int) ([]*model.ClipboardItem, int64, int, error) {
	offset := (page - 1) * size
	var items []*model.ClipboardItem
	var total int64

	// 获取符合条件的记录总数
	query := db.GetDB().Model(&model.ClipboardItem{})
	if channelID != "" {
		query = query.Where("channel_id = ?", channelID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	// 计算总页数
	totalPages := int(total / int64(size))
	if total%int64(size) > 0 {
		totalPages++
	}

	// 获取分页数据
	if err := query.Order("created_at DESC").Offset(offset).Limit(size).Find(&items).Error; err != nil {
		return nil, 0, 0, err
	}

	return items, total, totalPages, nil
}

// FindByType 按类型查找剪贴板项目
func (r *clipboardRepository) FindByType(contentType, channelID string, page, size int) ([]*model.ClipboardItem, int64, int, error) {
	offset := (page - 1) * size
	var items []*model.ClipboardItem
	var total int64

	// 构建查询
	query := db.GetDB().Model(&model.ClipboardItem{}).Where("type = ?", contentType)
	if channelID != "" {
		query = query.Where("channel_id = ?", channelID)
	}

	// 获取总记录数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	// 计算总页数
	totalPages := int(total / int64(size))
	if total%int64(size) > 0 {
		totalPages++
	}

	// 获取分页数据
	if err := query.Order("created_at DESC").Offset(offset).Limit(size).Find(&items).Error; err != nil {
		return nil, 0, 0, err
	}

	return items, total, totalPages, nil
}

// FindByDeviceType 按设备类型查找剪贴板项目
func (r *clipboardRepository) FindByDeviceType(deviceType, channelID string, page, size int) ([]*model.ClipboardItem, int64, int, error) {
	offset := (page - 1) * size
	var items []*model.ClipboardItem
	var total int64

	// 构建查询
	query := db.GetDB().Model(&model.ClipboardItem{}).Where("device_type = ?", deviceType)
	if channelID != "" {
		query = query.Where("channel_id = ?", channelID)
	}

	// 获取总记录数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	// 计算总页数
	totalPages := int(total / int64(size))
	if total%int64(size) > 0 {
		totalPages++
	}

	// 获取分页数据
	if err := query.Order("created_at DESC").Offset(offset).Limit(size).Find(&items).Error; err != nil {
		return nil, 0, 0, err
	}

	return items, total, totalPages, nil
}

// FindByTypeAndDeviceType 同时按内容类型和设备类型查找剪贴板项目
func (r *clipboardRepository) FindByTypeAndDeviceType(contentType, deviceType, channelID string, page, size int) ([]*model.ClipboardItem, int64, int, error) {
	var items []*model.ClipboardItem
	var total int64

	query := db.GetDB().Model(&model.ClipboardItem{})

	if channelID != "" {
		query = query.Where("channel_id = ?", channelID)
	}

	if contentType != "" {
		query = query.Where("type = ?", contentType)
	}

	if deviceType != "" {
		query = query.Where("device_type = ?", deviceType)
	}

	// 获取总记录数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	// 计算总页数
	totalPages := int(total / int64(size))
	if total%int64(size) > 0 {
		totalPages++
	}

	// 获取分页数据
	offset := (page - 1) * size
	if err := query.Order("created_at DESC").Offset(offset).Limit(size).Find(&items).Error; err != nil {
		return nil, 0, 0, err
	}

	return items, total, totalPages, nil
}

// FindFavorites 查找收藏的剪贴板项目
func (r *clipboardRepository) FindFavorites(channelID string, limit int) ([]*model.ClipboardItem, error) {
	var items []*model.ClipboardItem
	query := db.GetDB()

	if channelID != "" {
		query = query.Where("channel_id = ?", channelID)
	}

	err := query.Where("favorite = ?", true).
		Order("updated_at DESC").
		Limit(limit).
		Find(&items).Error
	return items, err
}

// Update 更新剪贴板项目
func (r *clipboardRepository) Update(id, channelID string, updates map[string]interface{}) error {
	result := db.GetDB().Model(&model.ClipboardItem{}).
		Where("id = ? AND channel_id = ?", id, channelID).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("item not found")
	}

	return nil
}

// Delete 删除剪贴板项目
func (r *clipboardRepository) Delete(id, channelID string) error {
	result := db.GetDB().Where("id = ? AND channel_id = ?", id, channelID).
		Delete(&model.ClipboardItem{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("item not found")
	}

	return nil
}

// DeleteByContentHash 基于内容哈希删除同通道下的重复项，保留指定项目。
// 使用 (channel_id, content_hash) 复合索引，避免全表扫描。
func (r *clipboardRepository) DeleteByContentHash(channelID, contentHash, keepID string) (int64, error) {
	if channelID == "" || contentHash == "" || keepID == "" {
		return 0, nil
	}

	result := db.GetDB().
		Where("channel_id = ? AND id <> ? AND content_hash = ?", channelID, keepID, contentHash).
		Delete(&model.ClipboardItem{})
	return result.RowsAffected, result.Error
}

// cleanupBatchSize 批量扫描时每批查询的记录数
const cleanupBatchSize = 1000

// CleanupDuplicateContents 清理同通道下已存在的重复内容，保留每组最新项目。
// 使用 keyset 分批扫描（不删除），seen map 跨批保留以正确处理跨批重复项，
// 最后统一批量删除，避免 OFFSET 边扫边删导致跳过记录。
// 未上线前不兼容旧数据：只处理非空 content_hash 记录。
func (r *clipboardRepository) CleanupDuplicateContents(channelID string) (int64, error) {
	if channelID == "" {
		return 0, nil
	}

	seen := make(map[string]struct{})
	duplicateIDs := make([]string, 0)
	var cursorCreatedAt time.Time
	var cursorID string

	for {
		candidates, lastCreatedAt, lastID, err := r.fetchDuplicateBatch(channelID, cursorCreatedAt, cursorID)
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

	// 统一分批删除，避免单条 SQL IN (...) 过长
	var totalDeleted int64
	for i := 0; i < len(duplicateIDs); i += cleanupBatchSize {
		end := i + cleanupBatchSize
		if end > len(duplicateIDs) {
			end = len(duplicateIDs)
		}
		result := db.GetDB().
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
func (r *clipboardRepository) fetchDuplicateBatch(channelID string, cursorCreatedAt time.Time, cursorID string) ([]duplicateCandidate, time.Time, string, error) {
	var rows []duplicateBatchResult

	query := db.GetDB().
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
// seen map 在调用方维护，跨批保留，确保跨批重复项也能被正确识别。
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
func (r *clipboardRepository) Count(channelID string) (int64, error) {
	var count int64
	query := db.GetDB().Model(&model.ClipboardItem{})

	if channelID != "" {
		query = query.Where("channel_id = ?", channelID)
	}

	err := query.Count(&count).Error
	return count, err
}

// CountByType 按类型统计剪贴板项目数量
func (r *clipboardRepository) CountByType(contentType, channelID string) (int64, error) {
	var count int64
	query := db.GetDB().Model(&model.ClipboardItem{}).Where("type = ?", contentType)

	if channelID != "" {
		query = query.Where("channel_id = ?", channelID)
	}

	err := query.Count(&count).Error
	return count, err
}

// SearchByKeyword 按关键词搜索剪贴板项目（支持标题和内容搜索）
func (r *clipboardRepository) SearchByKeyword(keyword, channelID string, page, size int) ([]*model.ClipboardItem, int64, int, error) {
	offset := (page - 1) * size
	var items []*model.ClipboardItem
	var total int64

	// 构建搜索查询 - 在标题和内容中搜索关键词
	searchPattern := "%" + keyword + "%"
	query := db.GetDB().Model(&model.ClipboardItem{}).Where(
		"(title LIKE ? OR content LIKE ?)",
		searchPattern, searchPattern,
	)

	if channelID != "" {
		query = query.Where("channel_id = ?", channelID)
	}

	// 获取总记录数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	// 计算总页数
	totalPages := int(total / int64(size))
	if total%int64(size) > 0 {
		totalPages++
	}

	// 获取分页数据，按相关度和时间排序，避免将搜索词拼接进SQL。
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
