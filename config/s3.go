package config

type S3 struct {
	Bucket    string
	Region    string
	AccessKey string
	SecretKey string
}

func NewS3FromEnv() S3 {
	return S3{
		Bucket:    getEnv("S3_BUCKET", "example-bucket"),
		Region:    getEnv("S3_REGION", "us-east-1"),
		AccessKey: getEnv("S3_ACCESS_KEY", "default-access"),
		SecretKey: getEnv("S3_SECRET_KEY", "default-secret"),
	}
}
