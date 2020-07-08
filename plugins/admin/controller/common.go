package controller

import (
	"bytes"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	c "github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/modules/language"
	"github.com/GoAdminGroup/go-admin/modules/menu"
	"github.com/GoAdminGroup/go-admin/modules/service"
	"github.com/GoAdminGroup/go-admin/plugins/admin/models"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/constant"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"github.com/GoAdminGroup/go-admin/template/types"
	template2 "html/template"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

type Handler struct {
	config        *c.Config
	captchaConfig map[string]string
	services      service.List
	conn          db.Connection
	routes        context.RouterMap
	generators    table.GeneratorList // map[string]Generator
	operations    []context.Node
	navButtons    types.Buttons
	operationLock sync.Mutex
}

// 判斷參數cfg後設置Handler(struct)並回傳
func New(cfg ...Config) *Handler {
	if len(cfg) == 0 {
		return &Handler{
			operations: make([]context.Node, 0),
			navButtons: make(types.Buttons, 0),
		}
	}
	return &Handler{
		config:     cfg[0].Config,
		services:   cfg[0].Services,
		conn:       cfg[0].Connection,
		generators: cfg[0].Generators,
		operations: make([]context.Node, 0),
		navButtons: make(types.Buttons, 0),
	}
}

type Config struct {
	Config     *c.Config
	Services   service.List
	Connection db.Connection
	Generators table.GeneratorList
}

// 將參數cfg(struct)的值設置至Handler(struct)
func (h *Handler) UpdateCfg(cfg Config) {
	h.config = cfg.Config
	h.services = cfg.Services
	h.conn = cfg.Connection
	h.generators = cfg.Generators
}

// 將參數cap設置至Handler.captchaConfig(驗證碼配置)
func (h *Handler) SetCaptcha(cap map[string]string) {
	h.captchaConfig = cap
}

func (h *Handler) SetRoutes(r context.RouterMap) {
	h.routes = r
}

// BaseTable也屬於Table(interface)
// 先透過參數prefix取得Table(interface)，接著判斷條件後將[]context.Node加入至Handler.operations後回傳
func (h *Handler) table(prefix string, ctx *context.Context) table.Table {
	// 透過參數prefix取得h.generators[prefix]的值(func(ctx *context.Context) Table)
	t := h.generators[prefix](ctx)

	// 建立Invoker(Struct)並透過參數ctx取得UserModel，並且取得該user的role、權限與可用menu，最後檢查用戶權限
	// GetConnection取得匹配的service.Service然後轉換成Connection(interface)類別
	authHandler := auth.Middleware(db.GetConnection(h.services))

	// BaseTable也屬於Table(interface)
	// GetInfo在plugins\admin\modules\table\table.go，為Table(interface)的方法
	// 將參數值設置至base.Info(InfoPanel(struct)).primaryKey中後回傳
	// Callbacks類別為[]context.Node(struct)
	for _, cb := range t.GetInfo().Callbacks {
		// ContextNodeNeedAuth = need_auth
		if cb.Value[constant.ContextNodeNeedAuth] == 1 {
			// 判斷條件後將參數(類別context.Node)添加至Handler.operations
			h.addOperation(context.Node{
				Path:     cb.Path,
				Method:   cb.Method,
				Handlers: append([]context.Handler{authHandler}, cb.Handlers...),
			})
		} else {
			h.addOperation(context.Node{Path: cb.Path, Method: cb.Method, Handlers: cb.Handlers})
		}
	}

	// GetForm在plugins\admin\modules\table\table.go，為Table(interface)的方法
	// 將參數值設置至base.Form(InfoPanel(struct)).primaryKey中後回傳
	for _, cb := range t.GetForm().Callbacks {
		if cb.Value[constant.ContextNodeNeedAuth] == 1 {
			// 判斷條件後將參數(類別context.Node)添加至Handler.operations
			h.addOperation(context.Node{
				Path:     cb.Path,
				Method:   cb.Method,
				Handlers: append([]context.Handler{authHandler}, cb.Handlers...),
			})
		} else {
			h.addOperation(context.Node{Path: cb.Path, Method: cb.Method, Handlers: cb.Handlers})
		}
	}
	return t
}

func (h *Handler) route(name string) context.Router {
	return h.routes.Get(name)
}

func (h *Handler) routePath(name string, value ...string) string {
	return h.routes.Get(name).GetURL(value...)
}

func (h *Handler) routePathWithPrefix(name string, prefix string) string {
	return h.routePath(name, "prefix", prefix)
}

// 判斷條件後將參數(類別context.Node)添加至Handler.operations
func (h *Handler) addOperation(nodes ...context.Node) {
	h.operationLock.Lock()
	defer h.operationLock.Unlock()
	// TODO: 避免重复增加，第一次加入后，后面大部分会存在重复情况，以下循环可以优化
	addNodes := make([]context.Node, 0)
	for _, node := range nodes {
		// 在Handler.operations([]context.Node)執行迴圈，如果條件符合參數path、method則回傳true
		// 代表Handler.operations裡已經存在，則不添加
		if h.searchOperation(node.Path, node.Method) {
			continue
		}
		addNodes = append(addNodes, node)
	}
	h.operations = append(h.operations, addNodes...)
}

func (h *Handler) AddNavButton(btns ...types.Button) {
	for _, btn := range btns {
		h.navButtons = append(h.navButtons, btn)
		h.addOperation(btn.GetAction().GetCallbacks())
	}
}

func (h *Handler) AddNavButtonFront(btn types.Button) {
	h.navButtons = append(types.Buttons{btn}, h.navButtons...)
	h.addOperation(btn.GetAction().GetCallbacks())
}

// 在Handler.operations([]context.Node)執行迴圈，如果條件符合參數path、method則回傳true
func (h *Handler) searchOperation(path, method string) bool {
	for _, node := range h.operations {
		if node.Path == path && node.Method == method {
			return true
		}
	}
	return false
}

func (h *Handler) OperationHandler(path string, ctx *context.Context) bool {
	for _, node := range h.operations {
		if node.Path == path {
			for _, handler := range node.Handlers {
				handler(ctx)
			}
			return true
		}
	}
	return false
}

// 將參數設置至ExecuteParam(struct)，接著將給定的數據(types.Page(struct))寫入buf(struct)並輸出HTML至Context.response.Body
func (h *Handler) HTML(ctx *context.Context, user models.UserModel, panel types.Panel, animation ...bool) {
	buf := h.Execute(ctx, user, panel, animation...)
	ctx.HTML(http.StatusOK, buf.String())
}

// 將參數設置至ExecuteParam(struct)，接著將給定的數據(types.Page(struct))寫入buf(struct)並回傳
func (h *Handler) Execute(ctx *context.Context, user models.UserModel, panel types.Panel, animation ...bool) *bytes.Buffer {
	tmpl, tmplName := aTemplate().GetTemplate(isPjax(ctx))

	return template.Execute(template.ExecuteParam{
		User:      user,
		TmplName:  tmplName,
		Tmpl:      tmpl,
		Panel:     panel,
		Config:    *h.config,
		Menu:      menu.GetGlobalMenu(user, h.conn).SetActiveClass(h.config.URLRemovePrefix(ctx.Path())),
		Animation: len(animation) > 0 && animation[0] || len(animation) == 0,
		Buttons:   h.navButtons.CheckPermission(user),
	})
}

func isInfoUrl(s string) bool {
	reg, _ := regexp.Compile("(.*?)info/(.*?)$")
	sub := reg.FindStringSubmatch(s)
	return len(sub) > 2 && !strings.Contains(sub[2], "/")
}

func isNewUrl(s string, p string) bool {
	reg, _ := regexp.Compile("(.*?)info/" + p + "/new")
	return reg.MatchString(s)
}

func isEditUrl(s string, p string) bool {
	reg, _ := regexp.Compile("(.*?)info/" + p + "/edit")
	return reg.MatchString(s)
}

func (h *Handler) authSrv() *auth.TokenService {
	return auth.GetTokenService(h.services.Get(auth.TokenServiceKey))
}

func aAlert() types.AlertAttribute {
	return aTemplate().Alert()
}

func aForm() types.FormAttribute {
	// aTemplate判斷templateMap(map[string]Template)的key鍵是否參數globalCfg.Theme，有則回傳Template(interface)
	// Form在template\components\base.go中
	// Form為Template(interface)的方法，建立FormAttribute(struct)並設置值後回傳
	return aTemplate().Form()
}

func aRow() types.RowAttribute {
	// aTemplate判斷templateMap(map[string]Template)的key鍵是否參數globalCfg.Theme，有則回傳Template(interface)
	// Row在template\components\base.go中
	// Row為Template(interface)的方法，建立RowAttribute(struct)並設置值後回傳
	return aTemplate().Row()
}

func aCol() types.ColAttribute {
	// aTemplate判斷templateMap(map[string]Template)的key鍵是否參數globalCfg.Theme，有則回傳Template(interface)
	// Col在template\components\base.go中
	// Col為Template(interface)的方法，建立ColAttribute(struct)並設置值後回傳
	return aTemplate().Col()
}

func aButton() types.ButtonAttribute {
	return aTemplate().Button()
}

func aTree() types.TreeAttribute {
	// aTemplate判斷templateMap(map[string]Template)的key鍵是否參數globalCfg.Theme，有則回傳Template(interface)
	// Tree在template\components\base.go中
	// Tree為Template(interface)的方法，設置TreeAttribute(struct也是interface)並回傳
	return aTemplate().Tree()
}

func aTable() types.TableAttribute {
	return aTemplate().Table()
}

func aDataTable() types.DataTableAttribute {
	return aTemplate().DataTable()
}

func aBox() types.BoxAttribute {
	// aTemplate判斷templateMap(map[string]Template)的key鍵是否參數globalCfg.Theme，有則回傳Template(interface)
	// Box在template\components\base.go中
	// Box為Template(interface)的方法，設置BoxAttribute(struct也是interface)並回傳
	return aTemplate().Box()
}

func aTab() types.TabsAttribute {
	return aTemplate().Tabs()
}

// 判斷templateMap(map[string]Template)的key鍵是否參數globalCfg.Theme，有則回傳Template(interface)
func aTemplate() template.Template {
	// 判斷templateMap(map[string]Template)的key鍵是否參數globalCfg.Theme，有則回傳Template(interface)
	// GetTheme回傳globalCfg.Theme
	return template.Get(c.GetTheme())
}

func isPjax(ctx *context.Context) bool {
	return ctx.IsPjax()
}

func formFooter(page string, isHideEdit, isHideNew, isHideReset bool) template2.HTML {
	col1 := aCol().SetSize(types.SizeMD(2)).GetContent()

	var (
		checkBoxs  template2.HTML
		checkBoxJS template2.HTML

		editCheckBox = template.HTML(`
			<label class="pull-right" style="margin: 5px 10px 0 0;">
                <input type="checkbox" class="continue_edit" style="position: absolute; opacity: 0;"> ` + language.Get("continue editing") + `
            </label>`)
		newCheckBox = template.HTML(`
			<label class="pull-right" style="margin: 5px 10px 0 0;">
                <input type="checkbox" class="continue_new" style="position: absolute; opacity: 0;"> ` + language.Get("continue creating") + `
            </label>`)

		editWithNewCheckBoxJs = template.HTML(`$('.continue_edit').iCheck({checkboxClass: 'icheckbox_minimal-blue'}).on('ifChanged', function (event) {
		if (this.checked) {
			$('.continue_new').iCheck('uncheck');
			$('input[name="` + form.PreviousKey + `"]').val(location.href)
		} else {
			$('input[name="` + form.PreviousKey + `"]').val(previous_url_goadmin)
		}
	});	`)

		newWithEditCheckBoxJs = template.HTML(`$('.continue_new').iCheck({checkboxClass: 'icheckbox_minimal-blue'}).on('ifChanged', function (event) {
		if (this.checked) {
			$('.continue_edit').iCheck('uncheck');
			$('input[name="` + form.PreviousKey + `"]').val(location.href.replace('/edit', '/new'))
		} else {
			$('input[name="` + form.PreviousKey + `"]').val(previous_url_goadmin)
		}
	});`)
	)

	if page == "edit" {
		if isHideNew {
			newCheckBox = ""
			newWithEditCheckBoxJs = ""
		}
		if isHideEdit {
			editCheckBox = ""
			editWithNewCheckBoxJs = ""
		}
		checkBoxs = editCheckBox + newCheckBox
		checkBoxJS = `<script>	
	let previous_url_goadmin = $('input[name="` + form.PreviousKey + `"]').attr("value")
	` + editWithNewCheckBoxJs + newWithEditCheckBoxJs + `
</script>
`
	} else if page == "edit_only" && !isHideEdit {
		checkBoxs = editCheckBox
		checkBoxJS = template.HTML(`	<script>
	let previous_url_goadmin = $('input[name="` + form.PreviousKey + `"]').attr("value")
	$('.continue_edit').iCheck({checkboxClass: 'icheckbox_minimal-blue'}).on('ifChanged', function (event) {
		if (this.checked) {
			$('input[name="` + form.PreviousKey + `"]').val(location.href)
		} else {
			$('input[name="` + form.PreviousKey + `"]').val(previous_url_goadmin)
		}
	});
</script>
`)
	} else if page == "new" && !isHideNew {
		checkBoxs = newCheckBox
		checkBoxJS = template.HTML(`	<script>
	let previous_url_goadmin = $('input[name="` + form.PreviousKey + `"]').attr("value")
	$('.continue_new').iCheck({checkboxClass: 'icheckbox_minimal-blue'}).on('ifChanged', function (event) {
		if (this.checked) {
			$('input[name="` + form.PreviousKey + `"]').val(location.href)
		} else {
			$('input[name="` + form.PreviousKey + `"]').val(previous_url_goadmin)
		}
	});
</script>
`)
	}

	btn1 := aButton().SetType("submit").
		SetContent(language.GetFromHtml("Save")).
		SetThemePrimary().
		SetOrientationRight().
		GetContent()
	btn2 := template.HTML("")
	if !isHideReset {
		btn2 = aButton().SetType("reset").
			SetContent(language.GetFromHtml("Reset")).
			SetThemeWarning().
			SetOrientationLeft().
			GetContent()
	}
	col2 := aCol().SetSize(types.SizeMD(8)).
		SetContent(btn1 + checkBoxs + btn2 + checkBoxJS).GetContent()
	return col1 + col2
}

func filterFormFooter(infoUrl string) template2.HTML {
	col1 := aCol().SetSize(types.SizeMD(2)).GetContent()
	btn1 := aButton().SetType("submit").
		SetContent(icon.Icon(icon.Search, 2) + language.GetFromHtml("search")).
		SetThemePrimary().
		SetSmallSize().
		SetOrientationLeft().
		SetLoadingText(icon.Icon(icon.Spinner, 1) + language.GetFromHtml("search")).
		GetContent()
	btn2 := aButton().SetType("reset").
		SetContent(icon.Icon(icon.Undo, 2) + language.GetFromHtml("reset")).
		SetThemeDefault().
		SetOrientationLeft().
		SetSmallSize().
		SetHref(infoUrl).
		SetMarginLeft(12).
		GetContent()
	col2 := aCol().SetSize(types.SizeMD(8)).
		SetContent(btn1 + btn2).GetContent()
	return col1 + col2
}

func formContent(form types.FormAttribute, isTab bool) template2.HTML {
	if isTab {
		return form.GetContent()
	}
	return aBox().
		SetHeader(form.GetDefaultBoxHeader()).
		WithHeadBorder().
		SetStyle(" ").
		SetBody(form.GetContent()).
		GetContent()
}

func detailContent(form types.FormAttribute, editUrl, deleteUrl string) template2.HTML {
	return aBox().
		SetHeader(form.GetDetailBoxHeader(editUrl, deleteUrl)).
		WithHeadBorder().
		SetBody(form.GetContent()).
		GetContent()
}

// menuFormContent(菜單表單內容)將符合BoxAttribute.TemplateList["box"](map[string]string)的值加入text(string)，接著將參數及功能添加給新的模板並解析為模板的主體
// 將參數compo寫入buffer(bytes.Buffer)中最後輸出HTML
func menuFormContent(form types.FormAttribute) template2.HTML {
	// aBox在plugins\admin\controller\common.go中
	// aBox設置BoxAttribute(是struct也是interface)
	// SetHeader、SetStyle、SetBody、GetContent、WithHeadBorder都為BoxAttribute的方法
	// 都是將參數值設置至BoxAttribute(struct)
	// GetContent先依判斷條件設置BoxAttribute.Style
	// 將符合BoxAttribute.TemplateList["box"](map[string]string)的值加入text(string)，接著將參數及功能添加給新的模板並解析為模板的主體
	// 將參數compo寫入buffer(bytes.Buffer)中最後輸出HTML
	return aBox().
		SetHeader(form.GetBoxHeaderNoButton()).
		SetStyle(" ").
		WithHeadBorder().
		SetBody(form.GetContent()).
		GetContent()
}
