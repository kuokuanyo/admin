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

// 盢把计砞竚TreeAttribute(struct)
func (compo *TreeAttribute) SetTree(value []menu.Item) types.TreeAttribute {
	compo.Tree = value
	return compo
}

// 盢把计砞竚TreeAttribute(struct)
func (compo *TreeAttribute) SetEditUrl(value string) types.TreeAttribute {
	compo.EditUrl = value
	return compo
}

// 盢把计砞竚TreeAttribute(struct)
func (compo *TreeAttribute) SetUrlPrefix(value string) types.TreeAttribute {
	compo.UrlPrefix = value
	return compo
}

// 盢把计砞竚TreeAttribute(struct)
func (compo *TreeAttribute) SetDeleteUrl(value string) types.TreeAttribute {
	compo.DeleteUrl = value
	return compo
}

// 盢把计砞竚TreeAttribute(struct)
func (compo *TreeAttribute) SetOrderUrl(value string) types.TreeAttribute {
	compo.OrderUrl = value
	return compo
}

// 盢才TreeAttribute.TemplateList["components/tree"](map[string]string)text(string)钡帝盢把计の睰倒穝家狾秆猂家狾砰
// 盢把计compo糶buffer(bytes.Buffer)い程块HTML
func (compo *TreeAttribute) GetContent() template.HTML {
	// template\components\composer.go
	// 盢才TreeAttribute.TemplateList["components/tree"](map[string]string)text(string)钡帝盢把计の睰倒穝家狾秆猂家狾砰
	// 盢把计compo糶buffer(bytes.Buffer)い程块HTML
	return ComposeHtml(compo.TemplateList, *compo, "tree")
}

// 盢才TreeAttribute.TemplateList["components/tree-header"](map[string]string)text(string)钡帝盢把计の睰倒穝家狾秆猂家狾砰
// 盢把计compo糶buffer(bytes.Buffer)い程块HTML
func (compo *TreeAttribute) GetTreeHeader() template.HTML {
	// template\components\composer.go
	// 盢才TreeAttribute.TemplateList["components/tree-header"](map[string]string)text(string)钡帝盢把计の睰倒穝家狾秆猂家狾砰
	// 盢把计compo糶buffer(bytes.Buffer)い程块HTML
	return ComposeHtml(compo.TemplateList, *compo, "tree-header")
}
