package amazon

import (
	"io"

	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/astaxie/beego"

	"github.com/astaxie/beego/logs"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func getCredentialsFromFile() *aws.Config {
	creds := credentials.NewSharedCredentials("conf/credentials.sh", "default")

	return &aws.Config{
		Region:      aws.String("us-east-2"),
		Credentials: creds,
	}
}

func SaveImageToS3(file io.Reader, filename string) (url string, err error) {
	// The session the S3 Uploader will use
	config := getCredentialsFromFile()
	sess := session.Must(session.NewSession(config))

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(beego.AppConfig.String("bucket")),
		Key:    aws.String(filename),
		Body:   file,
		//ContentType: aws.String("image/png"),
	})
	if err != nil {
		logs.Error(err)
		return "", err
	}
	return result.Location, nil
}
