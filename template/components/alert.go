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

// �N�Ѽƭȳ]�m��AlertAttribute(struct)
func (compo *AlertAttribute) SetTheme(value string) types.AlertAttribute {
	compo.Theme = value
	return compo
}

// �N�Ѽƭȳ]�m��AlertAttribute(struct)
func (compo *AlertAttribute) SetTitle(value template.HTML) types.AlertAttribute {
	compo.Title = value
	return compo
}

// �N�Ѽƭȳ]�m��AlertAttribute(struct)
func (compo *AlertAttribute) SetContent(value template.HTML) types.AlertAttribute {
	compo.Content = value
	return compo
}

// �����N�ѼƳ]�m��AlertAttribute(struct)��A���۱N�ŦXAlertAttribute.TemplateList["components/alert"](map[string]string)���ȥ[�Jtext(string)�A���۱N�ѼƤΥ\��K�[���s���ҪO�øѪR���ҪO���D��
// �N�Ѽ�compo�g�Jbuffer(bytes.Buffer)���̫��XHTML
func (compo *AlertAttribute) Warning(msg string) template.HTML {
	// SetTitle�BSetTheme�BSetContent�N�ѼƳ]�m��AlertAttribute(struct)��
	// GetContent�����N�ŦXAlertAttribute.TemplateList["components/alert"](map[string]string)���ȥ[�Jtext(string)�A���۱N�ѼƤΥ\��K�[���s���ҪO�øѪR���ҪO���D��
	// �N�Ѽ�compo�g�Jbuffer(bytes.Buffer)���̫��XHTML
	return compo.SetTitle(errors.MsgWithIcon).
		SetTheme("warning").
		SetContent(language.GetFromHtml(template.HTML(msg))).
		GetContent()
}

// �����N�ŦXAlertAttribute.TemplateList["components/alert"](map[string]string)���ȥ[�Jtext(string)�A���۱N�ѼƤΥ\��K�[���s���ҪO�øѪR���ҪO���D��
// �N�Ѽ�compo�g�Jbuffer(bytes.Buffer)���̫��XHTML
func (compo *AlertAttribute) GetContent() template.HTML {
	// �btemplate\components\composer.go
	// �����N�ŦXAlertAttribute.TemplateList["components/alert"](map[string]string)���ȥ[�Jtext(string)�A���۱N�ѼƤΥ\��K�[���s���ҪO�øѪR���ҪO���D��
	// �N�Ѽ�compo�g�Jbuffer(bytes.Buffer)���̫��XHTML
	return ComposeHtml(compo.TemplateList, *compo, "alert")
}
