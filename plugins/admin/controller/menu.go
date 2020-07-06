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
func (h *Handler) ShowMenu(ctx *context.Context) {
	// getMenuInfoPanel(���omenu��T���O)
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
func (h *Handler) ShowEditMenu(ctx *context.Context) {

	// �ˬdurl����id�Ѽ�(�]���O�n�s��Y��menu�A�ݭn�]�mid = ?)
	if ctx.Query("id") == "" {
		h.getMenuInfoPanel(ctx, template.Get(h.config.Theme).Alert().Warning(errors.WrongID))

		ctx.AddHeader("Content-Type", "text/html; charset=utf-8")
		ctx.AddHeader(constant.PjaxUrlHeader, h.routePath("menu"))
		return
	}

	model := h.table("menu", ctx)
	formInfo, err := model.GetDataWithId(parameter.BaseParam().WithPKs(ctx.Query("id")))

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
	// GetMenuDeleteParam�NContext.UserValue(map[string]interface{})[delete_menu_param]�����ഫ��MenuDeleteParam(struct)���O
	// MenuWithId�bplugins\admin\models\menu.go
	// MenuWithId�z�L�ѼƱNid�Ptablename(goadmin_menu)�]�m��MenuModel(struct)��^��
	// SetConn�N�Ѽ�h.conn�]�m��MenuModel.Base.Conn
	models.MenuWithId(guard.GetMenuDeleteParam(ctx).Id).SetConn(h.conn).Delete()
	response.OkWithMsg(ctx, language.Get("delete succeed"))
}

// EditMenu edit the menu of given id.
// �NContext.UserValue(map[string]interface{})[edit_menu_param]�����ഫ��MenuEditParam(struct)���O
// ���Ngoadmin_role_menu��ƪ�menu_id = MenuModel.Id����ƧR���A���ۦp�G���bmultipart/form-data���]�wroles[]�ȡA�ˬd�����N�Ѽ�roleId(role_id)�PMenuModel.Id(menu_id)�[�Jgoadmin_role_menu��ƪ�
// ���۱Ngoadmin_menu��ƪ����id = MenuModel.Id����Ƴz�L�Ѽ�(��multipart/form-data�]�m)��s
func (h *Handler) EditMenu(ctx *context.Context) {
	// �NContext.UserValue(map[string]interface{})[edit_menu_param]�����ഫ��MenuEditParam(struct)���O
	param := guard.GetMenuEditParam(ctx)

	// �P�_MenuNewParam.Alert�O�_�X�{ĵ�i(���O�ŭ�)
	if param.HasAlert() {
		h.getMenuInfoPanel(ctx, param.Alert)
		ctx.AddHeader("Content-Type", "text/html; charset=utf-8")
		ctx.AddHeader(constant.PjaxUrlHeader, h.routePath("menu"))
		return
	}

	// MenuWithId�z�L�ѼƱNparam.id�Ptablename(goadmin_menu)�]�m��MenuModel(struct)��^��
	// SetConn�N�Ѽ�h.conn�]�m��MenuModel.Base.Conn
	menuModel := models.MenuWithId(param.Id).SetConn(h.conn)

	// TODO: use transaction
	// DeleteRoles�R��goadmin_role_menu��ƪ�menu_id = MenuModel.Id���󪺸��
	deleteRolesErr := menuModel.DeleteRoles()
	if db.CheckError(deleteRolesErr, db.DELETE) {
		formInfo, _ := h.table("menu", ctx).GetDataWithId(parameter.BaseParam().WithPKs(param.Id))
		h.showEditMenu(ctx, formInfo, deleteRolesErr)
		ctx.AddHeader(constant.PjaxUrlHeader, h.routePath("menu"))
		return
	}

	// �p�Gmultipart/form-data���]�wroles[]��
	// AddRole���ˬdgoadmin_role_menu����A���۱N�Ѽ�roleId(role_id)�PMenuModel.Id(menu_id)�[�Jgoadmin_role_menu��ƪ�
	for _, roleId := range param.Roles {
		_, addRoleErr := menuModel.AddRole(roleId)
		if db.CheckError(addRoleErr, db.INSERT) {
			formInfo, _ := h.table("menu", ctx).GetDataWithId(parameter.BaseParam().WithPKs(param.Id))
			h.showEditMenu(ctx, formInfo, addRoleErr)
			ctx.AddHeader(constant.PjaxUrlHeader, h.routePath("menu"))
			return
		}
	}

	// Update�Ngoadmin_menu��ƪ����id = MenuModel.Id����Ƴz�L�Ѽ�(��multipart/form-data�]�m)��s
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
// �NContext.UserValue(map[string]interface{})[new_menu_param]�����ഫ��MenuNewParam(struct)���O�A���۱NMenuNewParam(struct)�ȷs�W�ܸ�ƪ�(MenuModel.Base.TableName(goadmin_menu))��
// �̫�p�Gmultipart/form-data���]�wroles[]�ȡA�ˬd�����N�Ѽ�roleId(role_id)�PMenuModel.Id(menu_id)�[�Jgoadmin_role_menu��ƪ�
func (h *Handler) NewMenu(ctx *context.Context) {
    // �NContext.UserValue(map[string]interface{})[new_menu_param]�����ഫ��MenuNewParam(struct)���O
	param := guard.GetMenuNewParam(ctx)

	// �P�_MenuNewParam.Alert�O�_�X�{ĵ�i(���O�ŭ�)
	if param.HasAlert() {
		h.getMenuInfoPanel(ctx, param.Alert)
		ctx.AddHeader("Content-Type", "text/html; charset=utf-8")
		// PjaxUrlHeader = X-PJAX-Url
		ctx.AddHeader(constant.PjaxUrlHeader, h.routePath("menu"))
		return
	}

	// �z�L�Ѽ�ctx�^�ǥثe�n�J���Τ�(Context.UserValue["user"])���ഫ��UserModel
	user := auth.Auth(ctx)

	// TODO: use transaction
	// Menu�NMenuModel(struct).Base.TableName�]�mgoadmin_menu��^��
	// SetConn�N�Ѽ�h.conn�]�m��MenuModel.Base.Conn
	// New�N�Ѽƭȷs�W�ܸ�ƪ�(MenuModel.Base.TableName(goadmin_menu))���A�̫�N�Ѽƭȳ��]�m�bMenuModel��
	menuModel, createErr := models.Menu().SetConn(h.conn).
	    // GetGlobalMenu�^�ǰѼ�user(struct)��Menu(�]�mmenuList�BmenuOption�BMaxOrder)
		New(param.Title, param.Icon, param.Uri, param.Header, param.ParentId, (menu.GetGlobalMenu(user, h.conn)).MaxOrder+1)

	if db.CheckError(createErr, db.INSERT) {
		h.showNewMenu(ctx, createErr)
		return
	}

	// �p�Gmultipart/form-data���]�wroles[]��
	// AddRole �ˬdgoadmin_role_menu��ƪ�̬O�_���ŦXrole_id = �Ѽ�roleId�Pmenu_id = MenuModel.Id������A���۱N�Ѽ�roleId(role_id)�PMenuModel.Id(menu_id)�[�Jgoadmin_role_menu��ƪ�
	for _, roleId := range param.Roles {
		_, addRoleErr := menuModel.AddRole(roleId)
		if db.CheckError(addRoleErr, db.INSERT) {
			h.showNewMenu(ctx, addRoleErr)
			return
		}
	}

	 // GetGlobalMenu�^�ǰѼ�user(struct)��Menu(�]�mmenuList�BmenuOption�BMaxOrder)
	 // AddMaxOrder�NMenu.MaxOrder+1
	menu.GetGlobalMenu(user, h.conn).AddMaxOrder()

	h.getMenuInfoPanel(ctx, "")
	ctx.AddHeader("Content-Type", "text/html; charset=utf-8")
	// PjaxUrlHeader = X-PJAX-Url
	ctx.AddHeader(constant.PjaxUrlHeader, h.routePath("menu"))
}

// MenuOrder change the order of menu items.
// ���omultipart/form-data����_order�Ѽƫ���menu����
func (h *Handler) MenuOrder(ctx *context.Context) {

	var data []map[string]interface{}
	// FormValue���omultipart/form-data����_order�Ѽƫ�ѽX��data([]map[string]interface{})
	_ = json.Unmarshal([]byte(ctx.FormValue("_order")), &data)

	// Menu�NMenuModel(struct).Base.TableName�]�mgoadmin_menu��^��
	// SetConn�N�Ѽ�con�]�m��MenuModel.Base.Conn
	// ResetOrder���menu������
	models.Menu().SetConn(h.conn).ResetOrder([]byte(ctx.FormValue("_order")))

	// �^��code�Bmsg
	response.Ok(ctx)
}

func (h *Handler) getMenuInfoPanel(ctx *context.Context, alert template2.HTML) {
	// �z�L�Ѽ�ctx�^�ǥثe�n�J���Τ�(Context.UserValue["user"])���ഫ��UserModel
	user := auth.Auth(ctx)

	// aTemplate�P�_templateMap(map[string]Template)��key��O�_�Ѽ�globalCfg.Theme�A���h�^��Template(interface)
	// Tree��Template(interface)����k
	tree := aTree().
		SetTree((menu.GetGlobalMenu(user, h.conn)).List).
		SetEditUrl(h.routePath("menu_edit_show")).
		SetUrlPrefix(h.config.Prefix()).
		SetDeleteUrl(h.routePath("menu_delete")).
		SetOrderUrl(h.routePath("menu_order")).
		GetContent()

	header := aTree().GetTreeHeader()
	box := aBox().SetHeader(header).SetBody(tree).GetContent()
	col1 := aCol().SetSize(types.SizeMD(6)).SetContent(box).GetContent()

	formInfo := h.table("menu", ctx).GetNewForm()

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

	col2 := aCol().SetSize(types.SizeMD(6)).SetContent(newForm).GetContent()

	row := aRow().SetContent(col1 + col2).GetContent()

	h.HTML(ctx, user, types.Panel{
		Content:     alert + row,
		Description: "Menus Manage",
		Title:       "Menus Manage",
	})
}
