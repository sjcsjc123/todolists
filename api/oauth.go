package api

import (
	"TodoLists/common/errorCode"
	"TodoLists/common/errorMsg"
	"TodoLists/common/result"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sjcsjc123/go-third-login/wxApplets"
	_ "github.com/sjcsjc123/go-third-login/wxApplets"
	_ "github.com/tidwall/gjson"
	"net/http"
)

//小程序登录传入的数据
type ConfirmLoginDto struct {
	Iv        string `json:"iv"`
	WxCode    string `json:"wx_code"`
	Token     string `json:"token"`     //账号登录认证
	Nickname  string `json:"nickname"`  //微信昵称
	Headimage string `json:"headimage"` //微信头像
	User      string `json:"user"`      //账号
	Passwd    string `json:"passwd"`    //密码
}

//账号信息
type User struct {
	Username string `json:"username"`
	Phone    string `json:"phone"`
}

/*微信小程序登录 返回值*/
type WXLoginResp struct {
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

type WxAccessTokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

/*type RequestBody struct {
	Code string `json:"code"`
}*/

/*type LoginReq struct {
	Code  string `json:"code"`
	Phone string `json:"phone"`
}*/

type RequestBody struct {
	WxLoginCode string `json:"wxLoginCode"`
	PhoneCode   string `json:"phoneCode"`
}

type Result struct {
	Phone  string
	Openid string
}

func WxLoginHandler(c *gin.Context) {
	var requestBody RequestBody
	err := c.ShouldBindWith(&requestBody, binding.JSON)
	if err != nil {
		result.JsonError(c, err)
	}
	phone, openid, err := wxApplets.Login(requestBody.WxLoginCode, requestBody.PhoneCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error(errorCode.GetPhoneError, errorMsg.GetPhoneErrorMsg))
		logger.Error("wx login fail:" + err.Error())
	}
	var resultData Result
	resultData.Phone = phone
	resultData.Openid = openid
	c.JSON(http.StatusOK, result.Success(resultData))
	/*var err error
	var loginReq LoginReq
	err = c.ShouldBindWith(&loginReq, binding.JSON)
	if err != nil {
		result.JsonError(c, err)
	}
	code := loginReq.Code
	phone := loginReq.Phone
	if code == "" {
		c.JSON(http.StatusBadRequest, result.Error(errorCode.WxCodeNull, errorMsg.WxCodeNullMsg))
		logger.Info("wxCode error: code null")
		return
	}
	wxLoginResp, err := wxLogin(code)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error(errorCode.WxLoginError, errorMsg.WxLoginErrorMsg))
		logger.Info("wx login error:" + err.Error())
		return
	}
	var user model.User
	err = database.DB.Model(&model.User{}).Where("phone = ?", phone).First(&user).Error
	if err != nil {
		result.MysqlError(c, err)
	}
	user.UserId = wxLoginResp.OpenId
	err = database.DB.Model(&model.User{}).Where("phone = ?", phone).Save(&user).Error
	if err != nil {
		result.MysqlError(c, err)
	}
	logger.Info("wx login success")
	c.Header("openid", wxLoginResp.OpenId)
	c.JSON(http.StatusOK, result.Success(wxLoginResp.OpenId))*/
}

/*根据wx_code返回token ...*/
/*func wxLogin(code string) (*WXLoginResp, error) {
	url := "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
	// 合成url, 这里的appId和secret是在微信公众平台上获取的
	url = fmt.Sprintf(url, config.Conf.GetString(constant.WechatAppId), config.Conf.GetString(constant.WechatSecret), code)
	// 创建http get请求
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// 解析http请求中body 数据到我们定义的结构体中
	wxResp := WXLoginResp{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&wxResp); err != nil {
		return nil, err
	}
	// 判断微信接口返回的是否是一个异常情况
	if wxResp.ErrCode != 0 {
		return nil, errors.New(fmt.Sprintf("ErrCode:%s  ErrMsg:%s", wxResp.ErrCode, wxResp.ErrMsg))
	}
	return &wxResp, nil
}

func GetPhone(c *gin.Context) {
	var appId = config.Conf.GetString(constant.WechatAppId)
	var secret = config.Conf.GetString(constant.WechatSecret)
	get, err := http.Get("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + appId + "&secret=" + secret)
	if err != nil {
		c.JSON(http.StatusOK, result.Error(errorCode.AccessTokenError, errorMsg.AccessTokenErrorMsg))
		logger.Error("get access token error:" + err.Error())
		return
	}
	decoder := json.NewDecoder(get.Body)
	var accessTokenResp WxAccessTokenResp
	if err := decoder.Decode(&accessTokenResp); err != nil {
		result.JsonError(c, err)
	}
	if accessTokenResp.ErrCode != 0 {
		PhoneError(c, err)
	}
	token := accessTokenResp.AccessToken
	code := c.GetHeader("code")
	values := url.Values{}
	values.Add("code", code)
	requestBody := RequestBody{
		Code: code,
	}
	js, err := json.MarshalIndent(&requestBody, "", "\t")
	if err != nil {
		result.JsonError(c, err)
	}
	request, err := http.NewRequest("post", "https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token="+token+"&", bytes.NewBuffer(js))
	if err != nil {
		PhoneError(c, err)
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		PhoneError(c, err)
	}
	body := response.Body
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(body)
	if err != nil {
		c.JSON(http.StatusOK, result.Error(errorCode.ReadCloserError, errorMsg.ReadCloserErrorMsg))
		logger.Error("read closer error:" + err.Error())
		return
	}
	phoneNumFormJson := gjson.Get(buf.String(), "phone_info.phoneNumber")
	phoneNum := phoneNumFormJson.String()
	errorCodeFormJson := gjson.Get(buf.String(), "errcode")
	errorMsgFormJson := gjson.Get(buf.String(), "errmsg")
	wxErrCode := errorCodeFormJson.Int()
	errMsg := errorMsgFormJson.String()
	if wxErrCode != 0 {
		c.JSON(http.StatusBadRequest, result.Error(int(wxErrCode), errMsg))
		logger.Error("get phone error:" + errMsg)
		return
	}
	var count int64
	err = database.DB.Model(&model.User{}).Where("phone = ?", phoneNum).Count(&count).Error
	if err != nil {
		result.MysqlError(c, err)
	}
	if count >= 1 {
		//已存入数据库
	} else {
		err = database.DB.Model(&model.User{}).Create(&model.User{
			Phone: phoneNum,
		}).Error
		if err != nil {
			result.MysqlError(c, err)
		}
	}
	logger.Info("parse phone success ...")
	c.JSON(http.StatusOK, result.Success(phoneNum))
}*/

func PhoneError(c *gin.Context, err error) {
	c.JSON(http.StatusOK, result.Error(errorCode.GetPhoneError, errorMsg.GetPhoneErrorMsg))
	logger.Error("get phone error:" + err.Error())
	return
}
