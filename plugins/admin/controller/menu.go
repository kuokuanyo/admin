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
// ���O�B�z�W�U�b����檺HTML�y�k�A�̫ᵲ�X�ÿ�XHTML
func (h *Handler) ShowMenu(ctx *context.Context) {
	// getMenuInfoPanel(���o����T���O)���O�B�z�W�U�b����檺HTML�y�k�A�̫ᵲ�X�ÿ�XHTML
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

	// �ˬdurl����id�Ѽ�(�]���O�n�s��Y��menu�A�ݭn�]�mid = ?)
	if ctx.Query("id") == "" {
		h.getMenuInfoPanel(ctx, template.Get(h.config.Theme).Alert().Warning(errors.WrongID))

		ctx.AddHeader("Content-Type", "text/html; charset=utf-8")
		ctx.AddHeader(constant.PjaxUrlHeader, h.routePath("menu"))
		return
	}

	// ���z�L�Ѽ�"menu"���oTable(interface)�A���ۧP�_�����N[]context.Node�[�J��Handler.operations��^��
	model := h.table("menu", ctx)

	// BaseParam�]�m��(���Ƥέ���Size)��Parameters(struct)�æ^��
	// WithPKs�N�Ѽ�(�h��string)���X�ó]�m��Parameters.Fields["__pk"]��^��
	// GetDataWithId�bplugins\admin\modules\table\default.go
	

	// �z�L�Ѽ�ctx�^�ǥثe�n�J���Τ�(Context.UserValue["user"])���ഫ��UserModel
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


	// aForm�bplugins\admin\controller\common.go��
	// aForm�]�mFormAttribute(�Ostruct�]�Ointerface)
	// �N�Ѽƭȳ]�m��FormFields(struct)
	// �P�_�����A�NFormFields�K�[��FormAttribute.ContentList([]FormFields)
    // ���۱N�ŦXFormAttribute.TemplateList["components/�h�ӰѼ�"](map[string]string)���ȥ[�Jtext(string)�A���۱N�ѼƤΥ\��K�[���s���ҪO�øѪR���ҪO���D��
	// �N�Ѽ�compo�g�Jbuffer(bytes.Buffer)���̫��XHTML�^��
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

// getMenuInfoPanel(���o����T���O)���O�B�z�W�U�b����檺HTML�y�k�A�̫ᵲ�X�ÿ�XHTML
func (h *Handler) getMenuInfoPanel(ctx *context.Context, alert template2.HTML) {
	// �z�L�Ѽ�ctx�^�ǥثe�n�J���Τ�(Context.UserValue["user"])���ഫ��UserModel
	user := auth.Auth(ctx)

	// aTree�bplugins\admin\controller\common.go��
	// aTree�P�_templateMap(map[string]Template)��key��O�_�Ѽ�globalCfg.Theme�A���h�^��Template(interface)
	// ���۳]�mTreeAttribute(struct�]�Ointerface)�æ^��
	// SetEditUrl�BSetUrlPrefix�BSetDeleteUrl�BSetOrderUrl�BGetContent����TreeAttribute����k
	// ���O�N�Ѽƭȳ]�m��TreeAttribute(struct)
	// GetContent�����N�ŦXcompo.TemplateList["components/tree"](map[string]string)���ȥ[�Jtext(string)�A���۱N�ѼƤΥ\��K�[���s���ҪO�øѪR���ҪO���D��
	// �N�Ѽ�compo�g�Jbuffer(bytes.Buffer)���̫��XHTML
	// tree��/admin/menu����ܪ��𪬹ϫe�ݻy�k(�M��HTML id="tree-model")
	tree := aTree().
		SetTree((menu.GetGlobalMenu(user, h.conn)).List).
		SetEditUrl(h.routePath("menu_edit_show")).
		SetUrlPrefix(h.config.Prefix()).
		SetDeleteUrl(h.routePath("menu_delete")).
		SetOrderUrl(h.routePath("menu_order")).
		GetContent()

	// GetTreeHeader��TreeAttribute����k
	// �����N�ŦXcompo.TemplateList["components/tree-header"](map[string]string)���ȥ[�Jtext(string)�A���۱N�ѼƤΥ\��K�[���s���ҪO�øѪR���ҪO���D��
	// �N�Ѽ�compo�g�Jbuffer(bytes.Buffer)���̫��XHTML
	// header��/admin/menu�����𪬹ϤW�����|�ӫ��s�e�ݻy�k
	// header���M��class="btn-group"
	header := aTree().GetTreeHeader()

	// aBox�bplugins\admin\controller\common.go��
	// aBox�]�mBoxAttribute(�Ostruct�]�Ointerface)
	// SetHeader�BSetBody�BGetContent����BoxAttribute����k
	// ���O�N�Ѽƭȳ]�m��BoxAttribute(struct)
	// GetContent���̧P�_����]�mBoxAttribute.Style
	// �N�ŦXBoxAttribute.TemplateList["box"](map[string]string)���ȥ[�Jtext(string)�A���۱N�ѼƤΥ\��K�[���s���ҪO�øѪR���ҪO���D��
	// �N�Ѽ�compo�g�Jbuffer(bytes.Buffer)���̫��XHTML
	// box���M��class="box box-"�A���Y��header�B���e��tree(�W�b����檺�y�k)
	box := aBox().SetHeader(header).SetBody(tree).GetContent()

	// aCol�bplugins\admin\controller\common.go��
	// aCol�]�mColAttribute(�Ostruct�]�Ointerface)
	// SetSize�BSetContent�BGetContent���OColAttribute����k
	// ���O�N�Ѽƭȳ]�m��ColAttribute(struct)
	// GetContent�N�ŦXColAttribute.TemplateList["col"](map[string]string)���ȥ[�Jtext(string)�A���۱N�ѼƤΥ\��K�[���s���ҪO�øѪR���ҪO���D��
	// �N�Ѽ�compo�g�Jbuffer(bytes.Buffer)���̫��XHTML
	// col1���M��class="col-md-6"�A���e��box(�W�b�����Ϫ�HTML�y�k)
	col1 := aCol().SetSize(types.SizeMD(6)).SetContent(box).GetContent()

	// BaseTable�]�ݩ�Table(interface)
	// table���z�L�Ѽ�"menu"���oTable(interface)�A���ۧP�_�����N[]context.Node�[�J��Handler.operations��^��
	// GetNewForm�bplugins\admin\modules\table\default.go
	// GetNewForm(���o�s���)�P�_����(TabGroups)��A�]�mFormInfo(struct)��æ^��
	formInfo := h.table("menu", ctx).GetNewForm()

	// aForm�bplugins\admin\controller\common.go��
	// aForm�]�mFormAttribute(�Ostruct�]�Ointerface)
	// �N�Ѽƭȳ]�m��FormFields(struct)
	// �P�_�����A�NFormFields�K�[��FormAttribute.ContentList([]FormFields)
    // ���۱N�ŦXFormAttribute.TemplateList["components/�h�ӰѼ�"](map[string]string)���ȥ[�Jtext(string)�A���۱N�ѼƤΥ\��K�[���s���ҪO�øѪR���ҪO���D��
	// �N�Ѽ�compo�g�Jbuffer(bytes.Buffer)���̫��XHTML�^��
	// menuFormContent(����椺�e)�N�ŦXBoxAttribute.TemplateList["box"](map[string]string)���ȥ[�Jtext(string)�A���۱N�ѼƤΥ\��K�[���s���ҪO�øѪR���ҪO���D��
	// �N�Ѽ�compo�g�Jbuffer(bytes.Buffer)���̫��XHTML
	// newForm���U�b���s�ت�檺HTML�y�k
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

	// aCol�bplugins\admin\controller\common.go��
	// aCol�]�mColAttribute(�Ostruct�]�Ointerface)
	// SetSize�BSetContent�BGetContent���OColAttribute����k
	// GetContent�N�ŦXColAttribute.TemplateList["col"](map[string]string)���ȥ[�Jtext(string)�A���۱N�ѼƤΥ\��K�[���s���ҪO�øѪR���ҪO���D��
	// �N�Ѽ�compo�g�Jbuffer(bytes.Buffer)���̫��XHTML
	// col2���M��class="col-md-6"�A���e��newForm(�U�b���s�ت�檺HTML�y�k)
	col2 := aCol().SetSize(types.SizeMD(6)).SetContent(newForm).GetContent()

	// aRow�bplugins\admin\controller\common.go��
	// aRow�]�mRowAttribute(�Ostruct�]�Ointerface)
	// �btemplate\components\composer.go
	// �����N�ŦXRowAttribute.TemplateList["components/row"](map[string]string)���ȥ[�Jtext(string)�A���۱N�ѼƤΥ\��K�[���s���ҪO�øѪR���ҪO���D��
	// �N�Ѽ�compo�g�Jbuffer(bytes.Buffer)���̫��XHTML
	// row���M��class="row"�A���e���W�U�b���Ҧ���檺HTML�y�k
	row := aRow().SetContent(col1 + col2).GetContent()

	// ��XHTML
	h.HTML(ctx, user, types.Panel{
		Content:     alert + row,
		Description: "Menus Manage",
		Title:       "Menus Manage",
	})
}
