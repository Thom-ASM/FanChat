package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	"github.com/aws/aws-sdk-go-v2/service/scheduler/types"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

type Event struct {
	ArtistName string `json:"artistName"`
	Location   string `json:"location"`
	EndTime    int64  `json:"endTime"`
	StartTime  int64  `json:"startTime"`
}

type App struct {
	youtube *youtube.Service
}

type BoundLiveStreamData struct {
	BroadcastID  string `dynamodbav:"broadcastId"`
	LivestreamID string `dynamodbav:"livestreamId"`
	RTMPAddr     string `dynamodbav:"rtmpAddr"`
	StreamKey    string `dynamodbav:"streamKey"`
	showId       string `dynamodbav:"showId"`
}

func insertBroadcast(service *youtube.Service, liveSteamName string, startTime time.Time) *youtube.LiveBroadcast {

	parts := []string{"snippet", "status"}

	// 2. Construct the LiveBroadcast struct (passed as a pointer)
	broadcast := &youtube.LiveBroadcast{
		Snippet: &youtube.LiveBroadcastSnippet{
			Title:              liveSteamName,
			ScheduledStartTime: startTime.Format(time.RFC3339),
		},
		Status: &youtube.LiveBroadcastStatus{
			PrivacyStatus: "unlisted",
		},
	}

	call := service.LiveBroadcasts.Insert(parts, broadcast)

	response, err := call.Do()

	if err != nil {

		log.Fatalf("failed to create broadcast : %v", err.Error())
	}

	return response
}

func createLivestream(service *youtube.Service) *youtube.LiveStream {

	parts := []string{"snippet", "cdn"}

	// 2. Construct the LiveBroadcast struct (passed as a pointer)
	broadcast := &youtube.LiveStream{
		Snippet: &youtube.LiveStreamSnippet{
			Title: "test livestream",
		},
		Cdn: &youtube.CdnSettings{

			FrameRate:     "30fps",
			Resolution:    "240p",
			IngestionType: "rtmp",
		},
	}

	call := service.LiveStreams.Insert(parts, broadcast)

	response, err := call.Do()

	if err != nil {

		log.Fatalf("failed to create livestream: %v", err.Error())
	}

	return response

}

func bindLiveStreamToBroadcast(service *youtube.Service, broadcastId string, livestreamId string) string {
	parts := []string{"id", "snippet", "contentDetails", "status"}

	call := service.LiveBroadcasts.Bind(broadcastId, parts)
	call = call.StreamId(livestreamId)

	response, err := call.Do()
	if err != nil {
		log.Fatalf("failed to bind livestream: %v", err)
	}

	fmt.Println(response.Id)
	return response.Id
}

func createHTTPClient(ctx context.Context, oauthConfig *oauth2.Config, awsConfig aws.Config) *http.Client {

	// fetch from aws
	svc := secretsmanager.NewFromConfig(awsConfig)

	secret, err := svc.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(""),
		VersionStage: aws.String("AWSCURRENT"),
	})

	if err != nil {
		log.Fatalf("failed to fetch secret %v", err)
	}

	if secret.SecretString == nil {
		log.Fatal("oauth secret contains no SecretString")
	}

	// create new oauthToken
	t := &oauth2.Token{}

	err = json.Unmarshal([]byte(*secret.SecretString), t)

	if err != nil {

		log.Fatalf("failed to decode oauthToken %v", err)

	}

	return oauthConfig.Client(ctx, t)

}

func genereateRTMPDetails(ytService *youtube.Service, livestreamName string, livestreamStartTime time.Time, showID string) BoundLiveStreamData {

	//create broadcast

	broadcast := insertBroadcast(ytService, livestreamName, livestreamStartTime)

	//create livestream

	liveStream := createLivestream(ytService)

	// bind both

	bindLiveStreamToBroadcast(ytService, broadcast.Id, liveStream.Id)

	// return data

	return BoundLiveStreamData{
		BroadcastID:  broadcast.Id,
		LivestreamID: liveStream.Id,
		RTMPAddr:     liveStream.Cdn.IngestionInfo.IngestionAddress,
		StreamKey:    liveStream.Cdn.IngestionInfo.StreamName,
		showId:       showID,
	}
}

func getSecretsManagerValue(ctx context.Context, smCLient *secretsmanager.Client, valueToSearch string) (*secretsmanager.GetSecretValueOutput, error) {
	secret, err := smCLient.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(valueToSearch),
		VersionStage: aws.String("AWSCURRENT"),
	})

	return secret, err
}

func handleRequest(ctx context.Context, event events.S3Event) error {

	//init aws service clients and misc config
	var DbTableName = os.Getenv("YOUTUBE_STREAM_TABLE_NAME")
	var YTOauthConfig = os.Getenv("OAUTH_CONFIG")

	ecsClusterArn := os.Getenv("ECS_CLUSTER_ARN")
	taskDefinitionArn := os.Getenv("ECS_TASK_DEFINITION_ARN")
	schedulerRoleArn := os.Getenv("SCHEDULER_ROLE_ARN")

	sdkConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Printf("failed to load default config: %s", err)
		return err
	}
	secretsManagerClient := secretsmanager.NewFromConfig(sdkConfig)
	dbCLient := dynamodb.NewFromConfig(sdkConfig)
	s3Client := s3.NewFromConfig(sdkConfig)
	schedulerClient := scheduler.NewFromConfig(sdkConfig)

	//Setup YT auth
	oauthConfig, err := getSecretsManagerValue(ctx, secretsManagerClient, YTOauthConfig)
	if err != nil {
		log.Fatalf("Failed to fetch oauth config from SM: %v", err)
	}
	if oauthConfig.SecretString == nil {
		log.Fatalf("oauthConfig is nil...")
	}
	googleOauthConfig, err := google.ConfigFromJSON([]byte(*oauthConfig.SecretString), youtube.YoutubeForceSslScope)

	if err != nil {
		log.Fatalf("failed to create oauthConfig %v", err)
	}

	httpClient := createHTTPClient(ctx, googleOauthConfig, sdkConfig)

	ytService, err := youtube.New(httpClient)

	if err != nil {

		log.Fatalf("failed to create youtube client: %v", err)
	}

	for _, record := range event.Records {
		bucket := record.S3.Bucket.Name
		key := record.S3.Object.URLDecodedKey
		file, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: &bucket,
			Key:    &key,
		})

		if err != nil {
			log.Printf("failed to get object from s3 bucket %s/%s: %s", bucket, key, err)
			return err
		}

		var upcomingEvents []Event

		json.NewDecoder(file.Body).Decode(&upcomingEvents)

		if err := json.NewDecoder(file.Body).Decode(&upcomingEvents); err != nil {
			return fmt.Errorf("failed to decode schedule document: %w", err)
		}

		// for each show, create a new broadcast/livestream/binding
		// then save it db table

		for _, item := range upcomingEvents {

			raw := fmt.Sprintf(
				"%s|%s|%d",
				item.ArtistName,
				item.Location,
				item.StartTime,
			)

			sum := sha256.Sum256([]byte(raw))
			showID := hex.EncodeToString(sum[:16])

			// we need to make idempotent
			// dbCLient.GetItem(ctx, &dynamodb.GetItemInput{Key: showID, TableName: &DbTableName, ProjectionExpression: aws.String("showID")})

			livestreamName := fmt.Sprintf("%s - %s", item.ArtistName, item.Location)

			startTime := time.Unix(item.StartTime, 0)

			// actual logic
			detailsToSave := genereateRTMPDetails(ytService, livestreamName, startTime, showID)

			// would ideally move this out to be a batch
			item, err := attributevalue.MarshalMap(detailsToSave)
			if err != nil {
				panic(err)
			}
			_, err = dbCLient.PutItem(ctx, &dynamodb.PutItemInput{
				TableName: aws.String(DbTableName), Item: item,
			})
			if err != nil {
				log.Printf("Couldn't add item to table. Here's why: %v\n", err)
			}

			isoTimestamp := startTime.UTC().Format("2006-01-02T15:04:05")

			schedulerClient.CreateSchedule(ctx, &scheduler.CreateScheduleInput{
				Name:                       aws.String(livestreamName),
				FlexibleTimeWindow:         &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
				ScheduleExpression:         aws.String(fmt.Sprintf("at(%s)", isoTimestamp)),
				ScheduleExpressionTimezone: aws.String("UTC"),
				Target: &types.Target{
					Arn:     aws.String(ecsClusterArn),
					RoleArn: aws.String(schedulerRoleArn),

					EcsParameters: &types.EcsParameters{
						TaskDefinitionArn: aws.String(taskDefinitionArn),
						LaunchType:        types.LaunchTypeFargate,
					},
				},
			})

		}

	}

	return nil
}

func main() {
	lambda.Start(handleRequest)
}
