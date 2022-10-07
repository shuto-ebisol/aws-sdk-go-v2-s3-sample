// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX - License - Identifier: Apache - 2.0
// Go V2 SDK examples: Common S3 actions

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func BucketOps(client s3.Client, name string) {
	fmt.Println("Create Presign client")
	presignClient := s3.NewPresignClient(&client)

	// 署名付きURL発行用パラメータ設定
	shopID := 2032
	reservationID := 1
	order := 1
	now := 0000000000000000
	extension := "png"
	key := fmt.Sprintf("contents/%d/%d/%d_%d.%s", shopID, reservationID, order, now, extension)
	presignParams := &s3.GetObjectInput{
		Bucket: aws.String(name),
		Key:    aws.String(key),
	}

	// Apply an expiration via an option function
	presignDuration := func(po *s3.PresignOptions) {
		po.Expires = 5 * time.Minute
	}

	presignResult, err := presignClient.PresignGetObject(context.TODO(), presignParams, presignDuration)
	if err != nil {
		panic("Couldn't get presigned URL for GetObject")
	}

	fmt.Printf("Presigned URL For object: %s\n", presignResult.URL)
}

func main() {
	//snippet-start:[s3.go-v2.s3_basics]
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           "http://localhost:4566",
			SigningRegion: "ap-northeast-1",
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("ap-northeast-1"),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	s3client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true // ローカルではこれをしないと、http://{{ バケット名 }}.{{ エンドポイント }}の形式でURLが組み立てられてしまう
	})

	myBucketName := "local-secure-reservation-images-xxxxxxxxxxxx"
	fmt.Printf("Bucket name: %v\n", myBucketName)

	BucketOps(*s3client, myBucketName)
}
