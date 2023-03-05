package main

import (
	"crypto/rsa"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/jwk"
	"runtime/debug"
	"time"
)

// Commit 代表 git commit hash
var Commit = func() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				return setting.Value
			}
		}
	}

	return ""
}()

// GenerateLineBotJWTToken 產生取得 lineBot channel access token 所需要的 JWT
/*
private key example
{
  "alg": "RS256",
  "d": "GaDzOmc4......",
  "dp": "WAByrYmh......",
  "dq": "WLwjYun0......",
  "e": "AQ......",
  "ext": true,
  "key_ops": [
    "sign"
  ],
  "kty": "RSA",
  "n": "vsbOUoFA......",
  "p": "5QJitCu9......",
  "q": "1ULfGui5......",
  "qi": "2cK4apee......"
}

public key example
{
  "alg": "RS256",
  "e": "AQ......",
  "ext": true,
  "key_ops": [
    "verify"
  ],
  "kty": "RSA",
  "n": "vsbOUoFA......"
}
*/
func GenerateLineBotJWTToken(privateKey string, channelId string, kid string) (string, error) {
	signKey := &rsa.PrivateKey{}
	err := jwk.ParseRawKey([]byte(privateKey), signKey)
	if err != nil {
		return "", err
	}
	myClaim := make(jwt.MapClaims)
	myClaim["iss"] = channelId
	myClaim["sub"] = channelId
	myClaim["aud"] = "https://api.line.me/" // TODO 看情況更換
	myClaim["exp"] = time.Now().Unix() + (30 * 60)
	myClaim["token_exp"] = 30 * 24 * 60 * 60
	tk := jwt.NewWithClaims(jwt.SigningMethodRS256, myClaim)
	tk.Header["kid"] = kid
	jwt, err := tk.SignedString(signKey)
	if err != nil {
		return "", err
	}
	return jwt, nil
}
