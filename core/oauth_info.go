package core

type CtxValues struct {
	m map[string]interface{}
}

type OauthInfo struct {
	UserId       string `json:"userId"`
	TenantId     string `json:"tenantI"`
	DepartmentId string `json:"departmentId"`
	CompanyId    string `json:"companyId"`
}

const OauthInfoKey = "oauthInfo"

func NewCtxValue(m map[string]interface{}) *CtxValues {
	return &CtxValues{
		m: m,
	}
}

func (v CtxValues) Get(key string) interface{} {
	return v.m[key]
}
