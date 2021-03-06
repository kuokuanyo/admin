package controller

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/errors"
	"github.com/GoAdminGroup/go-admin/modules/logger"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/constant"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/parameter"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/response"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	template2 "html/template"
	"regexp"
	"runtime/debug"
	"strings"
)

// GlobalDeferHandler is a global error handler of admin plugin.
// 印出訪問訊息在終端機上並記錄所有操作行為至資料表(goadmin_operation_log)中
func (h *Handler) GlobalDeferHandler(ctx *context.Context) {
	// 在modules\logger\logger.go
	// 印出訪問訊息在終端機上
	logger.Access(ctx)

	// 在plugins\admin\controller\operation_log.go
	// 記錄所有操作行為至資料表(goadmin_operation_log)中
	h.RecordOperationLog(ctx)

	if err := recover(); err != nil {
		logger.Error(err)
		logger.Error(string(debug.Stack()[:]))

		var (
			errMsg string
			ok     bool
			e      error
		)

		if errMsg, ok = err.(string); !ok {
			if e, ok = err.(error); ok {
				errMsg = e.Error()
			}
		}

		if errMsg == "" {
			errMsg = "system error"
		}

		// 判斷header裡包含accept:json
		if ctx.WantJSON() {
			response.Error(ctx, errMsg)
			return
		}

		if ok, _ = regexp.MatchString("/edit(.*)", ctx.Path()); ok {
			h.setFormWithReturnErrMessage(ctx, errMsg, "edit")
			return
		}
		if ok, _ = regexp.MatchString("/new(.*)", ctx.Path()); ok {
			h.setFormWithReturnErrMessage(ctx, errMsg, "new")
			return
		}
		// 將參數設置至ExecuteParam(struct)，接著將給定的數據(types.Page(struct))寫入buf(struct)並輸出HTML至Context.response.Body
		h.HTML(ctx, auth.Auth(ctx), types.Panel{
			Content:     aAlert().Warning(errMsg),
			Description: template2.HTML(errors.Msg),
			Title:       template2.HTML(errors.Msg),
		})
	}
}

func (h *Handler) setFormWithReturnErrMessage(ctx *context.Context, errMsg string, kind string) {

	var (
		formInfo table.FormInfo
		prefix   = ctx.Query(constant.PrefixKey)
		panel    = h.table(prefix, ctx)
		f        = panel.GetForm()
	)

	if kind == "edit" {
		id := ctx.Query("id")
		if id == "" {
			id = ctx.Request.MultipartForm.Value[panel.GetPrimaryKey().Name][0]
		}
		formInfo, _ = h.table(prefix, ctx).GetDataWithId(parameter.GetParam(ctx.Request.URL,
			panel.GetInfo().DefaultPageSize,
			panel.GetInfo().SortField,
			panel.GetInfo().GetSort()).WithPKs(id))
	} else {
		formInfo = panel.GetNewForm()
		formInfo.Title = f.Title
		formInfo.Description = f.Description
	}

	queryParam := parameter.GetParam(ctx.Request.URL, panel.GetInfo().DefaultPageSize,
		panel.GetInfo().SortField, panel.GetInfo().GetSort()).GetRouteParamStr()

	h.HTML(ctx, auth.Auth(ctx), types.Panel{
		Content: aAlert().Warning(errMsg) + formContent(aForm().
			SetContent(formInfo.FieldList).
			SetTabContents(formInfo.GroupFieldList).
			SetTabHeaders(formInfo.GroupFieldHeaders).
			SetTitle(template2.HTML(strings.Title(kind))).
			SetPrimaryKey(panel.GetPrimaryKey().Name).
			SetPrefix(h.config.PrefixFixSlash()).
			SetHiddenFields(map[string]string{
				form.TokenKey:    h.authSrv().AddToken(),
				form.PreviousKey: h.config.Url("/info/" + prefix + queryParam),
			}).
			SetUrl(h.config.Url("/"+kind+"/"+prefix)).
			SetOperationFooter(formFooter(kind, f.IsHideContinueEditCheckBox, f.IsHideContinueNewCheckBox,
				f.IsHideResetButton)).
			SetHeader(f.HeaderHtml).
			SetFooter(f.FooterHtml), len(formInfo.GroupFieldHeaders) > 0),
		Description: template2.HTML(formInfo.Description),
		Title:       template2.HTML(formInfo.Title),
	})

	ctx.AddHeader(constant.PjaxUrlHeader, h.config.Url("/info/"+prefix+"/"+kind+queryParam))
}
