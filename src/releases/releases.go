package releases

import (
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
)

var s3Client *s3.S3

type Release struct {
	Platform     string
	Type         string
	Time         int64
	Name         string
	MajorRelease int
	MinorRelease int
	Compilation  int
}

func GetReleasesByTypePlatform(platform string, releaseType string) ([]Release, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(os.Getenv("S3_BUCKET")),
		Prefix: aws.String(platform + "/" + releaseType),
	}

	result, err := s3Client.ListObjectsV2(input)

	if err != nil {
		return nil, err
	}

	releases := []Release{}
	for _, object := range result.Contents {
		keyParts := strings.Split(*object.Key, "/")
		prettyKey := keyParts[len(keyParts)-1]

		releaseParts := strings.Split(prettyKey, ".")

		majorReleaseParts := strings.Split(releaseParts[0], "-")

		majorRelease, err := strconv.Atoi(majorReleaseParts[len(majorReleaseParts)-1])

		if err != nil {
			return nil, err
		}

		minorRelease, err := strconv.Atoi(releaseParts[1])

		if err != nil {
			return nil, err
		}

		compilationRelease, err := strconv.Atoi(releaseParts[2])

		if err != nil {
			return nil, err
		}

		releases = append(releases, Release{
			Name:         prettyKey,
			Type:         releaseType,
			Platform:     platform,
			Time:         object.LastModified.Unix(),
			MajorRelease: majorRelease,
			MinorRelease: minorRelease,
			Compilation:  compilationRelease,
		})
	}

	sort.Slice(releases, func(i, j int) bool {
		return releases[i].Time >= releases[j].Time
	})

	return releases, nil
}

func GetDownloadLink(platform string, releaseType string, prettyName string) (string, error) {
	req, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET")),
		Key:    aws.String(platform + "/" + releaseType + "/" + prettyName),
	})

	urlStr, err := req.Presign(15 * time.Minute)

	if err != nil {
		return "", err
	}

	return urlStr, nil
}

func init() {
	godotenv.Load()

	log.Println("Starting releases ..")
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(os.Getenv("awsAccessKeyId"), os.Getenv("awsSecretAccessKey"), ""),
		Endpoint:         aws.String(os.Getenv("S3_ENDPOINT")),
		Region:           aws.String("GRA"),
		S3ForcePathStyle: aws.Bool(true),
	}

	newSession, err := session.NewSession(s3Config)

	if err != nil {
		log.Fatal(err)
	}

	s3Client = s3.New(newSession)
}
