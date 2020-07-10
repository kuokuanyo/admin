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

// �N�Ѽƭȳ]�m��RowAttribute(struct)
func (compo *RowAttribute) SetContent(value template.HTML) types.RowAttribute {
	compo.Content = value
	return compo
}

// �N�Ѽƭȳ]�m��RowAttribute(struct)
func (compo *RowAttribute) AddContent(value template.HTML) types.RowAttribute {
	compo.Content += value
	return compo
}

// �����N�ŦXTreeAttribute.TemplateList["components/tree-header"](map[string]string)���ȥ[�Jtext(string)
// ���۱N�Ѽ�compo�g�Jbuffer(bytes.Buffer)���̫��XHTML
func (compo *RowAttribute) GetContent() template.HTML {
	// �btemplate\components\composer.go
	// �����N�ŦXTreeAttribute.TemplateList["components/row"](map[string]string)���ȥ[�Jtext(string)�A���۱N�ѼƤΥ\��K�[���s���ҪO�øѪR���ҪO���D��
	// �N�Ѽ�compo�g�Jbuffer(bytes.Buffer)���̫��XHTML
	return ComposeHtml(compo.TemplateList, *compo, "row")
}
