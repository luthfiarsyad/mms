package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luthfiarsyad/mms/internal/infrastructure/security"
	"github.com/luthfiarsyad/mms/internal/interface/http/request"
	"github.com/luthfiarsyad/mms/internal/usecase"
	"golang.org/x/crypto/bcrypt"

	domain "github.com/luthfiarsyad/mms/internal/domain/user"
	mysqlrepo "github.com/luthfiarsyad/mms/internal/infrastructure/persistence/mysql"
)

// For simplicity we wire dependencies inside NewAuthHandler
type AuthHandler struct {
	usecase *usecase.AuthUsecase
}

func NewAuthHandler() *AuthHandler {
	pas := security.NewPasetoService()
	db := mysqlrepo.Get()
	ur := mysqlrepo.NewUserRepo(db)
	us := domain.NewService(ur)
	uc := usecase.NewAuthUsecase(us, pas)
	return &AuthHandler{usecase: uc}
}
func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password),
		bcrypt.DefaultCost)
	if err != nil {

		return
	}
	u := &domain.User{
		Name:  req.Name,
		Email: req.Email,
		// password set in usecase
	}
	if err := h.usecase.Register(c.Request.Context(), u, string(hashed)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": u.ID, "email": u.Email,
		"created_at": u.CreatedAt})
}
func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// passwordCheck function uses bcrypt.CompareHashAndPassword
	passwordCheck := func(hashed, plain string) error {
		return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	}
	token, err := h.usecase.Login(c.Request.Context(), req.Email, req.Password,
		passwordCheck)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_token": token, "token_type": "bearer",
		"expires_in": 24 * 3600})
}
