package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidToken       = errors.New("invalid token")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserLocked         = errors.New("user locked")
)

type Claims struct {
	UserID   string   `json:"user_id"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
	SchoolID string   `json:"school_id,omitempty"` // chỉ có giá trị khi user là SCHOOL_ADMIN
	jwt.RegisteredClaims
}

type Authenticator struct {
	Secret     string
	TTLSeconds int
}

// NewAuthenticator tạo mới authenticator
func NewAuthenticator(secret string, ttlMinutes int) *Authenticator {
	return &Authenticator{
		Secret:     secret,
		TTLSeconds: ttlMinutes * 60,
	}
}

// SignToken tạo JWT token bằng Authenticator.
// schoolID rỗng ("") cho SUPER_ADMIN/TEACHER/PARENT, có giá trị cho SCHOOL_ADMIN
func (a *Authenticator) SignToken(userID, email string, roles []string, schoolID string) (string, error) {
	return Sign(a.Secret, time.Duration(a.TTLSeconds)*time.Second, userID, email, roles, schoolID)
}

func VerifyPassword(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}

// Parse giải mã và xác thực JWT token

// Sign tạo ra một JWT token sử dụng thuật toán HS256 để ký với thông tin người dùng.
// Các tham số:
//   - secret: chuỗi bí mật dùng để ký token.
//   - ttl: thời gian sống của token (time-to-live).
//   - userID: ID của người dùng.
//   - email: email của người dùng.
//   - roles: danh sách vai trò của người dùng.
//
// Hàm sẽ tạo một struct Claims chứa thông tin người dùng và các trường chuẩn của JWT (IssuedAt, ExpiresAt, Subject).
// Sau đó, hàm tạo một token mới với các claims này, ký bằng thuật toán HS256 và trả về chuỗi token đã ký.
func Sign(secret string, ttl time.Duration, userID, email string, roles []string, schoolID string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Email:    email,
		Roles:    roles,
		SchoolID: schoolID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID, // định danh đối tượng mà token cấp cho, gán bằng userID để xác định token này thuộc về người dùng nào.
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)), // token hết hạn, cộng ttl vào thời điểm hiện tại (now)
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(secret))
}

// Parse kiểm tra và giải mã một JWT token sử dụng thuật toán HS256 với secret cho trước.
// Các tham số:
//   - secret: chuỗi bí mật dùng để xác thực token.
//   - tokenStr: chuỗi JWT token cần kiểm tra.
//
// Hàm sẽ xác thực token, kiểm tra tính hợp lệ và giải mã các claims (thông tin người dùng) từ token.
// Nếu token hợp lệ, trả về con trỏ tới struct Claims chứa thông tin người dùng.
// Nếu token không hợp lệ hoặc không giải mã được, trả về lỗi ErrInvalidToken.
func Parse(secret, tokenStr string) (*Claims, error) {
	t, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (any, error) {
		if token.Method == nil || token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil || !t.Valid {
		return nil, ErrInvalidToken
	}
	claims, ok := t.Claims.(*Claims)
	if !ok {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
