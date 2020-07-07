package components

import (
	"github.com/GoAdminGroup/go-admin/modules/menu"
	"github.com/GoAdminGroup/go-admin/template/types"
	"html/template"
)

type TreeAttribute struct {
	Name      string
	Tree      []menu.Item
	EditUrl   string
	DeleteUrl string
	UrlPrefix string
	OrderUrl  string
	types.Attribute
}

// �N�Ѽƭȳ]�m��TreeAttribute(struct)
func (compo *TreeAttribute) SetTree(value []menu.Item) types.TreeAttribute {
	compo.Tree = value
	return compo
}

// �N�Ѽƭȳ]�m��TreeAttribute(struct)
func (compo *TreeAttribute) SetEditUrl(value string) types.TreeAttribute {
	compo.EditUrl = value
	return compo
}

// �N�Ѽƭȳ]�m��TreeAttribute(struct)
func (compo *TreeAttribute) SetUrlPrefix(value string) types.TreeAttribute {
	compo.UrlPrefix = value
	return compo
}

// �N�Ѽƭȳ]�m��TreeAttribute(struct)
func (compo *TreeAttribute) SetDeleteUrl(value string) types.TreeAttribute {
	compo.DeleteUrl = value
	return compo
}

// �N�Ѽƭȳ]�m��TreeAttribute(struct)
func (compo *TreeAttribute) SetOrderUrl(value string) types.TreeAttribute {
	compo.OrderUrl = value
	return compo
}

// �����N�ŦXTreeAttribute.TemplateList["components/tree"](map[string]string)���ȥ[�Jtext(string)�A���۱N�ѼƤΥ\��K�[���s���ҪO�øѪR���ҪO���D��
// �N�Ѽ�compo�g�Jbuffer(bytes.Buffer)���̫��XHTML
func (compo *TreeAttribute) GetContent() template.HTML {
	// �btemplate\components\composer.go
	// �����N�ŦXTreeAttribute.TemplateList["components/tree"](map[string]string)���ȥ[�Jtext(string)�A���۱N�ѼƤΥ\��K�[���s���ҪO�øѪR���ҪO���D��
	// �N�Ѽ�compo�g�Jbuffer(bytes.Buffer)���̫��XHTML
	return ComposeHtml(compo.TemplateList, *compo, "tree")
}

// �����N�ŦXTreeAttribute.TemplateList["components/tree-header"](map[string]string)���ȥ[�Jtext(string)�A���۱N�ѼƤΥ\��K�[���s���ҪO�øѪR���ҪO���D��
// �N�Ѽ�compo�g�Jbuffer(bytes.Buffer)���̫��XHTML
func (compo *TreeAttribute) GetTreeHeader() template.HTML {
	// �btemplate\components\composer.go
	// �����N�ŦXTreeAttribute.TemplateList["components/tree-header"](map[string]string)���ȥ[�Jtext(string)�A���۱N�ѼƤΥ\��K�[���s���ҪO�øѪR���ҪO���D��
	// �N�Ѽ�compo�g�Jbuffer(bytes.Buffer)���̫��XHTML
	return ComposeHtml(compo.TemplateList, *compo, "tree-header")
}
