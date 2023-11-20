package controllers

import (
	"net/http"
	"net/mail"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"strings"

	"github.com/saffer4u/udharam/v2/initializers"
	"github.com/saffer4u/udharam/v2/models"
	"golang.org/x/crypto/bcrypt"
)

type User struct { // Not the model, see this as serializer
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func CreateUserResponse(userModel models.User) User {
	return User{
		ID:        userModel.ID,
		Email:     strings.ToLower(userModel.Email),
		FirstName: userModel.FirstName,
		LastName:  userModel.LastName,
	}
}

func Signup(c *gin.Context) {

	// Get the email/pass off req body
	var body struct {
		Email     string
		Password  string
		FirstName string
		LastName  string
	}

	if err := c.Bind(&body); err != nil {

		models.CreateResponse(false, http.StatusBadRequest, "Invailid request body", nil)
		c.JSON(models.Response.Code, models.Response)
		return
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		models.CreateResponse(false, http.StatusBadRequest, "Failed to encrypt password", nil)
		c.JSON(models.Response.Code, models.Response)
		return
	}

	// Create user
	user := models.User{Email: strings.ToLower(body.Email), Password: string(hash), FirstName: body.FirstName, LastName: body.LastName}
	_, errMail := mail.ParseAddress(user.Email)

	if errMail != nil {
		models.CreateResponse(false, http.StatusBadRequest, "Invailid email", nil)
		c.JSON(models.Response.Code, models.Response)
		return
	}

	result := initializers.DB.Create(&user)

	if result.Error != nil {
		// pqErr, ok := result.Error.(*pq.Error)
		// if ok && pqErr.Code == "23505" {
		// 	models.CreateResponse(false, http.StatusBadRequest, "Email Already exist please login", nil)
		// 	c.JSON(models.Response.Code, models.Response)
		// 	return
		// }
		models.CreateResponse(false, http.StatusBadRequest, "User creation faild. "+"{"+result.Error.Error()+"}", nil)
		c.JSON(models.Response.Code, models.Response)
		return
	}

	// Respond
	models.CreateResponse(true, http.StatusOK, "Account created successfully", CreateUserResponse(user))
	c.JSON(models.Response.Code, models.Response)
}

func Login(c *gin.Context) {
	// Get the email/passs off req body
	var body struct {
		Email    string
		Password string
	}

	if err := c.Bind(&body); err != nil {

		models.CreateResponse(false, http.StatusBadRequest, "Invalid request body", nil)
		c.JSON(models.Response.Code, models.Response)
		return
	}

	// Look up requested user
	var user models.User
	ref := initializers.DB.First(&user, "email = ?", body.Email)

	if ref.Error != nil {

		models.CreateResponse(false, http.StatusNotFound, "User not found please signup if you don't have account", nil)
		c.JSON(models.Response.Code, models.Response)
		return
	}

	// Compare sent in pass with saved user hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		models.CreateResponse(false, http.StatusBadRequest, "Invalid password", nil)
		c.JSON(models.Response.Code, models.Response)
		return
	}

	// Generate a jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":     user.ID,
		"expTime": time.Now().Add(time.Minute * 5).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECERET")))

	if err != nil {
		models.CreateResponse(false, http.StatusInternalServerError, "Faild to generate token", nil)
		c.JSON(models.Response.Code, models.Response)
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24, "", "", false, true)

	// Send it back
	models.CreateResponse(true, http.StatusOK, "Login successful", CreateUserResponse(user))
	c.JSON(models.Response.Code, models.Response)
}

func Logout(c *gin.Context) {
	// tokenString, err := c.Cookie("Authorization")
	// if err != nil {
	// 	models.CreateResponse(false, http.StatusUnauthorized, "Token not found", nil)
	// 	c.AbortWithStatusJSON(models.Response.Code, models.Response)
	// 	return
	// }
	// token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	// 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
	// 		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	// 	}

	// 	// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
	// 	return []byte(os.Getenv("SECERET")), nil
	// })

}
