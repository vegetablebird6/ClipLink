package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"strings"
	stdtime "time"

	"github.com/google/uuid"

	"github.com/xiaojiu/cliplink/internal/common/validation"
	"github.com/xiaojiu/cliplink/internal/domain/model"
	"github.com/xiaojiu/cliplink/internal/domain/repository"
	"github.com/xiaojiu/cliplink/internal/domain/service"
)

func computeContentHash(content string) string {
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return ""
	}
	hash := sha256.Sum256([]byte(trimmed))
	return hex.EncodeToString(hash[:])
}

type clipboardService struct {
	clipboardRepo repository.ClipboardRepository
	syncEventRepo repository.SyncEventRepository
	deviceRepo    repository.DeviceRepository
}

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

func (s *clipboardService) requireActorDevice(ctx context.Context, deviceID, channelID string) (*model.Device, error) {
	if deviceID == "" {
		return nil, model.ErrInvalidInput
	}

	device, err := s.deviceRepo.FindByIDAndChannel(ctx, deviceID, channelID)
	if err != nil || device == nil {
		return nil, model.ErrInvalidInput
	}

	return device, nil
}

func (s *clipboardService) recordSyncEvent(ctx context.Context, event *model.SyncEvent) {
	if err := s.syncEventRepo.Save(ctx, event); err != nil {
		log.Printf("[clipboard] record sync event failed: action=%s target_id=%s err=%v", event.Action, event.TargetID, err)
	}
}

func (s *clipboardService) CreateClipboard(ctx context.Context, in service.CreateClipboardInput) (*service.ClipboardItemOutput, error) {
	if !validation.IsValidClipboardType(in.Type) {
		return nil, model.ErrInvalidInput
	}
	if !validation.IsValidDeviceType(in.ActorDeviceType) {
		return nil, model.ErrInvalidInput
	}
	if !validation.IsValidContentFormat(in.ContentFormat) {
		return nil, model.ErrInvalidInput
	}

	device, err := s.requireActorDevice(ctx, in.ActorDeviceID, in.ChannelID)
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

	if err := s.clipboardRepo.Save(ctx, item); err != nil {
		return nil, err
	}

	if in.CleanDuplicates {
		if item.ContentHash != "" {
			if _, err := s.clipboardRepo.DeleteByContentHash(ctx, in.ChannelID, item.ContentHash, item.ID); err != nil {
				return nil, err
			}
		}
	}

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
	s.recordSyncEvent(ctx, syncEvent)

	return toClipboardItemOutput(item), nil
}

func (s *clipboardService) GetLatestClipboard(ctx context.Context, channelID string, limit int) ([]*service.ClipboardItemOutput, error) {
	items, err := s.clipboardRepo.FindLatest(ctx, channelID, limit)
	if err != nil {
		return nil, err
	}
	return toClipboardItemOutputs(items), nil
}

func (s *clipboardService) GetClipboardItem(ctx context.Context, id, channelID string) (*service.ClipboardItemOutput, error) {
	item, err := s.clipboardRepo.FindByID(ctx, id, channelID)
	if err != nil {
		return nil, err
	}
	return toClipboardItemOutput(item), nil
}

func (s *clipboardService) GetClipboardHistory(ctx context.Context, channelID string, afterCreatedAt *stdtime.Time, afterID *string, size int) ([]*service.ClipboardItemOutput, error) {
	items, err := s.clipboardRepo.FindWithKeyset(ctx, channelID, afterCreatedAt, afterID, size)
	if err != nil {
		return nil, err
	}
	return toClipboardItemOutputs(items), nil
}

func (s *clipboardService) DeleteClipboard(ctx context.Context, in service.DeleteClipboardInput) error {
	device, err := s.requireActorDevice(ctx, in.ActorDeviceID, in.ChannelID)
	if err != nil {
		return err
	}

	item, err := s.clipboardRepo.FindByID(ctx, in.ID, in.ChannelID)
	if err != nil {
		return err
	}

	if err := s.clipboardRepo.Delete(ctx, in.ID, in.ChannelID); err != nil {
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

	s.recordSyncEvent(ctx, syncEvent)
	return nil
}

func (s *clipboardService) UpdateClipboard(ctx context.Context, in service.UpdateClipboardInput) (*service.ClipboardItemOutput, error) {
	device, err := s.requireActorDevice(ctx, in.ActorDeviceID, in.ChannelID)
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

	p := newClipboardPatch()
	if in.Title != nil {
		p.withTitle(*in.Title)
	}
	if in.Content != nil {
		p.withContent(*in.Content)
	}
	if in.Type != nil {
		p.withType(*in.Type)
	}
	if in.DeviceType != nil {
		p.withDeviceType(*in.DeviceType)
	}
	if in.ContentHTML != nil {
		p.withContentHTML(*in.ContentHTML)
	}
	if in.ContentFormat != nil {
		p.withContentFormat(*in.ContentFormat)
	}
	updates := p.toMap()

	if err := s.clipboardRepo.Update(ctx, in.ID, in.ChannelID, updates); err != nil {
		return nil, err
	}

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

	s.recordSyncEvent(ctx, syncEvent)

	item, err := s.clipboardRepo.FindByID(ctx, in.ID, in.ChannelID)
	if err != nil {
		return nil, err
	}
	return toClipboardItemOutput(item), nil
}

func (s *clipboardService) SetFavorite(ctx context.Context, in service.SetFavoriteInput) (*service.ClipboardItemOutput, error) {
	device, err := s.requireActorDevice(ctx, in.ActorDeviceID, in.ChannelID)
	if err != nil {
		return nil, err
	}

	item, err := s.clipboardRepo.FindByID(ctx, in.ID, in.ChannelID)
	if err != nil {
		return nil, err
	}

	updates := newClipboardPatch().
		withFavorite(in.Favorite).
		toMap()

	if err := s.clipboardRepo.Update(ctx, in.ID, in.ChannelID, updates); err != nil {
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
	s.recordSyncEvent(ctx, syncEvent)

	updated, err := s.clipboardRepo.FindByID(ctx, in.ID, in.ChannelID)
	if err != nil {
		return nil, err
	}
	return toClipboardItemOutput(updated), nil
}

func (s *clipboardService) GetFavoriteClipboard(ctx context.Context, channelID string, afterCreatedAt *stdtime.Time, afterID *string, size int) ([]*service.ClipboardItemOutput, error) {
	items, err := s.clipboardRepo.FindFavorites(ctx, channelID, afterCreatedAt, afterID, size)
	if err != nil {
		return nil, err
	}
	return toClipboardItemOutputs(items), nil
}

func (s *clipboardService) GetClipboardByType(ctx context.Context, contentType string, channelID string, afterCreatedAt *stdtime.Time, afterID *string, size int) ([]*service.ClipboardItemOutput, error) {
	items, err := s.clipboardRepo.FindByType(ctx, contentType, channelID, afterCreatedAt, afterID, size)
	if err != nil {
		return nil, err
	}
	return toClipboardItemOutputs(items), nil
}

func (s *clipboardService) GetClipboardByDeviceType(ctx context.Context, deviceType string, channelID string, afterCreatedAt *stdtime.Time, afterID *string, size int) ([]*service.ClipboardItemOutput, error) {
	items, err := s.clipboardRepo.FindByDeviceType(ctx, deviceType, channelID, afterCreatedAt, afterID, size)
	if err != nil {
		return nil, err
	}
	return toClipboardItemOutputs(items), nil
}

func (s *clipboardService) GetClipboardByTypeAndDeviceType(ctx context.Context, contentType, deviceType string, channelID string, afterCreatedAt *stdtime.Time, afterID *string, size int) ([]*service.ClipboardItemOutput, error) {
	items, err := s.clipboardRepo.FindByTypeAndDeviceType(ctx, contentType, deviceType, channelID, afterCreatedAt, afterID, size)
	if err != nil {
		return nil, err
	}
	return toClipboardItemOutputs(items), nil
}

func (s *clipboardService) SearchClipboard(ctx context.Context, keyword, channelID string, page, size int) (items []*service.ClipboardItemOutput, total int64, totalPages int, err error) {
	if keyword == "" {
		return []*service.ClipboardItemOutput{}, 0, 0, nil
	}

	modelItems, total, totalPages, err := s.clipboardRepo.SearchByKeyword(ctx, keyword, channelID, page, size)
	if err != nil {
		return nil, 0, 0, err
	}
	return toClipboardItemOutputs(modelItems), total, totalPages, nil
}

func (s *clipboardService) CleanupDuplicateContents(ctx context.Context, channelID string) (int64, error) {
	return s.clipboardRepo.CleanupDuplicateContents(ctx, channelID)
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
