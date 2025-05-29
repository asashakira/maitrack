package utils

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func UploadImageToS3(imageURL string) error {
	ctx := context.Background()

	// Load AWS config
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}
	s3Client := s3.NewFromConfig(cfg)

	// Download the image
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ⚠️
		},
	}
	resp, err := httpClient.Get(imageURL)
	if err != nil {
		return fmt.Errorf("failed to fetch image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non-200 response: %d", resp.StatusCode)
	}

	// Read the body into a buffer
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read image data: %w", err)
	}
	bodyReader := bytes.NewReader(bodyBytes)

	// Extract file name from URL
	filename := filepath.Base(imageURL)

	// Define bucket and key
	bucket := "assets.maitrack.com"
	key := fmt.Sprintf("songs/%s", filename)

	// Get content type
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = detectContentTypeFromExt(filename)
	}

	// Upload to S3
	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        &bucket,
		Key:           &key,
		Body:          bodyReader,
		ContentType:   &contentType,
	})
	if err != nil {
		return fmt.Errorf("failed to upload to S3: %w", err)
	}

	fmt.Printf("Uploaded %s to s3://%s/%s\n", filename, bucket, key)
	return nil
}

func detectContentTypeFromExt(filename string) string {
	switch strings.ToLower(filepath.Ext(filename)) {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}
