package web

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type JWT struct {
	header    string
	Payload   []byte
	signature string
}

type JWTDefaultParams struct {
	Iss string `json:"iss,omitempty"`
	Exp int64  `json:"exp,omitempty"`
	Sub string `json:"sub,omitempty"`
	Aud string `json:"aud,omitempty"`
	Nbf int64  `json:"nbf,omitempty"`
	Iat int64  `json:"iat,omitempty"`
	Jti string `json:"jti,omitempty"`
}

func encodeBase64(data string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(data))
}

func generateSignature(key []byte, data []byte) (string, error) {
	hash := hmac.New(sha256.New, key)
	_, err := hash.Write(data)
	if err != nil {
		return "", err
	}
	return encodeBase64(string(hash.Sum(nil))), nil
}

func CreateToken(key []byte, payloadData any) (string, error) {
	header := `{"alg":"HS256","typ":"JWT"}`
	payload, jsonErr := json.Marshal(payloadData)
	if jsonErr != nil {
		return "", fmt.Errorf("load JSON parsing error")
	}
	encodedHeader := encodeBase64(header)
	encodedPayload := encodeBase64(string(payload))
	HeaderAndPayload := encodedHeader + "." + encodedPayload
	signature, err := generateSignature(key, []byte(HeaderAndPayload))
	if err != nil {
		return "", err
	}
	return HeaderAndPayload + "." + signature, nil
}

func parseJwt(token string, key []byte) (*JWT, error) {
	jwtParts := strings.Split(token, ".")
	if len(jwtParts) != 3 {
		return nil, fmt.Errorf("illegal token")
	}
	encodedHeader := jwtParts[0]
	encodedPayload := jwtParts[1]
	signature := jwtParts[2]

	confirmSignature, err := generateSignature(key, []byte(encodedHeader+"."+encodedPayload))
	if err != nil {
		return nil, fmt.Errorf("signature generation error")
	}
	if signature != confirmSignature {
		return nil, fmt.Errorf("token verification failed")
	}
	dstPayload, _ := base64.RawURLEncoding.DecodeString(encodedPayload)
	return &JWT{encodedHeader, dstPayload, signature}, nil
}

func JwtConfirm(key []byte, headerKey string, obj any) HandlerFunc {
	isMap := false
	if reflect.TypeOf(obj).Kind() == reflect.Map {
		isMap = true
	}
	hasExp := false
	if isMap {
		mp := obj.(map[string]any)
		if _, has := mp["exp"]; has {
			hasExp = true
		}
	} else {
		metaObj := reflect.ValueOf(obj).Elem()
		if metaObj.FieldByName("Exp") != (reflect.Value{}) {
			hasExp = true
		}
	}
	return func(c *Context) {
		token := c.GetHeader(headerKey)
		jwt, err := parseJwt(token, key)
		if err != nil {
			c.Fail(http.StatusUnauthorized, err.Error())
			return
		}
		err = json.Unmarshal(jwt.Payload, &obj)
		if err != nil {
			c.Fail(http.StatusInternalServerError, err.Error())
			return
		}
		if isMap {
			if hasExp && int64(obj.(map[string]any)["exp"].(float64)) < time.Now().Unix() {
				c.Fail(http.StatusUnauthorized, "session expiration")
				return
			}
		} else {
			if hasExp && reflect.ValueOf(obj).Elem().FieldByName("Exp").Int() < time.Now().Unix() {
				c.Fail(http.StatusUnauthorized, "session expiration")
				return
			}
		}
		c.SetExtra("jwt", obj)
		c.Next()
	}
}
