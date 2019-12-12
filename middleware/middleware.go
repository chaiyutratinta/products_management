package middleware

import (
	"fmt"
	"log"
	"net/http"
	"products_management/models"

	"github.com/dgrijalva/jwt-go"
)

func AuthMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		authorization := req.Header.Get("Authorization")
		token, err := jwt.ParseWithClaims(authorization, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {

			return []byte("a40449d2a209f1ba700c20da616b01a2f360b39f97152aa384e01f54ecab17571c5311e5f83108bc57fc94ddcc2ba12530edc2db5f6a57458c8d330d6317307e"), nil
		})

		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusUnauthorized)

			return
		}

		if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
			next.ServeHTTP(writer, req)
			fmt.Println(claims)

			return
		}

	})
}
