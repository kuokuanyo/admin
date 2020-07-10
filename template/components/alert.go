package components

import (
	"github.com/GoAdminGroup/go-admin/modules/errors"
	"github.com/GoAdminGroup/go-admin/modules/language"
	"github.com/GoAdminGroup/go-admin/template/types"
	"html/template"
)

type AlertAttribute struct {
	Name    string
	Theme   string
	Title   template.HTML
	Content template.HTML
	types.Attribute
}

// 盢把计砞竚AlertAttribute(struct)
func (compo *AlertAttribute) SetTheme(value string) types.AlertAttribute {
	compo.Theme = value
	return compo
}

// 盢把计砞竚AlertAttribute(struct)
func (compo *AlertAttribute) SetTitle(value template.HTML) types.AlertAttribute {
	compo.Title = value
	return compo
}

// 盢把计砞竚AlertAttribute(struct)
func (compo *AlertAttribute) SetContent(value template.HTML) types.AlertAttribute {
	compo.Content = value
	return compo
}

// 盢把计砞竚AlertAttribute(struct)钡帝盢才AlertAttribute.TemplateList["components/alert"](map[string]string)text(string)钡帝盢把计の睰倒穝家狾秆猂家狾砰
// 盢把计compo糶buffer(bytes.Buffer)い程块HTML
func (compo *AlertAttribute) Warning(msg string) template.HTML {
	// SetTitleSetThemeSetContent盢把计砞竚AlertAttribute(struct)
	// GetContent盢才AlertAttribute.TemplateList["components/alert"](map[string]string)text(string)钡帝盢把计の睰倒穝家狾秆猂家狾砰
	// 盢把计compo糶buffer(bytes.Buffer)い程块HTML
	return compo.SetTitle(errors.MsgWithIcon).
		SetTheme("warning").
		SetContent(language.GetFromHtml(template.HTML(msg))).
		GetContent()
}

// 盢才AlertAttribute.TemplateList["components/alert"](map[string]string)text(string)钡帝盢把计の睰倒穝家狾秆猂家狾砰
// 盢把计compo糶buffer(bytes.Buffer)い程块HTML
func (compo *AlertAttribute) GetContent() template.HTML {
	// template\components\composer.go
	// 盢才AlertAttribute.TemplateList["components/alert"](map[string]string)text(string)钡帝盢把计の睰倒穝家狾秆猂家狾砰
	// 盢把计compo糶buffer(bytes.Buffer)い程块HTML
	return ComposeHtml(compo.TemplateList, *compo, "alert")
}
