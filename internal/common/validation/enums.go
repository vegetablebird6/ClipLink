package validation

var validClipboardTypes = map[string]bool{
	"text": true, "link": true, "code": true,
	"password": true, "image": true, "file": true,
}

var validDeviceTypes = map[string]bool{
	"phone": true, "tablet": true, "desktop": true, "other": true,
}

var validContentFormats = map[string]bool{
	"plain": true, "html": true,
}

// IsValidClipboardType 校验剪贴板内容类型是否合法
func IsValidClipboardType(t string) bool {
	return validClipboardTypes[t]
}

// IsValidDeviceType 校验设备类型是否合法
func IsValidDeviceType(t string) bool {
	return validDeviceTypes[t]
}

// IsValidContentFormat 校验内容格式是否合法，空值视为默认 plain
func IsValidContentFormat(f string) bool {
	if f == "" {
		return true
	}
	return validContentFormats[f]
}
