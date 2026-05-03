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

// CreateClipboard 创建剪贴板条目
func (s *clipboardService) CreateClipboard(in service.CreateClipboardInput) (*service.ClipboardItemOutput, error) {
	if !validation.IsValidClipboardType(in.Type) {
		return nil, model.ErrInvalidInput
	}
	if !validation.IsValidDeviceType(in.ActorDeviceType) {
		return nil, model.ErrInvalidInput
	}
	if !validation.IsValidContentFormat(in.ContentFormat) {
		return nil, model.ErrInvalidInput
	}

	device, err := s.requireActorDevice(in.ActorDeviceID, in.ChannelID)
	if err != nil {
		return nil, err
	}

	item := &model.ClipboardItem{
		ID:            uuid.New().String(),
		Title:         in.Title,
		Content:       in.Content,
		Type:          in.Type,
		DeviceID:      in.ActorDeviceID,
		DeviceType:    in.ActorDeviceType,
		ChannelID:     in.ChannelID,
		CreatedAt:     stdtime.Now(),
		UpdatedAt:     stdtime.Now(),
		ContentHTML:   in.ContentHTML,
		ContentFormat: in.ContentFormat,
		ContentHash:   computeContentHash(in.Content),
	}

	if err := s.clipboardRepo.Save(item); err != nil {
		return nil, err
	}

	if in.CleanDuplicates {
		if item.ContentHash != "" {
			if _, err := s.clipboardRepo.DeleteByContentHash(in.ChannelID, item.ContentHash, item.ID); err != nil {
				return nil, err
			}
		}
	}

	// 记录同步事件
	contentSummary := in.Content
	if len(contentSummary) > 100 {
		contentSummary = contentSummary[:100]
	}
	syncEvent := &model.SyncEvent{
		Action:          model.ActionSync,
		Content:         contentSummary,
		ChannelID:       in.ChannelID,
		TargetType:      model.TargetTypeClipboard,
		TargetID:        item.ID,
		Summary:         "同步剪贴板内容",
		ActorDeviceID:   device.ID,
		ActorDeviceName: device.Name,
		ActorDeviceType: device.Type,
		CreatedAt:       stdtime.Now(),
	}
	_ = s.syncEventRepo.Save(syncEvent)

	return toClipboardItemOutput(item), nil
}

// GetLatestClipboard 获取最新的剪贴板条目
func (s *clipboardService) GetLatestClipboard(channelID string, limit int) ([]*service.ClipboardItemOutput, error) {
	items, err := s.clipboardRepo.FindLatest(channelID, limit)
	if err != nil {
		return nil, err
	}
	return toClipboardItemOutputs(items), nil
}

// GetClipboardItem 获取剪贴板条目
func (s *clipboardService) GetClipboardItem(id, channelID string) (*service.ClipboardItemOutput, error) {
	item, err := s.clipboardRepo.FindByID(id, channelID)
	if err != nil {
		return nil, err
	}
	return toClipboardItemOutput(item), nil
}

// GetClipboardHistory 获取剪贴板历史记录（keyset 游标分页）
func (s *clipboardService) GetClipboardHistory(channelID string, afterCreatedAt *stdtime.Time, afterID *string, size int) ([]*service.ClipboardItemOutput, error) {
	items, err := s.clipboardRepo.FindWithKeyset(channelID, afterCreatedAt, afterID, size)
	if err != nil {
		return nil, err
	}
	return toClipboardItemOutputs(items), nil
}

// DeleteClipboard 删除剪贴板条目
func (s *clipboardService) DeleteClipboard(in service.DeleteClipboardInput) error {
	device, err := s.requireActorDevice(in.ActorDeviceID, in.ChannelID)
	if err != nil {
		return err
	}

	item, err := s.clipboardRepo.FindByID(in.ID, in.ChannelID)
	if err != nil {
		return err
	}

	if err := s.clipboardRepo.Delete(in.ID, in.ChannelID); err != nil {
		return err
	}

	syncEvent := &model.SyncEvent{
		Action:          model.ActionDelete,
		Content:         "删除剪贴板内容: " + item.Type,
		ChannelID:       in.ChannelID,
		TargetType:      model.TargetTypeClipboard,
		TargetID:        in.ID,
		Summary:         "删除剪贴板内容",
		ActorDeviceID:   device.ID,
		ActorDeviceName: device.Name,
		ActorDeviceType: device.Type,
		CreatedAt:       stdtime.Now(),
	}

	return s.syncEventRepo.Save(syncEvent)
}

// UpdateClipboard 更新剪贴板条目（部分更新）
func (s *clipboardService) UpdateClipboard(in service.UpdateClipboardInput) (*service.ClipboardItemOutput, error) {
	device, err := s.requireActorDevice(in.ActorDeviceID, in.ChannelID)
	if err != nil {
		return nil, err
	}

	if in.Type != nil && !validation.IsValidClipboardType(*in.Type) {
		return nil, model.ErrInvalidInput
	}
	if in.DeviceType != nil && !validation.IsValidDeviceType(*in.DeviceType) {
		return nil, model.ErrInvalidInput
	}
	if in.ContentFormat != nil && !validation.IsValidContentFormat(*in.ContentFormat) {
		return nil, model.ErrInvalidInput
	}

	updates := map[string]any{
		"updated_at": stdtime.Now(),
	}
	if in.Title != nil {
		updates["title"] = *in.Title
	}
	if in.Content != nil {
		updates["content"] = *in.Content
		updates["content_hash"] = computeContentHash(*in.Content)
	}
	if in.Type != nil {
		updates["type"] = *in.Type
	}
	if in.DeviceType != nil {
		updates["device_type"] = *in.DeviceType
	}
	if in.ContentHTML != nil {
		updates["content_html"] = *in.ContentHTML
	}
	if in.ContentFormat != nil {
		updates["content_format"] = *in.ContentFormat
	}

	if err := s.clipboardRepo.Update(in.ID, in.ChannelID, updates); err != nil {
		return nil, err
	}

	// 记录同步事件
	contentType := ""
	if in.Type != nil {
		contentType = *in.Type
	}
	syncEvent := &model.SyncEvent{
		Action:          model.ActionUpdate,
		Content:         "更新剪贴板内容: " + contentType,
		ChannelID:       in.ChannelID,
		TargetType:      model.TargetTypeClipboard,
		TargetID:        in.ID,
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
	item, err := s.clipboardRepo.FindByID(in.ID, in.ChannelID)
	if err != nil {
		return nil, err
	}
	return toClipboardItemOutput(item), nil
}

// SetFavorite 设置收藏状态
func (s *clipboardService) SetFavorite(in service.SetFavoriteInput) (*service.ClipboardItemOutput, error) {
	device, err := s.requireActorDevice(in.ActorDeviceID, in.ChannelID)
	if err != nil {
		return nil, err
	}

	item, err := s.clipboardRepo.FindByID(in.ID, in.ChannelID)
	if err != nil {
		return nil, err
	}

	updates := map[string]any{
		"favorite":   in.Favorite,
		"updated_at": stdtime.Now(),
	}

	if err := s.clipboardRepo.Update(in.ID, in.ChannelID, updates); err != nil {
		return nil, err
	}

	syncEvent := &model.SyncEvent{
		Action:          model.ActionUpdate,
		Content:         item.Title,
		ChannelID:       in.ChannelID,
		TargetType:      model.TargetTypeClipboard,
		TargetID:        in.ID,
		Summary:         "切换收藏状态",
		ActorDeviceID:   device.ID,
		ActorDeviceName: device.Name,
		ActorDeviceType: device.Type,
		CreatedAt:       stdtime.Now(),
	}
	_ = s.syncEventRepo.Save(syncEvent)

	// 获取更新后的数据
	updated, err := s.clipboardRepo.FindByID(in.ID, in.ChannelID)
	if err != nil {
		return nil, err
	}
	return toClipboardItemOutput(updated), nil
}

// GetFavoriteClipboard 获取收藏的剪贴板条目（keyset 游标分页）
func (s *clipboardService) GetFavoriteClipboard(channelID string, afterCreatedAt *stdtime.Time, afterID *string, size int) ([]*service.ClipboardItemOutput, error) {
	items, err := s.clipboardRepo.FindFavorites(channelID, afterCreatedAt, afterID, size)
	if err != nil {
		return nil, err
	}
	return toClipboardItemOutputs(items), nil
}

// GetClipboardByType 按内容类型获取剪贴板历史记录（keyset 游标分页）
func (s *clipboardService) GetClipboardByType(contentType string, channelID string, afterCreatedAt *stdtime.Time, afterID *string, size int) ([]*service.ClipboardItemOutput, error) {
	items, err := s.clipboardRepo.FindByType(contentType, channelID, afterCreatedAt, afterID, size)
	if err != nil {
		return nil, err
	}
	return toClipboardItemOutputs(items), nil
}

// GetClipboardByDeviceType 按设备类型获取剪贴板历史记录（keyset 游标分页）
func (s *clipboardService) GetClipboardByDeviceType(deviceType string, channelID string, afterCreatedAt *stdtime.Time, afterID *string, size int) ([]*service.ClipboardItemOutput, error) {
	items, err := s.clipboardRepo.FindByDeviceType(deviceType, channelID, afterCreatedAt, afterID, size)
	if err != nil {
		return nil, err
	}
	return toClipboardItemOutputs(items), nil
}

// GetClipboardByTypeAndDeviceType 同时按内容类型和设备类型获取剪贴板历史记录（keyset 游标分页）
func (s *clipboardService) GetClipboardByTypeAndDeviceType(contentType, deviceType string, channelID string, afterCreatedAt *stdtime.Time, afterID *string, size int) ([]*service.ClipboardItemOutput, error) {
	items, err := s.clipboardRepo.FindByTypeAndDeviceType(contentType, deviceType, channelID, afterCreatedAt, afterID, size)
	if err != nil {
		return nil, err
	}
	return toClipboardItemOutputs(items), nil
}

// SearchClipboard 按关键词搜索剪贴板条目
func (s *clipboardService) SearchClipboard(keyword, channelID string, page, size int) (items []*service.ClipboardItemOutput, total int64, totalPages int, err error) {
	if keyword == "" {
		return []*service.ClipboardItemOutput{}, 0, 0, nil
	}

	modelItems, total, totalPages, err := s.clipboardRepo.SearchByKeyword(keyword, channelID, page, size)
	if err != nil {
		return nil, 0, 0, err
	}
	return toClipboardItemOutputs(modelItems), total, totalPages, nil
}

// CleanupDuplicateContents 清理重复内容
func (s *clipboardService) CleanupDuplicateContents(channelID string) (int64, error) {
	return s.clipboardRepo.CleanupDuplicateContents(channelID)
}

// --- model → output converters ---

func toClipboardItemOutput(item *model.ClipboardItem) *service.ClipboardItemOutput {
	if item == nil {
		return nil
	}
	return &service.ClipboardItemOutput{
		ID:            item.ID,
		ChannelID:     item.ChannelID,
		Content:       item.Content,
		ContentHTML:   item.ContentHTML,
		ContentFormat: item.ContentFormat,
		Type:          item.Type,
		Title:         item.Title,
		DeviceID:      item.DeviceID,
		DeviceType:    item.DeviceType,
		Favorite:      item.Favorite,
		CreatedAt:     item.CreatedAt,
		UpdatedAt:     item.UpdatedAt,
	}
}

func toClipboardItemOutputs(items []*model.ClipboardItem) []*service.ClipboardItemOutput {
	result := make([]*service.ClipboardItemOutput, 0, len(items))
	for _, item := range items {
		result = append(result, toClipboardItemOutput(item))
	}
	return result
}
