package cloudstorage

import (
	"context"
	"fmt"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Minio struct {
	conn            *minio.Client
	endpoint        string
	accessKeyID     string
	secretAccessKey string
	useSSL          bool
}

func NewMinioClient(endpoint string, accessKeyID string, secretAccessKey string, useSSL bool) (*Minio, error) {
	//  Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	return &Minio{
		conn:            minioClient,
		endpoint:        endpoint,
		accessKeyID:     accessKeyID,
		secretAccessKey: secretAccessKey,
		useSSL:          useSSL,
	}, err
}

func (mc *Minio) ListBuckets(ctx context.Context) {
	buckets, err := mc.conn.ListBuckets(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	for _, bucket := range buckets {
		fmt.Println(bucket.Name)
	}
}
