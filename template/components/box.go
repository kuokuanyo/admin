package components

import (
	"fmt"
	"github.com/GoAdminGroup/go-admin/template/types"
	"html/template"
)

type BoxAttribute struct {
	Name              string
	Header            template.HTML
	Body              template.HTML
	Footer            template.HTML
	Title             template.HTML
	Theme             string
	HeadBorder        string
	HeadColor         string
	SecondHeaderClass string
	SecondHeader      template.HTML
	SecondHeadBorder  string
	SecondHeadColor   string
	Style             template.HTMLAttr
	Padding           string
	types.Attribute
}

// 盢把计砞竚BoxAttribute(struct)
func (compo *BoxAttribute) SetTheme(value string) types.BoxAttribute {
	compo.Theme = value
	return compo
}

// 盢把计砞竚BoxAttribute(struct)
func (compo *BoxAttribute) SetHeader(value template.HTML) types.BoxAttribute {
	compo.Header = value
	return compo
}

// 盢把计砞竚BoxAttribute(struct)
func (compo *BoxAttribute) SetBody(value template.HTML) types.BoxAttribute {
	compo.Body = value
	return compo
}

// 盢把计砞竚BoxAttribute(struct)
func (compo *BoxAttribute) SetStyle(value template.HTMLAttr) types.BoxAttribute {
	compo.Style = value
	return compo
}

// 盢把计砞竚BoxAttribute(struct)
func (compo *BoxAttribute) SetFooter(value template.HTML) types.BoxAttribute {
	compo.Footer = value
	return compo
}

// 盢把计砞竚BoxAttribute(struct)
func (compo *BoxAttribute) SetTitle(value template.HTML) types.BoxAttribute {
	compo.Title = value
	return compo
}

// 盢把计砞竚BoxAttribute(struct)
func (compo *BoxAttribute) SetHeadColor(value string) types.BoxAttribute {
	compo.HeadColor = value
	return compo
}

// 盢with-border砞竚BoxAttribute.HeadBorder(struct)
func (compo *BoxAttribute) WithHeadBorder() types.BoxAttribute {
	compo.HeadBorder = "with-border"
	return compo
}

// 盢把计砞竚BoxAttribute(struct)
func (compo *BoxAttribute) SetSecondHeader(value template.HTML) types.BoxAttribute {
	compo.SecondHeader = value
	return compo
}

// 盢把计砞竚BoxAttribute(struct)
func (compo *BoxAttribute) SetSecondHeadColor(value string) types.BoxAttribute {
	compo.SecondHeadColor = value
	return compo
}

// 盢把计砞竚BoxAttribute(struct)
func (compo *BoxAttribute) SetSecondHeaderClass(value string) types.BoxAttribute {
	compo.SecondHeaderClass = value
	return compo
}

// 盢padding:0;砞竚BoxAttribute(struct).Padding
func (compo *BoxAttribute) SetNoPadding() types.BoxAttribute {
	compo.Padding = "padding:0;"
	return compo
}

// 盢with-border砞竚BoxAttribute(struct).SecondHeadBorder
func (compo *BoxAttribute) WithSecondHeadBorder() types.BoxAttribute {
	compo.SecondHeadBorder = "with-border"
	return compo
}

// ㄌ耞兵ン砞竚BoxAttribute.Style
// 盢才BoxAttribute.TemplateList["box"](map[string]string)text(string)钡帝盢把计の睰倒穝家狾秆猂家狾砰
// 盢把计compo糶buffer(bytes.Buffer)い程块HTML
func (compo *BoxAttribute) GetContent() template.HTML {

	// ㄌ耞兵ン砞竚BoxAttribute.Style
	if compo.Style == "" {
		compo.Style = template.HTMLAttr(fmt.Sprintf(`style="overflow-x: scroll;overflow-y: hidden;%s"`, compo.Padding))
	} else {
		compo.Style = template.HTMLAttr(fmt.Sprintf(`style="%s"`, string(compo.Style)+compo.Padding))
	}

	// template\components\composer.go
	// 盢才BoxAttribute.TemplateList["box"](map[string]string)text(string)钡帝盢把计の睰倒穝家狾秆猂家狾砰
	// 盢把计compo糶buffer(bytes.Buffer)い程块HTML
	return ComposeHtml(compo.TemplateList, *compo, "box")
}
