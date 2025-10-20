package ND_Dkey

import (
	"net/http/cookiejar"
	"fmt"
	"github.com/go-resty/resty/v2"
	"encoding/json"
	"crypto/tls"
	
)


var client_ptr *resty.Client

type Client struct{}

//同步数据源，建议发送前同步一下
func (c *Client ) DatasourceSync() error{

	data:=struct{
		Data string `json:"data"`
		Message string `json:"message"`
	}{}

	for _,v := range( sync_tenant_config){
		response_ptr,err:=client_ptr.R().SetFormData(map[string]string{"id":v}).Post(syncEndpoint)
		if err!=nil {
			return err
		}

		err=json.Unmarshal(response_ptr.Body(), &data)
		if err!=nil {
			return err
		}
		if data.Data !="同步成功"{
			return fmt.Errorf(data.Data+":"+data.Message)
		}
	}
	
	return nil

}

//发送令牌
func (c *Client) SendMFA(Mail string,expireDays string) error{
	res_ptr,err:=client_ptr.R().SetQueryParams(map[string]string{"field":"email","identityStoreId":sync_tenant_config["dataSource1"],"keyword":Mail,"limit":"15","precise":"false","start":"0","tenantId":tenantId}).Get(searchUser_endpoint)
	if err!=nil {
		return err
	}
	userProfile:=struct{
		Data struct{
			Data []struct{
				User struct{
					Id struct{
						Id string `json:"id"`
					}`json:"id"`
					EmailAddress struct{
						Address string `json:"address"`
					}`json:"emailAddress"`
				}`json:"user"`
			} `json:"data"`
		}`json:"data`
	}{}

	err=json.Unmarshal(res_ptr.Body(), &userProfile)
	if err!=nil {
		return err
	}
	uid:=userProfile.Data.Data[0].User.Id.Id
	SendMFARes_ptr,err:=client_ptr.R().SetFormData(map[string]string{"userId":uid,"expireInDays":"5","tokenExpireInDays":expireDays,"tokenBindingStartDelayHours":"0","toEmailAddress":userProfile.Data.Data[0].User.EmailAddress.Address,"generateNewToken":"false"}).Post(deliver_endpoint)
	if err!=nil {
		return nil
	}
	SendResult:= struct{
		Success bool `json:"success"`
	}{}
	err=json.Unmarshal(SendMFARes_ptr.Body(), &SendResult)
	if err!=nil {
		return err
	}
	if SendResult.Success!=true {


		return fmt.Errorf("发送失败！"+string(SendMFARes_ptr.Body()))
	}
	return nil
}

//初始化
func NDInit() (*Client,error){

	cookiejar,_:=cookiejar.New(nil)
	client_ptr=resty.New().SetDebug(true).SetCookieJar(cookiejar).SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true,})
	res_ptr,err:=client_ptr.R().SetFormData(map[string]string{"loginName":adminLoginName,"password":adminPassword}).Post(login_endpoint)
	if err!=nil {
		return nil,err
	}
	if (*res_ptr).StatusCode()!=200 {
		return nil,fmt.Errorf("登录失败，请检查账号密码")
	}
	_,err=client_ptr.R().SetQueryParams(map[string]string{"tenantId":tenantId}).Get(set_endpoint)
	if err!=nil{
		return nil,err
	}
	return &Client{},nil
}