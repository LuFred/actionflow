package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/lufred/golog"
	// 引入proto包
)

//RecoverMiddleware
func RecoverMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				s := debug.Stack()
				log.Errorf("[%s]:%v host:%s error:%s  stack :%s", r.Host, r.Method, r.URL.Path, err, string(s[:]))
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)

}
