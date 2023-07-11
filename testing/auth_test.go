package testing

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/mas-wig/ta-v1.0.4/controllers"
	"github.com/mas-wig/ta-v1.0.4/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestSignInUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	assert.NoError(t, err)

	gormDB = gormDB.Session(&gorm.Session{DryRun: true})
	authController := controllers.NewAuthController(gormDB)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	payload := &models.SignInInput{
		Email:    "wigtono24@gmail.com",
		Password: "12345678",
		Role:     "user",
	}

	c.Request, _ = http.NewRequest(http.MethodPost, "", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = ioutil.NopCloser(bytes.NewBufferString(
		fmt.Sprintf(`{"email":"%s","password":"%s","role":"%s"}`, payload.Email, payload.Password, payload.Role)))

	authController.SignInUser(c)

	assert.Equal(t, 403, w.Code)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
