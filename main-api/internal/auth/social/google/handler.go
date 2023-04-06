package google

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/yerlanov/go-tour/common/session"
	"github.com/yerlanov/go-tour/main-api/internal/config"
	"github.com/yerlanov/go-tour/main-api/internal/storage"
	"github.com/yerlanov/go-tour/main-api/internal/storage/user"
	"github.com/yerlanov/go-tour/main-api/pkg/oauth/google"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"time"
)

type Handler struct {
	config           config.Config
	googleAuthConfig *google.Auth
	userRepository   user.Repository
	session          session.Session
}

func NewHandler(config *config.Config, storage storage.Storage, session session.Session) *Handler {
	return &Handler{
		googleAuthConfig: google.LoadAuthConfig(google.Config{
			ClientID:     config.Social.Google.ClientID,
			ClientSecret: config.Social.Google.ClientSecret,
			RedirectURL:  config.Social.Google.RedirectURL,
		}),
		userRepository: user.NewRepository(storage),
		session:        session,
	}
}

func (h *Handler) RegisterRouter(router *gin.RouterGroup) {
	router.GET("/login", h.Login)
	router.GET("/callback", h.Callback)
}

func (h *Handler) Login(ctx *gin.Context) {
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

	state, err := generateRandomState(32)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
	}
	sessionContent.Values["state"] = state

	redirectParam := ctx.Query("redirect_url")
	if redirectParam != "" {
		sessionContent.Values["redirectUrl"] = redirectParam
	}

	err = h.session.Set(ctx, sessionContent)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	url := h.googleAuthConfig.Google.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *Handler) Callback(ctx *gin.Context) {
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

	queryState := ctx.Query("state")
	if queryState != sessionContent.Values["state"] {
		ctx.JSON(500, gin.H{"error": "state mismatch"})
		return
	}

	delete(sessionContent.Values, "state")

	redirectUrl, ok := sessionContent.Values["redirectUrl"].(string)
	if !ok {
		redirectUrl = "/"
	}

	code := ctx.Query("code")
	token, err := h.googleAuthConfig.Google.Exchange(ctx, code)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	client := h.googleAuthConfig.Google.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var userDto User
	err = json.Unmarshal(data, &userDto)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	id, err := h.saveUser(ctx, userDto, token)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	sessionContent.Values["userId"] = id
	err = h.session.Set(ctx, sessionContent)

	ctx.Redirect(http.StatusTemporaryRedirect, redirectUrl)
}

func (h *Handler) saveUser(ctx context.Context, userDto User, token *oauth2.Token) (*string, error) {
	id, err := h.userRepository.Upsert(ctx, user.User{
		ID:         primitive.ObjectID{},
		Email:      userDto.Email,
		FirstName:  userDto.GivenName,
		LastName:   userDto.FamilyName,
		Provider:   "google_oauth2",
		AccessKey:  token.AccessToken,
		RefreshKey: token.RefreshToken,
		ExpireAt:   token.Expiry,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	})
	if err != nil {
		return nil, err
	}
	return id, nil
}

func generateRandomState(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buffer), nil
}
