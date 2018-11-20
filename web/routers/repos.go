package routers

import (
	"github.com/chennqqi/gshark/models"
	"github.com/chennqqi/gshark/util/common"
	"github.com/chennqqi/gshark/vars"

	"gopkg.in/macaron.v1"

	"github.com/go-macaron/session"

	"strconv"
)

func ListRepos(ctx *macaron.Context, sess session.Store) {
	page := ctx.Params(":page")
	p, _ := strconv.Atoi(page)
	p, pre, next := common.GetPreAndNext(p)

	if sess.Get("admin") != nil {
		repos, pages, _ := models.ListReposPage(p)
		pageList := common.GetPageList(p, vars.PageStep, pages)

		ctx.Data["pages"] = pages
		ctx.Data["page"] = p
		ctx.Data["pre"] = pre
		ctx.Data["next"] = next
		ctx.Data["pageList"] = pageList
		ctx.Data["repos"] = repos
		ctx.HTML(200, "repos")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func EnableRepo(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		models.EnableRepoById(int64(Id))
		ctx.Redirect("/admin/repos/list/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DisableRepo(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		models.DisableRepoById(int64(Id))
		ctx.Redirect("/admin/repos/list/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DeleteAllRepo(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		models.DeleteAllRepos()
		ctx.Redirect("/admin/repos/list/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}
