package utils

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

func SaveUploadedFile(file *multipart.FileHeader) (string, error) {
	ext := filepath.Ext(file.Filename)
	// Membangun path tujuan file
	filename := uuid.New().String() + ext
	dst := filepath.Join("public", "assets", "uploads", "img", filename)
	// Membuat folder "static/img" jika belum ada
	err := os.MkdirAll(filepath.Dir(dst), os.ModePerm)
	if err != nil {
		return "", err
	}

	// Membuka file yang diunggah
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Membuat file baru di direktori tujuan
	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Menyalin data dari file yang diunggah ke file baru
	_, err = io.Copy(out, src)
	if err != nil {
		return "", err
	}

	// Mengembalikan URL lengkap file yang tersimpan
	return filepath.Join("/", dst), nil
}
