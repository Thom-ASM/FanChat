package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambda/handlertrace"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

const youtubeScope = "https://www.googleapis.com/auth/youtube.force-ssl"

type OAuthSecret struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
}

type Event struct {
	Action          string `json:"action"`
	EventID         string `json:"eventId"`
	Title           string `json:"title"`
	ScheduledStart  string `json:"scheduledStart"`
	ScheduledEnd    string `json:"scheduledEnd"`
	YouTubeStreamID string `json:"youtubeStreamId"`
}

type Result struct {
	EventID            string `json:"eventId"`
	YouTubeBroadcastID string `json:"youtubeBroadcastId"`
	YouTubeVideoID     string `json:"youtubeVideoId"`
	YouTubeLiveChatID  string `json:"youtubeLiveChatId"`
}

type App struct {
	youtube *youtube.Service
}

func handleRequest(ctx context.Context, event events.S3Event) error {

	// get the streams from the event,
	sdkConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Printf("failed to load default config: %s", err)
		return err
	}

	s3Client := s3.NewFromConfig(sdkConfig)

	for _, record := range event.Records {
		bucket := record.S3.Bucket.Name
		key := record.S3.Object.URLDecodedKey

		file, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: &bucket,
			Key:    &key,
		})

		var upcomingEvents []Event

		data := json.NewDecoder(file.Body).Decode(&upcomingEvents)

		// we only want the next 10 events

		var recentUpcoming []Event = upcomingEvents[0:9]

		//Compare against dynamo db

		if err != nil {
			log.Printf("error failed to fetch file from %s/%s: %s", bucket, key, err)
			return err
		}
	}

	//Compare against dynamo db

	//create streams as required

	//generate eventbridge schedule

	return nil
}

func main() {
	lambda.Start(handleRequest())
}
