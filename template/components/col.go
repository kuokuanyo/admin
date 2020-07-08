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

// �N�ѼƳ]�m��ColAttribute(struct)
func (compo *ColAttribute) SetContent(value template.HTML) types.ColAttribute {
	compo.Content = value
	return compo
}

// �N�ѼƳ]�m��ColAttribute(struct)
func (compo *ColAttribute) AddContent(value template.HTML) types.ColAttribute {
	compo.Content += value
	return compo
}

// �N�ѼƳ]�m��ColAttribute(struct)
func (compo *ColAttribute) SetSize(value types.S) types.ColAttribute {
	compo.Size = ""
	for key, size := range value {
		compo.Size += "col-" + key + "-" + size + " "
	}
	return compo
}

// ���̧P�_����]�mColAttribute.Style
// �N�ŦXColAttribute.TemplateList["box"](map[string]string)���ȥ[�Jtext(string)�A���۱N�ѼƤΥ\��K�[���s���ҪO�øѪR���ҪO���D��
// �N�Ѽ�compo�g�Jbuffer(bytes.Buffer)���̫��XHTML
func (compo *ColAttribute) GetContent() template.HTML {
	// �btemplate\components\composer.go
	// �����N�ŦXColAttribute.TemplateList["col"](map[string]string)���ȥ[�Jtext(string)�A���۱N�ѼƤΥ\��K�[���s���ҪO�øѪR���ҪO���D��
	// �N�Ѽ�compo�g�Jbuffer(bytes.Buffer)���̫��XHTML
	return ComposeHtml(compo.TemplateList, *compo, "col")
}
