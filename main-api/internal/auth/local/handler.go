package local

import (
	"github.com/gin-gonic/gin"
	"github.com/yerlanov/go-tour/common/session"
	"github.com/yerlanov/go-tour/main-api/internal/storage"
)

type Handler struct {
	session session.Session
	service Service
	storage storage.Storage
}

func NewHandler(session session.Session, storage storage.Storage) *Handler {
	return &Handler{
		session: session,
		service: NewService(storage, session),
	}
}

func (h *Handler) RegisterRouter(router *gin.RouterGroup) {
	router.GET("/logout", h.Logout)
	router.POST("/registration", h.Registration)
	router.POST("/login", h.Login)
}

func (h *Handler) Logout(ctx *gin.Context) {
	content, exists := ctx.Get("session")
	if !exists {
		ctx.AbortWithStatus(401)
		return
	}

	sess, ok := content.(session.Content)
	if !ok {
		ctx.AbortWithStatus(401)
		return
	}

	delete(sess.Values, "userId")

	err := h.session.Set(ctx, sess)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}

	ctx.Redirect(302, "/")
}

func (h *Handler) Registration(ctx *gin.Context) {
	var userCreateDto UserCreateDto
	err := ctx.ShouldBindJSON(&userCreateDto)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
	}

	user, err := h.service.Registration(ctx, userCreateDto)
	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, user)
}

func (h *Handler) Login(ctx *gin.Context) {
	var login Login
	err := ctx.ShouldBindJSON(&login)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
	}

	user, err := h.service.Login(ctx, login.Email, login.Password)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, user)
}
