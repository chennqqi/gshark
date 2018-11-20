package routers

import (
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	"github.com/chennqqi/gshark/models"
	"github.com/chennqqi/gshark/util/common"
	"github.com/chennqqi/gshark/vars"
	"gopkg.in/macaron.v1"
	"strconv"
	"strings"
)

func ListAssets(ctx *macaron.Context, sess session.Store) {
	page := ctx.Params(":page")
	p, _ := strconv.Atoi(page)
	p, pre, next := common.GetPreAndNext(p)

	if sess.Get("admin") != nil {
		assets, pages, _ := models.ListInputInfoPage(p)
		pageList := common.GetPageList(p, vars.PageStep, pages)

		ctx.Data["pages"] = pages
		ctx.Data["page"] = p
		ctx.Data["pre"] = pre
		ctx.Data["next"] = next
		ctx.Data["pageList"] = pageList
		ctx.Data["assets"] = assets
		ctx.HTML(200, "assets")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func NewAssets(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		ctx.HTML(200, "assets_new")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DoNewAssets(ctx *macaron.Context, sess session.Store) {
	ctx.Req.ParseForm()
	if sess.Get("admin") != nil {
		Type := strings.TrimSpace(ctx.Req.Form.Get("type"))
		content := strings.TrimSpace(ctx.Req.Form.Get("content"))
		desc := strings.TrimSpace(ctx.Req.Form.Get("desc"))
		assets := models.NewInputInfo(Type, content, desc)
		assets.Insert()
		ctx.Redirect("/admin/assets/list/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func EditAssets(ctx *macaron.Context, sess session.Store, x csrf.CSRF) {
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		assets, _, _ := models.GetInputInfoById(int64(Id))
		ctx.Data["csrf_token"] = x.GetToken()
		ctx.Data["assets"] = assets
		ctx.Data["user"] = sess.Get("admin")
		ctx.HTML(200, "assets_edit")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DoEditAssets(ctx *macaron.Context, sess session.Store) {
	ctx.Req.ParseForm()
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		Type := strings.TrimSpace(ctx.Req.Form.Get("type"))
		content := strings.TrimSpace(ctx.Req.Form.Get("content"))
		desc := strings.TrimSpace(ctx.Req.Form.Get("desc"))
		models.EditInputInfoById(int64(Id), Type, content, desc)
		ctx.Redirect("/admin/assets/list/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DeleteAssets(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		models.DeleteInputInfoById(int64(Id))
		ctx.Redirect("/admin/assets/list/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DeleteAllAssets(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		models.DeleteAllInputInfo()
		ctx.Redirect("/admin/assets/list/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}
