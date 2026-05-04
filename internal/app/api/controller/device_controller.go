package controller

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xiaojiu/cliplink/internal/app/api/dto"
	"github.com/xiaojiu/cliplink/internal/common/response"
	"github.com/xiaojiu/cliplink/internal/common/validation"
	"github.com/xiaojiu/cliplink/internal/domain/model"
	"github.com/xiaojiu/cliplink/internal/domain/service"
)

// DeviceController 设备控制器
type DeviceController struct {
	deviceService service.DeviceService
}

// NewDeviceController 创建新的设备控制器
func NewDeviceController(deviceService service.DeviceService) *DeviceController {
	return &DeviceController{
		deviceService: deviceService,
	}
}

// RegisterDevice 注册设备并关联到通道
func (c *DeviceController) RegisterDevice(ctx *gin.Context) {
	channelID, ok := ctx.Get("channelID")
	if !ok {
		response.BadRequestWithCode(ctx, "CHANNEL_ID_REQUIRED", "error.channel_id_required", "")
		return
	}

	var req dto.RegisterDeviceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequestWithCode(ctx, "INVALID_INPUT", "error.invalid_input", "")
		return
	}

	if !validation.IsValidDeviceType(req.DeviceType) {
		response.BadRequestWithCode(ctx, "INVALID_DEVICE_TYPE", "error.invalid_device_type", req.DeviceType)
		return
	}

	rctx := ctx.Request.Context()
	device, err := c.deviceService.RegisterDevice(rctx, req.DeviceName, req.DeviceType, req.DeviceID)
	if err != nil {
		log.Printf("[device] register failed: %v", err)
		response.Error(ctx, err)
		return
	}

	err = c.deviceService.AddDeviceToChannel(rctx, device.ID, channelID.(string))
	if err != nil {
		log.Printf("[device] add to channel failed: %v", err)
		response.Error(ctx, err)
		return
	}

	deviceDTO := dto.ToDeviceResponse(&service.DeviceChannelOutput{
		ID:        device.ID,
		Name:      device.Name,
		Type:      device.Type,
		ChannelID: channelID.(string),
		LastSeen:  device.LastSeen,
		IsOnline:  device.IsOnline,
		CreatedAt: device.CreatedAt,
		JoinedAt:  time.Now(),
	})

	response.Success(ctx, deviceDTO, "设备注册成功")
}

// GetDevices 获取通道下的所有设备
func (c *DeviceController) GetDevices(ctx *gin.Context) {
	channelID, ok := ctx.Get("channelID")
	if !ok {
		response.BadRequestWithCode(ctx, "CHANNEL_ID_REQUIRED", "error.channel_id_required", "")
		return
	}

	devices, err := c.deviceService.GetDevicesByChannel(ctx.Request.Context(), channelID.(string))
	if err != nil {
		log.Printf("[device] get devices failed: %v", err)
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, dto.ToDeviceResponseList(devices), "获取成功")
}

// GetDeviceByID 获取特定设备
func (c *DeviceController) GetDeviceByID(ctx *gin.Context) {
	channelID, ok := ctx.Get("channelID")
	if !ok {
		response.BadRequestWithCode(ctx, "CHANNEL_ID_REQUIRED", "error.channel_id_required", "")
		return
	}
	deviceID := ctx.Param("deviceID")

	device, err := c.deviceService.GetDeviceInChannel(ctx.Request.Context(), deviceID, channelID.(string))
	if err != nil {
		if err == model.ErrDeviceNotFound {
			response.NotFound(ctx, "device not found")
			return
		}
		log.Printf("[device] get device failed: %v", err)
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, dto.ToDeviceResponse(device), "获取成功")
}

// UpdateDeviceStatus 更新设备状态
func (c *DeviceController) UpdateDeviceStatus(ctx *gin.Context) {
	channelID, ok := ctx.Get("channelID")
	if !ok {
		response.BadRequestWithCode(ctx, "CHANNEL_ID_REQUIRED", "error.channel_id_required", "")
		return
	}
	deviceID := ctx.Param("deviceID")

	var req dto.UpdateDeviceStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequestWithCode(ctx, "INVALID_INPUT", "error.invalid_input", "")
		return
	}

	rctx := ctx.Request.Context()
	device, err := c.deviceService.UpdateDeviceStatus(rctx, deviceID, *req.IsOnline)
	if err != nil {
		log.Printf("[device] update status failed: %v", err)
		response.Error(ctx, err)
		return
	}

	err = c.deviceService.UpdateDeviceInChannel(rctx, deviceID, channelID.(string), *req.IsOnline)
	if err != nil {
		// 忽略通道关联错误
	}

	deviceDTO, err := c.deviceService.GetDeviceInChannel(rctx, deviceID, channelID.(string))
	if err != nil {
		deviceDTO = &service.DeviceChannelOutput{
			ID:        device.ID,
			Name:      device.Name,
			Type:      device.Type,
			ChannelID: channelID.(string),
			LastSeen:  device.LastSeen,
			IsOnline:  device.IsOnline,
			CreatedAt: device.CreatedAt,
		}
	}

	response.Success(ctx, dto.ToDeviceResponse(deviceDTO), "设备状态已更新")
}

// UpdateDeviceName 更新设备名称
func (c *DeviceController) UpdateDeviceName(ctx *gin.Context) {
	channelID, ok := ctx.Get("channelID")
	if !ok {
		response.BadRequestWithCode(ctx, "CHANNEL_ID_REQUIRED", "error.channel_id_required", "")
		return
	}
	deviceID := ctx.Param("deviceID")

	rctx := ctx.Request.Context()
	inChannel, err := c.deviceService.IsDeviceInChannel(rctx, deviceID, channelID.(string))
	if err != nil {
		log.Printf("[device] check in channel failed: %v", err)
		response.Error(ctx, err)
		return
	}
	if !inChannel {
		response.NotFound(ctx, "device not found in channel")
		return
	}

	var req dto.UpdateDeviceNameRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequestWithCode(ctx, "INVALID_INPUT", "error.invalid_input", "")
		return
	}

	if _, err := c.deviceService.UpdateDevice(rctx, deviceID, req.Name, ""); err != nil {
		log.Printf("[device] update name failed: %v", err)
		response.Error(ctx, err)
		return
	}

	deviceDTO, err := c.deviceService.GetDeviceInChannel(rctx, deviceID, channelID.(string))
	if err != nil {
		log.Printf("[device] get device after update failed: %v", err)
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, dto.ToDeviceResponse(deviceDTO), "设备名称已更新")
}

// RemoveDevice 移除设备
func (c *DeviceController) RemoveDevice(ctx *gin.Context) {
	channelID, ok := ctx.Get("channelID")
	if !ok {
		response.BadRequestWithCode(ctx, "CHANNEL_ID_REQUIRED", "error.channel_id_required", "")
		return
	}
	deviceID := ctx.Param("deviceID")

	err := c.deviceService.RemoveDeviceFromChannel(ctx.Request.Context(), deviceID, channelID.(string))
	if err != nil {
		log.Printf("[device] remove from channel failed: %v", err)
		response.Error(ctx, err)
		return
	}

	response.SuccessWithMessage(ctx, "device removed from channel")
}
