package ND_Dkey

// datasource 需要在宁盾后台，点击“外部用户”，触发的/am/controller/tenant/api/2/identity/externalIdentityStore/list?tenant=<tenant> 接口中查
var sync_tenant_config map[string]string = map[string]string{
	"dataSource1": "<dataSource1>",
	"dataSource2": "<dataSource2>",
}

const (
	//根据实际修改
	NDhost                   = "<http://NDhost:port>"
	adminLoginName    string = "<adminLoginName>"
	adminPassword     string = "<adminPassword>"
	tenantName        string = "<tenantName>"
	identityStoreName string = "<identityStoreName>"
	//tenantId 需要在宁盾后台，点击“外部用户”，触发的/am/controller/tenant/api/2/identity/externalIdentityStore/list?tenant=<tenant> 接口中查
	tenantId string = "<tenantId>"


	//以下无需修改
	set_endpoint        string = NDhost + "/am/controller/tenant/current/set"
	deliver_endpoint    string = NDhost + "/am/controller/tenant/api/2/token/mobile/deliver/singleWithEmailAddress"
	searchUser_endpoint string = NDhost + "/am/controller/tenant/api/2/identity/2/user/search"
	login_endpoint      string = NDhost + "/am/controller/admin/login"
	syncEndpoint        string = NDhost + "/am/controller/tenant/api/2/identity/externalIdentityStore/sync"
)


