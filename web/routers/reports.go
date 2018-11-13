package routers

import (
	"github.com/neal1991/gshark/models"
	"github.com/neal1991/gshark/vars"
	"gopkg.in/macaron.v1"

	"fmt"
	"github.com/go-macaron/session"
	"github.com/neal1991/gshark/util/common"
	"github.com/neal1991/gshark/util/githubsearch"
	"net/url"
	"strconv"
)

func GetDetailedReportById(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		codeResultDetail, _ := models.GetCodeResultDetailById(int64(Id))
		setUserInfoOfCodeResultDetail(codeResultDetail)
		ctx.Data["detailed_report"] = codeResultDetail
		ctx.HTML(200, "report_detail")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func setUserInfoOfCodeResultDetail(detail *models.CodeResultDetail) {
	gitClient, _, _ := githubsearch.GetGithubClient()
	user, resp, err := gitClient.GetUserInfo(*detail.OwnerName)
	if err == nil && resp.StatusCode == 200 {
		detail.OwnerURl = user.HTMLURL
		detail.Blog = user.Blog
		detail.Company = user.Company
		detail.Email = user.Email
		detail.OwnerCreatedAt = user.CreatedAt.Format("2006-01-02 15:04:05")
		detail.Type = user.Type
	} else {
		fmt.Println(err)
	}
}

func ListGithubSearchResultByStatus(ctx *macaron.Context, sess session.Store) {
	status, _ := strconv.Atoi(ctx.Params(":status"))
	ctx.SetCookie("status", strconv.Itoa(status))
	p := 0

	renderDataForGithubSearchResult(ctx, sess, p, status)
}

func ListGithubSearchResult(ctx *macaron.Context, sess session.Store) {
	page := ctx.Params(":page")
	p, _ := strconv.Atoi(page)
	status := 0

	if ctx.GetCookie("status") != "" {
		status, _ = strconv.Atoi(ctx.GetCookie("status"))
	}
	renderDataForGithubSearchResult(ctx, sess, p, status)
}

func renderDataForGithubSearchResult(ctx *macaron.Context, sess session.Store, p, status int) {

	if sess.Get("admin") != nil {
		p, pre, next := common.GetPreAndNext(p)
		reports, pages, count := models.ListGithubSearchResultPage(p, status)
		pageList := common.GetPageList(p, vars.PageStep, pages)
		lastPage := 0
		if len(pageList) >= 1 {
			lastPage = pageList[len(pageList)-1]
		}

		ctx.Data["reports"] = reports
		ctx.Data["pages"] = pages
		ctx.Data["page"] = p
		ctx.Data["pre"] = pre
		ctx.Data["next"] = next
		ctx.Data["pageList"] = pageList
		ctx.Data["status"] = status
		ctx.Data["count"] = count
		ctx.Data["lastPage"] = lastPage
		ctx.Data["role"] = sess.Get("user").(string)
		ctx.HTML(200, "report_github")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func getRefer(ctx *macaron.Context) string {
	refer := "/admin/reports/github/"
	if _, ok := ctx.Req.Header["Referer"]; len(ctx.Req.Header["Referer"]) > 0 && ok {
		u := ctx.Req.Header["Referer"][0]
		urlParsed, err := url.Parse(u)
		if err == nil {
			refer = urlParsed.RequestURI()
		}
	}
	return refer
}

func ConfirmReportById(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		models.ConfirmResultById(int64(Id))
		// redirect to reports which have been confirmed
		ctx.Redirect("/admin/reports/github/query/1")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func CancelReportById(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		models.CancelReportById(int64(Id))
		ctx.Redirect(getRefer(ctx))
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DisableRepoById(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		omitRepo := true
		has, result, err := models.GetReportById(int64(Id), omitRepo)
		if err == nil && has {
			models.DisableRepoByUrl(result.Repository.GetHTMLURL())
			models.CancelReportsByRepo(int64(Id))
		}
		models.CancelReportById(int64(Id))
		ctx.Redirect(getRefer(ctx))
	} else {
		ctx.Redirect("/admin/login/")
	}
}
