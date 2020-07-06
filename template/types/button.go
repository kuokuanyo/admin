package types

import (
	"github.com/GoAdminGroup/go-admin/modules/utils"
	"github.com/GoAdminGroup/go-admin/plugins/admin/models"
	"html/template"
	"net/url"
)

type Button interface {
	Content() (template.HTML, template.JS)
	GetAction() Action
	URL() string
	METHOD() string
	ID() string
}

// BaseButton(struct)也是Button(interface)
type BaseButton struct {
	Id, Url, Method string
	Title           template.HTML
	Action          Action
}

// 回傳空字串
func (b *BaseButton) Content() (template.HTML, template.JS) {
	return "", ""
}

// 回傳BaseButton.Action(interface)
func (b *BaseButton) GetAction() Action {
	return b.Action
}

// 回傳BaseButton.ID
func (b *BaseButton) ID() string {
	return b.Id
}

// 回傳BaseButton.Url
func (b *BaseButton) URL() string {
	return b.Url
}

// 回傳BaseButton.Method
func (b *BaseButton) METHOD() string {
	return b.Method
}

// 基本的BaseButton(struct)再加上顏色、圖標...等資訊
type DefaultButton struct {
	*BaseButton
	Color     template.HTML
	TextColor template.HTML
	Icon      string
	Direction template.HTML
}

func GetDefaultButton(title template.HTML, icon string, action Action, colors ...template.HTML) *DefaultButton {
	return defaultButton(title, "right", icon, action, colors...)
}

func defaultButton(title, direction template.HTML, icon string, action Action, colors ...template.HTML) *DefaultButton {
	id := btnUUID()
	action.SetBtnId(id)

	var color, textColor template.HTML
	if len(colors) > 0 {
		color = colors[0]
	}
	if len(colors) > 1 {
		textColor = colors[1]
	}
	node := action.GetCallbacks()
	return &DefaultButton{
		BaseButton: &BaseButton{
			Id:     id,
			Title:  title,
			Action: action,
			Url:    node.Path,
			Method: node.Method,
		},
		Color:     color,
		TextColor: textColor,
		Icon:      icon,
		Direction: direction,
	}
}

func GetColumnButton(title template.HTML, icon string, action Action, colors ...template.HTML) *DefaultButton {
	return defaultButton(title, "", icon, action, colors...)
}

func (b *DefaultButton) Content() (template.HTML, template.JS) {

	color := template.HTML("")
	if b.Color != template.HTML("") {
		color = template.HTML(`background-color:`) + b.Color + template.HTML(`;`)
	}
	textColor := template.HTML("")
	if b.TextColor != template.HTML("") {
		textColor = template.HTML(`color:`) + b.TextColor + template.HTML(`;`)
	}

	style := template.HTML("")
	addColor := color + textColor

	if addColor != template.HTML("") {
		style = template.HTML(`style="`) + addColor + template.HTML(`"`)
	}

	h := `<div class="btn-group pull-` + b.Direction + `" style="margin-right: 10px">
                <a ` + style + ` class="` + template.HTML(b.Id) + ` btn btn-sm btn-default ` + b.Action.BtnClass() + `" ` + b.Action.BtnAttribute() + `>
                    <i class="fa ` + template.HTML(b.Icon) + `"></i>&nbsp;&nbsp;` + b.Title + `
                </a>
        </div>` + b.Action.ExtContent()
	return h, b.Action.Js()
}

type ActionButton struct {
	*BaseButton
}

func GetActionButton(title template.HTML, action Action, ids ...string) *ActionButton {

	id := ""
	if len(ids) > 0 {
		id = ids[0]
	} else {
		id = "action-info-btn-" + utils.Uuid(10)
	}

	action.SetBtnId(id)
	node := action.GetCallbacks()

	return &ActionButton{
		BaseButton: &BaseButton{
			Id:     id,
			Title:  title,
			Action: action,
			Url:    node.Path,
			Method: node.Method,
		},
	}
}

func (b *ActionButton) Content() (template.HTML, template.JS) {
	h := template.HTML(`<li style="cursor: pointer;"><a data-id="{{.Id}}" class="`+template.HTML(b.Id)+`" `+b.Action.BtnAttribute()+`>`+b.Title+`</a></li>`) + b.Action.ExtContent()
	return h, b.Action.Js()
}

type Buttons []Button

func (b Buttons) Content() (template.HTML, template.JS) {
	h := template.HTML("")
	j := template.JS("")

	for _, btn := range b {
		hh, jj := btn.Content()
		h += hh
		j += jj
	}
	return h, j
}

func (b Buttons) FooterContent() template.HTML {
	footer := template.HTML("")

	for _, btn := range b {
		footer += btn.GetAction().FooterContent()
	}
	return footer
}

// 檢查權限，回傳Buttons([]Button(interface))
func (b Buttons) CheckPermission(user models.UserModel) Buttons {
	btns := make(Buttons, 0)
	for _, btn := range b {
		// 檢查權限(藉由url、method)
		if user.CheckPermissionByUrlMethod(btn.URL(), btn.METHOD(), url.Values{}) {
			btns = append(btns, btn)
		}
	}
	return btns
}

type NavButton struct {
	*BaseButton
	Icon string
}

// 回傳NavButton(struct)設置資訊
func GetNavButton(title template.HTML, icon string, action Action) *NavButton {

	//btnUUID回傳"info-btn-" + utils.Uuid(10)
	id := btnUUID()
	action.SetBtnId(id)
	node := action.GetCallbacks()

	return &NavButton{
		BaseButton: &BaseButton{
			Id:     id,
			Title:  title,
			Action: action,
			Url:    node.Path,
			Method: node.Method,
		},
		Icon: icon,
	}
}

func (n *NavButton) Content() (template.HTML, template.JS) {

	icon := template.HTML("")
	title := template.HTML("")

	if n.Icon != "" {
		icon = template.HTML(`<i class="fa ` + n.Icon + `"></i>`)
	}

	if n.Title != "" {
		title = `<span>` + n.Title + `</span>`
	}

	h := template.HTML(`<li>
    <a class="`+template.HTML(n.Id)+`" `+n.Action.BtnAttribute()+`>
      `+icon+`
      `+title+`
    </a>
</li>`) + n.Action.ExtContent()
	return h, n.Action.Js()
}
