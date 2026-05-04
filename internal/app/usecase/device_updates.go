package usecase

import "time"

// --- device table partial update ---

type devicePatch struct {
	name     *string
	typ      *string
	lastSeen *time.Time
	isOnline *bool
}

func newDevicePatch() *devicePatch { return &devicePatch{} }

func (p *devicePatch) withName(v string) *devicePatch      { p.name = &v; return p }
func (p *devicePatch) withType(v string) *devicePatch       { p.typ = &v; return p }
func (p *devicePatch) withLastSeen(v time.Time) *devicePatch { p.lastSeen = &v; return p }
func (p *devicePatch) withIsOnline(v bool) *devicePatch     { p.isOnline = &v; return p }

func (p *devicePatch) toMap() map[string]any {
	m := make(map[string]any)
	if p.name != nil {
		m["name"] = *p.name
	}
	if p.typ != nil {
		m["type"] = *p.typ
	}
	if p.lastSeen != nil {
		m["last_seen"] = *p.lastSeen
	}
	if p.isOnline != nil {
		m["is_online"] = *p.isOnline
	}
	m["updated_at"] = time.Now()
	return m
}

// --- device_channel table partial update ---

type deviceChannelPatch struct {
	isActive   *bool
	lastSeenAt *time.Time
}

func newDeviceChannelPatch() *deviceChannelPatch { return &deviceChannelPatch{} }

func (p *deviceChannelPatch) withIsActive(v bool) *deviceChannelPatch       { p.isActive = &v; return p }
func (p *deviceChannelPatch) withLastSeenAt(v time.Time) *deviceChannelPatch { p.lastSeenAt = &v; return p }

func (p *deviceChannelPatch) toMap() map[string]any {
	m := make(map[string]any)
	if p.isActive != nil {
		m["is_active"] = *p.isActive
	}
	if p.lastSeenAt != nil {
		m["last_seen_at"] = *p.lastSeenAt
	}
	m["updated_at"] = time.Now()
	return m
}
