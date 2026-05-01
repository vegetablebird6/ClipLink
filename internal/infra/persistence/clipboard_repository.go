package persistence

import (
	"errors"
	"strings"
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
	Content     string
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

// DeleteDuplicates 删除同通道下内容相同的重复项，保留指定项目。
// Deprecated: 使用 DeleteByContentHash 替代，性能更优。
func (r *clipboardRepository) DeleteDuplicates(channelID, content, keepID string) error {
	if channelID == "" || content == "" || keepID == "" {
		return nil
	}

	return db.GetDB().
		Where("channel_id = ? AND id <> ? AND TRIM(content) = ?", channelID, keepID, strings.TrimSpace(content)).
		Delete(&model.ClipboardItem{}).Error
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

// cleanupBatchSize 批量清理时每批查询的记录数
const cleanupBatchSize = 1000

// CleanupDuplicateContents 清理同通道下已存在的重复内容，保留每组最新项目。
// 采用分批查询 + 按 content_hash 分组策略，避免一次性全量加载到内存。
func (r *clipboardRepository) CleanupDuplicateContents(channelID string) (int64, error) {
	if channelID == "" {
		return 0, nil
	}

	var totalDeleted int64
	offset := 0

	for {
		var candidates []duplicateCandidate
		if err := db.GetDB().
			Model(&model.ClipboardItem{}).
			Select("id", "content", "content_hash", "created_at").
			Where("channel_id = ?", channelID).
			Order("created_at DESC").
			Offset(offset).Limit(cleanupBatchSize).
			Find(&candidates).Error; err != nil {
			return totalDeleted, err
		}

		if len(candidates) == 0 {
			break
		}

		duplicateIDs := r.findDuplicateIDs(candidates)
		if len(duplicateIDs) > 0 {
			result := db.GetDB().
				Where("id IN ?", duplicateIDs).
				Delete(&model.ClipboardItem{})
			if result.Error != nil {
				return totalDeleted, result.Error
			}
			totalDeleted += result.RowsAffected
		}

		if len(candidates) < cleanupBatchSize {
			break
		}
		offset += cleanupBatchSize
	}

	return totalDeleted, nil
}

// findDuplicateIDs 从一批候选记录中找出重复项的 ID。
// 优先使用 content_hash 分组，对空 hash 的旧记录回退到 TRIM(content) 比较。
func (r *clipboardRepository) findDuplicateIDs(candidates []duplicateCandidate) []string {
	duplicateIDs := make([]string, 0)

	// 第一轮：按 content_hash 分组（快速路径）
	seenHash := make(map[string]struct{}, len(candidates))
	for _, c := range candidates {
		if c.ContentHash == "" {
			continue
		}
		if _, exists := seenHash[c.ContentHash]; exists {
			duplicateIDs = append(duplicateIDs, c.ID)
			continue
		}
		seenHash[c.ContentHash] = struct{}{}
	}

	// 第二轮：处理空 hash 的旧记录，回退到 TRIM(content) 比较
	seenContent := make(map[string]struct{})
	for _, c := range candidates {
		if c.ContentHash != "" {
			continue
		}
		normalized := strings.TrimSpace(c.Content)
		if normalized == "" {
			continue
		}
		if _, exists := seenContent[normalized]; exists {
			duplicateIDs = append(duplicateIDs, c.ID)
			continue
		}
		seenContent[normalized] = struct{}{}
	}

	return duplicateIDs
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
