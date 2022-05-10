package middleware

import (
	"net/http"
	"strings"

	"github.com/lufred/golog"
	// 引入proto包
	"fmt"
)

//UrlProcessMiddleware
func UrlProcessMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Debugf("[%s]:%v host:%s remote addr:%s", r.Host, r.Method, r.URL.Path, r.RemoteAddr)
		r.URL.Path = strings.ToLower(r.URL.Path)
		//callback接口可以不传token
		if strings.Contains(r.URL.Path, "callback") {
			fmt.Println("r.URL.Path", r.URL.Path)
			r.Header.Set("Authorization", "oauth f9894bed32164a49a3452a0802c1b11b")
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)

}
