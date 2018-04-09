package image_upload

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	awsAccessKey = "AKIAJDQHCA75XMHWVOTQ"
	awsSecret    = "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"
	token        = ""

	S3_BUCKET_RESOURCE_PREFIX = `https://s3-us-west-2.amazonaws.com/forbotchinmayr/media/`
)

func UploadImage(fileName string) {
	creds := credentials.NewStaticCredentials(awsAccessKey, awsSecret, token)
	_, err := creds.Get()
	if err != nil {
		fmt.Printf("bad credentials: %s", err)
	}

	cfg := aws.NewConfig().WithRegion("us-west-2").WithCredentials(creds)
	svc := s3.New(session.New(), cfg)

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("err opening file %s: %s", fileName, err)
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()

	buffer := make([]byte, size)
	file.Read(buffer)
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)

	path := "/media/" + file.Name()
	params := &s3.PutObjectInput{
		Bucket:        aws.String("forbotchinmayr"),
		Key:           aws.String(path),
		Body:          fileBytes,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fileType),
		ACL:           aws.String(s3.ObjectCannedACLPublicRead),
	}

	resp, err := svc.PutObject(params)
	if err != nil {
		fmt.Printf("bad response: %s", err)
	}
	fmt.Printf("Image upload to S3 response %s", awsutil.StringValue(resp))
}
