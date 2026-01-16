package httpapi

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/middleware"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
)

type Handlers struct {
	UserRepo  *repo.UserRepo
	JWTSecret string
	TTL       time.Duration
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewRouter(userRepo *repo.UserRepo, jwtSecret string, ttlMinutes int) *gin.Engine {
	r := gin.Default()
	h := &Handlers{UserRepo: userRepo, JWTSecret: jwtSecret, TTL: time.Duration(ttlMinutes) * time.Minute}

	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })
	r.POST("/auth/login", h.Login)

	protected := r.Group("/")
	protected.Use(middleware.AuthJWT(jwtSecret))
	protected.GET("/me", h.Me)

	admin := r.Group("/admin")
	admin.Use(middleware.AuthJWT(jwtSecret), middleware.RequireRole("ADMIN"))
	admin.GET("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"pong": "admin"}) })

	return r
}

func (h *Handlers) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil || req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email/password required"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	u, err := h.UserRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(500, gin.H{"error": "db error"})
		return
	}
	if u.Status == "locked" {
		c.JSON(http.StatusForbidden, gin.H{"error": "user locked"})
		return
	}
	if !auth.VerifyPassword(u.PasswordHash, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	roles, err := h.UserRepo.RolesOfUser(ctx, u.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "db error"})
		return
	}

	token, err := auth.Sign(h.JWTSecret, h.TTL, u.ID.String(), u.Email, roles)
	if err != nil {
		c.JSON(500, gin.H{"error": "token error"})
		return
	}

	c.JSON(200, gin.H{
		"access_token": token,
		"token_type":   "Bearer",
		"expires_in":   int(h.TTL.Seconds()),
	})
}

func (h *Handlers) Me(c *gin.Context) {
	claimsAny, _ := c.Get(middleware.CtxClaims)
	claims := claimsAny.(*auth.Claims)
	c.JSON(200, gin.H{"user_id": claims.UserID, "email": claims.Email, "roles": claims.Roles})
}
