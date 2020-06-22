package aws

// func S3Upload(filePath string) (string, error) {

// 	sess, err := session.NewSession(&aws.Config{
// 		Credentials:      credentials.NewStaticCredentials(setting.S3AccessKey, setting.S3SecretKey, ""),
// 		Endpoint:         aws.String(setting.S3EndPoint),
// 		Region:           aws.String(setting.S3Region),
// 		DisableSSL:       aws.Bool(true),
// 		S3ForcePathStyle: aws.Bool(false), //virtual-host style方式，不要修改
// 	})

// 	if err != nil {
// 		fmt.Println(err)
// 		return "", err
// 	}

// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		fmt.Printf("Unable to open file %q", err)
// 		return "", err
// 	}

// 	defer file.Close()

// 	uploader := s3manager.NewUploader(sess)

// 	uploadOutput, err := uploader.Upload(&s3manager.UploadInput{
// 		Bucket: aws.String(setting.S3Bucket),
// 		Key:    aws.String(filePath),
// 		Body:   file,
// 	})
// 	if err != nil {
// 		fmt.Printf("Unable to upload %q to %q, %v", filePath, setting.S3Bucket, err)
// 		return "", err
// 	}

// 	return uploadOutput.Location, nil
// }
