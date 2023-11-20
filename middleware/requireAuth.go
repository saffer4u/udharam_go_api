package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/saffer4u/udharam/v2/initializers"
	"github.com/saffer4u/udharam/v2/models"
)

func RequireAuth(c *gin.Context) {
	// Get the cookie off req
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		models.CreateResponse(false, http.StatusUnauthorized, "Please Login", nil)
		c.AbortWithStatusJSON(models.Response.Code, models.Response)
		return
	}

	// Decode/validate it
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECERET")), nil
	})
	if err != nil {
		models.CreateResponse(false, http.StatusInternalServerError, "Unbale to parse token", err)
		c.AbortWithStatusJSON(models.Response.Code, models.Response)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		// Check the exp
		if float64(time.Now().Unix()) > claims["expTime"].(float64) {
			models.CreateResponse(false, http.StatusUnauthorized, "Session timeout", nil)
			c.AbortWithStatusJSON(models.Response.Code, models.Response)
			return
		}

		// find user with token sub
		var user models.User
		initializers.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			models.CreateResponse(false, http.StatusUnauthorized, "User not found", nil)
			c.AbortWithStatusJSON(models.Response.Code, models.Response)
			return

		}

		// attach to request
		c.Set("user", user)

		// Continue
		c.Next()

	} else {
		models.CreateResponse(false, http.StatusUnauthorized, "Token not vailid", nil)
		c.JSON(models.Response.Code, models.Response)
	}

}
