package utils

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func SaveUploadedFile(file *multipart.FileHeader) (string, error) {
	ext := filepath.Ext(file.Filename)
	filename := uuid.New().String() + ext
	dst := filepath.Join("public", "static", "img", filename)
	// Membuat folder "static/img" jika belum ada
	err := os.MkdirAll(filepath.Dir(dst), os.ModePerm)
	if err != nil {
		return "", err
	}

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
	newDst := strings.ReplaceAll(dst, "public", "")
	return filepath.Join("/", newDst), nil
}
