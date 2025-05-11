// internal/services/openai_service.go
package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/wizenheimer/swiftcal/internal/config"
	"github.com/wizenheimer/swiftcal/internal/models"
	"github.com/wizenheimer/swiftcal/pkg/logger"
	"github.com/wizenheimer/swiftcal/templates"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"go.uber.org/zap"
)

type OpenAIService struct {
	client *openai.Client
	config *config.Config
}

func NewOpenAIService(cfg *config.Config) *OpenAIService {
	client := openai.NewClient(option.WithAPIKey(cfg.OpenAIAPIKey))
	return &OpenAIService{
		client: &client,
		config: cfg,
	}
}

func (s *OpenAIService) ProcessEmail(ctx context.Context, emailContent, subject, from, date string) (*models.EventsResponse, *models.TimezoneResponse, error) {
	emailText := fmt.Sprintf("Date: %s\nSubject: %s\nFrom: %s\n%s", date, subject, from, emailContent)

	// Process events and timezone in parallel
	eventsChan := make(chan *models.EventsResponse, 1)
	timezoneChan := make(chan *models.TimezoneResponse, 1)
	errChan := make(chan error, 2)

	// Extract events
	go func() {
		events, err := s.extractEvents(ctx, emailText)
		if err != nil {
			errChan <- err
			return
		}
		eventsChan <- events
	}()

	// Extract timezone
	go func() {
		timezone, err := s.extractTimezone(ctx, emailText)
		if err != nil {
			errChan <- err
			return
		}
		timezoneChan <- timezone
	}()

	// Wait for both results
	var eventsResponse *models.EventsResponse
	var timezoneResponse *models.TimezoneResponse
	var errors []error

	for i := 0; i < 2; i++ {
		select {
		case events := <-eventsChan:
			eventsResponse = events
		case timezone := <-timezoneChan:
			timezoneResponse = timezone
		case err := <-errChan:
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		logger.GetLogger().Error("OpenAI processing failed", zap.Error(errors[0]))
		return nil, nil, errors[0]
	}

	// Apply timezone to all events
	if eventsResponse != nil && eventsResponse.Events != nil && timezoneResponse != nil {
		for i := range eventsResponse.Events {
			eventsResponse.Events[i].TimeZone = timezoneResponse.Timezone
		}
	}

	if timezoneResponse != nil && timezoneResponse.Timezone != nil {
		logger.GetLogger().Debug("Email processed by OpenAI",
			zap.Int("events_count", len(eventsResponse.Events)),
			zap.String("timezone", *timezoneResponse.Timezone),
			zap.String("reason", timezoneResponse.Reason),
		)
	} else {
		logger.GetLogger().Debug("Email processed by OpenAI",
			zap.Int("events_count", len(eventsResponse.Events)),
		)
	}

	return eventsResponse, timezoneResponse, nil
}

func (s *OpenAIService) extractEvents(ctx context.Context, emailText string) (*models.EventsResponse, error) {
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(templates.GetEventExtractionPrompt()),
		openai.UserMessage(emailText),
	}

	completion, err := s.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages:    messages,
		Model:       openai.ChatModelGPT4oMini,
		Temperature: openai.Float(0.1),
		MaxTokens:   openai.Int(4096),
	})

	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(completion.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	content := completion.Choices[0].Message.Content
	if content == "" {
		return nil, fmt.Errorf("empty response from OpenAI")
	}

	var eventsResponse models.EventsResponse
	if err := json.Unmarshal([]byte(content), &eventsResponse); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAI response: %w", err)
	}

	logger.GetLogger().Debug("OpenAI usage for event extraction", zap.Any("usage", completion.Usage))
	return &eventsResponse, nil
}

func (s *OpenAIService) extractTimezone(ctx context.Context, emailText string) (*models.TimezoneResponse, error) {
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(templates.GetTimezoneExtractionPrompt()),
		openai.UserMessage(emailText),
	}

	completion, err := s.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages:    messages,
		Model:       openai.ChatModelGPT4oMini,
		Temperature: openai.Float(0.1),
		MaxTokens:   openai.Int(1024),
	})

	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(completion.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	content := completion.Choices[0].Message.Content
	if content == "" {
		return nil, fmt.Errorf("empty response from OpenAI")
	}

	var timezoneResponse models.TimezoneResponse
	if err := json.Unmarshal([]byte(content), &timezoneResponse); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAI response: %w", err)
	}

	logger.GetLogger().Debug("OpenAI usage for timezone extraction", zap.Any("usage", completion.Usage))
	return &timezoneResponse, nil
}
