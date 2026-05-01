package validation

var validClipboardTypes = map[string]bool{
	"text": true, "link": true, "code": true,
	"password": true, "image": true, "file": true,
}

var validDeviceTypes = map[string]bool{
	"phone": true, "tablet": true, "desktop": true, "other": true,
}

// IsValidClipboardType 校验剪贴板内容类型是否合法
func IsValidClipboardType(t string) bool {
	return validClipboardTypes[t]
}

// IsValidDeviceType 校验设备类型是否合法
func IsValidDeviceType(t string) bool {
	return validDeviceTypes[t]
}
