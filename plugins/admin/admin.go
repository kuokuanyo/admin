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
	// GeneratorList摸map[string]GeneratorGenerator摸func(ctx *context.Context) Table
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
	// 盢把计services(map[string]Service)砞竚Admin.Base(struct)
	admin.InitBase(services)

	// 硓筁把计("config")眔で皌Service(interface)
	// GetService盢把计services.Get("config")锣传ΘService(struct)肚Service.C(Config struct)
	c := config.GetService(services.Get("config"))

	// 盢把计connc砞竚SystemTable(struct)肚
	st := table.NewSystemTable(admin.Conn, c)

	// GeneratorList摸map[string]GeneratorGenerator摸func(ctx *context.Context) Table
	// Combine硓筁把计耞GeneratorList竒Τ赣keyvalue狦ぃ玥赣龄籔
	admin.tableList.Combine(table.GeneratorList{
		"manager":        st.GetManagerTable,
		"permission":     st.GetPermissionTable,
		"roles":          st.GetRolesTable,
		"op":             st.GetOpTable,
		"menu":           st.GetMenuTable,
		"normal_manager": st.GetNormalManagerTable,
		"site":           st.GetSiteTable,
	})

	// 盢把计admin.Services, admin.Conn, admin.tableList砞竚Admin.guardian(struct)肚
	admin.guardian = guard.New(admin.Services, admin.Conn, admin.tableList)

	// 盢把计砞竚Config(struct)
	handlerCfg := controller.Config{
		Config:     c,
		Services:   services,
		Generators: admin.tableList,
		Connection: admin.Conn,
	}

	// 盢把计handlerCfg(struct)把计砞竚Admin.handler(struct)
	admin.handler.UpdateCfg(handlerCfg)

	// ﹍てrouter
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
// 砞竚Admin(STRUCT)?忙^肚
func NewAdmin(tableCfg ...table.GeneratorList) *Admin {
	return &Admin{
		tableList: make(table.GeneratorList).CombineAll(tableCfg),
		Base:      &plugins.Base{PlugName: "admin"},
		// 砞竚Handler(struct)?忙^肚
		handler:   controller.New(),
	}
}

// SetCaptcha set captcha driver.
// 盢把计captcha(喷靡絏)砞竚?蹵dmin.handler.captchaConfig(struct)
func (admin *Admin) SetCaptcha(captcha map[string]string) *Admin {
	// SetCaptcha?bplugins\admin\controller\common.go
	// 盢把计captcha砞竚?蹾andler.captchaConfig(喷靡絏皌竚)
	admin.handler.SetCaptcha(captcha)
	return admin
}

// AddGenerator add table model generator.
// 盢把计keyのgen(function)睰?[计??GeneratorList(map[string]Generator)
func (admin *Admin) AddGenerator(key string, g table.Generator) *Admin {
	// Add盢把计keyのg(function)睰?[?蹵dmin.tableList(map[string]Generator)
	admin.tableList.Add(key, g)
	return admin
}

// AddGenerators add table model generators.
// 硓筁把计gen(?h??)?P耞GeneratorList?w竒Τ赣key?Bvalue?A?p狦ぃ?s?b玥?[?J赣龄籔??Admin.tableList
func (admin *Admin) AddGenerators(gen ...table.GeneratorList) *Admin {
	// 硓筁把计gen(?h??)?P耞GeneratorList?w竒Τ赣key?Bvalue?A?p狦ぃ?s?b玥?[?J赣龄籔??
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
