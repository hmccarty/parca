package calendar

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

type jwtResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func buildJWT(serviceEmail string, keyData []byte) (string, error) {
	currTime := time.Now().UTC()
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss":   serviceEmail,
		"scope": "https://www.googleapis.com/auth/calendar",
		"aud":   "https://accounts.google.com/o/oauth2/token",
		"exp":   currTime.Add(time.Hour * time.Duration(1)).Unix(),
		"iat":   currTime.Unix(),
	})

	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		return "", fmt.Errorf("create: parse key: %w", err)
	}

	token, err := unsignedToken.SignedString(key)
	if err != nil {
		return "", err
	}

	return token, nil
}

func requestToken(serviceEmail string, keyData []byte) (string, time.Time, error) {
	token, err := buildJWT(serviceEmail, keyData)
	if err != nil {
		return "", time.Time{}, err
	}

	endpoint := "https://oauth2.googleapis.com/token"
	data := url.Values{}
	data.Set("grant_type", "urn:ietf:params:oauth:grant-type:jwt-bearer")
	data.Set("assertion", token)

	client := &http.Client{}
	r, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatal(err)
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := client.Do(r)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var j jwtResp
	if err := json.Unmarshal(body, &j); err != nil {
		log.Fatal(err)
	}

	expiration := time.Now().UTC().Add(time.Duration(j.ExpiresIn) * time.Second)
	return j.AccessToken, expiration, nil
}
