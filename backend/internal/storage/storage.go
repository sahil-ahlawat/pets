package storage

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"pet-platform/backend/internal/config"
)

type StorageService interface {
	GeneratePresignedUploadURL(filename, mimeType, category string) (*PresignedUpload, error)
	GetPublicURL(storageKey string) string
}

type PresignedUpload struct {
	UploadURL  string            `json:"uploadUrl"`
	FileURL    string            `json:"fileUrl"`
	StorageKey string            `json:"storageKey"`
	Headers    map[string]string `json:"headers"`
}

type S3Storage struct {
	cfg *config.Config
}

func NewStorageService(cfg *config.Config) StorageService {
	return &S3Storage{cfg: cfg}
}

func (s *S3Storage) GeneratePresignedUploadURL(filename, mimeType, category string) (*PresignedUpload, error) {
	ext := filepath.Ext(filename)
	storageKey := fmt.Sprintf("%s/%s%s", category, uuid.New().String(), ext)
	expires := time.Now().Add(15 * time.Minute).Unix()

	// Direct S3 / Cloudflare R2 presigned URL signer calculation
	rawStringToSign := fmt.Sprintf("PUT\n\n%s\n%d\n/%s/%s", mimeType, expires, s.cfg.S3Bucket, storageKey)
	mac := hmac.New(sha256.New, []byte(s.cfg.S3SecretKey))
	mac.Write([]byte(rawStringToSign))
	signature := hex.EncodeToString(mac.Sum(nil))

	uploadURL := fmt.Sprintf("%s/%s/%s?AWSAccessKeyId=%s&Signature=%s&Expires=%d",
		s.cfg.S3Endpoint,
		s.cfg.S3Bucket,
		storageKey,
		url.QueryEscape(s.cfg.S3AccessKey),
		url.QueryEscape(signature),
		expires,
	)

	fileURL := s.GetPublicURL(storageKey)

	return &PresignedUpload{
		UploadURL:  uploadURL,
		FileURL:    fileURL,
		StorageKey: storageKey,
		Headers: map[string]string{
			"Content-Type": mimeType,
		},
	}, nil
}

func (s *S3Storage) GetPublicURL(storageKey string) string {
	return fmt.Sprintf("%s/%s/%s", s.cfg.S3Endpoint, s.cfg.S3Bucket, storageKey)
}
