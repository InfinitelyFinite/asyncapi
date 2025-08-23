package apiserver

import (
	"asyncapi/config"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var signingMethod = jwt.SigningMethodHS256

type JwtManager struct {
	config *config.Config
}

func NewJwtManager(config *config.Config) *JwtManager {
	return &JwtManager{config}
}

type TokenPair struct {
	AccessToken  *jwt.Token
	RefreshToken *jwt.Token
}

type CustomClaims struct {
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

func (j *JwtManager) Parse(token string) (*jwt.Token, error) {
	parser := jwt.NewParser()
	jwtToken, err := parser.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if t.Method != signingMethod {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(j.config.JwtSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	return jwtToken, nil
}

func (j *JwtManager) IsAccessToken(token *jwt.Token) bool {
	jwtClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}
	if tokenType, ok := jwtClaims["token_type"]; ok {
		return tokenType == "access"
	}
	return false
}

func (j *JwtManager) GenerateTokenPair(userId uuid.UUID) (*TokenPair, error) {
	issuer := "http://" + j.config.ApiServerHost + ":" + j.config.ApiServerPort
	now := time.Now()

	jwtAccessToken := jwt.NewWithClaims(signingMethod, CustomClaims{
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userId.String(),
			Issuer:    issuer,
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute * 15)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	})

	key := []byte(j.config.JwtSecret)
	signedAccessToken, err := jwtAccessToken.SignedString(key)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// parsing because parse method populates a lot of the things in the token that we dont get when we generate it
	accessToken, err := j.Parse(signedAccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to parse access token: %w", err)
	}

	jwtRefreshToken := jwt.NewWithClaims(signingMethod, CustomClaims{
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userId.String(),
			Issuer:    issuer,
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24 * 30)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	})

	signedRefreshToken, err := jwtRefreshToken.SignedString(key)
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	refreshToken, err := j.Parse(signedRefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to parse refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
