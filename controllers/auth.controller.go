package controllers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	"gorm.io/gorm"

	"github.com/mas-wig/ta-v1.0.4/initializers"
	"github.com/mas-wig/ta-v1.0.4/models"
	"github.com/mas-wig/ta-v1.0.4/utils"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(DB *gorm.DB) AuthController {
	return AuthController{DB}
}

// Register
func (ac *AuthController) SignUpUser(ctx *gin.Context) {
	var payload *models.SignUpInput

	if err := ctx.ShouldBind(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail 1", "message": err.Error()})
		return
	}

	if payload.Password != payload.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail 2", "message": "Passwords do not match"})
		return
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error 1", "message": err.Error()})
		return
	}

	var photoURL string
	if payload.Photo != nil {
		url, err := utils.SaveUploadedFile(payload.Photo)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error 2", "message": err.Error()})
			return
		}
		photoURL = url
	}
	now := time.Now()

	newUser := models.User{
		Email:     strings.ToLower(payload.Email),
		Username:  strings.ToLower(payload.Username),
		Password:  hashedPassword,
		FullName:  payload.FullName,
		Gender:    payload.Gender,
		Address:   payload.Address,
		Verified:  false,
		Photo:     photoURL,
		Role:      "user",
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := ac.DB.Create(&newUser)

	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail 3", "message": "User with that email already exists"})
		return
	} else if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error 3", "message": "Something bad happened"})
		return
	}

	config, _ := initializers.LoadConfig(".")

	// Generate Verification Code
	code := randstr.String(20)

	verificationCode := utils.Encode(code)

	// Update User in Database
	newUser.VerificationCode = verificationCode
	ac.DB.Save(newUser)

	var firstName = newUser.Username

	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	emailData := utils.EmailData{
		URL:       config.ClientOrigin + "/api/auth/verifyemail/" + code,
		FirstName: firstName,
		Subject:   "Your account verification code",
	}

	utils.SendEmail(&newUser, &emailData)

	message := "We sent an email with a verification code to " + newUser.Email
	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": message})
}

// Login
func (ac *AuthController) SignInUser(ctx *gin.Context) {
	var payload *models.SignInInput
	if err := ctx.ShouldBind(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	var user models.User

	if err := ac.DB.Where("email = ? AND role = ?", payload.Email, payload.Role).First(&user).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	if !user.Verified {
		ctx.JSON(http.StatusForbidden, gin.H{"status": "fail", "message": "Please verify your email"})
		return
	}

	if err := utils.VerifyPassword(user.Password, payload.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	config, _ := initializers.LoadConfig(".")

	// Generate Token
	token, err := utils.GenerateToken(config.TokenExpiresIn, user.ID, config.TokenSecret)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	switch payload.Role {
	case "user":
		ctx.SetCookie("access_token", token, config.TokenMaxAge*60, "/", "localhost", false, true)
		ctx.Redirect(http.StatusFound, "/users/dashboard")
	case "admin":
		ctx.SetCookie("access_token", token, config.TokenMaxAge*60, "/", "localhost", false, true)
		ctx.Redirect(http.StatusFound, "/api/users/me")
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "role tidak ada"})
	}
}

// Logout
func (ac *AuthController) LogoutUser(ctx *gin.Context) {
	ctx.SetCookie("token", "", -1, "/", "localhost", false, true)
	ctx.Redirect(http.StatusFound, "/api/auth/login")
}

// Verify Email
func (ac *AuthController) VerifyEmail(ctx *gin.Context) {

	code := ctx.Params.ByName("verificationCode")
	verificationCode := utils.Encode(code)

	var updatedUser models.User
	result := ac.DB.First(&updatedUser, "verification_code = ?", verificationCode)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid verification code or user doesn't exists"})
		return
	}

	if updatedUser.Verified {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User already verified"})
		return
	}

	updatedUser.VerificationCode = ""
	updatedUser.Verified = true
	ac.DB.Save(&updatedUser)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Email verified successfully"})
}
