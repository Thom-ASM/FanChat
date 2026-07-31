package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
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

func main() {
	ctx := context.Background()

	app, err := newApp(ctx)
	if err != nil {
		log.Fatalf("initialise application: %v", err)
	}

	lambda.Start(app.handle)
}

func newApp(ctx context.Context) (*App, error) {
	secretARN := os.Getenv("YOUTUBE_OAUTH_SECRET_ARN")
	if secretARN == "" {
		return nil, errors.New("YOUTUBE_OAUTH_SECRET_ARN is required")
	}

	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("load AWS config: %w", err)
	}

	secretsClient := secretsmanager.NewFromConfig(awsConfig)

	output, err := secretsClient.GetSecretValue(
		ctx,
		&secretsmanager.GetSecretValueInput{
			SecretId: aws.String(secretARN),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("read YouTube OAuth secret: %w", err)
	}

	if output.SecretString == nil {
		return nil, errors.New("YouTube OAuth secret has no SecretString")
	}

	var secret OAuthSecret

	if err := json.Unmarshal(
		[]byte(*output.SecretString),
		&secret,
	); err != nil {
		return nil, fmt.Errorf("decode YouTube OAuth secret: %w", err)
	}

	if secret.ClientID == "" {
		return nil, errors.New("YouTube OAuth client_id is missing")
	}

	if secret.ClientSecret == "" {
		return nil, errors.New("YouTube OAuth client_secret is missing")
	}

	if secret.RefreshToken == "" {
		return nil, errors.New("YouTube OAuth refresh_token is missing")
	}

	oauthConfig := &oauth2.Config{
		ClientID:     secret.ClientID,
		ClientSecret: secret.ClientSecret,
		Endpoint:     google.Endpoint,
		Scopes: []string{
			youtubeScope,
		},
	}

	token := &oauth2.Token{
		RefreshToken: secret.RefreshToken,
	}

	httpClient := oauthConfig.Client(ctx, token)

	youtubeService, err := youtube.NewService(
		ctx,
		option.WithHTTPClient(httpClient),
	)
	if err != nil {
		return nil, fmt.Errorf("create YouTube service: %w", err)
	}

	return &App{
		youtube: youtubeService,
	}, nil
}

func (app *App) handle(
	ctx context.Context,
	event Event,
) (Result, error) {
	log.Printf(
		"handling action=%q eventId=%q",
		event.Action,
		event.EventID,
	)

	switch event.Action {
	case "prepare":
		return app.prepareBroadcast(ctx, event)

	default:
		return Result{}, fmt.Errorf(
			"unsupported action %q",
			event.Action,
		)
	}
}

func (app *App) prepareBroadcast(
	ctx context.Context,
	event Event,
) (Result, error) {
	if err := validatePrepareEvent(event); err != nil {
		return Result{}, err
	}

	start, err := time.Parse(
		time.RFC3339,
		event.ScheduledStart,
	)
	if err != nil {
		return Result{}, fmt.Errorf(
			"parse scheduledStart: %w",
			err,
		)
	}

	end, err := time.Parse(
		time.RFC3339,
		event.ScheduledEnd,
	)
	if err != nil {
		return Result{}, fmt.Errorf(
			"parse scheduledEnd: %w",
			err,
		)
	}

	if !end.After(start) {
		return Result{}, errors.New(
			"scheduledEnd must be after scheduledStart",
		)
	}

	broadcast := &youtube.LiveBroadcast{
		Snippet: &youtube.LiveBroadcastSnippet{
			Title:              event.Title,
			Description:        "FanChat event: " + event.EventID,
			ScheduledStartTime: start.UTC().Format(time.RFC3339),
			ScheduledEndTime:   end.UTC().Format(time.RFC3339),
		},

		Status: &youtube.LiveBroadcastStatus{
			PrivacyStatus:           "unlisted",
			SelfDeclaredMadeForKids: false,

			// False is otherwise omitted by the generated JSON
			// marshaller, so explicitly include it.
			ForceSendFields: []string{
				"SelfDeclaredMadeForKids",
			},
		},

		ContentDetails: &youtube.LiveBroadcastContentDetails{
			EnableEmbed:     true,
			EnableAutoStart: false,
			EnableAutoStop:  false,

			MonitorStream: &youtube.MonitorStreamInfo{
				EnableMonitorStream: boolPtr(false),
			},

			// These fields are plain bool values. Explicitly include
			// them even though their configured value is false.
			ForceSendFields: []string{
				"EnableAutoStart",
				"EnableAutoStop",
			},
		},
	}

	created, err := app.youtube.LiveBroadcasts.
		Insert(
			[]string{
				"id",
				"snippet",
				"status",
				"contentDetails",
			},
			broadcast,
		).
		Context(ctx).
		Do()
	if err != nil {
		return Result{}, fmt.Errorf(
			"create YouTube broadcast: %w",
			err,
		)
	}

	if created.Id == "" {
		return Result{}, errors.New(
			"YouTube created a broadcast without returning an ID",
		)
	}

	log.Printf(
		"created YouTube broadcast broadcastId=%q eventId=%q",
		created.Id,
		event.EventID,
	)

	_, err = app.youtube.LiveBroadcasts.
		Bind(
			created.Id,
			[]string{
				"id",
				"snippet",
				"status",
				"contentDetails",
			},
		).
		StreamId(event.YouTubeStreamID).
		Context(ctx).
		Do()
	if err != nil {
		log.Printf(
			"binding failed; deleting broadcast broadcastId=%q",
			created.Id,
		)

		deleteErr := app.youtube.LiveBroadcasts.
			Delete(created.Id).
			Context(ctx).
			Do()

		if deleteErr != nil {
			log.Printf(
				"failed to clean up broadcast broadcastId=%q error=%v",
				created.Id,
				deleteErr,
			)
		}

		return Result{}, fmt.Errorf(
			"bind broadcast %q to stream %q: %w",
			created.Id,
			event.YouTubeStreamID,
			err,
		)
	}

	var liveChatID string

	if created.Snippet != nil {
		liveChatID = created.Snippet.LiveChatId
	}

	log.Printf(
		"bound broadcast broadcastId=%q streamId=%q",
		created.Id,
		event.YouTubeStreamID,
	)

	return Result{
		EventID:            event.EventID,
		YouTubeBroadcastID: created.Id,
		YouTubeVideoID:     created.Id,
		YouTubeLiveChatID:  liveChatID,
	}, nil
}

func validatePrepareEvent(event Event) error {
	if event.EventID == "" {
		return errors.New("eventId is required")
	}

	if event.Title == "" {
		return errors.New("title is required")
	}

	if event.ScheduledStart == "" {
		return errors.New("scheduledStart is required")
	}

	if event.ScheduledEnd == "" {
		return errors.New("scheduledEnd is required")
	}

	if event.YouTubeStreamID == "" {
		return errors.New("youtubeStreamId is required")
	}

	return nil
}

func boolPtr(value bool) *bool {
	return &value
}
