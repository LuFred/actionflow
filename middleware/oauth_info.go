package middleware

import (
	"actionflow/core"
	"context"
	"net/http"
)

//OauthInfoMiddleware
func OauthInfoMiddleware(srv *core.Server, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		oauthInfo := core.OauthInfo{}
		oauthInfo.CompanyId = r.Header.Get("CompanyId")
		oauthInfo.DepartmentId = r.Header.Get("DepartmentId")
		oauthInfo.TenantId = r.Header.Get("TenantId")
		oauthInfo.UserId = r.Header.Get("UserId")

		v := core.NewCtxValue(map[string]interface{}{
			core.OauthInfoKey: oauthInfo,
		})
		r = r.WithContext(context.WithValue(r.Context(), core.OauthInfoKey, *v))
		next.ServeHTTP(w, r)
	})
}
