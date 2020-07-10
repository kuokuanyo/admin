package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/constant"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/modules/file"
	"github.com/GoAdminGroup/go-admin/modules/language"
	"github.com/GoAdminGroup/go-admin/modules/utils"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	form2 "github.com/GoAdminGroup/go-admin/template/types/form"
)

type FieldOption struct {
	Text          string            `json:"text"`
	Value         string            `json:"value"`
	TextHTML      template.HTML     `json:"-"`
	Selected      bool              `json:"-"`
	SelectedLabel template.HTML     `json:"-"`
	Extra         map[string]string `json:"-"`
}

type FieldOptions []FieldOption

// 判斷條件後將參數labels([]template.HTML)加入FieldOptions[k].SelectedLabel
func (fo FieldOptions) SetSelected(val interface{}, labels []template.HTML) FieldOptions {

	if valArr, ok := val.([]string); ok {
		for k := range fo {
			text := fo[k].Text
			if text == "" {
				text = string(fo[k].TextHTML)
			}
			fo[k].Selected = utils.InArray(valArr, fo[k].Value) || utils.InArray(valArr, text)
			if fo[k].Selected {
				fo[k].SelectedLabel = labels[0]
			} else {
				fo[k].SelectedLabel = labels[1]
			}
		}
	} else {
		for k := range fo {
			text := fo[k].Text
			if text == "" {
				text = string(fo[k].TextHTML)
			}
			fo[k].Selected = fo[k].Value == val || text == val
			if fo[k].Selected {
				fo[k].SelectedLabel = labels[0]
			} else {
				fo[k].SelectedLabel = labels[1]
			}
		}
	}

	return fo
}

// 對FieldOptions([]FieldOption)執行迴圈，判斷條件後將參數(html)設置至FieldOptions[k].SelectedLabel後回傳
func (fo FieldOptions) SetSelectedLabel(labels []template.HTML) FieldOptions {
	for k := range fo {
		if fo[k].Selected {
			fo[k].SelectedLabel = labels[0]
		} else {
			fo[k].SelectedLabel = labels[1]
		}
	}
	return fo
}

func (fo FieldOptions) Marshal() string {
	if len(fo) == 0 {
		return ""
	}
	eo, err := json.Marshal(fo)

	if err != nil {
		return ""
	}

	return string(eo)
}

type OptionInitFn func(val FieldModel) FieldOptions

type OptionTable struct {
	Table          string
	TextField      string
	ValueField     string
	QueryProcessFn OptionTableQueryProcessFn
	ProcessFn      OptionProcessFn
}

type OptionTableQueryProcessFn func(sql *db.SQL) *db.SQL

type OptionProcessFn func(options FieldOptions) FieldOptions

// FormField is the form field with different options.
type FormField struct {
	Field    string          `json:"field"`
	TypeName db.DatabaseType `json:"type_name"`
	Head     string          `json:"head"`
	FormType form2.Type      `json:"form_type"`

	Default                template.HTML `json:"default"`
	Value                  template.HTML `json:"value"`
	Value2                 string        `json:"value_2"`
	Options                FieldOptions  `json:"options"`
	DefaultOptionDelimiter string        `json:"default_option_delimiter"`
	Label                  template.HTML `json:"label"`

	Placeholder string `json:"placeholder"`

	CustomContent template.HTML `json:"custom_content"`
	CustomJs      template.JS   `json:"custom_js"`
	CustomCss     template.CSS  `json:"custom_css"`

	Editable    bool `json:"editable"`
	NotAllowAdd bool `json:"not_allow_add"`
	Must        bool `json:"must"`
	Hide        bool `json:"hide"`

	Width int `json:"width"`

	Joins Joins `json:"-"`

	Divider      bool   `json:"divider"`
	DividerTitle string `json:"divider_title"`

	HelpMsg template.HTML `json:"help_msg"`

	OptionExt    template.JS  `json:"option_ext"`
	OptionInitFn OptionInitFn `json:"-"`
	OptionTable  OptionTable  `json:"-"`

	FieldDisplay `json:"-"`
	PostFilterFn PostFieldFilterFn `json:"-"`
}

// 將FormField(struct)的值更新後回傳
func (f FormField) UpdateValue(id, val string, res map[string]interface{}, sqls ...*db.SQL) FormField {
	if f.FormType.IsSelect() {
		// type OptionInitFn func(val FieldModel) FieldOptions
		// SetSelectedLabel(設置所選標籤)對FieldOptions([]FieldOption)執行迴圈，判斷條件後將參數(html)設置至FieldOptions[k].SelectedLabel後回傳
		// SelectedLabel判斷條件後回傳[]template.HTML
		if len(f.Options) == 0 && f.OptionInitFn != nil {
			f.Options = f.OptionInitFn(FieldModel{
				ID:    id,
				Value: val,
				Row:   res,
			}).SetSelectedLabel(f.FormType.SelectedLabel())
		} else if len(f.Options) == 0 && f.OptionTable.Table != "" && len(sqls) > 0 && sqls[0] != nil {

			// Table將SQL(struct)資訊清除後將參數設置至SQL.TableName回傳
			// Select將參數設置至SQL(struct).Fields並且設置SQL(struct).Functions
			sqls[0].Table(f.OptionTable.Table).Select(f.OptionTable.ValueField, f.OptionTable.TextField)

			if f.OptionTable.QueryProcessFn != nil {
				f.OptionTable.QueryProcessFn(sqls[0])
			}

			// 返回所有符合查詢的結果
			queryRes, err := sqls[0].All()
			if err == nil {
				for _, item := range queryRes {
					f.Options = append(f.Options, FieldOption{
						Value: fmt.Sprintf("%v", item[f.OptionTable.ValueField]),
						Text:  fmt.Sprintf("%v", item[f.OptionTable.TextField]),
					})
				}
			}

			if f.OptionTable.ProcessFn != nil {
				f.Options = f.OptionTable.ProcessFn(f.Options)
			}

			// SetSelected判斷條件後將參數f.FormType.SelectedLabel()([]template.HTML)加入FieldOptions[k].SelectedLabel
			f.Options.SetSelected(f.ToDisplay(FieldModel{
				ID:    id,
				Value: val,
				Row:   res,
			}), f.FormType.SelectedLabel())

		} else {
			// SetSelected判斷條件後將參數f.FormType.SelectedLabel()([]template.HTML)加入FieldOptions[k].SelectedLabel
			f.Options.SetSelected(f.ToDisplay(FieldModel{
				ID:    id,
				Value: val,
				Row:   res,
			}), f.FormType.SelectedLabel())
		}
	} else {
		value := f.ToDisplay(FieldModel{
			ID:    id,
			Value: val,
			Row:   res,
		})
		if v, ok := value.(template.HTML); ok {
			f.Value = v
		} else {
			f.Value = template.HTML(value.(string))
		}
	}
	return f
}

// UpdateDefaultValue將預設值更新後回傳FormField(struct)
func (f FormField) UpdateDefaultValue(sqls ...*db.SQL) FormField {
	f.Value = f.Default // template.HTML

	// IsSelect在template\types\form\form.go
	// 判斷f.FormTypet(unit8)是否符合條件
	if f.FormType.IsSelect() {
		if len(f.Options) == 0 && f.OptionInitFn != nil {
			// type OptionInitFn func(val FieldModel) FieldOptions
			// SetSelectedLabel(設置所選標籤)對FieldOptions([]FieldOption)執行迴圈，判斷條件後將參數(html)設置至FieldOptions[k].SelectedLabel後回傳
			// SelectedLabel判斷條件後回傳[]template.HTML
			f.Options = f.OptionInitFn(FieldModel{
				ID:    "",
				Value: string(f.Value),
				Row:   make(map[string]interface{}),
			}).SetSelectedLabel(f.FormType.SelectedLabel())
		} else if len(f.Options) == 0 && f.OptionTable.Table != "" && len(sqls) > 0 && sqls[0] != nil {
			// Table將SQL(struct)資訊清除後將參數設置至SQL.TableName回傳
			// Select將參數設置至SQL(struct).Fields並且設置SQL(struct).Functions
			sqls[0].Table(f.OptionTable.Table).Select(f.OptionTable.ValueField, f.OptionTable.TextField)

			if f.OptionTable.QueryProcessFn != nil {
				f.OptionTable.QueryProcessFn(sqls[0])
			}

			// 返回所有符合查詢的結果
			res, err := sqls[0].All()

			if err == nil {
				for _, item := range res {
					f.Options = append(f.Options, FieldOption{
						Value: fmt.Sprintf("%v", item[f.OptionTable.ValueField]),
						Text:  fmt.Sprintf("%v", item[f.OptionTable.TextField]),
					})
				}
			}

			if f.OptionTable.ProcessFn != nil {
				f.Options = f.OptionTable.ProcessFn(f.Options)
			}
			// SetSelected判斷條件後將參數f.FormType.SelectedLabel()([]template.HTML)加入FieldOptions[k].SelectedLabel
			f.Options.SetSelected(f.ToDisplay(FieldModel{
				ID:    "",
				Value: string(f.Value),
				Row:   make(map[string]interface{}),
			}), f.FormType.SelectedLabel())

		} else {
			// SelectedLabel判斷條件後回傳[]template.HTML
			// SetSelected判斷條件後將參數f.FormType.SelectedLabel()([]template.HTML)加入FieldOptions[k].SelectedLabel
			f.Options.SetSelected(string(f.Value), f.FormType.SelectedLabel())
		}
	}
	return f
}

// FillCustomContent(填寫自定義內容)判斷條件後設置FormField回傳
func (f FormField) FillCustomContent() FormField {
	// TODO: optimize
	if f.CustomContent != "" {
		f.CustomContent = template.HTML(f.fillCustom(string(f.CustomContent)))
	}
	if f.CustomJs != "" {
		f.CustomJs = template.JS(f.fillCustom(string(f.CustomJs)))
	}
	if f.CustomCss != "" {
		f.CustomCss = template.CSS(f.fillCustom(string(f.CustomCss)))
	}
	return f
}

func (f FormField) fillCustom(src string) string {
	t := template.New("custom")
	t, _ = t.Parse(src)
	buf := new(bytes.Buffer)
	_ = t.Execute(buf, f)
	return buf.String()
}

// FormPanel
type FormPanel struct {
	FieldList         FormFields
	curFieldListIndex int

	// Warn: may be deprecated in the future.
	TabGroups  TabGroups
	TabHeaders TabHeaders

	Table       string
	Title       string
	Description string

	Validator    FormPostFn
	PostHook     FormPostFn
	PreProcessFn FormPreProcessFn

	Callbacks Callbacks

	primaryKey primaryKey

	UpdateFn FormPostFn
	InsertFn FormPostFn

	IsHideContinueEditCheckBox bool
	IsHideContinueNewCheckBox  bool
	IsHideResetButton          bool
	IsHideBackButton           bool

	processChains DisplayProcessFnChains

	HeaderHtml template.HTML
	FooterHtml template.HTML
}

func NewFormPanel() *FormPanel {
	return &FormPanel{
		curFieldListIndex: -1,
		Callbacks:         make(Callbacks, 0),
	}
}

func (f *FormPanel) AddLimitFilter(limit int) *FormPanel {
	f.processChains = addLimit(limit, f.processChains)
	return f
}

func (f *FormPanel) AddTrimSpaceFilter() *FormPanel {
	f.processChains = addTrimSpace(f.processChains)
	return f
}

func (f *FormPanel) AddSubstrFilter(start int, end int) *FormPanel {
	f.processChains = addSubstr(start, end, f.processChains)
	return f
}

func (f *FormPanel) AddToTitleFilter() *FormPanel {
	f.processChains = addToTitle(f.processChains)
	return f
}

func (f *FormPanel) AddToUpperFilter() *FormPanel {
	f.processChains = addToUpper(f.processChains)
	return f
}

func (f *FormPanel) AddToLowerFilter() *FormPanel {
	f.processChains = addToLower(f.processChains)
	return f
}

func (f *FormPanel) AddXssFilter() *FormPanel {
	f.processChains = addXssFilter(f.processChains)
	return f
}

func (f *FormPanel) AddXssJsFilter() *FormPanel {
	f.processChains = addXssJsFilter(f.processChains)
	return f
}

// 將參數name、type設置至FormPanel.primaryKey後回傳
func (f *FormPanel) SetPrimaryKey(name string, typ db.DatabaseType) *FormPanel {
	f.primaryKey = primaryKey{Name: name, Type: typ}
	return f
}

func (f *FormPanel) HideContinueEditCheckBox() *FormPanel {
	f.IsHideContinueEditCheckBox = true
	return f
}

func (f *FormPanel) HideContinueNewCheckBox() *FormPanel {
	f.IsHideContinueNewCheckBox = true
	return f
}

func (f *FormPanel) HideResetButton() *FormPanel {
	f.IsHideResetButton = true
	return f
}

func (f *FormPanel) HideBackButton() *FormPanel {
	f.IsHideBackButton = true
	return f
}

func (f *FormPanel) AddField(head, field string, filedType db.DatabaseType, formType form2.Type) *FormPanel {

	f.FieldList = append(f.FieldList, FormField{
		Head:        head,
		Field:       field,
		TypeName:    filedType,
		Editable:    true,
		Hide:        false,
		Placeholder: language.Get("input") + " " + head,
		FormType:    formType,
		FieldDisplay: FieldDisplay{
			Display: func(value FieldModel) interface{} {
				return value.Value
			},
			DisplayProcessChains: chooseDisplayProcessChains(f.processChains),
		},
	})
	f.curFieldListIndex++

	if formType.IsFile() {
		f.FieldOptionExt(map[string]interface{}{
			"overwriteInitial":     true,
			"initialPreviewAsData": true,
			"browseLabel":          language.Get("Browse"),
			"showRemove":           false,
			"previewClass":         "preview-" + field,
			"showUpload":           false,
			"allowedFileTypes":     []string{"image"},
		})
	}

	if formType.IsCode() {
		f.FieldList[f.curFieldListIndex].OptionExt = `
	theme = "monokai";
	font_size = 14;
	language = "html";
	options = {useWorker: false};
`
	}

	return f
}

// Field attribute setting functions
// ====================================================

func (f *FormPanel) FieldDisplay(filter FieldFilterFn) *FormPanel {
	f.FieldList[f.curFieldListIndex].Display = filter
	return f
}

func (f *FormPanel) SetTable(table string) *FormPanel {
	f.Table = table
	return f
}

func (f *FormPanel) FieldMust() *FormPanel {
	f.FieldList[f.curFieldListIndex].Must = true
	return f
}

func (f *FormPanel) FieldHide() *FormPanel {
	f.FieldList[f.curFieldListIndex].Hide = true
	return f
}

func (f *FormPanel) FieldPlaceholder(placeholder string) *FormPanel {
	f.FieldList[f.curFieldListIndex].Placeholder = placeholder
	return f
}

func (f *FormPanel) FieldWidth(width int) *FormPanel {
	f.FieldList[f.curFieldListIndex].Width = width
	return f
}

func (f *FormPanel) FieldDivider(title ...string) *FormPanel {
	f.FieldList[f.curFieldListIndex].Divider = true
	if len(title) > 0 {
		f.FieldList[f.curFieldListIndex].DividerTitle = title[0]
	}
	return f
}

func (f *FormPanel) FieldHelpMsg(s template.HTML) *FormPanel {
	f.FieldList[f.curFieldListIndex].HelpMsg = s
	return f
}

func (f *FormPanel) FieldOptionInitFn(fn OptionInitFn) *FormPanel {
	f.FieldList[f.curFieldListIndex].OptionInitFn = fn
	return f
}

func (f *FormPanel) FieldOptionExt(m map[string]interface{}) *FormPanel {

	if f.FieldList[f.curFieldListIndex].FormType.IsCode() {
		f.FieldList[f.curFieldListIndex].OptionExt = template.JS(fmt.Sprintf(`
	theme = "%s";
	font_size = %s;
	language = "%s";
	options = %s;
`, m["theme"], m["font_size"], m["language"], m["options"]))
		return f
	}

	s, _ := json.Marshal(m)

	if f.FieldList[f.curFieldListIndex].OptionExt != template.JS("") {
		ss := string(f.FieldList[f.curFieldListIndex].OptionExt)
		ss = strings.Replace(ss, "}", "", strings.Count(ss, "}"))
		ss = strings.TrimRight(ss, " ")
		ss += ","
		f.FieldList[f.curFieldListIndex].OptionExt = template.JS(ss) + template.JS(strings.Replace(string(s), "{", "", 1))
	} else {
		f.FieldList[f.curFieldListIndex].OptionExt = template.JS(string(s))
	}

	return f
}

func (f *FormPanel) FieldOptionExtJS(js template.JS) *FormPanel {
	f.FieldList[f.curFieldListIndex].OptionExt = js
	return f
}

func (f *FormPanel) FieldEnableFileUpload(data ...interface{}) *FormPanel {

	url := f.OperationURL("/file/upload")

	if len(data) > 0 {
		url = data[0].(string)
	}

	field := f.FieldList[f.curFieldListIndex].Field

	f.FieldList[f.curFieldListIndex].OptionExt = template.JS(fmt.Sprintf(`
	%seditor.customConfig.uploadImgServer = '%s';
	%seditor.customConfig.uploadImgMaxSize = 3 * 1024 * 1024;
	%seditor.customConfig.uploadImgMaxLength = 5;
	%seditor.customConfig.uploadFileName = 'file';
`, field, url, field, field, field))

	var fileUploadHandler context.Handler
	if len(data) > 1 {
		fileUploadHandler = data[1].(context.Handler)
	} else {
		fileUploadHandler = func(ctx *context.Context) {
			if len(ctx.Request.MultipartForm.File) == 0 {
				ctx.JSON(http.StatusOK, map[string]interface{}{
					"errno": 400,
				})
				return
			}

			err := file.GetFileEngine(config.GetFileUploadEngine().Name).Upload(ctx.Request.MultipartForm)
			if err != nil {
				ctx.JSON(http.StatusOK, map[string]interface{}{
					"errno": 500,
				})
				return
			}

			var imgPath = make([]string, len(ctx.Request.MultipartForm.Value["file"]))
			for i, path := range ctx.Request.MultipartForm.Value["file"] {
				imgPath[i] = config.GetStore().URL(path)
			}

			ctx.JSON(http.StatusOK, map[string]interface{}{
				"errno": 0,
				"data":  imgPath,
			})
		}
	}

	f.Callbacks = f.Callbacks.AddCallback(context.Node{
		Path:     url,
		Method:   "post",
		Value:    map[string]interface{}{constant.ContextNodeNeedAuth: 1},
		Handlers: []context.Handler{fileUploadHandler},
	})

	return f
}

func (f *FormPanel) FieldDefault(def string) *FormPanel {
	f.FieldList[f.curFieldListIndex].Default = template.HTML(def)
	return f
}

func (f *FormPanel) FieldNotAllowEdit() *FormPanel {
	f.FieldList[f.curFieldListIndex].Editable = false
	return f
}

func (f *FormPanel) FieldNotAllowAdd() *FormPanel {
	f.FieldList[f.curFieldListIndex].NotAllowAdd = true
	return f
}

func (f *FormPanel) FieldFormType(formType form2.Type) *FormPanel {
	f.FieldList[f.curFieldListIndex].FormType = formType
	return f
}

func (f *FormPanel) FieldValue(value string) *FormPanel {
	f.FieldList[f.curFieldListIndex].Value = template.HTML(value)
	return f
}

func (f *FormPanel) FieldOptionsFromTable(table, textFieldName, valueFieldName string, process ...OptionTableQueryProcessFn) *FormPanel {
	var fn OptionTableQueryProcessFn
	if len(process) > 0 {
		fn = process[0]
	}
	f.FieldList[f.curFieldListIndex].OptionTable = OptionTable{
		Table:          table,
		TextField:      textFieldName,
		ValueField:     valueFieldName,
		QueryProcessFn: fn,
	}
	return f
}

func (f *FormPanel) FieldOptionsTableProcessFn(fn OptionProcessFn) *FormPanel {
	f.FieldList[f.curFieldListIndex].OptionTable.ProcessFn = fn
	return f
}

func (f *FormPanel) FieldOptions(options FieldOptions) *FormPanel {
	f.FieldList[f.curFieldListIndex].Options = options
	return f
}

func (f *FormPanel) FieldDefaultOptionDelimiter(delimiter string) *FormPanel {
	f.FieldList[f.curFieldListIndex].DefaultOptionDelimiter = delimiter
	return f
}

func (f *FormPanel) FieldPostFilterFn(post PostFieldFilterFn) *FormPanel {
	f.FieldList[f.curFieldListIndex].PostFilterFn = post
	return f
}

func (f *FormPanel) FieldLimit(limit int) *FormPanel {
	f.FieldList[f.curFieldListIndex].DisplayProcessChains = f.FieldList[f.curFieldListIndex].AddLimit(limit)
	return f
}

func (f *FormPanel) FieldTrimSpace() *FormPanel {
	f.FieldList[f.curFieldListIndex].DisplayProcessChains = f.FieldList[f.curFieldListIndex].AddTrimSpace()
	return f
}

func (f *FormPanel) FieldSubstr(start int, end int) *FormPanel {
	f.FieldList[f.curFieldListIndex].DisplayProcessChains = f.FieldList[f.curFieldListIndex].AddSubstr(start, end)
	return f
}

func (f *FormPanel) FieldToTitle() *FormPanel {
	f.FieldList[f.curFieldListIndex].DisplayProcessChains = f.FieldList[f.curFieldListIndex].AddToTitle()
	return f
}

func (f *FormPanel) FieldToUpper() *FormPanel {
	f.FieldList[f.curFieldListIndex].DisplayProcessChains = f.FieldList[f.curFieldListIndex].AddToUpper()
	return f
}

func (f *FormPanel) FieldToLower() *FormPanel {
	f.FieldList[f.curFieldListIndex].DisplayProcessChains = f.FieldList[f.curFieldListIndex].AddToLower()
	return f
}

func (f *FormPanel) FieldXssFilter() *FormPanel {
	f.FieldList[f.curFieldListIndex].DisplayProcessChains = f.FieldList[f.curFieldListIndex].DisplayProcessChains.
		Add(func(s string) string {
			return html.EscapeString(s)
		})
	return f
}

func (f *FormPanel) FieldCustomContent(content template.HTML) *FormPanel {
	f.FieldList[f.curFieldListIndex].CustomContent = content
	return f
}

func (f *FormPanel) FieldCustomJs(js template.JS) *FormPanel {
	f.FieldList[f.curFieldListIndex].CustomJs = js
	return f
}

func (f *FormPanel) FieldCustomCss(css template.CSS) *FormPanel {
	f.FieldList[f.curFieldListIndex].CustomCss = css
	return f
}

func (f *FormPanel) FieldOnSearch(url string, handler Handler, delay ...int) *FormPanel {
	ext, callback := searchJS(f.FieldList[f.curFieldListIndex].OptionExt, f.OperationURL(url), handler, delay...)
	f.FieldList[f.curFieldListIndex].OptionExt = ext
	f.Callbacks = f.Callbacks.AddCallback(callback)
	return f
}

func (f *FormPanel) FieldOnChooseCustom(js template.HTML) *FormPanel {
	f.FooterHtml += chooseCustomJS(f.FieldList[f.curFieldListIndex].Field, js)
	return f
}

type LinkField struct {
	Field   string
	Value   template.HTML
	Hide    bool
	Disable bool
}

func (f *FormPanel) FieldOnChooseMap(m map[string]LinkField) *FormPanel {
	f.FooterHtml += chooseMapJS(f.FieldList[f.curFieldListIndex].Field, m)
	return f
}

func (f *FormPanel) FieldOnChoose(val, field string, value template.HTML) *FormPanel {
	f.FooterHtml += chooseJS(f.FieldList[f.curFieldListIndex].Field, field, val, value)
	return f
}

func (f *FormPanel) OperationURL(id string) string {
	return config.Url("/operation/" + utils.WrapURL(id))
}

func (f *FormPanel) FieldOnChooseAjax(field, url string, handler Handler) *FormPanel {
	js, callback := chooseAjax(f.FieldList[f.curFieldListIndex].Field, field, f.OperationURL(url), handler)
	f.FooterHtml += js
	f.Callbacks = f.Callbacks.AddCallback(callback)
	return f
}

func (f *FormPanel) FieldOnChooseHide(value string, field ...string) *FormPanel {
	f.FooterHtml += chooseHideJS(f.FieldList[f.curFieldListIndex].Field, value, field...)
	return f
}

func (f *FormPanel) FieldOnChooseShow(value string, field ...string) *FormPanel {
	f.FooterHtml += chooseShowJS(f.FieldList[f.curFieldListIndex].Field, value, field...)
	return f
}

func (f *FormPanel) FieldOnChooseDisable(value string, field ...string) *FormPanel {
	f.FooterHtml += chooseDisableJS(f.FieldList[f.curFieldListIndex].Field, value, field...)
	return f
}

func searchJS(ext template.JS, url string, handler Handler, delay ...int) (template.JS, context.Node) {
	delayStr := "500"
	if len(delay) > 0 {
		delayStr = strconv.Itoa(delay[0])
	}

	if ext != template.JS("") {
		s := string(ext)
		s = strings.Replace(s, "{", "", 1)
		s = utils.ReplaceNth(s, "}", "", strings.Count(s, "}"))
		s = strings.TrimRight(s, " ")
		s += ","
		ext = template.JS(s)
	}

	return template.JS(`{
		`) + ext + template.JS(`
		ajax: {
		    url: "`+url+`",
		    dataType: 'json',
		    data: function (params) {
			      var query = {
			        	search: params.term,
						page: params.page || 1
			      }
			      return query;
		    },
		    delay: `+delayStr+`,
		    processResults: function (data, params) {
			      return data.data;
	    	}
	  	}
	}`), context.Node{
			Path:     url,
			Method:   "get",
			Handlers: context.Handlers{handler.Wrap()},
			Value:    map[string]interface{}{constant.ContextNodeNeedAuth: 1},
		}
}

func chooseCustomJS(field string, js template.HTML) template.HTML {
	return `<script>
$(".` + template.HTML(field) + `").on("select2:select",function(e){
	` + js + `
})
</script>`
}

func chooseMapJS(field string, m map[string]LinkField) template.HTML {
	cm := template.HTML("")

	for val, obejct := range m {
		if obejct.Hide {
			cm += `if (e.params.data.text === "` + template.HTML(val) + `") {
		$("label[for='` + template.HTML(obejct.Field) + `']").parent().hide()
	} else {
		$("label[for='` + template.HTML(obejct.Field) + `']").parent().show()
	}`
		} else if obejct.Disable {
			cm += `if (e.params.data.text === "` + template.HTML(val) + `") {
		$("#` + template.HTML(obejct.Field) + `").prop('disabled', true);
	} else {
		$("#` + template.HTML(obejct.Field) + `").prop('disabled', false);
	}`
		} else {
			cm += `if (e.params.data.text === "` + template.HTML(val) + `") {
		if ($(".` + template.HTML(obejct.Field) + `").length > 0) {
			$(".` + template.HTML(obejct.Field) + `").val("` + obejct.Value + `").select2()
		} else {
			$("#` + template.HTML(obejct.Field) + `").val("` + obejct.Value + `")
		}	
	}`
		}
	}

	return `<script>
$(".` + template.HTML(field) + `").on("select2:select",function(e){
	` + cm + `
})
</script>`
}

func chooseJS(field, chooseField, val string, value template.HTML) template.HTML {
	return `<script>
$(".` + template.HTML(field) + `").on("select2:select",function(e){
	if (e.params.data.text === "` + template.HTML(val) + `") {
		if ($(".` + template.HTML(chooseField) + `").length > 0) {
			$(".` + template.HTML(chooseField) + `").val("` + value + `").select2()
		} else {
			$("#` + template.HTML(chooseField) + `").val("` + value + `")
		}	
	}
})
</script>`
}

func chooseAjax(field, chooseField, url string, handler Handler) (template.HTML, context.Node) {
	return `<script>

let updateBoxSelections = function(selectObj, new_opts) {
    selectObj.html('');
    new_opts.forEach(function (opt) {
      	selectObj.append($('<option value="'+opt["id"]+'">'+opt["text"]+'</option>'));
    });
	console.log("selectObj.bootstrapDualListbox", selectObj.bootstrapDualListbox)
    selectObj.bootstrapDualListbox('refresh', true);
}

let req = function(selectObj, box, event) {
	$.ajax({
		url:"` + template.HTML(url) + `",
		type: 'post',
		dataType: 'text',
		data: {
			'value':$(".` + template.HTML(field) + `").val(),
			'event': event
		},
		success: function (data)  {
			if (typeof (data) === "string") {
				data = JSON.parse(data);
			}
			if (data.code === 0) {
				if (selectObj.length > 0) {
					if (typeof(data.data) === "object") {
						if (box) {
							updateBoxSelections(selectObj, data.data)
						} else {
							if (typeof(selectObj.attr("multiple")) !== "undefined") {
								selectObj.html("");
							}
							selectObj.select2({
								data: data.data
							});
						}	
					} else {
						if (box) {
							selectObj.val(data.data).select2()
						} else {
							
						}
					}
				} else {
					$('#` + template.HTML(chooseField) + `').val(data.data);
				}
			} else {
				swal(data.msg, '', 'error');
			}
		},
		error:function(){
			alert('error')
		}
	})
}

if ($("label[for='` + template.HTML(field) + `']").next().find(".bootstrap-duallistbox-container").length === 0) {
	$(".` + template.HTML(field) + `").on("select2:select", function(e) {
		let id = '` + template.HTML(chooseField) + `'
		let selectObj = $("."+id)
		if (selectObj.length > 0) {
			selectObj.val("").select2()
			selectObj.html('<option value="" selected="selected"></option>')
		}
		req(selectObj, false, "select");
	})
	if (typeof($(".` + template.HTML(field) + `").attr("multiple")) !== "undefined") {
		$(".` + template.HTML(field) + `").on("select2:unselect",function(e){
			let id = '` + template.HTML(chooseField) + `'
			let selectObj = $("."+id)
			if (selectObj.length > 0) {
				selectObj.val("").select2()
				selectObj.html('<option value="" selected="selected"></option>')
			}
			req(selectObj, false, "unselect");
		})
	}
} else {
	let ` + template.HTML(field) + `_lastState = $(".` + template.HTML(field) + `").val();

	$(".` + template.HTML(field) + `").on('change',function (e) {
    	var newState = $(this).val();                     
		if ($(` + template.HTML(field) + `_lastState).not(newState).get().length > 0) {
			let id = '` + template.HTML(chooseField) + `'
			req($("."+id), true, "unselect");
		}
		if ($(newState).not(` + template.HTML(field) + `_lastState).get().length > 0) {
			let id = '` + template.HTML(chooseField) + `'
			req($("."+id), true, "select");
		}
    	` + template.HTML(field) + `_lastState = newState;
	})
}
</script>`, context.Node{
			Path:     url,
			Method:   "post",
			Handlers: context.Handlers{handler.Wrap()},
			Value:    map[string]interface{}{constant.ContextNodeNeedAuth: 1},
		}
}

func chooseHideJS(field, value string, chooseFields ...string) template.HTML {
	if len(chooseFields) == 0 {
		return ""
	}

	hideText := template.HTML("")
	showText := template.HTML("")

	for _, f := range chooseFields {
		hideText += `$("label[for='` + template.HTML(f) + `']").parent().hide()
`
		showText += `$("label[for='` + template.HTML(f) + `']").parent().show()
`
	}

	return `<script>
$(".` + template.HTML(field) + `").on("select2:select",function(e){
	if (e.params.data.text === "` + template.HTML(value) + `") {
		` + hideText + `
	} else {
		` + showText + `
	}
})
$(function(){
	let data = $(".` + template.HTML(field) + `").select2("data");
	let text = "";
	if (data.length > 0) {
		text = data[0].text;
	}
	if (text === "` + template.HTML(value) + `") {
		` + hideText + `
	}
})
</script>`
}

func chooseShowJS(field, value string, chooseFields ...string) template.HTML {
	if len(chooseFields) == 0 {
		return ""
	}

	hideText := template.HTML("")
	showText := template.HTML("")

	for _, f := range chooseFields {
		hideText += `$("label[for='` + template.HTML(f) + `']").parent().hide()
`
		showText += `$("label[for='` + template.HTML(f) + `']").parent().show()
`
	}

	return `<script>
$(".` + template.HTML(field) + `").on("select2:select",function(e){
	if (e.params.data.text === "` + template.HTML(value) + `") {
		` + showText + `
	} else {
		` + hideText + `
	}
})
$(function(){
	let data = $(".` + template.HTML(field) + `").select2("data");
	let text = "";
	if (data.length > 0) {
		text = data[0].text;
	}
	if (text !== "` + template.HTML(value) + `") {
		` + hideText + `
	}
})
</script>`
}

func chooseDisableJS(field, value string, chooseFields ...string) template.HTML {
	if len(chooseFields) == 0 {
		return ""
	}

	disableText := template.HTML("")
	enableText := template.HTML("")

	for _, f := range chooseFields {
		disableText += `$("#` + template.HTML(f) + `").prop('disabled', true);
`
		enableText += `$("#` + template.HTML(f) + `").prop('disabled', false);
`
	}

	return `<script>
$(".` + template.HTML(field) + `").on("select2:select",function(e){
	if (e.params.data.text === "` + template.HTML(value) + `") {
		` + disableText + `
	} else {
		` + enableText + `
	}
})
</script>`
}

// FormPanel attribute setting functions
// ====================================================

func (f *FormPanel) SetTitle(title string) *FormPanel {
	f.Title = title
	return f
}

func (f *FormPanel) SetTabGroups(groups TabGroups) *FormPanel {
	f.TabGroups = groups
	return f
}

func (f *FormPanel) SetTabHeaders(headers ...string) *FormPanel {
	f.TabHeaders = headers
	return f
}

func (f *FormPanel) SetDescription(desc string) *FormPanel {
	f.Description = desc
	return f
}

func (f *FormPanel) SetHeaderHtml(header template.HTML) *FormPanel {
	f.HeaderHtml += header
	return f
}

func (f *FormPanel) SetFooterHtml(footer template.HTML) *FormPanel {
	f.FooterHtml += footer
	return f
}

func (f *FormPanel) SetPostValidator(va FormPostFn) *FormPanel {
	f.Validator = va
	return f
}

func (f *FormPanel) SetPreProcessFn(fn FormPreProcessFn) *FormPanel {
	f.PreProcessFn = fn
	return f
}

func (f *FormPanel) SetPostHook(fn FormPostFn) *FormPanel {
	f.PostHook = fn
	return f
}

func (f *FormPanel) SetUpdateFn(fn FormPostFn) *FormPanel {
	f.UpdateFn = fn
	return f
}

func (f *FormPanel) SetInsertFn(fn FormPostFn) *FormPanel {
	f.InsertFn = fn
	return f
}

func (f *FormPanel) GroupFieldWithValue(pk, id string, columns []string, res map[string]interface{}, sql ...func() *db.SQL) ([]FormFields, []string) {
	var (
		groupFormList = make([]FormFields, 0)
		groupHeaders  = make([]string, 0)
		hasPK         = false
	)

	if len(f.TabGroups) > 0 {
		for key, value := range f.TabGroups {
			list := make(FormFields, len(value))
			for j := 0; j < len(value); j++ {
				for _, field := range f.FieldList {
					if value[j] == field.Field {
						if field.Field == pk {
							hasPK = true
						}
						rowValue := modules.AorB(modules.InArray(columns, field.Field) || len(columns) == 0,
							db.GetValueFromDatabaseType(field.TypeName, res[field.Field], len(columns) == 0).String(), "")
						if len(sql) > 0 {
							list[j] = field.UpdateValue(id, rowValue, res, sql[0]())
						} else {
							list[j] = field.UpdateValue(id, rowValue, res)
						}
						if list[j].FormType == form2.File && list[j].Value != template.HTML("") {
							list[j].Value2 = config.GetStore().URL(string(list[j].Value))
						}
						break
					}
				}
			}

			groupFormList = append(groupFormList, list.FillCustomContent())
			groupHeaders = append(groupHeaders, f.TabHeaders[key])
		}

		if len(groupFormList) > 0 && !hasPK {
			groupFormList[len(groupFormList)-1] = groupFormList[len(groupFormList)-1].Add(FormField{
				Head:  pk,
				Field: pk,
				Value: template.HTML(id),
				Hide:  true,
			})
		}
	}

	return groupFormList, groupHeaders
}

// GroupField(欄位分組)先判斷條件後處理FormField，最後將FormField與TabHeader加入至groupFormList與groupHeaders後回傳
func (f *FormPanel) GroupField(sql ...func() *db.SQL) ([]FormFields, []string) {
	var (
		groupFormList = make([]FormFields, 0)
		groupHeaders  = make([]string, 0)
	)

	// FormPanel.TabGroups [][]string
	// 判斷條件
	if len(f.TabGroups) > 0 {
		for key, value := range f.TabGroups {
			list := make(FormFields, 0)
			for i := 0; i < len(value); i++ {
				for _, v := range f.FieldList {
					if v.Field == value[i] {
						if !v.NotAllowAdd {
							v.Editable = true
							if len(sql) > 0 {
								// 在template\types\form.go
								// UpdateDefaultValue首先對FieldOptions([]FieldOption)執行迴圈，判斷條件後將參數(html)設置至FieldOptions[k].SelectedLabel後回傳
								// 最後判斷條件後將參數f.FormType.SelectedLabel()([]template.HTML)加入FieldOptions[k].SelectedLabel，回傳FormField
								// FillCustomContent(填寫自定義內容)對FormFields([]FormField)執行迴圈，判斷條件後設置FormField，最後回傳FormFields([]FormField)
								list = append(list, v.UpdateDefaultValue(sql[0]()).FillCustomContent())
							} else {
								list = append(list, v.UpdateDefaultValue())
							}
							break
						}
					}
				}
			}
			groupFormList = append(groupFormList, list)
			// TabHeaders []string
			groupHeaders = append(groupHeaders, f.TabHeaders[key])
		}
	}
	return groupFormList, groupHeaders
}

// FieldsWithValue(對帶值的欄位更新)對FormPanel.FieldList(FormFields)執行迴圈，分別更新FormField(struct)並加入FormFields後回傳
func (f *FormPanel) FieldsWithValue(pk, id string, columns []string, res map[string]interface{}, sql ...func() *db.SQL) FormFields {
	var (
		list  = make(FormFields, 0)
		hasPK = false
	)
	// field為表單上所有欄位資訊
	for _, field := range f.FieldList {
		// rowValue為該欄位的值
		rowValue := field.GetRawValue(columns, res[field.Field])

		// 編輯menu頁面時都field.FatherField都為空
		if field.FatherField != "" {
			f.FieldList.FindTableField(field.Field, field.FatherField).UpdateValue(id, rowValue, res, sql())
		} else if field.FormType.IsTable() {
			list = append(list, field)
		} else {
			// 將field(struct)的值都更新並加入list([]FormField)中
			list = append(list, *(field.UpdateValue(id, rowValue, res, sql())))
		}

		if field.Field == pk {
			hasPK = true
		}
	}

	// hasPK判斷是否有primary key
	if !hasPK {
		list = list.Add(FormField{
			Head:       pk,
			FieldClass: pk,
			Field:      pk,
			Value:      template.HTML(id),
			FormType:   form2.Default,
			Hide:       true,
		})
	}

	// FillCustomContent(填寫自定義內容)對FormFields([]FormField)執行迴圈，判斷條件後設置FormField，最後回傳FormFields([]FormField)
	return list.FillCustomContent()
}

// GetRawValue為取得該欄位的值
func (f *FormField) GetRawValue(columns []string, v interface{}) string {
	isJSON := len(columns) == 0
	// AorB判斷bool返回第二個(對)或第三個參數(錯)
	// GetValueFromDatabaseType在\modules\db\types.go
	return modules.AorB(isJSON || modules.InArray(columns, f.Field),
		// f.TypeName為該欄位類型(ex: INT)
		// GetValueFromDatabaseType(從資料庫類型取得值)從SQL或JSON取得值
		db.GetValueFromDatabaseType(f.TypeName, v, isJSON).String(), "")
}

// 判斷欄位是否允許添加，例如ID無法手動增加，接著將預設值更新後得到FormField(struct)並加入FormFields中，最後回傳FormFields
func (f *FormPanel) FieldsWithDefaultValue(sql ...func() *db.SQL) FormFields {
	var list = make(FormFields, 0)
	for _, v := range f.FieldList {
		// 判斷欄位是否允許添加，例如ID、建立時間、更新時間欄位無法手動增加
		if v.allowAdd() {
			v.Editable = true
			if v.FatherField != "" {
				if len(sql) > 0 {
					f.FieldList.FindTableField(v.Field, v.FatherField).UpdateDefaultValue(sql[0]())
				} else {
					f.FieldList.FindTableField(v.Field, v.FatherField).UpdateDefaultValue(nil)
				}
			} else if v.FormType.IsTable() {
				list = append(list, v)
			} else {
				if len(sql) > 0 {
					// 新增菜單時(/menu/new)會執行
					// UpdateDefaultValue將預設值更新後回傳FormField(struct)
					list = append(list, *(v.UpdateDefaultValue(sql[0]())))
				} else {
					list = append(list, *(v.UpdateDefaultValue(nil)))
				}
			}
		}
	}
	// FillCustomContent(填寫自定義內容)對FormFields([]FormField)執行迴圈，判斷條件後設置FormField，最後回傳FormFields([]FormField)
	// 新增菜單api(/menu/new)不會執行下列動作(自定義)，所以list不變
	return list.FillCustomContent().RemoveNotShow()
}

type FormPreProcessFn func(values form.Values) form.Values

type FormPostFn func(values form.Values) error

type FormFields []FormField

type GroupFormFields []FormFields
type GroupFieldHeaders []string

func (f FormFields) Copy() FormFields {
	formList := make(FormFields, len(f))
	copy(formList, f)
	for i := 0; i < len(formList); i++ {
		formList[i].Options = make(FieldOptions, len(f[i].Options))
		for j := 0; j < len(f[i].Options); j++ {
			formList[i].Options[j] = FieldOption{
				Value:    f[i].Options[j].Value,
				Text:     f[i].Options[j].Text,
				TextHTML: f[i].Options[j].TextHTML,
				Selected: f[i].Options[j].Selected,
			}
		}
	}
	return formList
}

func (f FormFields) FindByFieldName(field string) FormField {
	for i := 0; i < len(f); i++ {
		if f[i].Field == field {
			return f[i]
		}
	}
	return FormField{}
}

// FillCustomContent(填寫自定義內容)對FormFields([]FormField)執行迴圈，判斷條件後設置FormField，最後回傳FormFields([]FormField)
func (f FormFields) FillCustomContent() FormFields {
	for i := range f {
		// 判斷是否是自定義
		// 新增菜單api(/menu/new)不會執行下列動作(自定義)
		if f[i].FormType.IsCustom() {
			// FillCustomContent(填寫自定義內容)判斷條件後設置FormField回傳
			f[i] = f[i].FillCustomContent()
		}
	}
	return f
}

func (f FormFields) Add(field FormField) FormFields {
	return append(f, field)
}
