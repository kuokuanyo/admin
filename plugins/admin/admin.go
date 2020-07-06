package admin

import (
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/service"
	"github.com/GoAdminGroup/go-admin/plugins"
	"github.com/GoAdminGroup/go-admin/plugins/admin/controller"
	"github.com/GoAdminGroup/go-admin/plugins/admin/models"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/guard"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	_ "github.com/GoAdminGroup/go-admin/template/types/display"
)

// Admin is a GoAdmin plugin.
type Admin struct {
	*plugins.Base
	// plugins\admin\modules\table\table.go
	// GeneratorList���O��map[string]Generator�AGenerator���O��func(ctx *context.Context) Table
	tableList table.GeneratorList 
	// plugins\admin\modules\guard
	guardian  *guard.Guard
	// plugins\admin\controller
	handler   *controller.Handler
}

// InitPlugin implements Plugin.InitPlugin.
// TODO: find a better way to manage the dependencies
func (admin *Admin) InitPlugin(services service.List) {

	// DO NOT DELETE
	// �N�Ѽ�services(map[string]Service)�]�m��Admin.Base(struct)
	admin.InitBase(services)

	// �z�L�Ѽ�("config")���o�ǰt��Service(interface)
	// GetService�N�Ѽ�services.Get("config")�ഫ��Service(struct)��^��Service.C(Config struct)
	c := config.GetService(services.Get("config"))

	// �N�Ѽ�conn�Bc�]�m��SystemTable(struct)��^��
	st := table.NewSystemTable(admin.Conn, c)

	// GeneratorList���O��map[string]Generator�AGenerator���O��func(ctx *context.Context) Table
	// Combine�z�L�ѼƧP�_GeneratorList�w�g����key�Bvalue�A�p�G���s�b�h�[�J����P��
	admin.tableList.Combine(table.GeneratorList{
		"manager":        st.GetManagerTable,
		"permission":     st.GetPermissionTable,
		"roles":          st.GetRolesTable,
		"op":             st.GetOpTable,
		"menu":           st.GetMenuTable,
		"normal_manager": st.GetNormalManagerTable,
		"site":           st.GetSiteTable,
	})

	// �N�Ѽ�admin.Services, admin.Conn, admin.tableList�]�mAdmin.guardian(struct)��^��
	admin.guardian = guard.New(admin.Services, admin.Conn, admin.tableList)

	// �N�ѼƳ]�m��Config(struct)
	handlerCfg := controller.Config{
		Config:     c,
		Services:   services,
		Generators: admin.tableList,
		Connection: admin.Conn,
	}

	// �N�Ѽ�handlerCfg(struct)�ѼƳ]�m��Admin.handler(struct)
	admin.handler.UpdateCfg(handlerCfg)

	// ��l��router
	admin.initRouter()
	admin.handler.SetRoutes(admin.App.Routers)
	admin.handler.AddNavButton(admin.UI.NavButtons...)

	// init site setting
	site := models.Site().SetConn(admin.Conn)
	site.Init(c.ToMap())
	_ = c.Update(site.AllToMap())

	table.SetServices(services)
}

// NewAdmin return the global Admin plugin.
// �]�mAdmin(STRUCT)?æ^��
func NewAdmin(tableCfg ...table.GeneratorList) *Admin {
	return &Admin{
		tableList: make(table.GeneratorList).CombineAll(tableCfg),
		Base:      &plugins.Base{PlugName: "admin"},
		// �]�mHandler(struct)?æ^��
		handler:   controller.New(),
	}
}

// SetCaptcha set captcha driver.
// �N�Ѽ�captcha(���ҽX)�]�m?�Admin.handler.captchaConfig(struct)
func (admin *Admin) SetCaptcha(captcha map[string]string) *Admin {
	// SetCaptcha?bplugins\admin\controller\common.go
	// �N�Ѽ�captcha�]�m?�Handler.captchaConfig(���ҽX�t�m)
	admin.handler.SetCaptcha(captcha)
	return admin
}

// AddGenerator add table model generator.
// �N�Ѽ�key��gen(function)�K?[��??��GeneratorList(map[string]Generator)
func (admin *Admin) AddGenerator(key string, g table.Generator) *Admin {
	// Add�N�Ѽ�key��g(function)�K?[?�Admin.tableList(map[string]Generator)
	admin.tableList.Add(key, g)
	return admin
}

// AddGenerators add table model generators.
// �z�L�Ѽ�gen(?h??)?P�_GeneratorList?w�g����key?Bvalue?A?p�G��?s?b�h?[?J����P??��Admin.tableList
func (admin *Admin) AddGenerators(gen ...table.GeneratorList) *Admin {
	// �z�L�Ѽ�gen(?h??)?P�_GeneratorList?w�g����key?Bvalue?A?p�G��?s?b�h?[?J����P??
	admin.tableList.CombineAll(gen)
	return admin
}

// AddGlobalDisplayProcessFn call types.AddGlobalDisplayProcessFn
func (admin *Admin) AddGlobalDisplayProcessFn(f types.DisplayProcessFn) *Admin {
	types.AddGlobalDisplayProcessFn(f)
	return admin
}

// AddDisplayFilterLimit call types.AddDisplayFilterLimit
func (admin *Admin) AddDisplayFilterLimit(limit int) *Admin {
	types.AddLimit(limit)
	return admin
}

// AddDisplayFilterTrimSpace call types.AddDisplayFilterTrimSpace
func (admin *Admin) AddDisplayFilterTrimSpace() *Admin {
	types.AddTrimSpace()
	return admin
}

// AddDisplayFilterSubstr call types.AddDisplayFilterSubstr
func (admin *Admin) AddDisplayFilterSubstr(start int, end int) *Admin {
	types.AddSubstr(start, end)
	return admin
}

// AddDisplayFilterToTitle call types.AddDisplayFilterToTitle
func (admin *Admin) AddDisplayFilterToTitle() *Admin {
	types.AddToTitle()
	return admin
}

// AddDisplayFilterToUpper call types.AddDisplayFilterToUpper
func (admin *Admin) AddDisplayFilterToUpper() *Admin {
	types.AddToUpper()
	return admin
}

// AddDisplayFilterToLower call types.AddDisplayFilterToLower
func (admin *Admin) AddDisplayFilterToLower() *Admin {
	types.AddToUpper()
	return admin
}

// AddDisplayFilterXssFilter call types.AddDisplayFilterXssFilter
func (admin *Admin) AddDisplayFilterXssFilter() *Admin {
	types.AddXssFilter()
	return admin
}

// AddDisplayFilterXssJsFilter call types.AddDisplayFilterXssJsFilter
func (admin *Admin) AddDisplayFilterXssJsFilter() *Admin {
	types.AddXssJsFilter()
	return admin
}
