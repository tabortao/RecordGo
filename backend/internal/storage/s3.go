package storage

import (
    "context"
    "errors"
    "io"
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/credentials"
    awsconfig "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
    "time"
)

type S3Storage struct {
    client *s3.Client
    bucket string
    publicBase string
    makePublic bool
    usePathStyle bool
}

func NewS3(cfgRegion, endpoint, accessKey, secretKey, bucket string, forcePath bool, publicBase string) (*S3Storage, error) {
    var awsCfg aws.Config
    var err error
    if accessKey != "" && secretKey != "" {
        awsCfg, err = awsconfig.LoadDefaultConfig(context.Background(),
            awsconfig.WithRegion(cfgRegion),
            awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
        )
    } else {
        awsCfg, err = awsconfig.LoadDefaultConfig(context.Background(), awsconfig.WithRegion(cfgRegion))
    }
    if err != nil { return nil, err }
    if endpoint != "" {
        awsCfg.EndpointResolverWithOptions = aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
            return aws.Endpoint{URL: endpoint, SigningRegion: cfgRegion, HostnameImmutable: true}, nil
        })
    }
    cli := s3.NewFromConfig(awsCfg, func(o *s3.Options) { o.UsePathStyle = forcePath })
    return &S3Storage{client: cli, bucket: bucket, publicBase: publicBase, makePublic: publicBase != "", usePathStyle: forcePath}, nil
}

func (s *S3Storage) DriverName() string { return "s3" }

func (s *S3Storage) PresignPut(ctx context.Context, key string, contentType string, expireSeconds int64) (string, map[string]string, error) {
    p := s3.NewPresignClient(s.client)
    in := &s3.PutObjectInput{Bucket: &s.bucket, Key: &key, ContentType: &contentType}
    if s.makePublic {
        acl := s3types.ObjectCannedACLPublicRead
        in.ACL = acl
    }
    dur := time.Duration(expireSeconds) * time.Second
    out, err := p.PresignPutObject(ctx, in, func(po *s3.PresignOptions) { po.Expires = dur })
    if err != nil { return "", nil, err }
    headers := map[string]string{"Content-Type": contentType}
    if s.makePublic { headers["x-amz-acl"] = string(s3types.ObjectCannedACLPublicRead) }
    return out.URL, headers, nil
}

func (s *S3Storage) PresignGet(ctx context.Context, key string, expireSeconds int64) (string, error) {
    if s.publicBase != "" {
        // MinIO/path-style 公共访问：需要在路径中包含桶名
        if s.usePathStyle {
            return s.publicBase + "/" + s.bucket + "/" + key, nil
        }
        // 虚拟主机风格（不常用于 MinIO）：假定 publicBase 已是桶域名
        return s.publicBase + "/" + key, nil
    }
    p := s3.NewPresignClient(s.client)
    in := &s3.GetObjectInput{Bucket: &s.bucket, Key: &key}
    dur := time.Duration(expireSeconds) * time.Second
    out, err := p.PresignGetObject(ctx, in, func(po *s3.PresignOptions) { po.Expires = dur })
    if err != nil { return "", err }
    return out.URL, nil
}

func (s *S3Storage) Put(ctx context.Context, key string, r io.Reader, contentType string) error {
    return errors.New("not supported")
}

func (s *S3Storage) PublicURL(key string) (string, bool) {
    if s.publicBase == "" { return "", false }
    return s.publicBase + "/" + key, true
}