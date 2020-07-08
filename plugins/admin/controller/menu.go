package controller

import (
	"encoding/json"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/modules/errors"
	"github.com/GoAdminGroup/go-admin/modules/language"
	"github.com/GoAdminGroup/go-admin/modules/menu"
	"github.com/GoAdminGroup/go-admin/plugins/admin/models"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/constant"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/guard"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/parameter"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/response"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
	template2 "html/template"
)

// ShowMenu show menu info page.
// 分別處理上下半部表單的HTML語法，最後結合並輸出HTML
func (h *Handler) ShowMenu(ctx *context.Context) {
	// getMenuInfoPanel(取得菜單資訊面板)分別處理上下半部表單的HTML語法，最後結合並輸出HTML
	h.getMenuInfoPanel(ctx, "")
}

// ShowNewMenu show new menu page.
func (h *Handler) ShowNewMenu(ctx *context.Context) {
	h.showNewMenu(ctx, nil)
}

func (h *Handler) showNewMenu(ctx *context.Context, err error) {
	panel := h.table("menu", ctx)

	formInfo := panel.GetNewForm()

	user := auth.Auth(ctx)

	var alert template2.HTML

	if err != nil {
		alert = aAlert().Warning(err.Error())
	}

	h.HTML(ctx, user, types.Panel{
		Content: alert + formContent(aForm().
			SetContent(formInfo.FieldList).
			SetTabContents(formInfo.GroupFieldList).
			SetTabHeaders(formInfo.GroupFieldHeaders).
			SetPrefix(h.config.PrefixFixSlash()).
			SetPrimaryKey(panel.GetPrimaryKey().Name).
			SetUrl(h.routePath("menu_edit")).
			SetHiddenFields(map[string]string{
				form2.TokenKey:    h.authSrv().AddToken(),
				form2.PreviousKey: h.routePath("menu"),
			}).
			SetOperationFooter(formFooter("new", false, false, false)), false),
		Description: template2.HTML(panel.GetForm().Description),
		Title:       template2.HTML(panel.GetForm().Title),
	})
}

// ShowEditMenu show edit menu page.
// 
func (h *Handler) ShowEditMenu(ctx *context.Context) {

	// 檢查url中的id參數(因為是要編輯某個menu，需要設置id = ?)
	if ctx.Query("id") == "" {
		h.getMenuInfoPanel(ctx, template.Get(h.config.Theme).Alert().Warning(errors.WrongID))

		ctx.AddHeader("Content-Type", "text/html; charset=utf-8")
		ctx.AddHeader(constant.PjaxUrlHeader, h.routePath("menu"))
		return
	}

	// 先透過參數"menu"取得Table(interface)，接著判斷條件後將[]context.Node加入至Handler.operations後回傳
	model := h.table("menu", ctx)

	// BaseParam設置值(頁數及頁數Size)至Parameters(struct)並回傳
	// WithPKs將參數(多個string)結合並設置至Parameters.Fields["__pk"]後回傳
	// GetDataWithId在plugins\admin\modules\table\default.go
	

	// 透過參數ctx回傳目前登入的用戶(Context.UserValue["user"])並轉換成UserModel
	user := auth.Auth(ctx)

	if err != nil {
		h.HTML(ctx, user, types.Panel{
			Content:     aAlert().Warning(err.Error()),
			Description: template2.HTML(model.GetForm().Description),
			Title:       template2.HTML(model.GetForm().Title),
		})
		return
	}

	h.showEditMenu(ctx, formInfo, nil)
}

func (h *Handler) showEditMenu(ctx *context.Context, formInfo table.FormInfo, err error) {

	var alert template2.HTML

	if err != nil {
		alert = aAlert().Warning(err.Error())
	}


	// aForm在plugins\admin\controller\common.go中
	// aForm設置FormAttribute(是struct也是interface)
	// 將參數值設置至FormFields(struct)
	// 判斷條件後，將FormFields添加至FormAttribute.ContentList([]FormFields)
    // 接著將符合FormAttribute.TemplateList["components/多個參數"](map[string]string)的值加入text(string)，接著將參數及功能添加給新的模板並解析為模板的主體
	// 將參數compo寫入buffer(bytes.Buffer)中最後輸出HTML回傳
	h.HTML(ctx, auth.Auth(ctx), types.Panel{
		Content: alert + formContent(aForm().
			SetContent(formInfo.FieldList).
			SetTabContents(formInfo.GroupFieldList).
			SetTabHeaders(formInfo.GroupFieldHeaders).
			SetPrefix(h.config.PrefixFixSlash()).
			SetPrimaryKey(h.table("menu", ctx).GetPrimaryKey().Name).
			SetUrl(h.routePath("menu_edit")).
			SetOperationFooter(formFooter("edit", false, false, false)).
			SetHiddenFields(map[string]string{
				form2.TokenKey:    h.authSrv().AddToken(),
				form2.PreviousKey: h.routePath("menu"),
			}), false),
		Description: template2.HTML(formInfo.Description),
		Title:       template2.HTML(formInfo.Title),
	})
	return
}

// DeleteMenu delete the menu of given id.
func (h *Handler) DeleteMenu(ctx *context.Context) {
	// GetMenuDeleteParam將Context.UserValue(map[string]interface{})[delete_menu_param]的值轉換成MenuDeleteParam(struct)類別
	// MenuWithId在plugins\admin\models\menu.go
	// MenuWithId透過參數將id與tablename(goadmin_menu)設置至MenuModel(struct)後回傳
	// SetConn將參數h.conn設置至MenuModel.Base.Conn
	models.MenuWithId(guard.GetMenuDeleteParam(ctx).Id).SetConn(h.conn).Delete()
	response.OkWithMsg(ctx, language.Get("delete succeed"))
}

// EditMenu edit the menu of given id.
// 將Context.UserValue(map[string]interface{})[edit_menu_param]的值轉換成MenuEditParam(struct)類別
// 先將goadmin_role_menu資料表中menu_id = MenuModel.Id的資料刪除，接著如果有在multipart/form-data有設定roles[]值，檢查條件後將參數roleId(role_id)與MenuModel.Id(menu_id)加入goadmin_role_menu資料表
// 接著將goadmin_menu資料表條件為id = MenuModel.Id的資料透過參數(由multipart/form-data設置)更新
func (h *Handler) EditMenu(ctx *context.Context) {
	// 將Context.UserValue(map[string]interface{})[edit_menu_param]的值轉換成MenuEditParam(struct)類別
	param := guard.GetMenuEditParam(ctx)

	// 判斷MenuNewParam.Alert是否出現警告(不是空值)
	if param.HasAlert() {
		h.getMenuInfoPanel(ctx, param.Alert)
		ctx.AddHeader("Content-Type", "text/html; charset=utf-8")
		ctx.AddHeader(constant.PjaxUrlHeader, h.routePath("menu"))
		return
	}

	// MenuWithId透過參數將param.id與tablename(goadmin_menu)設置至MenuModel(struct)後回傳
	// SetConn將參數h.conn設置至MenuModel.Base.Conn
	menuModel := models.MenuWithId(param.Id).SetConn(h.conn)

	// TODO: use transaction
	// DeleteRoles刪除goadmin_role_menu資料表中menu_id = MenuModel.Id條件的資料
	deleteRolesErr := menuModel.DeleteRoles()
	if db.CheckError(deleteRolesErr, db.DELETE) {
		formInfo, _ := h.table("menu", ctx).GetDataWithId(parameter.BaseParam().WithPKs(param.Id))
		h.showEditMenu(ctx, formInfo, deleteRolesErr)
		ctx.AddHeader(constant.PjaxUrlHeader, h.routePath("menu"))
		return
	}

	// 如果multipart/form-data有設定roles[]值
	// AddRole先檢查goadmin_role_menu條件，接著將參數roleId(role_id)與MenuModel.Id(menu_id)加入goadmin_role_menu資料表
	for _, roleId := range param.Roles {
		_, addRoleErr := menuModel.AddRole(roleId)
		if db.CheckError(addRoleErr, db.INSERT) {
			formInfo, _ := h.table("menu", ctx).GetDataWithId(parameter.BaseParam().WithPKs(param.Id))
			h.showEditMenu(ctx, formInfo, addRoleErr)
			ctx.AddHeader(constant.PjaxUrlHeader, h.routePath("menu"))
			return
		}
	}

	// Update將goadmin_menu資料表條件為id = MenuModel.Id的資料透過參數(由multipart/form-data設置)更新
	_, updateErr := menuModel.Update(param.Title, param.Icon, param.Uri, param.Header, param.ParentId)

	if db.CheckError(updateErr, db.UPDATE) {
		formInfo, _ := h.table("menu", ctx).GetDataWithId(parameter.BaseParam().WithPKs(param.Id))
		h.showEditMenu(ctx, formInfo, updateErr)
		ctx.AddHeader(constant.PjaxUrlHeader, h.routePath("menu"))
		return
	}

	h.getMenuInfoPanel(ctx, "")
	ctx.AddHeader("Content-Type", "text/html; charset=utf-8")
	// PjaxUrlHeader = X-PJAX-Url
	ctx.AddHeader(constant.PjaxUrlHeader, h.routePath("menu"))
}

// NewMenu create a new menu item.
// 將Context.UserValue(map[string]interface{})[new_menu_param]的值轉換成MenuNewParam(struct)類別，接著將MenuNewParam(struct)值新增至資料表(MenuModel.Base.TableName(goadmin_menu))中
// 最後如果multipart/form-data有設定roles[]值，檢查條件後將參數roleId(role_id)與MenuModel.Id(menu_id)加入goadmin_role_menu資料表
func (h *Handler) NewMenu(ctx *context.Context) {
    // 將Context.UserValue(map[string]interface{})[new_menu_param]的值轉換成MenuNewParam(struct)類別
	param := guard.GetMenuNewParam(ctx)

	// 判斷MenuNewParam.Alert是否出現警告(不是空值)
	if param.HasAlert() {
		h.getMenuInfoPanel(ctx, param.Alert)
		ctx.AddHeader("Content-Type", "text/html; charset=utf-8")
		// PjaxUrlHeader = X-PJAX-Url
		ctx.AddHeader(constant.PjaxUrlHeader, h.routePath("menu"))
		return
	}

	// 透過參數ctx回傳目前登入的用戶(Context.UserValue["user"])並轉換成UserModel
	user := auth.Auth(ctx)

	// TODO: use transaction
	// Menu將MenuModel(struct).Base.TableName設置goadmin_menu後回傳
	// SetConn將參數h.conn設置至MenuModel.Base.Conn
	// New將參數值新增至資料表(MenuModel.Base.TableName(goadmin_menu))中，最後將參數值都設置在MenuModel中
	menuModel, createErr := models.Menu().SetConn(h.conn).
	    // GetGlobalMenu回傳參數user(struct)的Menu(設置menuList、menuOption、MaxOrder)
		New(param.Title, param.Icon, param.Uri, param.Header, param.ParentId, (menu.GetGlobalMenu(user, h.conn)).MaxOrder+1)

	if db.CheckError(createErr, db.INSERT) {
		h.showNewMenu(ctx, createErr)
		return
	}

	// 如果multipart/form-data有設定roles[]值
	// AddRole 檢查goadmin_role_menu資料表裡是否有符合role_id = 參數roleId與menu_id = MenuModel.Id的條件，接著將參數roleId(role_id)與MenuModel.Id(menu_id)加入goadmin_role_menu資料表
	for _, roleId := range param.Roles {
		_, addRoleErr := menuModel.AddRole(roleId)
		if db.CheckError(addRoleErr, db.INSERT) {
			h.showNewMenu(ctx, addRoleErr)
			return
		}
	}

	 // GetGlobalMenu回傳參數user(struct)的Menu(設置menuList、menuOption、MaxOrder)
	 // AddMaxOrder將Menu.MaxOrder+1
	menu.GetGlobalMenu(user, h.conn).AddMaxOrder()

	h.getMenuInfoPanel(ctx, "")
	ctx.AddHeader("Content-Type", "text/html; charset=utf-8")
	// PjaxUrlHeader = X-PJAX-Url
	ctx.AddHeader(constant.PjaxUrlHeader, h.routePath("menu"))
}

// MenuOrder change the order of menu items.
// 取得multipart/form-data中的_order參數後更改menu順序
func (h *Handler) MenuOrder(ctx *context.Context) {

	var data []map[string]interface{}
	// FormValue取得multipart/form-data中的_order參數後解碼至data([]map[string]interface{})
	_ = json.Unmarshal([]byte(ctx.FormValue("_order")), &data)

	// Menu將MenuModel(struct).Base.TableName設置goadmin_menu後回傳
	// SetConn將參數con設置至MenuModel.Base.Conn
	// ResetOrder更改menu的順序
	models.Menu().SetConn(h.conn).ResetOrder([]byte(ctx.FormValue("_order")))

	// 回傳code、msg
	response.Ok(ctx)
}

// getMenuInfoPanel(取得菜單資訊面板)分別處理上下半部表單的HTML語法，最後結合並輸出HTML
func (h *Handler) getMenuInfoPanel(ctx *context.Context, alert template2.HTML) {
	// 透過參數ctx回傳目前登入的用戶(Context.UserValue["user"])並轉換成UserModel
	user := auth.Auth(ctx)

	// aTree在plugins\admin\controller\common.go中
	// aTree判斷templateMap(map[string]Template)的key鍵是否參數globalCfg.Theme，有則回傳Template(interface)
	// 接著設置TreeAttribute(struct也是interface)並回傳
	// SetEditUrl、SetUrlPrefix、SetDeleteUrl、SetOrderUrl、GetContent都為TreeAttribute的方法
	// 都是將參數值設置至TreeAttribute(struct)
	// GetContent首先將符合compo.TemplateList["components/tree"](map[string]string)的值加入text(string)，接著將參數及功能添加給新的模板並解析為模板的主體
	// 將參數compo寫入buffer(bytes.Buffer)中最後輸出HTML
	// tree為/admin/menu中顯示的樹狀圖前端語法(尋找HTML id="tree-model")
	tree := aTree().
		SetTree((menu.GetGlobalMenu(user, h.conn)).List).
		SetEditUrl(h.routePath("menu_edit_show")).
		SetUrlPrefix(h.config.Prefix()).
		SetDeleteUrl(h.routePath("menu_delete")).
		SetOrderUrl(h.routePath("menu_order")).
		GetContent()

	// GetTreeHeader為TreeAttribute的方法
	// 首先將符合compo.TemplateList["components/tree-header"](map[string]string)的值加入text(string)，接著將參數及功能添加給新的模板並解析為模板的主體
	// 將參數compo寫入buffer(bytes.Buffer)中最後輸出HTML
	// header為/admin/menu中為樹狀圖上面的四個按鈕前端語法
	// header為尋找class="btn-group"
	header := aTree().GetTreeHeader()

	// aBox在plugins\admin\controller\common.go中
	// aBox設置BoxAttribute(是struct也是interface)
	// SetHeader、SetBody、GetContent都為BoxAttribute的方法
	// 都是將參數值設置至BoxAttribute(struct)
	// GetContent先依判斷條件設置BoxAttribute.Style
	// 將符合BoxAttribute.TemplateList["box"](map[string]string)的值加入text(string)，接著將參數及功能添加給新的模板並解析為模板的主體
	// 將參數compo寫入buffer(bytes.Buffer)中最後輸出HTML
	// box為尋找class="box box-"，標頭為header、內容為tree(上半部表單的語法)
	box := aBox().SetHeader(header).SetBody(tree).GetContent()

	// aCol在plugins\admin\controller\common.go中
	// aCol設置ColAttribute(是struct也是interface)
	// SetSize、SetContent、GetContent都是ColAttribute的方法
	// 都是將參數值設置至ColAttribute(struct)
	// GetContent將符合ColAttribute.TemplateList["col"](map[string]string)的值加入text(string)，接著將參數及功能添加給新的模板並解析為模板的主體
	// 將參數compo寫入buffer(bytes.Buffer)中最後輸出HTML
	// col1為尋找class="col-md-6"，內容為box(上半部表單圖的HTML語法)
	col1 := aCol().SetSize(types.SizeMD(6)).SetContent(box).GetContent()

	// BaseTable也屬於Table(interface)
	// table先透過參數"menu"取得Table(interface)，接著判斷條件後將[]context.Node加入至Handler.operations後回傳
	// GetNewForm在plugins\admin\modules\table\default.go
	// GetNewForm(取得新表單)判斷條件(TabGroups)後，設置FormInfo(struct)後並回傳
	formInfo := h.table("menu", ctx).GetNewForm()

	// aForm在plugins\admin\controller\common.go中
	// aForm設置FormAttribute(是struct也是interface)
	// 將參數值設置至FormFields(struct)
	// 判斷條件後，將FormFields添加至FormAttribute.ContentList([]FormFields)
    // 接著將符合FormAttribute.TemplateList["components/多個參數"](map[string]string)的值加入text(string)，接著將參數及功能添加給新的模板並解析為模板的主體
	// 將參數compo寫入buffer(bytes.Buffer)中最後輸出HTML回傳
	// menuFormContent(菜單表單內容)將符合BoxAttribute.TemplateList["box"](map[string]string)的值加入text(string)，接著將參數及功能添加給新的模板並解析為模板的主體
	// 將參數compo寫入buffer(bytes.Buffer)中最後輸出HTML
	// newForm為下半部新建表單的HTML語法
	newForm := menuFormContent(aForm().
		SetPrefix(h.config.PrefixFixSlash()).
		SetUrl(h.routePath("menu_new")).
		SetPrimaryKey(h.table("menu", ctx).GetPrimaryKey().Name).
		SetHiddenFields(map[string]string{
			form2.TokenKey:    h.authSrv().AddToken(),
			form2.PreviousKey: h.routePath("menu"),
		}).
		SetOperationFooter(formFooter("menu", false, false, false)).
		SetTitle("New").
		SetContent(formInfo.FieldList).
		SetTabContents(formInfo.GroupFieldList).
		SetTabHeaders(formInfo.GroupFieldHeaders))

	// aCol在plugins\admin\controller\common.go中
	// aCol設置ColAttribute(是struct也是interface)
	// SetSize、SetContent、GetContent都是ColAttribute的方法
	// GetContent將符合ColAttribute.TemplateList["col"](map[string]string)的值加入text(string)，接著將參數及功能添加給新的模板並解析為模板的主體
	// 將參數compo寫入buffer(bytes.Buffer)中最後輸出HTML
	// col2為尋找class="col-md-6"，內容為newForm(下半部新建表單的HTML語法)
	col2 := aCol().SetSize(types.SizeMD(6)).SetContent(newForm).GetContent()

	// aRow在plugins\admin\controller\common.go中
	// aRow設置RowAttribute(是struct也是interface)
	// 在template\components\composer.go
	// 首先將符合RowAttribute.TemplateList["components/row"](map[string]string)的值加入text(string)，接著將參數及功能添加給新的模板並解析為模板的主體
	// 將參數compo寫入buffer(bytes.Buffer)中最後輸出HTML
	// row為尋找class="row"，內容為上下半部所有表單的HTML語法
	row := aRow().SetContent(col1 + col2).GetContent()

	// 輸出HTML
	h.HTML(ctx, user, types.Panel{
		Content:     alert + row,
		Description: "Menus Manage",
		Title:       "Menus Manage",
	})
}
