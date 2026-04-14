package handler

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kingwrcy/moments/db"
	"github.com/kingwrcy/moments/vo"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserHandler struct {
	base BaseHandler
}

type loginSuccessDTO struct {
	Token    string `json:"token,omitempty"`    // token
	Username string `json:"username,omitempty"` //用户名
	Id       int32  `json:"id,omitempty"`       //用户ID
}

func NewUserHandler(injector do.Injector) *UserHandler {
	return &UserHandler{do.MustInvoke[BaseHandler](injector)}
}

// Login godoc
//
//	@Tags		User
//	@Summary	用户登录
//	@Accept		json
//	@Produce	json
//	@Param		object	body		vo.LoginReq	true	"用户登录"
//	@Success	200		{object}	loginSuccessDTO
//	@Router		/api/user/login [post]
func (u UserHandler) Login(c echo.Context) error {
	var req vo.LoginReq
	err := c.Bind(&req)
	if err != nil {
		return FailResp(c, ParamError)
	}

	var user db.User
	err = u.base.db.Where("username = ?", req.Username).First(&user).Error
	if err != nil {
		return FailRespWithMsg(c, Fail, "用户不存在或密码不正确")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		return FailRespWithMsg(c, Fail, "用户不存在或密码不正确")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"userId":   user.Id,
	})

	tokenString, err := token.SignedString([]byte(u.base.cfg.JwtKey))
	if err != nil {
		u.base.log.Error().Msgf("生成jwt token异常:%s", err)
		return FailRespWithMsg(c, Fail, "登录异常")
	}
	return SuccessResp(c, loginSuccessDTO{
		Token:    tokenString,
		Username: user.Username,
		Id:       user.Id,
	})
}

// Reg godoc
//
//	@Tags		User
//	@Summary	用户注册
//	@Accept		json
//	@Produce	json
//	@Param		object	body	vo.RegReq	true	"用户注册"
//	@Success	200
//	@Router		/api/user/reg [post]
func (u UserHandler) Reg(c echo.Context) error {
	var (
		req         vo.RegReq
		count       int64
		user        db.User
		now         = time.Now()
		sysConfig   db.SysConfig
		sysConfigVO vo.FullSysConfigVO
	)

	u.base.db.First(&sysConfig)
	_ = json.Unmarshal([]byte(sysConfig.Content), &sysConfigVO)

	if !sysConfigVO.EnableRegister {
		return FailRespWithMsg(c, Fail, "当前未开启注册用户")
	}

	err := c.Bind(&req)
	if err != nil {
		return FailResp(c, ParamError)
	}

	if len(req.Username) < 3 {
		return FailRespWithMsg(c, Fail, "用户名最少3个字符")
	}
	if req.Password != req.RepeatPassword {
		return FailRespWithMsg(c, Fail, "两次密码不一致")
	}
	u.base.db.Table("User").Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		return FailRespWithMsg(c, Fail, "用户名已存在")
	}
	user.Username = req.Username
	pwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		u.base.log.Error().Msgf("密码加密异常:%s", err)
		return FailRespWithMsg(c, Fail, "密码加密异常")
	}
	user.Password = string(pwd)
	user.CreatedAt = &now
	user.UpdatedAt = &now
	user.Nickname = req.Username
	user.AvatarUrl = "/avatar.webp"
	user.Slogan = "修道者，逆天而行，注定要一生孤独。"
	user.CoverUrl = "/cover.webp"
	if err := u.base.db.Save(&user).Error; err != nil {
		u.base.log.Error().Msgf("注册用户异常:%s", err)
		return FailRespWithMsg(c, Fail, "注册用户异常")
	}
	return SuccessResp(c, h{})
}

// ProfileForUser godoc
//
//	@Tags		User
//	@Summary	获取指定用户信息
//	@Accept		json
//	@Produce	json
//	@param		string	path		string	true	"用户名"
//	@Success	200		{object}	db.User
//	@Router		/api/user/profile/{username} [post]
func (u UserHandler) ProfileForUser(c echo.Context) error {
	username := c.Param("username")
	var user db.User
	u.base.db.Select("username", "nickname", "slogan", "id", "avatarUrl", "coverUrl", "email").Find(&user, "username = ?", username)
	return SuccessResp(c, user)
}

// Profile godoc
//
//	@Tags			User
//	@Summary		获取用户信息
//	@Description	当前如果已经登录了,获取当前用户信息,否则获取管理员的用户信息
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	db.User
//	@Router			/user/profile [post]
func (u UserHandler) Profile(c echo.Context) error {

	context := c.(CustomContext)
	currentUser := context.CurrentUser()
	if currentUser == nil {
		u.base.db.Select("username", "nickname", "slogan", "id", "avatarUrl", "coverUrl", "email").First(&currentUser)
	}

	return SuccessResp(c, currentUser)
}

// SaveProfile godoc
//
//	@Tags		User
//	@Summary	保存用户信息
//	@Accept		json
//	@Produce	json
//	@Param		object		body	vo.ProfileReq	true	"保存用户信息"
//	@Param		x-api-token	header	string			true	"登录TOKEN"
//	@Success	200
//	@Router		/api/user/saveProfile [post]
func (u UserHandler) SaveProfile(c echo.Context) error {
	var (
		req  vo.ProfileReq
		user db.User
	)
	err := c.Bind(&req)
	if err != nil {
		return FailResp(c, ParamError)
	}
	context := c.(CustomContext)
	currentUser := context.CurrentUser()
	if currentUser == nil {
		return FailResp(c, TokenMissing)
	}
	u.base.db.Find(&user, currentUser.Id)
	if req.Password != "" {
		password, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
		if err != nil {
			return FailResp(c, Fail)
		}
		user.Password = string(password)
	}
	user.Nickname = req.Nickname
	user.AvatarUrl = req.AvatarUrl
	user.Slogan = req.Slogan
	user.CoverUrl = req.CoverUrl
	user.Email = req.Email

	if err := u.base.db.Save(&user).Error; err != nil {
		return FailResp(c, Fail)
	}
	return SuccessResp(c, h{})
}

// ListUsers godoc
//
//	@Tags		User
//	@Summary	管理员获取用户列表
//	@Accept		json
//	@Produce	json
//	@Param		x-api-token	header	string	true	"登录TOKEN"
//	@Success	200
//	@Router		/api/user/list [post]
func (u UserHandler) ListUsers(c echo.Context) error {
	context := c.(CustomContext)
	currentUser := context.CurrentUser()
	if currentUser == nil || currentUser.Id != 1 {
		return FailRespWithMsg(c, Fail, "没有权限")
	}

	var users []vo.AdminUserVO
	if err := u.base.db.Table("User").
		Select("id", "username", "nickname", "avatarUrl", "slogan", "coverUrl", "email", "createdAt", "updatedAt").
		Order("id asc").
		Find(&users).Error; err != nil {
		return FailRespWithMsg(c, Fail, "读取用户列表失败")
	}

	return SuccessResp(c, users)
}

// AdminSaveUser godoc
//
//	@Tags		User
//	@Summary	管理员保存用户资料
//	@Accept		json
//	@Produce	json
//	@Param		object		body	vo.AdminUserSaveReq	true	"管理员保存用户资料"
//	@Param		x-api-token	header	string				true	"登录TOKEN"
//	@Success	200
//	@Router		/api/user/adminSave [post]
func (u UserHandler) AdminSaveUser(c echo.Context) error {
	var req vo.AdminUserSaveReq
	if err := c.Bind(&req); err != nil {
		return FailResp(c, ParamError)
	}

	context := c.(CustomContext)
	currentUser := context.CurrentUser()
	if currentUser == nil || currentUser.Id != 1 {
		return FailRespWithMsg(c, Fail, "没有权限")
	}

	var user db.User
	if err := u.base.db.First(&user, req.ID).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return FailRespWithMsg(c, Fail, "用户不存在")
	}

	user.Nickname = req.Nickname
	user.Slogan = req.Slogan
	user.Email = req.Email

	if req.Password != "" {
		password, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
		if err != nil {
			return FailRespWithMsg(c, Fail, "密码加密失败")
		}
		user.Password = string(password)
	}

	if err := u.base.db.Save(&user).Error; err != nil {
		return FailRespWithMsg(c, Fail, "保存用户失败")
	}

	return SuccessResp(c, h{})
}

// DeleteUser godoc
//
//	@Tags		User
//	@Summary	管理员删除用户
//	@Accept		json
//	@Produce	json
//	@Param		id			query	int		true	"用户ID"
//	@Param		x-api-token	header	string	true	"登录TOKEN"
//	@Success	200
//	@Router		/api/user/delete [post]
func (u UserHandler) DeleteUser(c echo.Context) error {
	context := c.(CustomContext)
	currentUser := context.CurrentUser()
	if currentUser == nil || currentUser.Id != 1 {
		return FailRespWithMsg(c, Fail, "没有权限")
	}

	id, err := strconv.Atoi(c.QueryParam("id"))
	if err != nil || id <= 0 {
		return FailResp(c, ParamError)
	}
	if int32(id) == currentUser.Id || id == 1 {
		return FailRespWithMsg(c, Fail, "管理员账号不允许删除")
	}

	var user db.User
	if err = u.base.db.First(&user, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return FailRespWithMsg(c, Fail, "用户不存在")
	}

	var memoCount int64
	u.base.db.Table("Memo").Where("userId = ?", id).Count(&memoCount)
	if memoCount > 0 {
		return FailRespWithMsg(c, Fail, "该用户仍有动态内容，暂不支持直接删除")
	}

	var commentCount int64
	u.base.db.Table("Comment").Where("author = ?", id).Count(&commentCount)
	if commentCount > 0 {
		return FailRespWithMsg(c, Fail, "该用户仍有评论内容，暂不支持直接删除")
	}

	if err = u.base.db.Delete(&user).Error; err != nil {
		return FailRespWithMsg(c, Fail, "删除用户失败")
	}

	return SuccessResp(c, h{})
}
