package config

import (
	"fmt"
	"os"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Mengambil data dari Variables Railway
	// Pastikan nama variabel di Railway (DB_USER, DB_PASS, dll) sudah diisi
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Pengecekan agar aplikasi memberi tahu jika ada variabel yang lupa diisi di Railway
	if dbHost == "" || dbUser == "" {
		panic("Error: Variabel database belum lengkap di Settings -> Variables Railway!")
	}

	// DSN untuk membuat database jika belum ada
	serverDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort)
	serverDB, err := gorm.Open(mysql.Open(serverDSN), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Gagal menyambung ke server database: %v", err))
	}

	// Membuat database jika belum ada
	if err := serverDB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", sanitizeIdentifier(dbName))).Error; err != nil {
		panic(fmt.Sprintf("Gagal membuat database: %v", err))
	}

	// DSN untuk koneksi utama ke database yang sudah dibuat/ada
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(fmt.Sprintf("Gagal menyambung ke database: %v", err))
	}

	DB = database
	fmt.Println("Database terhubung dengan sukses!")
}

func sanitizeIdentifier(value string) string {
	var builder strings.Builder
	for _, char := range value {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '_' {
			builder.WriteRune(char)
		}
	}

	if builder.Len() == 0 {
		return "risa_resa_cookies"
	}

	return builder.String()
}
