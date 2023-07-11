package controllers

import (
	"encoding/base64"
	"log"
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
		ctx.HTML(http.StatusBadRequest, "400.html", nil)
		return
	}

	if payload.Password != payload.PasswordConfirm {
		ctx.HTML(http.StatusBadRequest, "400.html", nil)
		return
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		ctx.HTML(http.StatusBadGateway, "500.html", nil)
		return
	}

	var photoURL string
	if payload.Photo != nil {
		url, err := utils.SaveUploadedFile(payload.Photo)
		if err != nil {
			ctx.HTML(http.StatusBadGateway, "500.html", nil)
			return
		}
		photoURL = url
	}
	now := time.Now()

	var adminCount int64
	ac.DB.Model(&models.User{}).Where("role = ?", "admin").Count(&adminCount)

	var roleUser = "admin"
	if adminCount >= 1 {
		roleUser = "user"
	}

	newUser := models.User{
		Email:              strings.ToLower(payload.Email),
		Username:           strings.ToLower(payload.Username),
		Password:           hashedPassword,
		FullName:           payload.FullName,
		Gender:             payload.Gender,
		Address:            payload.Address,
		Verified:           false,
		Photo:              photoURL,
		Role:               roleUser,
		PasswordResetToken: "",
		PasswordResetAt:    now,
		Acc:                false,
		CreatedAt:          now,
		UpdatedAt:          now,
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
	code := randstr.String(20)
	verificationCode := utils.Encode(code)
	newUser.VerificationCode = verificationCode
	var firstName = newUser.Username

	ac.DB.Save(newUser)

	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	emailData := utils.EmailData{
		URL:       config.ClientOrigin + "/api/auth/verifyemail/" + code,
		FirstName: firstName,
		Subject:   "Link verifikasi akun anda..",
	}
	utils.SendEmail(&newUser, &emailData)

	ctx.Redirect(http.StatusFound, "/login")
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
		ctx.Redirect(http.StatusFound, "/admin/dashboard")
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "role tidak ada"})
	}
}

// Logout
func (ac *AuthController) LogoutUser(ctx *gin.Context) {
	ctx.SetCookie("token", "", -1, "/", "localhost", false, true)
	ctx.Redirect(http.StatusFound, "/login")
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

	ctx.HTML(http.StatusOK, "verifikasi.html", nil)
}

func (ac *AuthController) ForgotPassword(ctx *gin.Context) {
	var payload *models.ForgotPasswordInput

	if err := ctx.ShouldBind(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var user models.User
	result := ac.DB.First(&user, "email = ?", strings.ToLower(payload.Email))

	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	if !user.Verified {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Account not verified"})
		return
	}

	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load config", err)
	}

	resetToken := randstr.String(20)

	passwordResetToken := utils.Encode(resetToken)
	user.PasswordResetToken = passwordResetToken
	user.PasswordResetAt = time.Now().Add(time.Minute * 60)

	ac.DB.Save(&user)

	var firstName = user.FullName
	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	emailData := utils.EmailData{
		URL:       config.ClientOrigin + "/api/auth/resettoken/" + resetToken,
		FirstName: firstName,
		Subject:   "Password reset token anda (valid untuk 10min)",
	}
	utils.SendEmail(&user, &emailData)
	ctx.Redirect(http.StatusFound, "/login")
}

func (ac *AuthController) GetResetPasswordToken(ctx *gin.Context) {
	resetToken := ctx.Params.ByName("newresetToken")
	var updatedUser models.User
	result := ac.DB.First(&updatedUser, "password_reset_token=?", base64.StdEncoding.EncodeToString([]byte(resetToken)))
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid verification code or user doesn't exists"})
		return
	}
	ctx.HTML(http.StatusFound, "reset-password.html", gin.H{"ResetToken": resetToken})
}

func (ac *AuthController) ResetPassword(ctx *gin.Context) {
	var payload *models.ResetPasswordInput
	resetToken := ctx.Params.ByName("resetToken")

	if err := ctx.ShouldBind(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if payload.Password != payload.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords do not match"})
		return
	}

	hashedPassword, _ := utils.HashPassword(payload.Password)

	passwordResetToken := utils.Encode(resetToken)

	var updatedUser models.User
	result := ac.DB.First(&updatedUser, "password_reset_token = ? AND password_reset_at > ?", passwordResetToken, time.Now())
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "The reset token is invalid or has expired"})
		return
	}

	updatedUser.Password = hashedPassword
	updatedUser.PasswordResetToken = ""
	ac.DB.Save(&updatedUser)

	time.Sleep(3 * time.Second)
	ctx.Redirect(http.StatusFound, "/login")
	ctx.SetCookie("token", "", -1, "/login", "localhost", false, true)
}

// TODO: Keknya susah kalo tak implentasi sekerang fitur ginian

// func (ac *AuthController) SignUpWithGoogle() {
//
// }
