package components

import (
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/menu"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"html/template"
)

type Base struct {
	Attribute types.Attribute
}

// �إ�BoxAttribute(struct)�ó]�m�ȫ�^��
func (b Base) Box() types.BoxAttribute {
	return &BoxAttribute{
		Name:       "box",
		Header:     template.HTML(""),
		Body:       template.HTML(""),
		Footer:     template.HTML(""),
		Title:      "",
		HeadBorder: "",
		Attribute:  b.Attribute,
	}
}

// �إ�ColAttribute(struct)�ó]�m�ȫ�^��
func (b Base) Col() types.ColAttribute {
	return &ColAttribute{
		Name:      "col",
		Size:      "col-md-2",
		Content:   "",
		Attribute: b.Attribute,
	}
}

// �إ�FormAttribute(struct)�ó]�m�ȫ�^��
func (b Base) Form() types.FormAttribute {
	return &FormAttribute{
		Name:         "form",
		Content:      []types.FormField{},
		Url:          "/",
		Method:       "post",
		HiddenFields: make(map[string]string),
		Layout:       form.LayoutDefault,
		Title:        "edit",
		Attribute:    b.Attribute,
		CdnUrl:       config.GetAssetUrl(),
		HeadWidth:    2,
		InputWidth:   8,
	}
}

// �إ�ImgAttribute(struct)�ó]�m�ȫ�^��
func (b Base) Image() types.ImgAttribute {
	return &ImgAttribute{
		Name:      "image",
		Width:     "50",
		Height:    "50",
		Src:       "",
		Attribute: b.Attribute,
	}
}

// �إ�TabsAttribute(struct)�ó]�m�ȫ�^��
func (b Base) Tabs() types.TabsAttribute {
	return &TabsAttribute{
		Name:      "tabs",
		Attribute: b.Attribute,
	}
}

// �إ�AlertAttribute(struct)�ó]�m�ȫ�^��
func (b Base) Alert() types.AlertAttribute {
	return &AlertAttribute{
		Name:      "alert",
		Attribute: b.Attribute,
	}
}

// �إ�LabelAttribute(struct)�ó]�m�ȫ�^��
func (b Base) Label() types.LabelAttribute {
	return &LabelAttribute{
		Name:      "label",
		Type:      "",
		Content:   "",
		Attribute: b.Attribute,
	}
}

// �إ�LinkAttribute(struct)�ó]�m�ȫ�^��
func (b Base) Link() types.LinkAttribute {
	return &LinkAttribute{
		Name:      "link",
		Content:   "",
		Attribute: b.Attribute,
	}
}

// �إ�PopupAttribute(struct)�ó]�m�ȫ�^��
func (b Base) Popup() types.PopupAttribute {
	return &PopupAttribute{
		Name:      "popup",
		Attribute: b.Attribute,
	}
}

// �إ�PaginatorAttribute(struct)�ó]�m�ȫ�^��
func (b Base) Paginator() types.PaginatorAttribute {
	return &PaginatorAttribute{
		Name:      "paginator",
		Attribute: b.Attribute,
	}
}

// �إ�RowAttribute(struct)�ó]�m�ȫ�^��
func (b Base) Row() types.RowAttribute {
	return &RowAttribute{
		Name:      "row",
		Content:   "",
		Attribute: b.Attribute,
	}
}

// �إ�ButtonAttribute(struct)�ó]�m�ȫ�^��
func (b Base) Button() types.ButtonAttribute {
	return &ButtonAttribute{
		Name:      "button",
		Content:   "",
		Href:      "",
		Attribute: b.Attribute,
	}
}

// �إ�TableAttribute(struct)�ó]�m�ȫ�^��
func (b Base) Table() types.TableAttribute {
	return &TableAttribute{
		Name:      "table",
		Thead:     make(types.Thead, 0),
		InfoList:  make([]map[string]types.InfoItem, 0),
		Type:      "normal",
		Style:     "hover",
		Layout:    "auto",
		Attribute: b.Attribute,
	}
}

// �إ�DataTableAttribute(struct)�ó]�m�ȫ�^��
func (b Base) DataTable() types.DataTableAttribute {
	return &DataTableAttribute{
		TableAttribute: *(b.Table().
			SetType("data-table").(*TableAttribute)),
		EditUrl:   "",
		NewUrl:    "",
		Style:     "hover",
		Attribute: b.Attribute,
	}
}

// �إ�TreeAttribute(struct)�ó]�m�ȫ�^��
func (b Base) Tree() types.TreeAttribute {
	return &TreeAttribute{
		Name:      "tree",
		Tree:      []menu.Item{},
		Attribute: b.Attribute,
	}
}
