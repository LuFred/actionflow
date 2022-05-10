package core

import (
	"actionflow/pkg/ljson"
	"errors"
	"fmt"

	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
)

//ParseParam 参数注入
func ParseParam(dest interface{}, r *http.Request) error {
	value := reflect.ValueOf(dest)
	// json.Unmarshal returns errors for these
	if value.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to StructScan destination")
	}
	if value.IsNil() {
		return errors.New("nil pointer passed to StructScan destination")
	}
	rDest := value.Elem()
	if !rDest.CanSet() {
		return errors.New("can not modify dest element")
	}
	types := rDest.Type()
	var contentType string
	jData := new(ljson.JSONData)
	var err error
	if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		contentType = "application/json"
		if r.ContentLength > 0 {
			jData, err = ljson.NewJSONReader(r.Body)
			if err != nil {
				return fmt.Errorf("Parameter error: %s", err.Error())
			}
		}
	}
	for i := 0; i < types.NumField(); i++ {
		field := types.Field(i)
		if jsonName, ok := field.Tag.Lookup("json"); ok {
			values := rDest.FieldByName(field.Name)
			var transport string
			transport, _ = field.Tag.Lookup("transport")
			var jValue reflect.Value

			switch contentType {
			case "application/json":
				jValue, err = GetReflectValueByJSON(r, rDest.FieldByName(field.Name).Type(), jsonName, transport, jData)
			default:
				jValue, err = GetReflectValueByFORM(r, rDest.FieldByName(field.Name).Type(), jsonName, transport)
			}
			if err != nil {
				if err == ljson.ErrorNotFound {
					if isRequired, ok := field.Tag.Lookup("required"); ok {
						if strings.TrimSpace(strings.ToLower(isRequired)) == "true" {
							return fmt.Errorf("The parameter '%s' can not be empty", jsonName)
						}
					}
				} else {
					return fmt.Errorf("Parameter %s error: %s", jsonName, err.Error())
				}
			}
			values.Set(jValue)
		}
	}

	return nil
}

/*GetReflectValueByJSON get parameter value
参数获取途径主要分3种
path:
	chi库
query：
	url.query
body：
	content-Type:application/json
	自定义json解析
*/
func GetReflectValueByJSON(r *http.Request, destType reflect.Type, jsonName, transport string, jData *ljson.JSONData) (reflect.Value, error) {
	var jValue reflect.Value
	var err error
	switch strings.ToLower(transport) {
	case "path":
		val := getValueByPath(r, jsonName)
		jValue, err = bind(destType, val)
	case "query":
		val := r.URL.Query().Get(jsonName)
		jValue, err = bind(destType, val)
	case "body":
		fallthrough
	default:
		jValue, err = jData.Get(jsonName, destType)
	}
	return jValue, err
}

/*GetReflectValueByFORM get parameter value
参数获取途径主要分3种
path:
	chi库
query：
	url.query
body：
	content-Type:application/x-www-form-urlencoded
	r.Form
*/
func GetReflectValueByFORM(r *http.Request, destType reflect.Type, jsonName, transport string) (reflect.Value, error) {
	var jValue reflect.Value
	var err error
	if r.Form == nil {
		r.ParseForm()
	}
	switch strings.ToLower(transport) {
	case "path":
		val := getValueByPath(r, jsonName)
		jValue, err = bind(destType, val)
	case "query":
		val := r.URL.Query().Get(jsonName)
		jValue, err = bind(destType, val)
	case "body":
		fallthrough
	default:
		sVal := r.Form.Get(jsonName)
		jValue, err = bind(destType, sVal)
	}
	return jValue, err
}

//getValueByPath 获取path中的参数
func getValueByPath(r *http.Request, key string) string {
	val := chi.URLParam(r, key)
	return val
}
func bind(typ reflect.Type, val string) (reflect.Value, error) {
	rv := reflect.Zero(typ)
	var err error
	if len(val) == 0 {
		return rv, ljson.ErrorNotFound
	}
	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		rv, err = bindInt(val, typ)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		rv, err = bindUint(val, typ)
	case reflect.Float32, reflect.Float64:
		rv, err = bindFloat(val, typ)
	case reflect.String:
		rv, err = bindString(val, typ)
	case reflect.Bool:
		rv, err = bindBool(val, typ)
	}
	return rv, err
}
func bindInt(val string, typ reflect.Type) (reflect.Value, error) {
	iv, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return reflect.Zero(typ), err
	}
	pValue := reflect.New(typ)
	pValue.Elem().SetInt(iv)
	return pValue.Elem(), nil
}

func bindUint(val string, typ reflect.Type) (reflect.Value, error) {
	uv, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return reflect.Zero(typ), err
	}
	pValue := reflect.New(typ)
	pValue.Elem().SetUint(uv)
	return pValue.Elem(), nil

}

func bindFloat(val string, typ reflect.Type) (reflect.Value, error) {
	fv, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return reflect.Zero(typ), err
	}
	pValue := reflect.New(typ)
	pValue.Elem().SetFloat(fv)
	return pValue.Elem(), nil
}

func bindString(val string, typ reflect.Type) (reflect.Value, error) {
	return reflect.ValueOf(val), nil
}

func bindBool(val string, typ reflect.Type) (reflect.Value, error) {
	val = strings.TrimSpace(strings.ToLower(val))
	switch val {
	case "true", "on", "1":
		return reflect.ValueOf(true), nil
	}
	return reflect.ValueOf(false), nil
}
