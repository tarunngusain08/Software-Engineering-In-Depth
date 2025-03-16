package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	// MinIO server details
	endpoint := "localhost:9000" // Ensure this is the correct API port
	accessKeyID := "minioadmin"
	secretAccessKey := "minioadmin"
	bucketName := "inferenceservice"
	objectName := "testfile.txt"
	filePath := "testfile.txt"
	contentType := "text/plain"

	// Initialize MinIO client
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}

	// Check if the MinIO server is reachable and bucket exists
	exists, err := minioClient.BucketExists(context.Background(), bucketName)
	if err != nil {
		log.Fatalf("Failed to connect to MinIO server: %v", err)
	}

	// Create the bucket if it does not exist
	if !exists {
		err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalf("Failed to create bucket: %v", err)
		}
		fmt.Printf("Successfully created bucket %s\n", bucketName)
	}

	// Create a sample file for upload
	err = os.WriteFile(filePath, []byte("Hello, this is a test file!"), 0644)
	if err != nil {
		log.Fatalf("Failed to create test file: %v", err)
	}

	// Upload the file
	ctx := context.Background()
	uploadInfo, err := minioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Successfully uploaded %s (size: %d bytes)\n", uploadInfo.Key, uploadInfo.Size)
}
