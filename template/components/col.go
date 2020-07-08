package components

import (
	"github.com/GoAdminGroup/go-admin/template/types"
	"html/template"
)

type ColAttribute struct {
	Name    string
	Content template.HTML
	Size    string
	types.Attribute
}

// 盢把计砞竚ColAttribute(struct)
func (compo *ColAttribute) SetContent(value template.HTML) types.ColAttribute {
	compo.Content = value
	return compo
}

// 盢把计砞竚ColAttribute(struct)
func (compo *ColAttribute) AddContent(value template.HTML) types.ColAttribute {
	compo.Content += value
	return compo
}

// 盢把计砞竚ColAttribute(struct)
func (compo *ColAttribute) SetSize(value types.S) types.ColAttribute {
	compo.Size = ""
	for key, size := range value {
		compo.Size += "col-" + key + "-" + size + " "
	}
	return compo
}

// ㄌ耞兵ン砞竚ColAttribute.Style
// 盢才ColAttribute.TemplateList["box"](map[string]string)text(string)钡帝盢把计の睰倒穝家狾秆猂家狾砰
// 盢把计compo糶buffer(bytes.Buffer)い程块HTML
func (compo *ColAttribute) GetContent() template.HTML {
	// template\components\composer.go
	// 盢才ColAttribute.TemplateList["col"](map[string]string)text(string)钡帝盢把计の睰倒穝家狾秆猂家狾砰
	// 盢把计compo糶buffer(bytes.Buffer)い程块HTML
	return ComposeHtml(compo.TemplateList, *compo, "col")
}
