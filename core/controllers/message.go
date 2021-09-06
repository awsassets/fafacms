package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/hunterhug/fafacms/core/model"
	"github.com/hunterhug/fafacms/core/util"
	log "github.com/hunterhug/golog"
	"math"
)

type ListMessageRequest struct {
	MessageId       int64    `json:"message_id"`
	MessageType     int      `json:"message_type" validate:"oneof=-1 0 1 2 3 4 5 6 7 8 9 10 11"`
	ReceiveUserId   int64    `json:"receive_user_id"`
	ChanelUserId    int64    `json:"chanel_user_id"`
	ReceiveStatus   int      `json:"receive_status" validate:"oneof=-1 0 1 2"`
	SendStatus      int      `json:"send_status" validate:"oneof=-1 0 1"`
	GlobalMessageId int64    `json:"global_message_id"`
	CreateTimeBegin int64    `json:"create_time_begin"`
	CreateTimeEnd   int64    `json:"create_time_end"`
	Sort            []string `json:"sort"`
	PageHelp
}

type ListMessageResponse struct {
	Messages      []model.Message               `json:"messages"`
	Comments      map[int64]model.Comment       `json:"comments"`
	Contents      map[int64]model.ContentHelper `json:"contents"`
	ExtraUsers    map[int64]model.UserHelper    `json:"extra_users"`
	ExtraComments map[int64]model.CommentHelper `json:"extra_comments"`
	UnRead        map[string]int                `json:"un_read"`
	PageHelp
}

func ListMessageHelper(c *gin.Context, isAdmin bool) {
	resp := new(Resp)

	respResult := new(ListMessageResponse)
	req := new(ListMessageRequest)
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
		log.Errorf("ListMessageHelper err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	var yourUserId int64 = 0
	var all = true

	if !isAdmin {
		uu, err := GetUserSession(c)
		if err != nil {
			log.Errorf("ListMessageHelper err: %s", err.Error())
			resp.Error = Error(GetUserSessionError, err.Error())
			return
		}

		yourUserId = uu.Id
		all = false

		// global message insert
		go model.InsertGlobalMessageToUser(yourUserId)
	}

	// new query list session
	session := model.FaFaRdb.Client.NewSession()
	defer session.Close()

	// group list where prepare
	session.Table(new(model.Message)).Where("1=1")

	if req.MessageId != 0 {
		session.And("id=?", req.MessageId)
	}

	if req.MessageType != -1 {
		session.And("message_type=?", req.MessageType)
	}

	if req.ReceiveStatus != -1 {
		session.And("receive_status=?", req.ReceiveStatus)
	}

	if req.CreateTimeBegin > 0 {
		session.And("create_time>=?", req.CreateTimeBegin)
	}

	if req.CreateTimeEnd > 0 {
		session.And("create_time<?", req.CreateTimeEnd)
	}

	if !all {
		// search chanel not delete message
		if req.ChanelUserId != 0 {
			session.And("send_status=?", 0).And("private_chanel=?", model.GetChanelName(yourUserId, req.ChanelUserId))
		} else {
			// search your message which not delete
			session.And("receive_user_id=?", yourUserId).And("receive_status!=?", 2)
		}
	} else {
		if req.ReceiveUserId != 0 {
			if req.ChanelUserId != 0 {
				// admin can search all
				if req.SendStatus != -1 {
					session.And("send_status=?", req.SendStatus)
				}
				session.And("private_chanel=?", model.GetChanelName(req.ReceiveUserId, req.ChanelUserId))
			} else {
				session.And("receive_user_id=?", req.ReceiveUserId)
			}
		}

	}

	if req.GlobalMessageId != 0 {
		session.And("global_message_id=?", req.GlobalMessageId)
	}

	// count unread messages
	countMap, err := model.GroupCount(yourUserId)
	if err != nil {
		log.Errorf("ListMessageHelper err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	// if count>0 start list
	cs := make([]model.Message, 0)
	p := &req.PageHelp

	// sql build
	p.build(session, req.Sort, model.MessageSortName)

	// do query
	total, err := session.FindAndCount(&cs)
	if err != nil {
		log.Errorf("ListMessageHelper err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	contentIds := make(map[int64]struct{})
	userIds := make(map[int64]struct{})
	commentIds := make(map[int64]struct{})
	for _, v := range cs {
		// content ids collect
		if v.ContentId != 0 {
			contentIds[v.ContentId] = struct{}{}
		}

		switch v.MessageType {
		case model.MessageTypeCommentForContent, model.MessageTypeCommentForComment:
			// if comment is anonymous must hide user id
			if v.CommentAnonymous == 1 && v.CommentIsYourSelf == 0 {
				if !all {
					v.UserId = 0
				}
			}
		}

		// user id collect
		if v.UserId != 0 {
			userIds[v.UserId] = struct{}{}
		}

		// comment id collect
		if v.CommentId != 0 {
			commentIds[v.CommentId] = struct{}{}
		}
	}

	// get all none delete comment
	comments, err := model.GetComment(util.MapToArray(commentIds), all)
	if err != nil {
		log.Errorf("ListMessageHelper err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	// extra comment id collect
	commentIds2 := make(map[int64]struct{})
	for k := range comments {
		commentIds2[k] = struct{}{}
	}

	// get extra comment info
	comments2, user2, err := model.GetCommentAndCommentUser(util.MapToArray(commentIds2), all, util.MapToArray(userIds), yourUserId)
	if err != nil {
		log.Errorf("ListMessageHelper err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	// get content base info
	contents, err := model.GetContentHelper(util.MapToArray(contentIds), all, yourUserId)
	if err != nil {
		log.Errorf("ListMessageHelper err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	// result
	respResult.Messages = cs
	respResult.UnRead = countMap
	respResult.Comments = comments
	respResult.ExtraComments = comments2
	respResult.ExtraUsers = user2
	respResult.Contents = contents
	p.Pages = int(math.Ceil(float64(total) / float64(p.Limit)))
	p.Total = int(total)
	respResult.PageHelp = *p
	resp.Data = respResult
	resp.Flag = true
}

func ListMessage(c *gin.Context) {
	ListMessageHelper(c, false)
}

func ListAllMessage(c *gin.Context) {
	ListMessageHelper(c, true)
}

type ReadMessageRequest struct {
	Ids []int64 `json:"ids"`
}

func ReadMessage(c *gin.Context) {
	resp := new(Resp)
	req := new(ReadMessageRequest)
	defer func() {
		JSONL(c, 200, req, resp)
	}()

	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	if len(req.Ids) == 0 {
		log.Errorf("ReadMessage err: %s", "ids empty")
		resp.Error = Error(ParasError, "ids empty")
		return
	}

	uu, err := GetUserSession(c)
	if err != nil {
		log.Errorf("ReadMessage err: %s", err.Error())
		resp.Error = Error(GetUserSessionError, err.Error())
		return
	}

	m := new(model.Message)
	m.ReceiveUserId = uu.Id
	m.ReceiveStatus = 1
	err = m.ReceiveUpdate(req.Ids)
	if err != nil {
		log.Errorf("ReadMessage err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	resp.Flag = true
}

type DeleteMessageRequest struct {
	Ids []int64 `json:"ids"`
}

func DeleteMessage(c *gin.Context) {
	resp := new(Resp)
	req := new(DeleteMessageRequest)
	defer func() {
		JSONL(c, 200, req, resp)
	}()

	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	if len(req.Ids) == 0 {
		log.Errorf("DeleteMessage err: %s", "ids empty")
		resp.Error = Error(ParasError, "ids empty")
		return
	}
	uu, err := GetUserSession(c)
	if err != nil {
		log.Errorf("DeleteMessage err: %s", err.Error())
		resp.Error = Error(GetUserSessionError, err.Error())
		return
	}

	m := new(model.Message)
	m.ReceiveUserId = uu.Id
	m.ReceiveStatus = 2
	err = m.ReceiveUpdate(req.Ids)
	if err != nil {
		log.Errorf("DeleteMessage err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	resp.Flag = true
}

type CreateGlobalMessageRequest struct {
	UserIds   []int64 `json:"user_ids"`
	AllPeople bool    `json:"all_people"`
	Message   string  `json:"message"`
	RightNow  bool    `json:"right_now"` // valid in all_people true
}

func CreateGlobalMessage(c *gin.Context) {
	resp := new(Resp)
	req := new(CreateGlobalMessageRequest)
	defer func() {
		JSONL(c, 200, req, resp)
	}()

	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	if !req.AllPeople {
		// alone sent to all user_ids, should not empty
		if len(req.UserIds) == 0 {
			log.Errorf("CreateGlobalMessage err: %s", "user_ids empty")
			resp.Error = Error(ParasError, "user_ids empty")
			return
		}

		// user should all exit
		if !model.UserAllExist(req.UserIds) {
			log.Errorf("CreateGlobalMessage err: %s", "user_ids not right")
			resp.Error = Error(ParasError, "user_ids not right")
			return
		}

	}

	if !req.AllPeople {
		for _, v := range req.UserIds {
			m := new(model.Message)
			m.MessageType = model.MessageTypeGlobal
			m.SendMessage = req.Message
			m.ReceiveUserId = v
			err := m.Insert()
			if err != nil {
				log.Errorf("CreateGlobalMessage err: %s", err.Error())
				resp.Error = Error(DBError, err.Error())
				return
			}
		}

		resp.Flag = true
		return
	}

	// insert global table
	gm := new(model.GlobalMessage)
	gm.SendMessage = req.Message

	// how much user now
	uCount, err := model.UserCount()
	if err != nil {
		log.Errorf("CreateGlobalMessage err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	gm.Total = uCount
	if req.RightNow {
		gm.Status = 1
	}

	err = gm.Insert()
	if err != nil {
		log.Errorf("CreateGlobalMessage err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	resp.Flag = true
}

type ListGlobalMessageRequest struct {
	Id              int64    `json:"id"`
	CreateTimeBegin int64    `json:"create_time_begin"`
	CreateTimeEnd   int64    `json:"create_time_end"`
	Status          int      `json:"status" validate:"oneof=-1 0 1 2"`
	Sort            []string `json:"sort"`
	PageHelp
}

type ListGlobalMessageResponse struct {
	Message []model.GlobalMessage `json:"message"`
	PageHelp
}

func ListGlobalMessage(c *gin.Context) {
	resp := new(Resp)
	respResult := new(ListGlobalMessageResponse)

	req := new(ListGlobalMessageRequest)
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
		log.Errorf("ListGlobalMessage err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	// new query list session
	session := model.FaFaRdb.Client.NewSession()
	defer session.Close()

	// group list where prepare
	session.Table(new(model.GlobalMessage)).Where("1=1")

	// query prepare
	if req.Id != 0 {
		session.And("id=?", req.Id)
	}

	if req.Status != -1 {
		session.And("status=?", req.Status)
	}

	if req.CreateTimeBegin > 0 {
		session.And("create_time>=?", req.CreateTimeBegin)
	}

	if req.CreateTimeEnd > 0 {
		session.And("create_time<?", req.CreateTimeEnd)
	}

	cs := make([]model.GlobalMessage, 0)
	p := &req.PageHelp

	// sql build
	p.build(session, req.Sort, model.GlobalMessageSortName)

	// do query
	total, err := session.FindAndCount(&cs)
	if err != nil {
		log.Errorf("ListGlobalMessage err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	// result
	respResult.Message = cs
	p.Pages = int(math.Ceil(float64(total) / float64(p.Limit)))
	p.Total = int(total)
	respResult.PageHelp = *p
	resp.Data = respResult
	resp.Flag = true
}

type UpdateGlobalMessageRequest struct {
	Id     int64 `json:"id" validate:"required"`
	Status int   `json:"status" validate:"oneof=0 1 2"`
}

func UpdateGlobalMessage(c *gin.Context) {
	resp := new(Resp)
	req := new(UpdateGlobalMessageRequest)
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
		log.Errorf("UpdateGlobalMessage err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	n := new(model.GlobalMessage)
	n.Id = req.Id

	exist, err := n.Get()
	if err != nil {
		log.Errorf("UpdateGlobalMessage err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	if !exist {
		log.Errorf("UpdateGlobalMessage err: %s", "global message not found")
		resp.Error = Error(GlobalMessageNotFound, "")
		return
	}

	after := new(model.GlobalMessage)
	after.Id = n.Id
	after.Status = req.Status

	_, err = after.Update()
	if err != nil {
		log.Errorf("UpdateGlobalMessage err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	resp.Flag = true
}

type SendPrivateMessageRequest struct {
	UserId  int64  `json:"user_id"`
	Message string `json:"message"`
}

func SendPrivateMessage(c *gin.Context) {
	resp := new(Resp)
	req := new(SendPrivateMessageRequest)
	defer func() {
		JSONL(c, 200, req, resp)
	}()

	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	if req.UserId == 0 {
		log.Errorf("SendPrivateMessage err: %s", "user_id empty")
		resp.Error = Error(ParasError, "user_id empty")
		return
	}

	targetUser := new(model.User)
	targetUser.Id = req.UserId
	ok, err := targetUser.GetRaw()
	if err != nil {
		log.Errorf("SendPrivateMessage err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	if !ok {
		log.Errorf("SendPrivateMessage err: %s", "user not found")
		resp.Error = Error(UserNotFound, "")
		return
	}

	if targetUser.Status == 0 {
		log.Errorf("SendPrivateMessage err: %s", "user not activate")
		resp.Error = Error(UserNotActivate, "")
		return
	}

	uu, err := GetUserSession(c)
	if err != nil {
		log.Errorf("SendPrivateMessage err: %s", err.Error())
		resp.Error = Error(GetUserSessionError, err.Error())
		return
	}

	err = model.Private(uu.Id, targetUser.Id, req.Message)
	if err != nil {
		log.Errorf("SendPrivateMessage err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	resp.Flag = true
}

type DeletePrivateMessageRequest struct {
	Ids []int64 `json:"ids"`
}

func DeletePrivateMessage(c *gin.Context) {
	resp := new(Resp)
	req := new(DeletePrivateMessageRequest)
	defer func() {
		JSONL(c, 200, req, resp)
	}()

	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	if len(req.Ids) == 0 {
		log.Errorf("DeletePrivateMessage err: %s", "ids empty")
		resp.Error = Error(ParasError, "ids empty")
		return
	}
	uu, err := GetUserSession(c)
	if err != nil {
		log.Errorf("DeletePrivateMessage err: %s", err.Error())
		resp.Error = Error(GetUserSessionError, err.Error())
		return
	}

	m := new(model.Message)
	m.SendUserId = uu.Id
	m.SendStatus = 1
	err = m.SendUpdate(req.Ids)
	if err != nil {
		log.Errorf("DeletePrivateMessage err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	resp.Flag = true
}
