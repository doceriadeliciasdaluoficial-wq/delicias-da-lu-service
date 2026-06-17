package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"time"

	"cloud.google.com/go/storage"
)

const (
	bucketName = "delicias-da-lu-images"
	// URL format for signed URLs
	signedURLExpiry = time.Hour
)

type StorageRepository interface {
	UploadFile(ctx context.Context, objectPath string, file *multipart.FileHeader) (string, error)
	UploadBytes(ctx context.Context, objectPath string, data []byte, contentType string) (string, error)
	GetSignedURL(ctx context.Context, objectPath string) (string, error)
	DeleteFile(ctx context.Context, objectPath string) error
	FileExists(ctx context.Context, objectPath string) (bool, error)
}

type storageRepositoryImpl struct {
	client *storage.Client
}

func NewStorageRepository(client *storage.Client) StorageRepository {
	return &storageRepositoryImpl{client: client}
}

// UploadFile uploads a multipart file to Cloud Storage
func (r *storageRepositoryImpl) UploadFile(ctx context.Context, objectPath string, file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Read file into memory
	data, err := io.ReadAll(src)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return r.UploadBytes(ctx, objectPath, data, file.Header.Get("Content-Type"))
}

// UploadBytes uploads raw bytes to Cloud Storage
func (r *storageRepositoryImpl) UploadBytes(ctx context.Context, objectPath string, data []byte, contentType string) (string, error) {
	wctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	obj := r.client.Bucket(bucketName).Object(objectPath)
	w := obj.NewWriter(wctx)

	w.ContentType = contentType
	w.Metadata = map[string]string{
		"uploaded": time.Now().Format(time.RFC3339),
	}

	if _, err := io.Copy(w, bytes.NewReader(data)); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	if err := w.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	// Make object public (if bucket is configured for public read)
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(signedURLExpiry),
	}

	url, err := r.client.Bucket(bucketName).SignedURL(objectPath, opts)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}

	return url, nil
}

// GetSignedURL generates a signed URL for accessing a file
func (r *storageRepositoryImpl) GetSignedURL(ctx context.Context, objectPath string) (string, error) {
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(signedURLExpiry),
	}

	url, err := r.client.Bucket(bucketName).SignedURL(objectPath, opts)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}

	return url, nil
}

// DeleteFile deletes a file from Cloud Storage
func (r *storageRepositoryImpl) DeleteFile(ctx context.Context, objectPath string) error {
	if err := r.client.Bucket(bucketName).Object(objectPath).Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// FileExists checks if a file exists in Cloud Storage
func (r *storageRepositoryImpl) FileExists(ctx context.Context, objectPath string) (bool, error) {
	_, err := r.client.Bucket(bucketName).Object(objectPath).Attrs(ctx)
	if err == storage.ErrObjectNotExist {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
