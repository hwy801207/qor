package admin

import (
	"github.com/qor/qor"
	"github.com/qor/qor/roles"

	"net/http"
	"path"
	"regexp"
	"strings"
)

func (admin *Admin) generateContext(w http.ResponseWriter, r *http.Request) *qor.Context {
	var currentUser *qor.CurrentUser
	context := qor.Context{Config: admin.Config}
	if admin.auth != nil {
		// FIXME
		currentUser = admin.auth.GetCurrentUser(&context)
	}
	context.Roles = roles.MatchedRoles(r, currentUser)
	return &context
}

func (admin *Admin) AddToMux(prefix string, mux *http.ServeMux) {
	// format "/admin" to "/admin/"
	// the trail "/" will match under domain, refer function pathMatch in net/http/server.go
	prefix = regexp.MustCompile("//(//)*").ReplaceAllString("/"+prefix+"/", "/")
	admin.Prefix = prefix

	mux.HandleFunc(strings.TrimRight(prefix, "/"), func(w http.ResponseWriter, r *http.Request) {
		admin.Dashboard(w, r, admin.generateContext(w, r))
	})

	pathMatch := regexp.MustCompile(path.Join(prefix, `(\w+)(?:/(\w+))?[^/]*/?$`))
	mux.HandleFunc(prefix, func(w http.ResponseWriter, r *http.Request) {
		var isIndexURL, isShowURL bool
		context := admin.generateContext(w, r)

		// 128 MB
		r.ParseMultipartForm(32 << 22)
		if len(r.Form["_method"]) > 0 {
			r.Method = strings.ToUpper(r.Form["_method"][0])
		}

		matches := pathMatch.FindStringSubmatch(r.URL.Path)
		if len(matches) == 0 {
			admin.Dashboard(w, r, admin.generateContext(w, r))
			return
		}

		if _, ok := admin.Resources[matches[1]]; matches[1] != "" && ok {
			isIndexURL = true
			context.ResourceName = matches[1]

			if matches[2] != "" { // "/admin/user/1234"
				context.ResourceID = matches[2]
				isIndexURL = false
				isShowURL = true
			}
		}

		switch {
		case r.Method == "GET" && isIndexURL:
			admin.Index(w, r, context)
		case r.Method == "GET" && isShowURL && context.ResourceID == "new":
			admin.New(w, r, context)
		case r.Method == "GET" && isShowURL:
			admin.Show(w, r, context)
		case r.Method == "POST" && isShowURL:
			admin.Update(w, r, context)
		case r.Method == "POST" && isIndexURL:
			admin.Create(w, r, context)
		case r.Method == "DELETE" && isShowURL:
			admin.Delete(w, r, context)
		default:
			http.NotFound(w, r)
		}
	})
}
