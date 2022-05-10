package ljson

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Value struct {
	data   interface{}
	exists bool
}

//JSONData json数据对象
type JSONData struct {
	m    map[string]*Value
	data interface{}
}

var (
	//ErrorNotFound 不存在指定键的值
	ErrorNotFound     = errors.New("The value of the specified key does not exist")
	ErrNotArray       = errors.New("Not an array")
	ErrNotNumber      = errors.New("not a number")
	ErrNotBool        = errors.New("no bool")
	ErrNotStruct      = errors.New("not an object")
	ErrNotObjectArray = errors.New("not an object array")
	ErrNotString      = errors.New("not a string")
	ErrNotTime        = errors.New("not a time")
)

//NewJSONReader 从reader流中获取带数据对象
func NewJSONReader(r io.Reader) (*JSONData, error) {
	d := new(JSONData)
	nd := json.NewDecoder(r)
	nd.UseNumber()
	err := nd.Decode(&d.data)
	if err != nil {
		return nil, err
	}
	err = d.init()
	if err != nil {
		return nil, err
	}
	return d, nil
}

//init data初始化到m集合中
func (j *JSONData) init() error {
	valid := false
	switch j.data.(type) {
	case map[string]interface{}:
		valid = true
	}
	if valid {
		j.m = make(map[string]*Value)
		for k, v := range j.data.(map[string]interface{}) {
			j.m[k] = &Value{v, true}
		}
		return nil
	}
	return fmt.Errorf("data is not a json object")
}

func (j JSONData) Get(key string, typ reflect.Type) (reflect.Value, error) {
	val := j.m[key]
	if val == nil {
		return reflect.Zero(typ), ErrorNotFound
	}
	return Bind(typ, j.m[key])
}
func bindSlice(vals []*Value, typ reflect.Type) (reflect.Value, error) {
	sv := reflect.MakeSlice(typ, 0, len(vals))
	sliceValueType := typ.Elem()
	for _, v := range vals {
		childValue, err := Bind(sliceValueType, v)
		if err != nil {
			return reflect.Zero(typ), err
		}
		sv = reflect.Append(sv, childValue)
	}

	return sv, nil
}

func Bind(typ reflect.Type, val *Value) (reflect.Value, error) {
	rv := reflect.Zero(typ)
	if val == nil {
		return rv, nil
	}

	switch typ.Kind() {
	case reflect.Ptr:
		vs, err := Bind(typ.Elem(), val)
		if err != nil {
			return reflect.Zero(typ), err
		}
		return vs.Addr(), nil
	case reflect.Slice:
		vs, err := val.Slice()
		if err != nil {
			return reflect.Zero(typ), err
		}
		return bindSlice(vs, typ)
	case reflect.Struct:
		//时间类型特殊处理
		var err error
		if reflect.Zero(typ).Type().Name() == "Time" {
			s, err := val.TimeString()
			if err != nil {
				return reflect.Zero(typ), err
			}
			return bindTime(s, typ), err
		}
		s, err := val.Struct()
		if err != nil {
			return reflect.Zero(typ), err
		}
		return bindStruct(s, typ)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := val.Number()
		if err != nil {
			return rv, err
		}
		rv = bindInt(string(n), typ)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := val.Number()
		if err != nil {
			return rv, err
		}
		rv = bindUint(string(n), typ)
	case reflect.Float32, reflect.Float64:
		n, err := val.Number()
		if err != nil {
			return rv, err
		}
		rv = bindFloat(string(n), typ)
	case reflect.String:
		s, err := val.String()
		if err != nil {
			return rv, err
		}
		rv = bindString(s, typ)
	case reflect.Bool:
		b, err := val.Boolean()
		if err != nil {
			return rv, err
		}
		rv = reflect.ValueOf(b)
	}
	return rv, nil
}
func (v *Value) Slice() ([]*Value, error) {
	var valid bool

	switch v.data.(type) {
	case []interface{}:
		valid = true
	}
	var slice []*Value
	if valid {
		for _, item := range v.data.([]interface{}) {
			val := Value{item, true}
			slice = append(slice, &val)
		}
		return slice, nil

	}
	return slice, ErrNotArray
}
func (v *Value) String() (string, error) {
	var valid bool

	// Check the type of this data
	switch v.data.(type) {
	case string:
		valid = true
		break
	}

	if valid {
		return v.data.(string), nil
	}

	return "", ErrNotString
}
func (v *Value) TimeString() (time.Time, error) {
	var valid bool
	// Check the type of this data
	switch v.data.(type) {
	case string:
		valid = true
		break
	default:
		base, _ := time.ParseInLocation("2006-01-02 15:04:05", "2006-01-02 15:04:05", time.Local)
		return base, nil
	}
	if valid {
		strTime := v.data.(string)
		reg := `^(\d{4}-\d{2}-\d{2}[T| ]{0,1}\d{2}:\d{2}:\d{2,3})`
		r, _ := regexp.Compile(reg)
		strTime = r.FindString(strTime)
		if strTime == "" {
			base, _ := time.ParseInLocation("2006-01-02 15:04:05", "2006-01-02 15:04:05", time.Local)
			return base, fmt.Errorf("%s foramt=2006-01-02 15:04:05", ErrNotTime)
		}
		strTime = strings.Replace(strTime, "T", " ", -1)
		t, err := time.ParseInLocation("2006-01-02 15:04:05", strTime, time.Local)
		if err != nil {
			return t, ErrNotTime
		}
		return t, nil
	}
	base, _ := time.ParseInLocation("2006-01-02 15:04:05", "2006-01-02 15:04:05", time.Local)
	return base, ErrNotTime
}
func (v *Value) Struct() (*JSONData, error) {
	var valid bool
	// Check the type of this data
	switch v.data.(type) {
	case map[string]interface{}:
		valid = true
		break

	}
	if valid {
		jData := new(JSONData)
		m := make(map[string]*Value)
		for key, element := range v.data.(map[string]interface{}) {
			m[key] = &Value{element, true}
		}
		jData.data = v.data
		jData.m = m

		return jData, nil
	}

	return nil, ErrNotStruct
}

// Attempts to typecast the current value into a number.
// Returns error if the current value is not a json number.
// Example:
//		ageNumber, err := ageValue.Number()
func (v *Value) Number() (json.Number, error) {
	var valid bool
	// Check the type of this data
	switch v.data.(type) {
	case json.Number:
		valid = true
		break
	}
	if valid {
		return v.data.(json.Number), nil
	}
	return "", ErrNotNumber
}

// Attempts to typecast the current value into a bool.
// Returns error if the current value is not a json boolean.
// Example:
//		ageBool, err := ageValue.Boolean()
func (v *Value) Boolean() (bool, error) {
	var valid bool
	// Check the type of this data
	switch v.data.(type) {
	case bool:
		valid = true
		break
	}
	if valid {
		return v.data.(bool), nil
	}
	return false, ErrNotBool
}

func bindStruct(jData *JSONData, typ reflect.Type) (reflect.Value, error) {
	np := reflect.New(typ)
	nv := np.Elem()
	fCount := nv.NumField()
	sType := nv.Type()
	for i := 0; i < fCount; i++ {
		field := sType.Field(i)
		if jsonName, ok := field.Tag.Lookup("json"); ok {
			val := jData.m[jsonName]
			var v reflect.Value
			var err error
			if val == nil {
				v = reflect.Zero(field.Type)
			} else {
				v, err = Bind(field.Type, val)
			}
			if err != nil {
				return reflect.Zero(typ), err
			}
			nv.FieldByName(field.Name).Set(v)
		}
	}
	return nv, nil
}
func bindTime(val time.Time, typ reflect.Type) reflect.Value {
	pValue := reflect.New(typ)
	pValue.Elem().Set(reflect.ValueOf(val))
	return pValue.Elem()
}
func bindInt(val string, typ reflect.Type) reflect.Value {
	iv, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return reflect.Zero(typ)
	}
	pValue := reflect.New(typ)
	pValue.Elem().SetInt(iv)
	return pValue.Elem()
}

func bindUint(val string, typ reflect.Type) reflect.Value {
	uv, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return reflect.Zero(typ)
	}
	pValue := reflect.New(typ)
	pValue.Elem().SetUint(uv)
	return pValue.Elem()
}

func bindFloat(val string, typ reflect.Type) reflect.Value {
	fv, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return reflect.Zero(typ)
	}
	pValue := reflect.New(typ)
	pValue.Elem().SetFloat(fv)
	return pValue.Elem()
}
func bindBool(val string, typ reflect.Type) reflect.Value {
	val = strings.TrimSpace(strings.ToLower(val))
	switch val {
	case "true", "on", "1":
		return reflect.ValueOf(true)
	}
	return reflect.ValueOf(false)
}

func bindString(val string, typ reflect.Type) reflect.Value {
	return reflect.ValueOf(val)
}
