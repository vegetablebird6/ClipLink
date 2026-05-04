package usecase

import "time"

// clipboardPatch 剪贴板表部分更新
type clipboardPatch struct {
	title         *string
	content       *string
	typ           *string
	deviceType    *string
	contentHTML   *string
	contentFormat *string
	favorite      *bool
}

func newClipboardPatch() *clipboardPatch { return &clipboardPatch{} }

func (p *clipboardPatch) withTitle(v string) *clipboardPatch         { p.title = &v; return p }
func (p *clipboardPatch) withContent(v string) *clipboardPatch       { p.content = &v; return p }
func (p *clipboardPatch) withType(v string) *clipboardPatch          { p.typ = &v; return p }
func (p *clipboardPatch) withDeviceType(v string) *clipboardPatch    { p.deviceType = &v; return p }
func (p *clipboardPatch) withContentHTML(v string) *clipboardPatch   { p.contentHTML = &v; return p }
func (p *clipboardPatch) withContentFormat(v string) *clipboardPatch { p.contentFormat = &v; return p }
func (p *clipboardPatch) withFavorite(v bool) *clipboardPatch        { p.favorite = &v; return p }

func (p *clipboardPatch) toMap() map[string]any {
	m := make(map[string]any)
	if p.title != nil {
		m["title"] = *p.title
	}
	if p.content != nil {
		m["content"] = *p.content
		m["content_hash"] = computeContentHash(*p.content)
	}
	if p.typ != nil {
		m["type"] = *p.typ
	}
	if p.deviceType != nil {
		m["device_type"] = *p.deviceType
	}
	if p.contentHTML != nil {
		m["content_html"] = *p.contentHTML
	}
	if p.contentFormat != nil {
		m["content_format"] = *p.contentFormat
	}
	if p.favorite != nil {
		m["favorite"] = *p.favorite
	}
	m["updated_at"] = time.Now()
	return m
}
