package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/hunterhug/fafacms/core/config"
	"github.com/hunterhug/fafacms/core/model"
	"github.com/hunterhug/fafacms/core/session"
	"github.com/hunterhug/fafacms/core/util"
	"github.com/hunterhug/fafacms/core/util/mail"
	log "github.com/hunterhug/golog"
	"math"
	"strings"
	"time"
)

type RegisterUserRequest struct {
	Name          string `json:"name" validate:"required,alphanumunicode"`
	NickName      string `json:"nick_name" validate:"required"`
	Email         string `json:"email" validate:"required,email"`
	WeChat        string `json:"wechat" validate:"omitempty,alphanumunicode"`
	WeiBo         string `json:"weibo" validate:"omitempty,url"`
	Github        string `json:"github" validate:"omitempty,url"`
	QQ            string `json:"qq" validate:"omitempty,numeric"`
	Password      string `json:"password" validate:"alphanumunicode"`
	RePassword    string `json:"repassword" validate:"eqfield=Password"`
	Gender        int    `json:"gender" validate:"oneof=0 1 2"`
	ShortDescribe string `json:"short_describe"`
	Describe      string `json:"describe"`
	ImagePath     string `json:"image_path"`
}

// RegisterUser User register, anyone can use email register
func RegisterUser(c *gin.Context) {
	resp := new(Resp)
	req := new(RegisterUserRequest)
	defer func() {
		JSONL(c, 200, req, resp)
	}()

	// if close register direct return
	if config.FaFaConfig.DefaultConfig.CloseRegister {
		resp.Error = Error(CloseRegisterError, "")
		return
	}

	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	var validate = validator.New()
	err := validate.Struct(req)
	if err != nil {
		log.Errorf("RegisterUser err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	// name can not repeat and prefix with @
	u := new(model.User)
	if strings.Contains(req.Name, "@") {
		log.Errorf("RegisterUser err: %s", "@ can not be")
		resp.Error = Error(ParasError, "@ can not be")
		return
	}

	if req.NickName == model.AnonymousUser {
		log.Errorf("RegisterUser err: %s", "can not be anonymous name")
		resp.Error = Error(NickNameAlreadyBeUsed, "")
		return
	}

	u.Name = req.Name
	repeat, err := u.IsNameRepeat()
	if err != nil {
		log.Errorf("RegisterUser err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	if repeat {
		log.Errorf("RegisterUser err: %s", "name already use by other")
		resp.Error = Error(UserNameAlreadyBeUsed, "")
		return
	}

	// nickname also must unique
	u.NickName = req.NickName
	repeat, err = u.IsNickNameRepeat()
	if err != nil {
		log.Errorf("RegisterUser err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	if repeat {
		log.Errorf("RegisterUser err: %s", "nickname already use by other")
		resp.Error = Error(NickNameAlreadyBeUsed, "")
		return
	}

	// email also
	u.Email = req.Email
	repeat, err = u.IsEmailRepeat()
	if err != nil {
		log.Errorf("RegisterUser err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	if repeat {
		log.Errorf("RegisterUser err: %s", "email already use by other")
		resp.Error = Error(EmailAlreadyBeUsed, "")
		return
	}

	// activate code gen
	u.ActivateCode = util.GetGUID()
	u.ActivateCodeExpired = time.Now().Add(5 * time.Minute).Unix()
	u.ShortDescribe = req.ShortDescribe
	u.Describe = req.Describe
	u.Password = req.Password
	u.Gender = req.Gender
	u.WeChat = req.WeChat
	u.QQ = req.QQ
	u.Github = req.Github
	u.WeiBo = req.WeiBo

	// send email
	mm := new(mail.Message)
	mm.Sender = config.FaFaConfig.MailConfig
	mm.To = u.Email
	mm.ToName = u.NickName
	mm.Body = fmt.Sprintf(mm.Body, "Register", u.ActivateCode)
	err = mm.Sent()
	if err != nil {
		log.Errorf("RegisterUser err:%s", err.Error())
		resp.Error = Error(EmailSendError, err.Error())
		return
	}

	err = u.InsertOne()
	if err != nil {
		log.Errorf("RegisterUser err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	// if debug will return some info
	if AuthDebug {
		resp.Data = u
	}

	resp.Flag = true
}

// CreateUser Create user, admin url
func CreateUser(c *gin.Context) {
	resp := new(Resp)
	req := new(RegisterUserRequest)
	defer func() {
		JSONL(c, 200, req, resp)
	}()

	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	var validate = validator.New()
	err := validate.Struct(req)
	if err != nil {
		log.Errorf("CreateUser err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	u := new(model.User)
	if strings.Contains(req.Name, "@") {
		log.Errorf("CreateUser err: %s", "@ can not be")
		resp.Error = Error(ParasError, "@ can not be")
		return
	}

	if req.NickName == model.AnonymousUser {
		log.Errorf("CreateUser err: %s", "can not be anonymous name")
		resp.Error = Error(NickNameAlreadyBeUsed, "")
		return
	}

	u.Name = req.Name
	repeat, err := u.IsNameRepeat()
	if err != nil {
		log.Errorf("CreateUser err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	if repeat {
		log.Errorf("CreateUser err: %s", "name already use by other")
		resp.Error = Error(UserNameAlreadyBeUsed, "")
		return
	}

	u.NickName = req.NickName
	repeat, err = u.IsNickNameRepeat()
	if err != nil {
		log.Errorf("CreateUser err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	if repeat {
		log.Errorf("CreateUser err: %s", "nickname already use by other")
		resp.Error = Error(NickNameAlreadyBeUsed, "")
		return
	}

	// email check
	u.Email = req.Email
	repeat, err = u.IsEmailRepeat()
	if err != nil {
		log.Errorf("CreateUser err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	if repeat {
		log.Errorf("CreateUser err: %s", "email already use by other")
		resp.Error = Error(EmailAlreadyBeUsed, "")
		return
	}

	// if image not empty
	if req.ImagePath != "" {
		u.HeadPhoto = req.ImagePath
		p := new(model.File)
		p.Url = req.ImagePath
		ok, err := p.Exist()
		if err != nil {
			log.Errorf("CreateUser err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}

		if !ok {
			log.Errorf("CreateUser err: image not exist")
			resp.Error = Error(FileCanNotBeFound, "")
			return
		}
	}

	u.Describe = req.Describe
	u.ShortDescribe = req.ShortDescribe
	u.NickName = req.NickName
	u.Password = req.Password
	u.Gender = req.Gender
	u.WeChat = req.WeChat
	u.QQ = req.QQ
	u.Github = req.Github
	u.WeiBo = req.WeiBo

	// default is activate
	u.Status = 1
	u.ActivateTime = time.Now().Unix()
	err = u.InsertOne()
	if err != nil {
		log.Errorf("CreateUser err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	resp.Flag = true
	resp.Data = u
}

type ActivateUserRequest struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required"`
}

// ActivateUser Activate by oneself
func ActivateUser(c *gin.Context) {
	resp := new(Resp)
	req := new(ActivateUserRequest)
	defer func() {
		JSONL(c, 200, req, resp)
	}()

	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	var validate = validator.New()
	err := validate.Struct(req)
	if err != nil {
		log.Errorf("ActivateUser err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	// email and activate code must together
	u := new(model.User)
	u.ActivateCode = req.Code
	u.Email = req.Email

	// whether exist
	exist, err := u.IsActivateCodeExist()
	if err != nil {
		log.Errorf("ActivateUser err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	if !exist {
		log.Errorf("ActivateUser err:%s", "not exist code")
		resp.Error = Error(ActivateCodeWrong, "not exist code")
		return
	}

	// has been activated direct return
	if u.Status != 0 {
		resp.Flag = true
		return
	}

	// activate code expired, must gen again
	if u.ActivateCodeExpired < time.Now().Unix() {
		log.Errorf("ActivateUser err:%s", "code expired")
		resp.Error = Error(ActivateCodeExpired, "")
		return
	} else {
		u.Status = 1
		err = u.UpdateActivateStatus()
		if err != nil {
			log.Errorf("ActivateUser err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}

		// activate success will soon set session
		token, err := SetUserSession(u.Id)
		if err != nil {
			log.Errorf("ActivateUser err:%s", err.Error())
			resp.Error = Error(SetUserSessionError, err.Error())
			return
		}

		// return token
		resp.Data = token
	}

	resp.Flag = true
}

type ResendActivateCodeToUserRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ResendActivateCodeToUser Activate code expire must resent email ang get new one
func ResendActivateCodeToUser(c *gin.Context) {
	resp := new(Resp)
	req := new(ResendActivateCodeToUserRequest)
	defer func() {
		JSONL(c, 200, req, resp)
	}()

	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	var validate = validator.New()
	err := validate.Struct(req)
	if err != nil {
		log.Errorf("ResendActivateCodeToUser err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	// get user info by email
	u := new(model.User)
	u.Email = req.Email
	ok, err := u.GetUserByEmail()
	if err != nil {
		log.Errorf("ResendActivateCodeToUser err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	if !ok {
		log.Errorf("ResendActivateCodeToUser err:%s", "email not found")
		resp.Error = Error(EmailNotFound, "")
		return
	}

	if u.Status != 0 {
		resp.Flag = true
		return
	} else if u.ActivateCodeExpired > time.Now().Unix() {
		// can not gen a new code because expire time not reach
		log.Errorf("ResendUser err:%s", "code not expired")
		resp.Error = Error(ActivateCodeNotExpired, "")
		return
	}

	// update activate code, expire after 5 min
	err = u.UpdateActivateCode()
	if err != nil {
		log.Errorf("ResendUser err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	// send email
	mm := new(mail.Message)
	mm.Sender = config.FaFaConfig.MailConfig
	mm.To = u.Email
	mm.ToName = u.NickName
	mm.Body = fmt.Sprintf(mm.Body, "Register", u.ActivateCode)
	err = mm.Sent()
	if err != nil {
		log.Errorf("ResendUser err:%s", err.Error())
		resp.Error = Error(EmailSendError, err.Error())
		return
	}
	resp.Flag = true
}

type ForgetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ForgetPasswordOfUser When user want to modify password or forget password, can gen a code to change password
func ForgetPasswordOfUser(c *gin.Context) {
	resp := new(Resp)
	req := new(ForgetPasswordRequest)
	defer func() {
		JSONL(c, 200, req, resp)
	}()

	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	var validate = validator.New()
	err := validate.Struct(req)
	if err != nil {
		log.Errorf("RegisterUser err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	u := new(model.User)
	u.Email = req.Email
	ok, err := u.GetUserByEmail()
	if err != nil {
		log.Errorf("ForgetPassword err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	if !ok {
		log.Errorf("ForgetPassword err:%s", "email not found")
		resp.Error = Error(EmailNotFound, "")
		return
	}

	// only code expired can gen a new one again
	if u.ResetCodeExpired < time.Now().Unix() {
		// code is valid in 5 min
		err = u.UpdateCode()
		if err != nil {
			log.Errorf("ForgetPassword comerr:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}

		// send email
		mm := new(mail.Message)
		mm.Sender = config.FaFaConfig.MailConfig
		mm.To = u.Email
		mm.ToName = u.NickName
		mm.Body = fmt.Sprintf(mm.Body, "Forget Password", u.ResetCode)
		err = mm.Sent()
		if err != nil {
			log.Errorf("ForgetPassword err:%s", err.Error())
			resp.Error = Error(EmailSendError, err.Error())
			return
		}

	} else {
		log.Errorf("ForgetPassword err:%s", "reset code expired time not reach")
		resp.Error = Error(ResetCodeExpiredTimeNotReach, "")
		return
	}

	resp.Flag = true
}

type ChangePasswordRequest struct {
	Email      string `json:"email" validate:"required,email"`
	Code       string `json:"code" validate:"required"`
	Password   string `json:"password" validate:"alphanumunicode"`
	RePassword string `json:"repassword" validate:"eqfield=Password"`
}

// ChangePasswordOfUser Change password by a forged password email code
func ChangePasswordOfUser(c *gin.Context) {
	resp := new(Resp)
	req := new(ChangePasswordRequest)
	defer func() {
		JSONL(c, 200, req, resp)
	}()

	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	var validate = validator.New()
	err := validate.Struct(req)
	if err != nil {
		log.Errorf("ChangePassword err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	u := new(model.User)
	u.Email = req.Email
	ok, err := u.GetUserByEmail()
	if err != nil {
		log.Errorf("ChangePassword err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	if !ok {
		log.Errorf("ChangePassword err:%s", "email not found")
		resp.Error = Error(EmailNotFound, "")
		return
	}

	// rest code is the same can change
	if u.ResetCode == req.Code {
		u.Password = req.Password
		err = u.UpdatePassword()
		if err != nil {
			log.Errorf("ChangePassword err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}
	} else {
		log.Errorf("ChangePassword err:%s", "reset code wrong")
		resp.Error = Error(RestCodeWrong, "")
		return
	}

	// after change password, session will delete all
	DeleteUserAllSession(u.Id)
	resp.Flag = true
}

type UpdateUserRequest struct {
	NickName      string `json:"nick_name" validate:"omitempty"`
	WeChat        string `json:"wechat" validate:"omitempty,alphanumunicode"`
	WeiBo         string `json:"weibo" validate:"omitempty,url"`
	Github        string `json:"github" validate:"omitempty,url"`
	QQ            string `json:"qq" validate:"omitempty,numeric"`
	Gender        int    `json:"gender" validate:"oneof=0 1 2"`
	Describe      string `json:"describe"`
	ShortDescribe string `json:"short_describe"`
	ImagePath     string `json:"image_path"`
}

func UpdateUser(c *gin.Context) {
	resp := new(Resp)
	req := new(UpdateUserRequest)
	defer func() {
		JSONL(c, 200, req, resp)
	}()

	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	var validate = validator.New()
	err := validate.Struct(req)
	if err != nil {
		log.Errorf("UpdateUser err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	// get oneself info
	uu, err := GetUserSession(c)
	if err != nil {
		log.Errorf("UpdateUser err: %s", err.Error())
		resp.Error = Error(GetUserSessionError, err.Error())
		return
	}

	uuu := new(model.User)
	uuu.Id = uu.Id
	ok, err := uuu.GetRaw()
	if err != nil {
		log.Errorf("UpdateUser err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	if !ok {
		log.Errorf("UpdateUser err: %s", "user not found")
		resp.Error = Error(UserNotFound, "")
		return
	}

	u := new(model.User)
	u.Id = uu.Id
	if req.ImagePath != "" && req.ImagePath != uuu.HeadPhoto {
		u.HeadPhoto = req.ImagePath
		p := new(model.File)
		p.Url = req.ImagePath
		ok, err := p.Exist()
		if err != nil {
			log.Errorf("UpdateUser err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}

		if !ok {
			log.Errorf("UpdateUser err: image not exist")
			resp.Error = Error(FileCanNotBeFound, "")
			return
		}
	}

	if req.NickName == model.AnonymousUser {
		log.Errorf("UpdateUser err: %s", "can not be anonymous name")
		resp.Error = Error(NickNameAlreadyBeUsed, "")
		return
	}

	// nickname can change 2 times one month
	if req.NickName != "" && req.NickName != uuu.NickName {
		if uuu.NickNameUpdateTime != 0 {
			passTime := time.Now().Unix() - uuu.NickNameUpdateTime
			if passTime < 15*24*3600 {
				log.Errorf("UpdateUser err: %s", "nickname can not change time not reach")
				resp.Error = Error(NickNameCanNotChangeForTimeNotReach, fmt.Sprintf("remain %d days", passTime/(24*3600)))
				return
			}
		}
		u.NickName = req.NickName
		u.NickNameUpdateTime = time.Now().Unix()
		repeat, err := u.IsNickNameRepeat()
		if err != nil {
			log.Errorf("UpdateUser err: %s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}
		if repeat {
			log.Errorf("UpdateUser err: %s", "nickname already use by other")
			resp.Error = Error(NickNameAlreadyBeUsed, "")
			return
		}
	}

	u.Describe = req.Describe
	u.ShortDescribe = req.ShortDescribe
	u.Gender = req.Gender
	u.WeChat = req.WeChat
	u.QQ = req.QQ
	u.Github = req.Github
	u.WeiBo = req.WeiBo
	err = u.UpdateInfo()
	if err != nil {
		log.Errorf("UpdateUser err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	err = session.Mgr.RefreshUser([]string{fmt.Sprintf("%d", u.Id)}, SessionExpireTime)
	if err != nil {
		log.Errorf("UpdateUser err:%s", err.Error())
		resp.Error = Error(RefreshUserCacheError, err.Error())
		return
	}

	resp.Flag = true
	resp.Data = u
}

// TakeUser Take oneself user info
func TakeUser(c *gin.Context) {
	resp := new(Resp)
	defer func() {
		JSONL(c, 200, nil, resp)
	}()

	// just get from session
	u, err := GetUserSession(c)
	if err != nil {
		log.Errorf("TakeUser err:%s", err.Error())
		resp.Error = Error(GetUserSessionError, err.Error())
		return
	}

	user := new(model.User)
	user.Id = u.Id
	exist, err := model.FaFaRdb.Client.Get(user)
	if err != nil {
		log.Errorf("TakeUser err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	if !exist {
		log.Errorf("TakeUser err:%s", "user  not found")
		resp.Error = Error(UserNotFound, "")
		return
	}

	v := user
	p := People{}
	p.Id = v.Id
	p.ShortDescribe = v.ShortDescribe
	p.Describe = v.Describe
	p.CreateTime = GetSecond2DateTimes(v.CreateTime)
	p.CreateTimeInt = v.CreateTime

	if v.Status == 2 {
		p.IsInBlack = true
	}

	p.UpdateTimeInt = v.UpdateTime
	if v.UpdateTime > 0 {
		p.UpdateTime = GetSecond2DateTimes(v.UpdateTime)
	}

	p.LoginTimeInt = v.LoginTime
	if v.LoginTime > 0 {
		p.LoginTime = GetSecond2DateTimes(v.LoginTime)
	}

	p.LoginIp = v.LoginIp

	p.NickNameUpdateTimeInt = v.NickNameUpdateTime
	if v.NickNameUpdateTime > 0 {
		p.NickNameUpdateTime = GetSecond2DateTimes(v.NickNameUpdateTime)
	}

	p.ActivateTimeInt = v.ActivateTime
	if v.ActivateTime > 0 {
		p.ActivateTime = GetSecond2DateTimes(v.ActivateTime)
	}
	p.Email = v.Email
	p.Github = v.Github
	p.Name = v.Name
	p.NickName = v.NickName
	p.HeadPhoto = v.HeadPhoto
	p.QQ = v.QQ
	p.WeChat = v.WeChat
	p.WeiBo = v.WeiBo
	p.Gender = v.Gender
	p.FollowingNum = v.FollowingNum
	p.FollowedNum = v.FollowedNum
	p.ContentNum = v.ContentNum
	p.ContentCoolNum = v.ContentCoolNum
	p.IsVip = v.Vip == 1
	resp.Flag = true
	resp.Data = p
}

type ListUserRequest struct {
	Id              int      `json:"id"`
	Name            string   `json:"name"`
	CreateTimeBegin int64    `json:"create_time_begin"`
	CreateTimeEnd   int64    `json:"create_time_end"`
	UpdateTimeBegin int64    `json:"update_time_begin"`
	UpdateTimeEnd   int64    `json:"update_time_end"`
	Sort            []string `json:"sort"`
	Email           string   `json:"email" validate:"omitempty,email"`
	WeChat          string   `json:"wechat" validate:"omitempty,alphanumunicode"`
	WeiBo           string   `json:"weibo" validate:"omitempty,url"`
	Github          string   `json:"github" validate:"omitempty,url"`
	QQ              string   `json:"qq" validate:"omitempty,numeric"`
	Gender          int      `json:"gender" validate:"oneof=-1 0 1 2"`
	Status          int      `json:"status" validate:"oneof=-1 0 1 2"`
	Vip             int      `json:"vip" validate:"oneof=-1 0 1"`
	PageHelp
}

type ListUserResponse struct {
	Users []model.User `json:"users"`
	PageHelp
}

// ListUser List user, admin url
func ListUser(c *gin.Context) {
	resp := new(Resp)

	respResult := new(ListUserResponse)
	req := new(ListUserRequest)
	defer func() {
		JSONL(c, 200, req, resp)
	}()

	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	var validate = validator.New()
	err := validate.Struct(req)
	if err != nil {
		log.Errorf("ListUser err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	// new query list session
	sess := model.FaFaRdb.Client.NewSession()
	defer sess.Close()

	// group list where prepare
	sess.Table(new(model.User)).Where("1=1")

	// query prepare
	if req.Id != 0 {
		sess.And("id=?", req.Id)
	}
	if req.Name != "" {
		sess.And("name=?", req.Name)
	}

	if req.Status != -1 {
		sess.And("status=?", req.Status)
	}

	if req.Gender != -1 {
		sess.And("gender=?", req.Gender)
	}

	if req.Vip != -1 {
		sess.And("vip=?", req.Vip)
	}

	if req.QQ != "" {
		sess.And("q_q=?", req.QQ)
	}

	if req.Email != "" {
		sess.And("email=?", req.Email)
	}

	if req.Github != "" {
		sess.And("github=?", req.Github)
	}

	if req.WeiBo != "" {
		sess.And("wei_bo=?", req.WeiBo)
	}
	if req.WeChat != "" {
		sess.And("we_chat=?", req.WeChat)
	}

	if req.CreateTimeBegin > 0 {
		sess.And("create_time>=?", req.CreateTimeBegin)
	}

	if req.CreateTimeEnd > 0 {
		sess.And("create_time<?", req.CreateTimeEnd)
	}

	if req.UpdateTimeBegin > 0 {
		sess.And("update_time>=?", req.UpdateTimeBegin)
	}

	if req.UpdateTimeEnd > 0 {
		sess.And("update_time<?", req.UpdateTimeEnd)
	}

	users := make([]model.User, 0)
	p := &req.PageHelp

	// sql build
	p.build(sess, req.Sort, model.UserSortName)

	// do query
	total, err := sess.FindAndCount(&users)
	if err != nil {
		log.Errorf("ListUser err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	// result
	respResult.Users = users
	p.Pages = int(math.Ceil(float64(total) / float64(p.Limit)))
	p.Total = int(total)
	respResult.PageHelp = *p
	resp.Data = respResult
	resp.Flag = true
}

type ListGroupUserRequest struct {
	GroupId int `json:"group_id" validate:"required"`
}

type ListGroupUserResponse struct {
	Users []model.User `json:"users"`
}

// ListGroupUser List the users of group
func ListGroupUser(c *gin.Context) {
	resp := new(Resp)

	respResult := new(ListGroupUserResponse)
	req := new(ListGroupUserRequest)
	defer func() {
		JSONL(c, 200, req, resp)
	}()

	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	var validate = validator.New()
	err := validate.Struct(req)
	if err != nil {
		log.Errorf("ListGroupUser err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	// new query list session
	sess := model.FaFaRdb.Client.NewSession()
	defer sess.Close()

	users := make([]model.User, 0)

	// group list where prepare
	err = sess.Table(new(model.User)).Where("group_id=?", req.GroupId).Find(&users)
	if err != nil {
		log.Errorf("ListUser err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	respResult.Users = users
	resp.Data = respResult
	resp.Flag = true
}

type AssignGroupRequest struct {
	GroupId      int64   `json:"group_id"`
	GroupRelease int     `json:"group_release"`
	Users        []int64 `json:"users"`
}

// AssignGroupToUser Assign user to a group, every user can only have less than one group
func AssignGroupToUser(c *gin.Context) {
	resp := new(Resp)
	req := new(AssignGroupRequest)
	defer func() {
		JSONL(c, 200, req, resp)
	}()

	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	if len(req.Users) == 0 {
		log.Errorf("AssignGroupToUser err:%s", "users empty")
		resp.Error = Error(ParasError, "users empty")
		return
	}

	// release the user of group, user will not belong to any group
	if req.GroupRelease == 1 {
		u := new(model.User)
		num, err := model.FaFaRdb.Client.Cols("group_id").In("id", req.Users).Update(u)
		if err != nil {
			log.Errorf("AssignGroupToUser err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}

		uidList := make([]string, 0)
		for _, v := range req.Users {
			uidList = append(uidList, fmt.Sprintf("%d", v))
		}

		err = session.Mgr.RefreshUser(uidList, SessionExpireTime)
		if err != nil {
			log.Errorf("AssignGroupToUser err:%s", err.Error())
			resp.Error = Error(RefreshUserCacheError, err.Error())
			return
		}
		resp.Data = num
	} else {
		if req.GroupId == 0 {
			log.Errorf("AssignGroupToUser err:%s", "group id empty")
			resp.Error = Error(ParasError, "group_id empty")
			return
		}

		g := new(model.Group)
		g.Id = req.GroupId
		exist, err := g.GetById()
		if err != nil {
			log.Errorf("AssignGroupToUser err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}

		if !exist {
			log.Errorf("AssignGroupToUser err:%s", "group not found")
			resp.Error = Error(GroupNotFound, "")
			return
		}

		u := new(model.User)
		u.GroupId = req.GroupId
		num, err := model.FaFaRdb.Client.Cols("group_id").In("id", req.Users).Update(u)
		if err != nil {
			log.Errorf("AssignGroupToUser err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}

		uidList := make([]string, 0)
		for _, v := range req.Users {
			uidList = append(uidList, fmt.Sprintf("%d", v))
		}

		err = session.Mgr.RefreshUser(uidList, SessionExpireTime)
		if err != nil {
			log.Errorf("AssignGroupToUser err:%s", err.Error())
			resp.Error = Error(RefreshUserCacheError, err.Error())
			return
		}
		resp.Data = num
	}

	resp.Flag = true
}

type UpdateUserAdminRequest struct {
	Id       int64  `json:"id" validate:"required"`
	NickName string `json:"nick_name" validate:"omitempty"`
	Password string `json:"password,omitempty"`
	Status   int    `json:"status" validate:"oneof=0 1 2"` // o nothing 1 activate 2 ban
	Vip      int    `json:"vip" validate:"oneof=0 1 2"`    // 1 become vip, 2 no vip
}

// UpdateUserAdmin Update user info, admin url. Can change user password, black one user, change nickname etc.
func UpdateUserAdmin(c *gin.Context) {
	resp := new(Resp)
	req := new(UpdateUserAdminRequest)
	defer func() {
		JSONL(c, 200, req, resp)
	}()

	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	var validate = validator.New()
	err := validate.Struct(req)
	if err != nil {
		log.Errorf("UpdateUserAdmin err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	uu := new(model.User)
	uu.Id = req.Id
	ok, err := uu.GetRaw()
	if err != nil {
		log.Errorf("UpdateUserAdmin err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	if !ok {
		log.Errorf("UpdateUserAdmin err: %s", "user not found")
		resp.Error = Error(UserNotFound, "")
		return
	}

	u := new(model.User)
	if req.NickName == model.AnonymousUser {
		log.Errorf("UpdateUserAdmin err: %s", "can not be anonymous name")
		resp.Error = Error(NickNameAlreadyBeUsed, "")
		return
	}

	// admin can change nickname no more limit
	if req.NickName != "" && req.NickName != uu.NickName {
		u.NickName = req.NickName
		repeat, err := u.IsNickNameRepeat()
		if err != nil {
			log.Errorf("UpdateUserAdmin err: %s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}
		if repeat {
			log.Errorf("UpdateUserAdmin err: %s", "nickname already use by other")
			resp.Error = Error(NickNameAlreadyBeUsed, "")
			return
		}
	}
	u.Id = req.Id
	u.Password = req.Password

	// change user status, 1->2, 2->1
	u.Status = req.Status

	// vip change
	u.Vip = uu.Vip
	if req.Vip == 1 {
		u.Vip = 1
	} else if req.Vip == 2 {
		u.Vip = 0
	}
	err = u.UpdateInfoMustVip()
	if err != nil {
		log.Errorf("UpdateUserAdmin err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	err = session.Mgr.RefreshUser([]string{fmt.Sprintf("%d", u.Id)}, SessionExpireTime)
	if err != nil {
		log.Errorf("UpdateUserAdmin err:%s", err.Error())
		resp.Error = Error(RefreshUserCacheError, err.Error())
		return
	}
	resp.Data = u
	resp.Flag = true
}
