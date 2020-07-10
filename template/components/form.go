package components

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/language"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

type FormAttribute struct {
	Name            string
	Header          template.HTML
	Content         types.FormFields
	ContentList     []types.FormFields
	Layout          form.Layout
	TabContents     []types.FormFields
	TabHeaders      []string
	Footer          template.HTML
	Url             string
	Method          string
	PrimaryKey      string
	HeadWidth       int
	InputWidth      int
	HiddenFields    map[string]string
	Title           template.HTML
	OperationFooter template.HTML
	Prefix          string
	CdnUrl          string
	types.Attribute
}

// �N�Ѽƭȳ]�m��FormAttribute(struct)
func (compo *FormAttribute) SetHeader(value template.HTML) types.FormAttribute {
	compo.Header = value
	return compo
}

// �N�Ѽƭȳ]�m��FormAttribute(struct)
func (compo *FormAttribute) SetPrimaryKey(value string) types.FormAttribute {
	compo.PrimaryKey = value
	return compo
}

// �N�Ѽƭȳ]�m��FormAttribute(struct)
func (compo *FormAttribute) SetContent(value types.FormFields) types.FormAttribute {
	compo.Content = value
	return compo
}

// �N�Ѽƭȳ]�m��FormAttribute(struct)
func (compo *FormAttribute) SetTabContents(value []types.FormFields) types.FormAttribute {
	compo.TabContents = value
	return compo
}

// �N�Ѽƭȳ]�m��FormAttribute(struct)
func (compo *FormAttribute) SetTabHeaders(value []string) types.FormAttribute {
	compo.TabHeaders = value
	return compo
}

// �N�Ѽƭȳ]�m��FormAttribute(struct)
func (compo *FormAttribute) SetHeadWidth(width int) types.FormAttribute {
	compo.HeadWidth = width
	return compo
}

// �N�Ѽƭȳ]�m��FormAttribute(struct)
func (compo *FormAttribute) SetInputWidth(width int) types.FormAttribute {
	compo.InputWidth = width
	return compo
}

// �N�Ѽƭȳ]�m��FormAttribute(struct)
func (compo *FormAttribute) SetFooter(value template.HTML) types.FormAttribute {
	compo.Footer = value
	return compo
}

// �N�Ѽƭȳ]�m��FormAttribute(struct)
func (compo *FormAttribute) SetLayout(layout form.Layout) types.FormAttribute {
	compo.Layout = layout
	return compo
}

// �N�Ѽƭȳ]�m��FormAttribute(struct)
func (compo *FormAttribute) SetPrefix(value string) types.FormAttribute {
	compo.Prefix = value
	return compo
}

// �N�Ѽƭȳ]�m��FormAttribute(struct)
func (compo *FormAttribute) SetUrl(value string) types.FormAttribute {
	compo.Url = value
	return compo
}

// �N�Ѽƭȳ]�m��FormAttribute(struct)
func (compo *FormAttribute) SetHiddenFields(fields map[string]string) types.FormAttribute {
	compo.HiddenFields = fields
	return compo
}



// �N�Ѽƭȳ]�m��FormAttribute(struct)
func (compo *FormAttribute) SetMethod(value string) types.FormAttribute {
	compo.Method = value
	return compo
}

// �N�Ѽƭȳ]�m��FormAttribute(struct)
func (compo *FormAttribute) SetTitle(value template.HTML) types.FormAttribute {
	compo.Title = value
	return compo
}

// GetDefaultBoxHeader�P�_����(�O�_���ê�^����)����oclass="box-title"HTML�y�k
func (compo *FormAttribute) GetDefaultBoxHeader(hideBack bool) template.HTML {
	if hideBack {
		return template.HTML(fmt.Sprintf(`<h3 class="box-title">%s</h3>`, language.GetFromHtml(compo.Title)))
	}

	return template.HTML(fmt.Sprintf(`<h3 class="box-title">%s</h3>
            <div class="box-tools">
                <div class="btn-group pull-right" style="margin-right: 10px">
                    <a href='%s' class="btn btn-sm btn-default form-history-back"><i
                                class="fa fa-arrow-left"></i>??%s</a>
                </div>
            </div>`, language.GetFromHtml(compo.Title), compo.HiddenFields[form2.PreviousKey], language.Get("Back")))
}

// GetDetailBoxHeader(���o�Ӹ`��BoxHeader)�N�s��ΧR����HTML�y�k�B�z��[�Jclass="box-title��HTML�y�k��
func (compo *FormAttribute) GetDetailBoxHeader(editUrl, deleteUrl string) template.HTML {

	var (
		editBtn   string
		deleteBtn string
	)

	if editUrl != "" {
		editBtn = fmt.Sprintf(`
                <div class="btn-group pull-right" style="margin-right: 10px">
                    <a href='%s' class="btn btn-sm btn-primary"><i
                                class="fa fa-edit"></i>??%s</a>
                </div>`, editUrl, language.Get("Edit"))
	}

	if deleteUrl != "" {
		deleteBtn = fmt.Sprintf(`
                <div class="btn-group pull-right" style="margin-right: 10px">
                    <a href='javascript:;' class="btn btn-sm btn-danger delete-btn"><i
                                class="fa fa-trash"></i>??%s</a>
                </div>`, language.Get("Delete"))
	}

	return template.HTML(`<h3 class="box-title">`) + language.GetFromHtml(compo.Title) + template.HTML(`</h3>
            <div class="box-tools">
				`+deleteBtn+editBtn+`
                <div class="btn-group pull-right" style="margin-right: 10px">
                    <a href='`+compo.HiddenFields[form2.PreviousKey]+`' class="btn btn-sm btn-default form-history-back"><i
                                class="fa fa-arrow-left"></i>??`+language.Get("Back")+`</a>
                </div>
            </div>`)
}

// GetBoxHeaderNoButton(���oBoxHeader���n���s))��HTML�y�k
func (compo *FormAttribute) GetBoxHeaderNoButton() template.HTML {
	return template.HTML(fmt.Sprintf(`<h3 class="box-title">%s</h3>`, language.GetFromHtml(compo.Title)))
}

// �N�Ѽƭȳ]�m��FormAttribute(struct)
func (compo *FormAttribute) SetOperationFooter(value template.HTML) types.FormAttribute {
	compo.OperationFooter = value
	return compo
}

// �P�_�����NFormFields�K�[��FormAttribute.ContentList([]FormFields)
// �����N�ŦXTreeAttribute.TemplateList["components/�h��"](map[string]string)���ȥ[�Jtext(string)
// ���۱N�Ѽ�compo�g�Jbuffer(bytes.Buffer)���̫��XHTML
func (compo *FormAttribute) GetContent() template.HTML {
	// GetAssetUrl return globalCfg.AssetUrl
	compo.CdnUrl = config.GetAssetUrl()

	// Col�bbtemplate\types\form\form.go
	// type Layout uint8
	if col := compo.Layout.Col(); col > 0 {
		compo.ContentList = make([]types.FormFields, col)
		index := 0
		for i := 0; i < len(compo.Content); i++ {
			ii := index % col
			// �NFormFields�K?[?�FormAttribute.ContentList([]FormFields)
			compo.ContentList[ii] = append(compo.ContentList[ii], compo.Content[i])
			if i < len(compo.Content)-1 {
				if strings.Contains(compo.Content[i+1].Field, "__goadmin_operator__") {
					compo.ContentList[ii] = append(compo.ContentList[ii], compo.Content[i+1])
					i++
				}
			}
			index++
		}
	}

	// �btemplate\components\composer.go
	// �����N�ŦXTreeAttribute.TemplateList["components/�h��"](map[string]string)���ȥ[�Jtext(string)�A���۱N�ѼƤΥ\��K�[���s���ҪO�øѪR���ҪO���D��
	// �N�Ѽ�compo�g�Jbuffer(bytes.Buffer)���̫��XHTML
	return ComposeHtml(compo.TemplateList, *compo, "form",
		"form/default", "form/file", "form/multi_file", "form/textarea", "form/custom",
		"form/selectbox", "form/text", "form/radio", "form/switch",
		"form/password", "form/code", "form/select", "form/singleselect", "form/datetime_range",
		"form/richtext", "form/iconpicker", "form/datetime", "form/number", "form/number_range",
		"form/email", "form/url", "form/ip", "form/color", "form/currency", "form_components",
		"form_layout_default", "form_layout_two_col", "form_layout_tab", "form_components_layout", "form_layout_flow")
}
