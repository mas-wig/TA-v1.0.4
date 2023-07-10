package testing

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mas-wig/ta-v1.0.4/controllers"
	"github.com/mas-wig/ta-v1.0.4/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestEncode(t *testing.T) {
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
	absensiController := controllers.NewAbsensiController(gormDB)

	w := httptest.NewRecorder()
	r := gin.Default()
	r.POST("/api/absen/create", absensiController.CreateAbsensi) // Menggunakan /api/absen/create sebagai endpoint

	payload := &models.CreatePresensi{
		Hari:           "Senin",
		Lokasi:         "Kelas A",
		Kehadiran:      "Hadir",
		TanggalWaktu:   "2023-06-30 08:00:00",
		InformasiMedis: "Sehat",
	}

	form := url.Values{}
	form.Add("hari", payload.Hari)
	form.Add("lokasi", payload.Lokasi)
	form.Add("kehadiran", payload.Kehadiran)
	form.Add("date", payload.TanggalWaktu)
	form.Add("catatankesehatan", payload.InformasiMedis)

	req, err := http.NewRequest(http.MethodPost, "/api/absen/create", strings.NewReader(form.Encode())) // Menggunakan /api/absen/create sebagai endpoint
	assert.NoError(t, err)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Request.PostForm = form

	now := time.Now()
	currentUser := models.User{
		ID:        "ea3d0806-287b-49fc-84ad-b538ac452d40",
		Email:     "wigtono24@gmail.com",
		FullName:  "wigtono",
		Password:  "$2a$10$6u.BsHRzkAkrtGz1o69O6..RDnPenCDBf0xZ2I09YN7FZE6NPylCK",
		Username:  "wigtono",
		Gender:    "pria",
		Address:   "efdt",
		Role:      "user",
		Photo:     "/static/img/7f25281d-dd08-4461-bfc8-e30c0b2c338b.png",
		Acc:       false,
		CreatedAt: now,
		UpdatedAt: now,
	}
	c.Set("currentUser", currentUser)

	r.ServeHTTP(w, req)

	assert.Equal(t, 500, w.Code)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDecodeByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to initialize database mock: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	router := gin.Default()
	absensiController := controllers.NewAbsensiController(gormDB)

	absenID := "ea3d0806-287b-49fc-84ad-b538ac452d40"
	currentUser := models.User{
		ID:       absenID,
		FullName: "John Doe",
	}

	encodePresensi := models.EncodePresensi{
		ID:             absenID,
		NamaSiswa:      base64.StdEncoding.EncodeToString([]byte("John Doe")),
		IDSiswa:        uuid.New(),
		Hari:           base64.StdEncoding.EncodeToString([]byte("Monday")),
		Lokasi:         base64.StdEncoding.EncodeToString([]byte("Classroom 101")),
		TanggalWaktu:   base64.StdEncoding.EncodeToString([]byte("2023-06-30 09:00")),
		Kehadiran:      base64.StdEncoding.EncodeToString([]byte("Present")),
		InformasiMedis: base64.StdEncoding.EncodeToString([]byte("No specific medical information")),
	}

	mock.ExpectQuery("SELECT (.+) FROM encode_presensis").WithArgs(absenID).
		WillReturnRows(mock.NewRows([]string{"ID", "NamaSiswa", "IDSiswa", "Hari", "Lokasi", "TanggalWaktu", "Kehadiran", "InformasiMedis"}).
			AddRow(encodePresensi.ID, encodePresensi.NamaSiswa, encodePresensi.IDSiswa, encodePresensi.Hari, encodePresensi.Lokasi, encodePresensi.TanggalWaktu, encodePresensi.Kehadiran, encodePresensi.InformasiMedis))

	// Set up the HTTP request for the successful scenario
	request, _ := http.NewRequest("GET", "/api/absen/decode/"+absenID, nil)
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.Use(func(c *gin.Context) {
		c.Set("currentUser", currentUser)
		c.Next()
	})

	router.ServeHTTP(recorder, request)

	assert.Equal(t, 404, recorder.Code)
	notFoundAbsenID := absenID

	mock.ExpectQuery("SELECT (.+) FROM encode_presensis").WithArgs(notFoundAbsenID).
		WillReturnError(gorm.ErrRecordNotFound)

	request, _ = http.NewRequest("GET", "/api/absen/decode/"+notFoundAbsenID, nil)
	request.Header.Set("Content-Type", "application/json")
	recorder = httptest.NewRecorder()
	router.GET("/api/absen/decode/:absenId", absensiController.DecodeByID)

	router.Use(func(c *gin.Context) {
		c.Set("currentUser", currentUser)
		c.Next()
	})

	router.ServeHTTP(recorder, request)
	assert.Equal(t, 400, recorder.Code)
}
