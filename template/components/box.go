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

// �N�ѼƳ]�m��BoxAttribute(struct)
func (compo *BoxAttribute) SetTheme(value string) types.BoxAttribute {
	compo.Theme = value
	return compo
}

// �N�ѼƳ]�m��BoxAttribute(struct)
func (compo *BoxAttribute) SetHeader(value template.HTML) types.BoxAttribute {
	compo.Header = value
	return compo
}

// �N�ѼƳ]�m��BoxAttribute(struct)
func (compo *BoxAttribute) SetBody(value template.HTML) types.BoxAttribute {
	compo.Body = value
	return compo
}

// �N�ѼƳ]�m��BoxAttribute(struct)
func (compo *BoxAttribute) SetStyle(value template.HTMLAttr) types.BoxAttribute {
	compo.Style = value
	return compo
}

// �N�ѼƳ]�m��BoxAttribute(struct)
func (compo *BoxAttribute) SetFooter(value template.HTML) types.BoxAttribute {
	compo.Footer = value
	return compo
}

// �N�ѼƳ]�m��BoxAttribute(struct)
func (compo *BoxAttribute) SetTitle(value template.HTML) types.BoxAttribute {
	compo.Title = value
	return compo
}

// �N�ѼƳ]�m��BoxAttribute(struct)
func (compo *BoxAttribute) SetHeadColor(value string) types.BoxAttribute {
	compo.HeadColor = value
	return compo
}

// �Nwith-border�]�m��BoxAttribute.HeadBorder(struct)
func (compo *BoxAttribute) WithHeadBorder() types.BoxAttribute {
	compo.HeadBorder = "with-border"
	return compo
}

// �N�ѼƳ]�m��BoxAttribute(struct)
func (compo *BoxAttribute) SetSecondHeader(value template.HTML) types.BoxAttribute {
	compo.SecondHeader = value
	return compo
}

// �N�ѼƳ]�m��BoxAttribute(struct)
func (compo *BoxAttribute) SetSecondHeadColor(value string) types.BoxAttribute {
	compo.SecondHeadColor = value
	return compo
}

// �N�ѼƳ]�m��BoxAttribute(struct)
func (compo *BoxAttribute) SetSecondHeaderClass(value string) types.BoxAttribute {
	compo.SecondHeaderClass = value
	return compo
}

// �Npadding:0;�]�m��BoxAttribute(struct).Padding
func (compo *BoxAttribute) SetNoPadding() types.BoxAttribute {
	compo.Padding = "padding:0;"
	return compo
}

// �Nwith-border�]�m��BoxAttribute(struct).SecondHeadBorder
func (compo *BoxAttribute) WithSecondHeadBorder() types.BoxAttribute {
	compo.SecondHeadBorder = "with-border"
	return compo
}

// ���̧P�_����]�mBoxAttribute.Style
// �N�ŦXBoxAttribute.TemplateList["box"](map[string]string)���ȥ[�Jtext(string)�A���۱N�ѼƤΥ\��K�[���s���ҪO�øѪR���ҪO���D��
// �N�Ѽ�compo�g�Jbuffer(bytes.Buffer)���̫��XHTML
func (compo *BoxAttribute) GetContent() template.HTML {

	// �̧P�_����]�mBoxAttribute.Style
	if compo.Style == "" {
		compo.Style = template.HTMLAttr(fmt.Sprintf(`style="overflow-x: scroll;overflow-y: hidden;%s"`, compo.Padding))
	} else {
		compo.Style = template.HTMLAttr(fmt.Sprintf(`style="%s"`, string(compo.Style)+compo.Padding))
	}

	// �btemplate\components\composer.go
	// �����N�ŦXBoxAttribute.TemplateList["box"](map[string]string)���ȥ[�Jtext(string)�A���۱N�ѼƤΥ\��K�[���s���ҪO�øѪR���ҪO���D��
	// �N�Ѽ�compo�g�Jbuffer(bytes.Buffer)���̫��XHTML
	return ComposeHtml(compo.TemplateList, *compo, "box")
}
