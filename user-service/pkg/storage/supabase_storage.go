package storage

import (
	"context"
	"fmt"
	_ "io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"micro-warehouse/user-service/configs"

	storage_go "github.com/supabase-community/storage-go"
)

type SupabaseInterface interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader, folder string) (*UploadResult, error)
}

type SupabaseStorage struct {
	client *storage_go.Client
	cfg    configs.Config
}

// UploadFile implements [SupabaseInterface].
func (s *SupabaseStorage) UploadFile(ctx context.Context, file *multipart.FileHeader, folder string) (*UploadResult, error) {
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// generate unique filename
	ext := filepath.Ext((file.Filename))
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%s_%d%s", strings.TrimSuffix(file.Filename, ext), timestamp, ext)

	// create file path
	filepath := fmt.Sprintf("%s/%s", folder, filename)

	// use the simpler implementation with proper Content-Type
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		switch strings.ToLower(ext) {
		case ".jpg", ".jpeg":
			contentType = "image/jpeg"
		case ".png":
			contentType = "image/png"
		case ".webp":
			contentType = "image/webp"
		case ".svg":
			contentType = "image/svg+xml"
		default:
			contentType = "application/octet-stream"
		}
	}

	// create client with proper Content-Type
	storageURL := s.cfg.Supabase.URL
	client := storage_go.NewClient(storageURL, s.cfg.Supabase.Key, map[string]string{
		"Content-Type": contentType,
	})

	// upload file
	_, err = client.UploadFile(s.cfg.Supabase.Bucket, filepath, src)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to supabase: %w", err)
	}

	// get public url
	publicUrl := client.GetPublicUrl(s.cfg.Supabase.Bucket, filepath)
	
	return &UploadResult{
		URL: publicUrl.SignedURL,
		Path: filepath,
		Filename: filename,
	}, nil
}

type UploadResult struct {
	URL      string `json:"url"`
	Path     string `json:"path"`
	Filename string `json:"filename"`
}

// const supabaseStoragePath = "/storage/v1"

func NewSupabaseStorage(cfg configs.Config) SupabaseInterface {
	storageURL := cfg.Supabase.URL
	client := storage_go.NewClient(storageURL, cfg.Supabase.Key, nil)
	return &SupabaseStorage{
		client: client,
		cfg:    cfg,
	}
}
