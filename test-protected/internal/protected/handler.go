package protected

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-tour/test-protected/internal/session"
	"github.com/go-tour/test-protected/internal/storage"
	"github.com/go-tour/test-protected/internal/storage/user"
	"net/http"
)

type Handler struct {
	userRepository user.Repository
}

func NewHandler(storage storage.Storage) *Handler {
	return &Handler{
		userRepository: user.NewRepository(storage),
	}
}

func (h *Handler) RegisterRouter(router *gin.RouterGroup) {
	router.GET("", h.Protected)
}

func (h *Handler) Protected(ctx *gin.Context) {
	value, exists := ctx.Get("session")
	if !exists {
		ctx.JSON(500, gin.H{"error": "session not found"})
		return
	}

	sessionContent, ok := value.(session.Content)
	if !ok {
		// Handle cases when the session content has an unexpected type
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Session content has an unexpected type"})
		return
	}

	if _, ok := sessionContent.Values["userId"]; !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User is not logged in"})
		return
	}

	fmt.Println(sessionContent.Values["userId"].(string))

	user, err := h.userRepository.GetByID(ctx, sessionContent.Values["userId"].(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}
