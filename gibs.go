package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gocql/gocql"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Gibs struct {
	UserUuid     gocql.UUID
	ResponseCode int
	ares         ApiResponse
}

type JwtClaims struct {
	UserUuid string `json:"userUuid"`
	jwt.StandardClaims
}

type JwtToken struct {
	Token string `json:"token"`
}

func (gibs *Gibs) encodeResponse(w http.ResponseWriter) {
	// TODO: we need to decide at which point content-type should be defaulted
	// to json, if it's defaulted to JSON at all
	//w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(gibs.ResponseCode)
	if err := json.NewEncoder(w).Encode(gibs.ares); err != nil {
		panic(errors.Wrap(err, "failed json encode"))
	}
}

func createJwtOnSignIn(user User) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userUuid": user.UserUuid,
		// TODO: should there be more claims here?
	})
	// TODO: make secret code
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		panic(errors.Wrap(err, "failed to sign token"))
	}
	return tokenString
}

func (gibs *Gibs) checkAndGetJwtClaims(authHeader string, access int) error {
	if authHeader == "" {
		if access != P_PUBLIC {
			return errors.New("Not authorized.")
		}
		return nil
	}

	token, err := jwt.ParseWithClaims(authHeader, &JwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

	if err != nil {
		return err
		//msg := err.Error()
	}

	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		//fmt.Printf("%v %v", claims.Foo, claims.StandardClaims.ExpiresAt)
		msg := "Token is valid with userUuid " + claims.UserUuid
		gibs.UserUuid, err = gocql.ParseUUID(claims.UserUuid)
		if err != nil {
			return errors.Wrap(err, "failure decoding gibs userUuid")
		}
		fmt.Println(msg)
	} else {
		return errors.New("Token is invalid. Please try relogging.")
	}

	return nil
}

func unpackReqIntoGibs(w http.ResponseWriter, r *http.Request, access int) Gibs {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	var gibs Gibs
	gibs.ares.Body = make(map[string]interface{})
	log.Info("unpacking req")
	gibs.ResponseCode = 200

	authHeader := r.Header.Get("Authorization")
	err := gibs.checkAndGetJwtClaims(authHeader, access)
	if err != nil {
		gibs.ares.Info = ERROR
		gibs.ares.appendFlash(err.Error(), "error")
		gibs.encodeResponse(w)
		// TODO: make sure request is killed when this panic occurs
		panic(errors.Wrap(err, "failed to authenticate gibs"))
	}

	switch access {
	case P_PUBLIC:
	case P_USER:
		if gibs.UserUuid.String() == EmptyUuidString {
			//panic(errors.New("Must be logged in to view this page"))
			log.Info("User must be logged in for this route")
			gibs.ares.Info = REJECTED
			gibs.ares.appendFlash("Must be logged in", "error")
		}
	case P_PUBLIC_MEM_SESSION:
	default:
		//err = errors.New("Access value not found")
	}

	log.Info("done unpacking")
	//gibs.ares.printFlashMsgs()

	return gibs
}

func (gibs *Gibs) sendAuthError(w http.ResponseWriter) {
	gibs.ares.Info = REJECTED
	gibs.ares.appendFlash("Authorization failure.", "error")
	log.Info("unauthorized attempt by gibs.UserUuid: " + gibs.UserUuid.String())
	gibs.encodeResponse(w)
}
