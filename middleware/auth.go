package middleware

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"go-menu/resource/user"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// Auth0Config Auth0の設定情報
type Auth0Config struct {
	Domain   string
	Audience string
}

// JWKSResponse Auth0 JWKS レスポンス構造体
type JWKSResponse struct {
	Keys []JWK `json:"keys"`
}

// JWK JSON Web Key 構造体
type JWK struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// AuthMiddleware Auth0 JWT トークン検証ミドルウェア
func AuthMiddleware(userDriver user.UserDriver, auth0Config Auth0Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authorization ヘッダーからトークンを抽出
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Bearer トークンフォーマットのチェック
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// JWT トークンを解析してクレームを取得
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// アルゴリズムの検証
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// Auth0の公開鍵を取得
			kid, ok := token.Header["kid"].(string)
			if !ok {
				return nil, errors.New("kid header is required")
			}

			publicKey, err := getAuth0PublicKey(auth0Config.Domain, kid)
			if err != nil {
				return nil, err
			}

			return publicKey, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid token: " + err.Error(),
			})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Token is not valid",
			})
			c.Abort()
			return
		}

		// クレームの取得
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid token claims",
			})
			c.Abort()
			return
		}

		// トークンの有効期限チェック
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": "Token has expired",
				})
				c.Abort()
				return
			}
		}

		// Audience の検証
		if aud, ok := claims["aud"].([]interface{}); ok {
			audienceValid := false
			for _, audience := range aud {
				if audStr, ok := audience.(string); ok && audStr == auth0Config.Audience {
					audienceValid = true
					break
				}
			}
			if !audienceValid {
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": "Invalid audience",
				})
				c.Abort()
				return
			}
		} else if aud, ok := claims["aud"].(string); ok {
			if aud != auth0Config.Audience {
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": "Invalid audience",
				})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Audience claim is required",
			})
			c.Abort()
			return
		}

		// Auth0 sub の取得
		auth0Sub, ok := claims["sub"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Subject claim is required",
			})
			c.Abort()
			return
		}

		// データベースからユーザーを取得
		user, err := userDriver.GetUserByAuth0Sub(auth0Sub)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusForbidden, gin.H{
					"message": "User not found",
				})
				c.Abort()
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to get user: " + err.Error(),
			})
			c.Abort()
			return
		}

		// コンテキストにユーザー情報を設定
		c.Set("userID", user.UserID)
		c.Set("auth0Sub", auth0Sub)

		c.Next()
	}
}

// getAuth0PublicKey Auth0から公開鍵を取得
func getAuth0PublicKey(domain, kid string) (*rsa.PublicKey, error) {
	jwksURL := fmt.Sprintf("https://%s/.well-known/jwks.json", domain)

	resp, err := http.Get(jwksURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var jwks JWKSResponse
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, err
	}

	for _, key := range jwks.Keys {
		if key.Kid == kid {
			return convertJWKToRSAPublicKey(key)
		}
	}

	return nil, errors.New("unable to find appropriate key")
}

// convertJWKToRSAPublicKey JWKをRSA公開鍵に変換
func convertJWKToRSAPublicKey(jwk JWK) (*rsa.PublicKey, error) {
	// n (modulus) をデコード
	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, err
	}

	// e (exponent) をデコード
	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, err
	}

	// big.Int に変換
	n := big.NewInt(0).SetBytes(nBytes)

	// exponent を int に変換
	var e int
	for i, b := range eBytes {
		e += int(b) << (8 * (len(eBytes) - 1 - i))
	}

	return &rsa.PublicKey{
		N: n,
		E: e,
	}, nil
}

// NewAuth0Config 環境変数からAuth0設定を作成
func NewAuth0Config() Auth0Config {
	return Auth0Config{
		Domain:   os.Getenv("AUTH0_DOMAIN"),
		Audience: os.Getenv("AUTH0_AUDIENCE"),
	}
}
