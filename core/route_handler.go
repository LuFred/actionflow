package core

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
)

//HandlerTransport imnplements Transport by forwarding context to a handler
type HandlerTransport struct {
	Ctx context.Context
}

// HandlerOption sets a parameter for the InjectBaseHandler
type HandlerOption func(h *HandlerTransport)

//CallHandler 请求函数注入
func CallHandler(s *Server, dto interface{}, h HandlerInterface, mappingMethod string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lg := s.Logger()
		handleValue := reflect.ValueOf(h)
		if handleValue.Type().Kind() != reflect.Ptr {
			panic("must pass a pointer, not a value, to h destination")
		}

		newHandle := reflect.New(handleValue.Type().Elem())

		controllerCtx := NewContext(r, w, s)
		ctxOauthInfo := r.Context().Value(OauthInfoKey)
		if ctxOauthInfo != nil {
			err := initOauthInfo(controllerCtx, ctxOauthInfo)
			if err != nil {
				lg.Error(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		if val := newHandle.MethodByName("Init"); val.IsValid() {
			val.Call([]reflect.Value{reflect.ValueOf(controllerCtx)})
		}

		if val := newHandle.MethodByName(mappingMethod); val.IsValid() {
			if dto != nil {
				newDto := reflect.New(reflect.TypeOf(dto).Elem())
				err := ParseParam(newDto.Interface(), r)
				if err != nil {
					w.Header().Add("Content-Type", "application/json; charset=utf-8")
					w.WriteHeader(http.StatusBadRequest)
					em := fmt.Sprintf(`{"message":"%s"}`, err.Error())
					w.Write([]byte(em))
					return
				}
				if newDto.IsNil() {
					val.Call(nil)
				} else {
					val.Call([]reflect.Value{newDto})
				}
			} else {
				val.Call(nil)
			}
		} else {
			t := reflect.Indirect(newHandle).Type()
			lg.Error(fmt.Sprintf("%s method doesn't exist in the controller %s", mappingMethod, t.Name()))
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func initOauthInfo(httpContext *Context, ctxOauthInfo interface{}) error {
	ctxValue, ok := ctxOauthInfo.(CtxValues)
	if ok {
		_oauthInfo := ctxValue.Get(OauthInfoKey)
		oi := _oauthInfo.(OauthInfo)
		httpContext.setOauthInfo(&oi)
	} else {
		return fmt.Errorf("don't know the type")
	}

	return nil
}
