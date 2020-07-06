package admin

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/response"
	"github.com/GoAdminGroup/go-admin/template"
)

// initRouter initialize the router and return the context.
// ��l��router
func (admin *Admin) initRouter() *Admin {
	// �^�Ƿs��App(struct)�A�Ū�
	app := context.NewApp()

	// �N�Ѽ�config.Prefix()�Badmin.globalErrorHandler(���~�B�z�{��)�s�W��RouterGroup(struct)
	// Prefix�^��globalCfg(Config struct).prefix
	// globalErrorHandler�P�_�O�_�N���I���������j��Context.handlers[ctx.index](ctx)
	// �̫�L�X�X�ݰT���b�׺ݾ��W�ðO���Ҧ��ާ@�欰�ܸ�ƪ�(goadmin_operation_log)��
	route := app.Group(config.Prefix(), admin.globalErrorHandler)

	// auth
	// GetLoginUrl�^��globalCfg.LoginUrl
	// ShowLogin�bplugins\admin\controller\auth.go
	// ShowLogin�P�_map[string]Component(interface)�O�_���Ѽ�login(key)���ȡA���۰���template�Ndata�g�Jbuf�ÿ�XHTML
	route.GET(config.GetLoginUrl(), admin.handler.ShowLogin)

	// ���J��username�Bpassword�������ҫ���ouser��role�Bpermission�Υi��menu�A�̫��s��ƪ�(goadmin_users)���K�X��(�[�K)
	route.POST("/signin", admin.handler.Auth)

	// auto install
	// plugins\admin\controller\install.go
	// �إ�buffer(bytes.Buffer)�ÿ�XHTML
	route.GET("/install", admin.handler.ShowInstall)
	// �ˬd��Ʈw�s�u�ѼƬO�_���T
	route.POST("/install/database/check", admin.handler.CheckDatabase)

	// �B�z�e�ݪ��ɮ�
	// Get�P�_templateMap(map[string]Template)��key��O�_�Ѽ�theme�A���h�^��Template(interface)
	// GetTheme�^��globalCfg.Theme
	// GetAssetList��Template(interface)����k
	for _, path := range template.Get(config.GetTheme()).GetAssetList() {
		route.GET("/assets"+path, admin.handler.Assets)
	}
	// GetComponentAsset�ˬdcompMap(map[string]Component)������@�@�[�J�}�C([]string)��
	for _, path := range template.GetComponentAsset() {
		route.GET("/assets"+path, admin.handler.Assets)
	}

	// �N�Ѽ�"/"�Bauth.middleware(admin.Conn)�s�W��RouterGroup(struct)
	// Middleware�إ�Invoker(Struct)�óz�L�Ѽ�ctx���oUserModel�A�åB���o��user��role�B�v���P�i��menu�A�̫��ˬd�Τ��v��
	// authRoute�ݭn����user��role�B�v���P�i��menu�A�̫��ˬd�Τ��v��
	authRoute := route.Group("/", auth.Middleware(admin.Conn))

	// auth
	// �n�X�òM��cookie��^��n�J����
	authRoute.GET("/logout", admin.handler.Logout)


	// menus
	// �ݭn���Ѽ�id = ?
	// MenuDelete�d��url���Ѽ�id���ȫ�Nid�]�m��MenuDeleteParam(struct)�A���۱N�ȳ]�m��Context.UserValue[delete_menu_param]���A�̫����j��Context.handlers[ctx.index](ctx)
	// DeleteMenu�R������MenuModel.id����ơA���F�R��goadmin_menu���~�٭n�R��goadmin_role_menu���
	// �p�GMenuModel.id�O��L��檺���šA�]�����R��
	authRoute.POST("/menu/delete", admin.guardian.MenuDelete, admin.handler.DeleteMenu).Name("menu_delete")

	// MenuNew�bplugins\admin\modules\guard\menu_new.go
	// MenuNew�ǥѰѼƨ��omultipart/form-data���]�m���ȡA��������token�ñNmultipart/form-data��key�Bvalue�ȳ]�m��Context.UserValue[new_menu_param]�A�̫����j��Context.handlers[ctx.index](ctx)
    // NewMenu�NContext.UserValue(map[string]interface{})[new_menu_param]�����ഫ��MenuNewParam(struct)���O�A���۱NMenuNewParam(struct)�ȷs�W�ܸ�ƪ�(MenuModel.Base.TableName(goadmin_menu))��
    // �̫�p�Gmultipart/form-data���]�wroles[]�ȡA�ˬd�����N�Ѽ�roleId(role_id)�PMenuModel.Id(menu_id)�[�Jgoadmin_role_menu��ƪ�
	authRoute.POST("/menu/new", admin.guardian.MenuNew, admin.handler.NewMenu).Name("menu_new")

	// MenuEdit�bplugins\admin\modules\guard\menu_edit.go��
	// MenuEdit�ǥѰѼƨ��omultipart/form-data���]�m���ȡA��������token�ñNmultipart/form-data��key�Bvalue�ȳ]�m��Context.UserValue[edit_menu_param]�A�̫����j��Context.handlers[ctx.index](ctx)
	// EditMenu�NContext.UserValue(map[string]interface{})[edit_menu_param]�����ഫ��MenuEditParam(struct)���O
	// ���Ngoadmin_role_menu��ƪ�menu_id = MenuModel.Id����ƧR���A���ۦp�G���bmultipart/form-data���]�wroles[]�ȡA�ˬd�����N�Ѽ�roleId(role_id)�PMenuModel.Id(menu_id)�[�Jgoadmin_role_menu��ƪ�
	// ���۱Ngoadmin_menu��ƪ����id = MenuModel.Id����Ƴz�L�Ѽ�(��multipart/form-data�]�m)��s
	authRoute.POST("/menu/edit", admin.guardian.MenuEdit, admin.handler.EditMenu).Name("menu_edit")

	// ���omultipart/form-data����_order�Ѽƫ���menu����
	// -----------�٤����Dorder�ѼƦp��]�m-----------
	authRoute.POST("/menu/order", admin.handler.MenuOrder).Name("menu_order")


	authRoute.GET("/menu", admin.handler.ShowMenu).Name("menu")
	authRoute.GET("/menu/edit/show", admin.handler.ShowEditMenu).Name("menu_edit_show")
	authRoute.GET("/menu/new", admin.handler.ShowNewMenu).Name("menu_new_show")

	// Group�N�Ѽ�"/"�Bauth.middleware(admin.Conn)�Badmin.guardian.CheckPrefix�s�W��RouterGroup(struct)
	// CheckPrefix�bplugins\admin\modules\guard\guard.go
	// CheckPrefix�d��url�̪��Ѽ�(__prefix)�A�p�GGuard.tableList�s�b��prefix(key)�h����j��
	// authPrefixRoute�ݭn����user��role�B�v���P�i��menu�A�̫��ˬd�Τ��v���A�H�άd��url�̪��Ѽ�(__prefix)
	authPrefixRoute := route.Group("/", auth.Middleware(admin.Conn), admin.guardian.CheckPrefix)

	// add delete modify query
	authPrefixRoute.GET("/info/:__prefix/detail", admin.handler.ShowDetail).Name("detail")
	authPrefixRoute.GET("/info/:__prefix/edit", admin.guardian.ShowForm, admin.handler.ShowForm).Name("show_edit")
	authPrefixRoute.GET("/info/:__prefix/new", admin.guardian.ShowNewForm, admin.handler.ShowNewForm).Name("show_new")
	authPrefixRoute.POST("/edit/:__prefix", admin.guardian.EditForm, admin.handler.EditForm).Name("edit")
	authPrefixRoute.POST("/new/:__prefix", admin.guardian.NewForm, admin.handler.NewForm).Name("new")
	authPrefixRoute.POST("/delete/:__prefix", admin.guardian.Delete, admin.handler.Delete).Name("delete")
	authPrefixRoute.POST("/export/:__prefix", admin.guardian.Export, admin.handler.Export).Name("export")
	authPrefixRoute.GET("/info/:__prefix", admin.handler.ShowInfo).Name("info")

	authPrefixRoute.POST("/update/:__prefix", admin.guardian.Update, admin.handler.Update).Name("update")

	authRoute.GET("/application/info", admin.handler.SystemInfo)

	route.ANY("/operation/:__goadmin_op_id", auth.Middleware(admin.Conn), admin.handler.Operation)

	if config.GetOpenAdminApi() {
		
		// crud json apis
		apiRoute := route.Group("/api", auth.Middleware(admin.Conn), admin.guardian.CheckPrefix)
		apiRoute.GET("/list/:__prefix", admin.handler.ApiList).Name("api_info")
		apiRoute.GET("/detail/:__prefix", admin.handler.ApiDetail).Name("api_detail")
		apiRoute.POST("/delete/:__prefix", admin.guardian.Delete, admin.handler.Delete).Name("api_delete")
		apiRoute.POST("/update/:__prefix", admin.guardian.EditForm, admin.handler.ApiUpdate).Name("api_edit")
		apiRoute.GET("/update/form/:__prefix", admin.guardian.ShowForm, admin.handler.ApiUpdateForm).Name("api_show_edit")
		apiRoute.POST("/create/:__prefix", admin.guardian.NewForm, admin.handler.ApiCreate).Name("api_new")
		apiRoute.GET("/create/form/:__prefix", admin.guardian.ShowNewForm, admin.handler.ApiCreateForm).Name("api_show_new")
		apiRoute.POST("/export/:__prefix", admin.guardian.Export, admin.handler.Export).Name("api_export")

		apiRoute.POST("/update/:__prefix", admin.guardian.Update, admin.handler.Update).Name("api_update")
	}

	admin.App = app
	return admin
}

// globalErrorHandler(���~�B�z�{��)
// �P�_�O�_�N���I���������j��Context.handlers[ctx.index](ctx)
// �̫�L�X�X�ݰT���b�׺ݾ��W�ðO���Ҧ��ާ@�欰�ܸ�ƪ�(goadmin_operation_log)��
func (admin *Admin) globalErrorHandler(ctx *context.Context) {
	// �bplugins\admin\controller\handler.go
	// �L�X�X�ݰT���b�׺ݾ��W�ðO���Ҧ��ާ@�欰�ܸ�ƪ�(goadmin_operation_log)��
	defer admin.handler.GlobalDeferHandler(ctx)

	// OffLineHandler�bplugins\admin\modules\response\response.go
	// OffLineHandler(���u�B�z�{��)�Ofunc(ctx *context.Context)
	// OffLineHandler�P�_���I�O�_�n�����A�p�n�����A�P�_method�O�_��get�H��header�̥]�taccept:html���XHTML
	response.OffLineHandler(ctx)

	// ����j��Context.handlers[ctx.index](ctx)
	ctx.Next()
}
