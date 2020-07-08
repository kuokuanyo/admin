package components

import (
	"github.com/GoAdminGroup/go-admin/template/types"
	"html/template"
)

type RowAttribute struct {
	Name    string
	Content template.HTML
	types.Attribute
}

// 盢把计砞竚RowAttribute(struct)
func (compo *RowAttribute) SetContent(value template.HTML) types.RowAttribute {
	compo.Content = value
	return compo
}

// 盢把计砞竚RowAttribute(struct)
func (compo *RowAttribute) AddContent(value template.HTML) types.RowAttribute {
	compo.Content += value
	return compo
}

// template\components\composer.go
// 盢才RowAttribute.TemplateList["components/row"](map[string]string)text(string)钡帝盢把计の睰倒穝家狾秆猂家狾砰
// 盢把计compo糶buffer(bytes.Buffer)い程块HTML
func (compo *RowAttribute) GetContent() template.HTML {
	// template\components\composer.go
	// 盢才RowAttribute.TemplateList["components/row"](map[string]string)text(string)钡帝盢把计の睰倒穝家狾秆猂家狾砰
	// 盢把计compo糶buffer(bytes.Buffer)い程块HTML
	return ComposeHtml(compo.TemplateList, *compo, "row")
}
