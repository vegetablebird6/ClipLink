package usecase

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	stdtime "time"

	"github.com/google/uuid"

	"github.com/xiaojiu/cliplink/internal/common/validation"
	"github.com/xiaojiu/cliplink/internal/domain/model"
	"github.com/xiaojiu/cliplink/internal/domain/repository"
	"github.com/xiaojiu/cliplink/internal/domain/service"
)

// computeContentHash 计算内容的 SHA-256 哈希。
// 去重语义：基于 trim 后的纯文本 content 计算，同文本不同 HTML 样式视为同一内容。
// 空字符串返回空，不参与去重。
func computeContentHash(content string) string {
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return ""
	}
	hash := sha256.Sum256([]byte(trimmed))
	return hex.EncodeToString(hash[:])
}

// clipboardService 剪贴板服务实现
type clipboardService struct {
	clipboardRepo repository.ClipboardRepository
	syncEventRepo repository.SyncEventRepository
	deviceRepo    repository.DeviceRepository
}

// NewClipboardService 创建新的剪贴板服务
func NewClipboardService(
	clipboardRepo repository.ClipboardRepository,
	syncEventRepo repository.SyncEventRepository,
	deviceRepo repository.DeviceRepository,
) service.ClipboardService {
	return &clipboardService{
		clipboardRepo: clipboardRepo,
		syncEventRepo: syncEventRepo,
		deviceRepo:    deviceRepo,
	}
}

// requireActorDevice 验证并返回执行操作的设备
func (s *clipboardService) requireActorDevice(deviceID, channelID string) (*model.Device, error) {
	if deviceID == "" {
		return nil, model.ErrInvalidInput
	}

	device, err := s.deviceRepo.FindByIDAndChannel(deviceID, channelID)
	if err != nil || device == nil {
		return nil, model.ErrInvalidInput
	}

	return device, nil
}

func deviceName(device *model.Device) string {
	if device == nil {
		return ""
	}
	return device.Name
}

func getDeviceTypeStr(device *model.Device) string {
	if device == nil {
		return ""
	}
	return device.Type
}

// SaveClipboard 保存剪贴板项目
func (s *clipboardService) SaveClipboard(title, content, contentType, deviceID, deviceType, channelID string, cleanDuplicates bool, contentHTML, contentFormat string) (*model.ClipboardItem, error) {
	if !validation.IsValidClipboardType(contentType) {
		return nil, model.ErrInvalidInput
	}
	if !validation.IsValidDeviceType(deviceType) {
		return nil, model.ErrInvalidInput
	}
	if !validation.IsValidContentFormat(contentFormat) {
		return nil, model.ErrInvalidInput
	}

	// 验证执行者设备（必须在写入同步事件前校验）
	device, err := s.requireActorDevice(deviceID, channelID)
	if err != nil {
		return nil, err
	}

	item := &model.ClipboardItem{
		ID:            uuid.New().String(),
		Title:         title,
		Content:       content,
		Type:          contentType,
		DeviceID:      deviceID,
		DeviceType:    deviceType,
		ChannelID:     channelID,
		CreatedAt:     stdtime.Now(),
		UpdatedAt:     stdtime.Now(),
		ContentHTML:   contentHTML,
		ContentFormat: contentFormat,
		ContentHash:   computeContentHash(content),
	}

	// 保存到数据库
	if err := s.clipboardRepo.Save(item); err != nil {
		return nil, err
	}

	if cleanDuplicates {
		// 优先用 hash 索引快速删除已有 hash 的重复项
		if item.ContentHash != "" {
			if _, err := s.clipboardRepo.DeleteByContentHash(channelID, item.ContentHash, item.ID); err != nil {
				return nil, err
			}
		}
	}

	// 记录同步事件
	contentSummary := content
	if len(contentSummary) > 100 {
		contentSummary = contentSummary[:100]
	}
	syncEvent := &model.SyncEvent{
		Action:          model.ActionSync,
		Content:         contentSummary,
		ChannelID:       channelID,
		TargetType:      model.TargetTypeClipboard,
		TargetID:        item.ID,
		Summary:         "同步剪贴板内容",
		ActorDeviceID:   device.ID,
		ActorDeviceName: device.Name,
		ActorDeviceType: device.Type,
		CreatedAt:       stdtime.Now(),
	}
	_ = s.syncEventRepo.Save(syncEvent)

	return item, nil
}

// GetLatestClipboard 获取最新的剪贴板项目
func (s *clipboardService) GetLatestClipboard(channelID string, limit int) ([]*model.ClipboardItem, error) {
	return s.clipboardRepo.FindLatest(channelID, limit)
}

// GetClipboardItem 获取剪贴板项目
func (s *clipboardService) GetClipboardItem(id string, channelID string) (*model.ClipboardItem, error) {
	return s.clipboardRepo.FindByID(id, channelID)
}

// GetClipboardHistory 获取剪贴板历史记录（keyset 游标分页）
func (s *clipboardService) GetClipboardHistory(channelID string, afterCreatedAt *stdtime.Time, afterID *string, size int) ([]*model.ClipboardItem, error) {
	return s.clipboardRepo.FindWithKeyset(channelID, afterCreatedAt, afterID, size)
}

// DeleteClipboard 删除剪贴板项目（actorDeviceID 为执行删除操作的设备）
func (s *clipboardService) DeleteClipboard(id string, channelID string, actorDeviceID string) error {
	// 验证执行者设备
	device, err := s.requireActorDevice(actorDeviceID, channelID)
	if err != nil {
		return err
	}

	// 记录同步历史（先获取 item 信息）
	item, err := s.clipboardRepo.FindByID(id, channelID)
	if err != nil {
		return err
	}

	// 从数据库删除
	if err := s.clipboardRepo.Delete(id, channelID); err != nil {
		return err
	}

	// 记录同步事件（用操作者设备快照作为 actor）
	syncEvent := &model.SyncEvent{
		Action:          model.ActionDelete,
		Content:         "删除剪贴板内容: " + item.Type,
		ChannelID:       channelID,
		TargetType:      model.TargetTypeClipboard,
		TargetID:        id,
		Summary:         "删除剪贴板内容",
		ActorDeviceID:   device.ID,
		ActorDeviceName: device.Name,
		ActorDeviceType: device.Type,
		CreatedAt:       stdtime.Now(),
	}

	return s.syncEventRepo.Save(syncEvent)
}

// UpdateClipboard 更新剪贴板项目（部分更新）
func (s *clipboardService) UpdateClipboard(id, channelID, actorDeviceID string, input *service.UpdateClipboardInput) (*model.ClipboardItem, error) {
	// 验证执行者设备
	device, err := s.requireActorDevice(actorDeviceID, channelID)
	if err != nil {
		return nil, err
	}

	// 校验字段（如果提供了）
	if input.Type != nil && !validation.IsValidClipboardType(*input.Type) {
		return nil, model.ErrInvalidInput
	}
	if input.DeviceType != nil && !validation.IsValidDeviceType(*input.DeviceType) {
		return nil, model.ErrInvalidInput
	}
	if input.ContentFormat != nil && !validation.IsValidContentFormat(*input.ContentFormat) {
		return nil, model.ErrInvalidInput
	}

	// 只更新请求中显式出现的字段
	updates := map[string]any{
		"updated_at": stdtime.Now(),
	}
	if input.Title != nil {
		updates["title"] = *input.Title
	}
	if input.Content != nil {
		updates["content"] = *input.Content
		updates["content_hash"] = computeContentHash(*input.Content)
	}
	if input.Type != nil {
		updates["type"] = *input.Type
	}
	if input.DeviceType != nil {
		updates["device_type"] = *input.DeviceType
	}
	if input.ContentHTML != nil {
		updates["content_html"] = *input.ContentHTML
	}
	if input.ContentFormat != nil {
		updates["content_format"] = *input.ContentFormat
	}

	// 更新到数据库
	if err := s.clipboardRepo.Update(id, channelID, updates); err != nil {
		return nil, err
	}

	// 记录同步事件
	contentType := ""
	if input.Type != nil {
		contentType = *input.Type
	}
	syncEvent := &model.SyncEvent{
		Action:          model.ActionUpdate,
		Content:         "更新剪贴板内容: " + contentType,
		ChannelID:       channelID,
		TargetType:      model.TargetTypeClipboard,
		TargetID:        id,
		Summary:         "更新剪贴板内容",
		ActorDeviceID:   device.ID,
		ActorDeviceName: device.Name,
		ActorDeviceType: device.Type,
		CreatedAt:       stdtime.Now(),
	}

	if err := s.syncEventRepo.Save(syncEvent); err != nil {
		return nil, err
	}

	// 获取更新后的数据
	return s.clipboardRepo.FindByID(id, channelID)
}

// ToggleFavorite 切换收藏状态（actorDeviceID 为执行操作的设备，必须属于 channel）
func (s *clipboardService) ToggleFavorite(id string, isFavorite bool, channelID string, actorDeviceID string) (*model.ClipboardItem, error) {
	// 验证 actor 设备（必须在写业务数据前校验）
	device, err := s.requireActorDevice(actorDeviceID, channelID)
	if err != nil {
		return nil, err
	}

	// 获取当前项目
	item, err := s.clipboardRepo.FindByID(id, channelID)
	if err != nil {
		return nil, err
	}

	// 设置收藏状态为指定值
	updates := map[string]any{
		"favorite":   isFavorite,
		"updated_at": stdtime.Now(),
	}

	// 更新到数据库
	if err := s.clipboardRepo.Update(id, channelID, updates); err != nil {
		return nil, err
	}

	// 记录同步事件
	syncEvent := &model.SyncEvent{
		Action:          model.ActionUpdate,
		Content:         item.Title,
		ChannelID:       channelID,
		TargetType:      model.TargetTypeClipboard,
		TargetID:        id,
		Summary:         "切换收藏状态",
		ActorDeviceID:   device.ID,
		ActorDeviceName: device.Name,
		ActorDeviceType: device.Type,
		CreatedAt:       stdtime.Now(),
	}
	// 忽略同步事件保存错误，不影响主流程
	_ = s.syncEventRepo.Save(syncEvent)

	// 获取更新后的数据
	return s.clipboardRepo.FindByID(id, channelID)
}

// GetFavoriteClipboard 获取收藏的剪贴板项目
func (s *clipboardService) GetFavoriteClipboard(channelID string, limit int) ([]*model.ClipboardItem, error) {
	return s.clipboardRepo.FindFavorites(channelID, limit)
}

// GetClipboardByType 按内容类型获取剪贴板历史记录（keyset 游标分页）
func (s *clipboardService) GetClipboardByType(contentType string, channelID string, afterCreatedAt *stdtime.Time, afterID *string, size int) ([]*model.ClipboardItem, error) {
	return s.clipboardRepo.FindByType(contentType, channelID, afterCreatedAt, afterID, size)
}

// GetClipboardByDeviceType 按设备类型获取剪贴板历史记录（keyset 游标分页）
func (s *clipboardService) GetClipboardByDeviceType(deviceType string, channelID string, afterCreatedAt *stdtime.Time, afterID *string, size int) ([]*model.ClipboardItem, error) {
	return s.clipboardRepo.FindByDeviceType(deviceType, channelID, afterCreatedAt, afterID, size)
}

// GetClipboardByTypeAndDeviceType 同时按内容类型和设备类型获取剪贴板历史记录（keyset 游标分页）
func (s *clipboardService) GetClipboardByTypeAndDeviceType(contentType, deviceType string, channelID string, afterCreatedAt *stdtime.Time, afterID *string, size int) ([]*model.ClipboardItem, error) {
	return s.clipboardRepo.FindByTypeAndDeviceType(contentType, deviceType, channelID, afterCreatedAt, afterID, size)
}

// SearchClipboard 按关键词搜索剪贴板项目
func (s *clipboardService) SearchClipboard(keyword, channelID string, page, size int) (items []*model.ClipboardItem, total int64, totalPages int, err error) {
	// 验证关键词不为空
	if keyword == "" {
		return []*model.ClipboardItem{}, 0, 0, nil
	}

	// 调用仓库层搜索方法
	return s.clipboardRepo.SearchByKeyword(keyword, channelID, page, size)
}

// CleanupDuplicateContents 清理同一通道下已存在的重复剪贴板内容。
func (s *clipboardService) CleanupDuplicateContents(channelID string) (int64, error) {
	return s.clipboardRepo.CleanupDuplicateContents(channelID)
}
