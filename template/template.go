// Copyright 2019 GoAdmin Core Team. All rights reserved.
// Use of this source code is governed by a Apache-2.0 style
// license that can be found in the LICENSE file.

package template

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"path"
	"plugin"
	"strings"
	"sync"

	c "github.com/GoAdminGroup/go-admin/modules/config"
	e "github.com/GoAdminGroup/go-admin/modules/errors"
	"github.com/GoAdminGroup/go-admin/modules/language"
	"github.com/GoAdminGroup/go-admin/modules/logger"
	"github.com/GoAdminGroup/go-admin/modules/menu"
	"github.com/GoAdminGroup/go-admin/modules/system"
	"github.com/GoAdminGroup/go-admin/modules/utils"
	"github.com/GoAdminGroup/go-admin/plugins/admin/models"
	"github.com/GoAdminGroup/go-admin/template/login"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// Template is the interface which contains methods of ui components.
// It will be used in the plugins for custom the ui.
// ui�ե󪺤�k�A�N�bplugins���۩w�qui
type Template interface {
	Name() string

	// Components

	// layout
	Col() types.ColAttribute
	Row() types.RowAttribute

	// form and table
	Form() types.FormAttribute
	Table() types.TableAttribute
	DataTable() types.DataTableAttribute

	Tree() types.TreeAttribute
	Tabs() types.TabsAttribute
	Alert() types.AlertAttribute
	Link() types.LinkAttribute

	Paginator() types.PaginatorAttribute
	Popup() types.PopupAttribute
	Box() types.BoxAttribute

	Label() types.LabelAttribute
	Image() types.ImgAttribute

	Button() types.ButtonAttribute

	// Builder methods
	GetTmplList() map[string]string
	GetAssetList() []string
	GetAssetImportHTML(exceptComponents ...string) template.HTML
	GetAsset(string) ([]byte, error)
	GetTemplate(bool) (*template.Template, string)
	GetVersion() string
	GetRequirements() []string
}

const (
	CompCol       = "col"
	CompRow       = "row"
	CompForm      = "form"
	CompTable     = "table"
	CompDataTable = "datatable"
	CompTree      = "tree"
	CompTabs      = "tabs"
	CompAlert     = "alert"
	CompLink      = "link"
	CompPaginator = "paginator"
	CompPopup     = "popup"
	CompBox       = "box"
	CompLabel     = "label"
	CompImage     = "image"
	CompButton    = "button"
)

func HTML(s string) template.HTML {
	return template.HTML(s)
}

func CSS(s string) template.CSS {
	return template.CSS(s)
}

func JS(s string) template.JS {
	return template.JS(s)
}

// The templateMap contains templates registered.
var templateMap = make(map[string]Template)

// Get the template interface by theme name. If the
// name is not found, it panics.
// �P�_templateMap(map[string]Template)��key��O�_�Ѽ�theme�A���h�^��Template(interface)
func Get(theme string) Template {
	if temp, ok := templateMap[theme]; ok {
		return temp
	}
	panic("wrong theme name")
}

// Get the default template with the theme name set with the global config.
// If the name is not found, it panics.
// �p�G�D�D�W�٤w�g�q�L�����t�m�A���o�w�]��Template(interface)
func Default() Template {
	// �p�D�D�W�٤w�g�q�L�t�m�^��true
	// GetTheme�^��globalCfg.Theme(�bmodules\config\config.go)
	if temp, ok := templateMap[c.GetTheme()]; ok {
		return temp
	}
	panic("wrong theme name")
}

var (
	templateMu sync.Mutex
	compMu     sync.Mutex
)

// Add makes a template available by the provided theme name.
// If Add is called twice with the same name or if template is nil,
// it panics.
func Add(name string, temp Template) {
	templateMu.Lock()
	defer templateMu.Unlock()
	if temp == nil {
		panic("template is nil")
	}
	if _, dup := templateMap[name]; dup {
		panic("add template twice " + name)
	}
	templateMap[name] = temp
}

func CheckRequirements() (bool, bool) {
	if !CheckThemeRequirements() {
		return false, true
	}
	if !utils.InArray(DefaultThemeNames, Default().Name()) {
		return true, true
	}
	return true, VersionCompare(Default().GetVersion(), system.RequireThemeVersion()[Default().Name()])
}

func CheckThemeRequirements() bool {
	return VersionCompare(system.Version(), Default().GetRequirements())
}

func VersionCompare(toCompare string, versions []string) bool {
	for _, v := range versions {
		if v == toCompare || utils.CompareVersion(v, toCompare) {
			return true
		}
	}
	return false
}

var DefaultThemeNames = []string{"adminlte", "sword"}

func Themes() []string {
	names := make([]string, len(templateMap))
	i := 0
	for k := range templateMap {
		names[i] = k
		i++
	}
	return names
}

func AddFromPlugin(name string, mod string) {

	plug, err := plugin.Open(mod)
	if err != nil {
		logger.Error("AddFromPlugin err", err)
		panic(err)
	}

	tempPlugin, err := plug.Lookup(strings.Title(name))
	if err != nil {
		logger.Error("AddFromPlugin err", err)
		panic(err)
	}

	var temp Template
	temp, ok := tempPlugin.(Template)
	if !ok {
		logger.Error("AddFromPlugin err: unexpected type from module symbol")
		panic(errors.New("AddFromPlugin err: unexpected type from module symbol"))
	}

	Add(name, temp)
}

// Component is the interface which stand for a ui component.
type Component interface {
	// GetTemplate return a *template.Template and a given key.
	GetTemplate() (*template.Template, string)

	// GetAssetList return the assets url suffix used in the component.
	// example:
	//
	// {{.UrlPrefix}}/assets/login/css/bootstrap.min.css => login/css/bootstrap.min.css
	//
	// See:
	// https://github.com/GoAdminGroup/go-admin/blob/master/template/login/theme1.tmpl#L32
	// https://github.com/GoAdminGroup/go-admin/blob/master/template/login/list.go
	GetAssetList() []string

	// GetAsset return the asset content according to the corresponding url suffix.
	// Asset content is recommended to use the tool go-bindata to generate.
	//
	// See: http://github.com/jteeuwen/go-bindata
	GetAsset(string) ([]byte, error)

	GetContent() template.HTML

	IsAPage() bool

	GetName() string
}

// GetLoginComponent�]�mLogin(struct)�æ^��
// Login(struct)�]�OComponent(interface)
var compMap = map[string]Component{
	"login": login.GetLoginComponent(),
}

// GetComp gets the component by registered name. If the
// name is not found, it panics.
// �P�_map[string]Component�O�_���Ѽ�name(key)���ȡA�����ܫh�^��Component(interface)
func GetComp(name string) Component {
	// Component(interface)
	if comp, ok := compMap[name]; ok {
		return comp
	}
	panic("wrong component name")
}

// �ˬdcompMap(map[string]Component)������@�@�[�J�}�C([]string)��
func GetComponentAsset() []string {
	assets := make([]string, 0)
	for _, comp := range compMap {
		assets = append(assets, comp.GetAssetList()...)
	}
	return assets
}

// �ˬdcompMap(map[string]Component)������O�_�ŦX����å[�J�}�C([]string)��
func GetComponentAssetWithinPage() []string {
	assets := make([]string, 0)
	for _, comp := range compMap {
		if !comp.IsAPage() {
			assets = append(assets, comp.GetAssetList()...)
		}
	}
	return assets
}

// �B�zasset��æ^��HTML�y�k
func GetComponentAssetImportHTML() (res template.HTML) {
	// Default()���o�w�]��template(�D�D�W�٤w�g�q�L�����t�m)
	// GetExcludeThemeComponents(�bmodules\config\config.go)�A���oglobalCfg.ExcludeThemeComponents([]string)
	// GetAssetImportHTML(Template(interface)����k)
	// res���ҨϥΪ�js�y��
	res = Default().GetAssetImportHTML(c.GetExcludeThemeComponents()...)
	// �b�������������asset
	// �ˬdmap[string]Component����O�_�ŦX����å[�J�}�C([]string)��
	assets := GetComponentAssetWithinPage()
	for i := 0; i < len(assets); i++ {
		// �z�L�Ѽ�assets[i]�P�_css��js�ɮסA���oHTML
		res += getHTMLFromAssetUrl(assets[i])
	}
	return
}

// �z�L�Ѽ�s�P�_css��js�ɮסA���oHTML
func getHTMLFromAssetUrl(s string) template.HTML {
	switch path.Ext(s) {
	case ".css":
		return template.HTML(`<link rel="stylesheet" href="` + c.GetAssetUrl() + c.Url("/assets"+s) + `">`)
	case ".js":
		return template.HTML(`<script src="` + c.GetAssetUrl() + c.Url("/assets"+s) + `"></script>`)
	default:
		return ""
	}
}

// ��map[string]Component�j��A��C�@��Component(interface)����GetAsset��k
func GetAsset(path string) ([]byte, error) {
	for _, comp := range compMap {
		res, err := comp.GetAsset(path)
		if err == nil {
			return res, err
		}
	}
	return nil, errors.New(path + " not found")
}

// AddComp makes a component available by the provided name.
// If Add is called twice with the same name or if component is nil,
// it panics.
func AddComp(comp Component) {
	compMu.Lock()
	defer compMu.Unlock()
	if comp == nil {
		panic("component is nil")
	}
	if _, dup := compMap[comp.GetName()]; dup {
		panic("add component twice " + comp.GetName())
	}
	compMap[comp.GetName()] = comp
}

// AddLoginComp add the specified login component.
func AddLoginComp(comp Component) {
	compMu.Lock()
	defer compMu.Unlock()
	compMap["login"] = comp
}

// SetComp makes a component available by the provided name.
// If the value corresponding to the key is empty or if component is nil,
// it panics.
func SetComp(name string, comp Component) {
	compMu.Lock()
	defer compMu.Unlock()
	if comp == nil {
		panic("component is nil")
	}
	if _, dup := compMap[name]; dup {
		compMap[name] = comp
	}
}

type ExecuteParam struct {
	User       models.UserModel
	Tmpl       *template.Template
	TmplName   string
	Panel      types.Panel
	Config     c.Config
	Menu       *menu.Menu
	Animation  bool
	Buttons    types.Buttons
	NoCompress bool
}

// �N���w���ƾ�(types.Page(struct))�g�Jbuf(struct)�æ^��
func Execute(param ExecuteParam) *bytes.Buffer {

	buf := new(bytes.Buffer)
	// ExecuteTemplate��html/template�M��
	// ExecuteTemplate�N���w���ƾ�(�ĤT�ӰѼ�)�g�J�Ѽ�buf
	// NewPageParam(struct)�btemplate\types\page.go��
	// NewPage�NNewPageParam(struct)���ȳ]�m��Page(struct)�æ^��
	err := param.Tmpl.ExecuteTemplate(buf, param.TmplName,
		types.NewPage(types.NewPageParam{
			User:    param.User,
			Menu:    param.Menu,
			Panel:   param.Panel.GetContent(append([]bool{param.Config.IsProductionEnvironment() && (!param.NoCompress)}, param.Animation)...),
			Assets:  GetComponentAssetImportHTML(),
			Buttons: param.Buttons,
		}))
	if err != nil {
		fmt.Println("Execute err", err)
	}
	return buf
}

// �z�L�Ѽ�msg�]�mPanel(struct)
func WarningPanel(msg string) types.Panel {
	return types.Panel{
		// Default()���o�w�]��template(�D�D�W�٤w�g�q�L�����t�m)
		// Alert��Template(interface)����k
		Content:     Default().Alert().Warning(msg),
		Description: template.HTML(e.Msg),
		Title:       template.HTML(e.Msg),
	}
}

var DefaultFuncMap = template.FuncMap{
	"lang":     language.Get,
	"langHtml": language.GetFromHtml,
	"link": func(cdnUrl, prefixUrl, assetsUrl string) string {
		if cdnUrl == "" {
			return prefixUrl + assetsUrl
		}
		return cdnUrl + assetsUrl
	},
	"isLinkUrl": func(s string) bool {
		return (len(s) > 7 && s[:7] == "http://") || (len(s) > 8 && s[:8] == "https://")
	},
	"render": func(s, old, repl template.HTML) template.HTML {
		return template.HTML(strings.Replace(string(s), string(old), string(repl), -1))
	},
	"renderJS": func(s template.JS, old, repl template.HTML) template.JS {
		return template.JS(strings.Replace(string(s), string(old), string(repl), -1))
	},
	"divide": func(a, b int) int {
		return a / b
	},
	"renderRowDataHTML": func(id, content template.HTML) template.HTML {
		return template.HTML(types.ParseTableDataTmplWithID(id, string(content)))
	},
	"renderRowDataJS": func(id template.HTML, content template.JS) template.JS {
		return template.JS(types.ParseTableDataTmplWithID(id, string(content)))
	},
	"js": func(s interface{}) template.JS {
		if ss, ok := s.(string); ok {
			return template.JS(ss)
		}
		if ss, ok := s.(template.HTML); ok {
			return template.JS(ss)
		}
		return ""
	},
}

type BaseComponent struct{}

func (b BaseComponent) GetAssetList() []string               { return make([]string, 0) }
func (b BaseComponent) GetAsset(name string) ([]byte, error) { return nil, nil }