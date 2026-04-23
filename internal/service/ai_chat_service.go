package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/LuuDinhTheTai/tzone/infrastructure/configuration"
	"github.com/LuuDinhTheTai/tzone/internal/dto"
)

type AIChatService struct {
	apiKey        string
	model         string
	youtubeAPIKey string
	publicBaseURL string
	bucket        string
	httpClient    *http.Client
	deviceIndex   map[string]catalogDevice
	catalog       []catalogDevice
}

type catalogDevice struct {
	ID        string
	BrandName string
	ModelName string
	ImageURL  string
	OS        string
	Chipset   string
	Memory    string
	Battery   string
	Price     string
	searchDoc string
}

type phoneCatalogFileBrand struct {
	BrandName string                   `json:"brand_name"`
	Devices   []phoneCatalogFileDevice `json:"devices"`
}

type phoneCatalogFileDevice struct {
	ID             phoneCatalogObjectID           `json:"_id"`
	ModelName      string                         `json:"model_name"`
	ImageURL       string                         `json:"imageUrl"`
	Specifications phoneCatalogFileSpecifications `json:"specifications"`
}

type phoneCatalogObjectID struct {
	OID string `json:"$oid"`
}

type phoneCatalogFileSpecifications struct {
	Platform phoneCatalogPlatform `json:"Platform"`
	Memory   phoneCatalogMemory   `json:"Memory"`
	Battery  phoneCatalogBattery  `json:"Battery"`
	Misc     phoneCatalogMisc     `json:"Misc"`
}

type phoneCatalogPlatform struct {
	OS      string `json:"OS"`
	Chipset string `json:"Chipset"`
}

type phoneCatalogMemory struct {
	Internal string `json:"Internal"`
}

type phoneCatalogBattery struct {
	Type string `json:"Type"`
}

type phoneCatalogMisc struct {
	Price string `json:"price"`
}

type geminiRequest struct {
	Contents         []geminiContent        `json:"contents"`
	GenerationConfig geminiGenerationConfig `json:"generationConfig"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiGenerationConfig struct {
	Temperature      float64 `json:"temperature"`
	ResponseMimeType string  `json:"responseMimeType"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

type geminiStructuredReply struct {
	Reply     string   `json:"reply"`
	DeviceIDs []string `json:"device_ids"`
}

type geminiVideoReply struct {
	PhoneName   string `json:"phoneName"`
	VideoTitle  string `json:"videoTitle"`
	YoutubeUrl  string `json:"youtubeUrl"`
	ChannelName string `json:"channelName"`
}

type geminiVideoReplyRow struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type scoredCandidate struct {
	Device catalogDevice
	Score  int
}

type GeminiAPIError struct {
	StatusCode int
	Body       string
}

func (e *GeminiAPIError) Error() string {
	if strings.TrimSpace(e.Body) == "" {
		return fmt.Sprintf("gemini request failed with status %d", e.StatusCode)
	}
	return fmt.Sprintf("gemini request failed with status %d: %s", e.StatusCode, e.Body)
}

func (e *GeminiAPIError) FriendlyMessage() string {
	switch e.StatusCode {
	case http.StatusTooManyRequests:
		return "The AI service is temporarily rate-limited or out of quota. Please try again in a few minutes."
	case http.StatusUnauthorized:
		return "The AI service key is invalid or expired. Please check GEMINI_API_KEY."
	case http.StatusForbidden:
		return "The AI service does not have permission to use this model. Please check your Gemini API access."
	case http.StatusNotFound:
		return "The configured Gemini model is not available. Please update GEMINI_MODEL or use a supported model."
	default:
		return "Unable to generate AI recommendations right now. Please try again later."
	}
}

func NewAIChatService(cfg configuration.AIConfig) (*AIChatService, error) {
	catalog, err := loadCatalogFromFile(cfg.PhoneDataPath, cfg.MinioPublicBaseURL, cfg.MinioBucket)
	if err != nil {
		return nil, err
	}

	deviceIndex := make(map[string]catalogDevice, len(catalog))
	for _, device := range catalog {
		deviceIndex[device.ID] = device
	}

	return &AIChatService{
		apiKey:        strings.TrimSpace(cfg.GeminiAPIKey),
		model:         strings.TrimSpace(cfg.GeminiModel),
		youtubeAPIKey: strings.TrimSpace(cfg.YouTubeAPIKey),
		publicBaseURL: strings.TrimSpace(cfg.MinioPublicBaseURL),
		bucket:        strings.TrimSpace(cfg.MinioBucket),
		httpClient:    &http.Client{Timeout: 20 * time.Second},
		deviceIndex:   deviceIndex,
		catalog:       catalog,
	}, nil
}

func loadCatalogFromFile(path string, publicBaseURL string, bucket string) ([]catalogDevice, error) {
	if strings.TrimSpace(path) == "" {
		return nil, fmt.Errorf("phone catalog path is empty")
	}

	resolvedPath, err := resolveCatalogPath(path)
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(resolvedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read phone catalog %s: %w", resolvedPath, err)
	}

	var brands []phoneCatalogFileBrand
	if err := json.Unmarshal(content, &brands); err != nil {
		return nil, fmt.Errorf("failed to parse phone catalog %s: %w", resolvedPath, err)
	}

	catalog := make([]catalogDevice, 0)
	for _, brand := range brands {
		for _, item := range brand.Devices {
			id := strings.TrimSpace(item.ID.OID)
			if id == "" {
				continue
			}

			device := catalogDevice{
				ID:        id,
				BrandName: strings.TrimSpace(brand.BrandName),
				ModelName: strings.TrimSpace(item.ModelName),
				ImageURL:  normalizeCatalogImageURL(item.ImageURL, publicBaseURL, bucket),
				OS:        strings.TrimSpace(item.Specifications.Platform.OS),
				Chipset:   strings.TrimSpace(item.Specifications.Platform.Chipset),
				Memory:    strings.TrimSpace(item.Specifications.Memory.Internal),
				Battery:   strings.TrimSpace(item.Specifications.Battery.Type),
				Price:     strings.TrimSpace(item.Specifications.Misc.Price),
			}
			device.searchDoc = strings.ToLower(strings.Join([]string{
				device.BrandName,
				device.ModelName,
				device.OS,
				device.Chipset,
				device.Memory,
				device.Battery,
				device.Price,
			}, " "))

			catalog = append(catalog, device)
		}
	}

	if len(catalog) == 0 {
		return nil, fmt.Errorf("phone catalog %s does not contain devices", resolvedPath)
	}

	return catalog, nil
}

func normalizeCatalogImageURL(imageURL string, publicBaseURL string, bucket string) string {
	raw := strings.TrimSpace(imageURL)
	if raw == "" {
		return ""
	}
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		return raw
	}

	normalizedPath := strings.TrimLeft(strings.ReplaceAll(raw, "\\", "/"), "/")
	if normalizedPath == "" {
		return ""
	}

	if publicBaseURL != "" && bucket != "" {
		return fmt.Sprintf("%s/%s/%s", strings.TrimRight(publicBaseURL, "/"), strings.Trim(bucket, "/"), normalizedPath)
	}

	return "/" + normalizedPath
}

func resolveCatalogPath(path string) (string, error) {
	cleanPath := filepath.Clean(strings.TrimSpace(path))
	if filepath.IsAbs(cleanPath) {
		if fileExists(cleanPath) {
			return cleanPath, nil
		}
		return "", fmt.Errorf("failed to read phone catalog %s: %w", cleanPath, os.ErrNotExist)
	}

	tryPaths := make([]string, 0, 8)

	if wd, err := os.Getwd(); err == nil {
		tryPaths = append(tryPaths,
			filepath.Join(wd, cleanPath),
			filepath.Join(wd, "..", cleanPath),
			filepath.Join(wd, "..", "..", cleanPath),
		)
	}

	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		tryPaths = append(tryPaths,
			filepath.Join(exeDir, cleanPath),
			filepath.Join(exeDir, "..", cleanPath),
		)
	}

	tryPaths = append(tryPaths, cleanPath)

	seen := map[string]struct{}{}
	for _, candidate := range tryPaths {
		absCandidate, err := filepath.Abs(candidate)
		if err != nil {
			continue
		}
		if _, ok := seen[absCandidate]; ok {
			continue
		}
		seen[absCandidate] = struct{}{}
		if fileExists(absCandidate) {
			return absCandidate, nil
		}
	}

	return "", fmt.Errorf("failed to read phone catalog %s: %w", cleanPath, os.ErrNotExist)
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func (s *AIChatService) Recommend(ctx context.Context, req dto.AIChatRecommendRequest) (*dto.AIChatRecommendResponse, error) {
	if strings.TrimSpace(req.Message) == "" {
		return nil, fmt.Errorf("message is required")
	}

	candidates := s.pickCandidates(req.Message, 40)
	structured := &geminiStructuredReply{}
	if s.apiKey != "" {
		modelReply, err := s.askGemini(ctx, req.Message, candidates, req.Limit)
		if err != nil {
			var geminiErr *GeminiAPIError
			if errors.As(err, &geminiErr) {
				switch geminiErr.StatusCode {
				case http.StatusTooManyRequests, http.StatusUnauthorized, http.StatusForbidden:
					return nil, geminiErr
				}
			}
			log.Printf("⚠️ Gemini call failed, using local fallback: %v", err)
		} else {
			structured = modelReply
		}
	} else {
		log.Printf("⚠️ GEMINI_API_KEY missing, using local fallback recommendations")
	}

	cards := s.buildCards(structured.DeviceIDs, req.Limit)
	if len(cards) == 0 {
		for _, candidate := range candidates {
			cards = append(cards, toCard(candidate))
			if len(cards) >= req.Limit {
				break
			}
		}
	}

	reply := strings.TrimSpace(structured.Reply)
	if reply == "" {
		reply = s.buildFallbackReply(cards)
	}

	return &dto.AIChatRecommendResponse{
		Reply:   reply,
		Devices: cards,
	}, nil
}

func (s *AIChatService) FindVideoReviews(ctx context.Context, req dto.AIVideoReviewRequest) (*dto.AIVideoReviewResponse, error) {
	deviceName := strings.TrimSpace(req.DeviceName)
	if deviceName == "" {
		return nil, fmt.Errorf("device_name is required")
	}

	fallbackVideos := buildSingleYouTubeFallbackVideo(deviceName)

	if s.youtubeAPIKey == "" {
		return &dto.AIVideoReviewResponse{
			Reply:  "I could not fetch a direct review video right now. Here is one YouTube search link you can open now.",
			Videos: fallbackVideos,
		}, nil
	}

	query := deviceName + " review"
	videos, err := s.fetchYouTubeVideos(ctx, query, req.Limit)
	if err != nil {
		log.Printf("⚠️ YouTube API call failed: %v", err)
		return &dto.AIVideoReviewResponse{
			Reply:  "I could not fetch a direct review video at the moment. Here is one YouTube search link you can open now.",
			Videos: fallbackVideos,
		}, nil
	}

	reply := "Here is a YouTube review link for this device."
	if len(videos) == 0 {
		videos = fallbackVideos
		reply = "I could not verify a direct review video right now. Here is one YouTube search link for this device."
	} else {
		reply = fmt.Sprintf("Found a YouTube review for %s", deviceName)
	}

	return &dto.AIVideoReviewResponse{Reply: reply, Videos: videos}, nil
}

func (s *AIChatService) fetchYouTubeVideos(ctx context.Context, query string, limit int) ([]dto.AIVideoReview, error) {
	if s.youtubeAPIKey == "" {
		return nil, fmt.Errorf("youtube API key is missing")
	}

	if limit <= 0 {
		limit = 1
	}

	searchURL := fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?part=snippet&q=%s&type=video&maxResults=%d&key=%s", url.QueryEscape(query), limit, s.youtubeAPIKey)
	
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create youtube request: %w", err)
	}

	res, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call youtube api: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("youtube api returned status %d: %s", res.StatusCode, string(body))
	}

	var ytRes struct {
		Items []struct {
			Id struct {
				VideoId string `json:"videoId"`
			} `json:"id"`
			Snippet struct {
				Title string `json:"title"`
			} `json:"snippet"`
		} `json:"items"`
	}

	if err := json.NewDecoder(res.Body).Decode(&ytRes); err != nil {
		return nil, fmt.Errorf("failed to decode youtube response: %w", err)
	}

	var videos []dto.AIVideoReview
	for _, item := range ytRes.Items {
		if item.Id.VideoId != "" {
			videos = append(videos, dto.AIVideoReview{
				Title: item.Snippet.Title,
				URL:   "https://www.youtube.com/watch?v=" + item.Id.VideoId,
			})
		}
	}

	return videos, nil
}

func (s *AIChatService) buildFallbackReply(cards []dto.RecommendedDeviceCard) string {
	if len(cards) == 0 {
		return "I could not find a strong match yet. Please share more details about your budget, brand, operating system, camera, or battery needs."
	}

	names := make([]string, 0, len(cards))
	for _, card := range cards {
		fullName := strings.TrimSpace(strings.TrimSpace(card.BrandName) + " " + strings.TrimSpace(card.ModelName))
		if fullName == "" {
			continue
		}
		names = append(names, fullName)
		if len(names) >= 3 {
			break
		}
	}

	if len(names) == 0 {
		return "I found several devices that match your request in the current catalog."
	}

	return fmt.Sprintf("I recommend %s. You can open the cards below to see the details.", strings.Join(names, ", "))
}

func (s *AIChatService) pickCandidates(message string, max int) []catalogDevice {
	query := strings.ToLower(strings.TrimSpace(message))
	if query == "" {
		if len(s.catalog) <= max {
			return append([]catalogDevice(nil), s.catalog...)
		}
		return append([]catalogDevice(nil), s.catalog[:max]...)
	}

	tokens := strings.Fields(query)
	scored := make([]scoredCandidate, 0, len(s.catalog))
	for _, device := range s.catalog {
		score := 0
		for _, token := range tokens {
			if len(token) < 2 {
				continue
			}
			if strings.Contains(device.searchDoc, token) {
				score++
			}
		}
		if score > 0 {
			scored = append(scored, scoredCandidate{Device: device, Score: score})
		}
	}

	if len(scored) == 0 {
		if len(s.catalog) <= max {
			return append([]catalogDevice(nil), s.catalog...)
		}
		return append([]catalogDevice(nil), s.catalog[:max]...)
	}

	sort.SliceStable(scored, func(i, j int) bool {
		if scored[i].Score == scored[j].Score {
			return scored[i].Device.ModelName < scored[j].Device.ModelName
		}
		return scored[i].Score > scored[j].Score
	})

	if len(scored) > max {
		scored = scored[:max]
	}

	result := make([]catalogDevice, 0, len(scored))
	for _, item := range scored {
		result = append(result, item.Device)
	}
	return result
}

func (s *AIChatService) askGemini(ctx context.Context, message string, candidates []catalogDevice, limit int) (*geminiStructuredReply, error) {
	modelsToTry := buildGeminiModelCandidates(s.model)

	catalogLines := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		catalogLines = append(catalogLines, fmt.Sprintf(
			"- id=%s | brand=%s | model=%s | os=%s | chipset=%s | memory=%s | battery=%s | price=%s",
			candidate.ID,
			candidate.BrandName,
			candidate.ModelName,
			candidate.OS,
			candidate.Chipset,
			candidate.Memory,
			candidate.Battery,
			candidate.Price,
		))
	}

	prompt := fmt.Sprintf(`You are a phone recommendation assistant for TZone.
Only recommend devices from the catalog below.
Do not recommend any product outside the catalog.
Return ONLY valid JSON following this schema:
{"reply":"string","device_ids":["id1","id2"]}
- reply: concise, friendly, and in English.
- device_ids: at most %d items, and each id must exist in the catalog.
- If nothing matches, return an empty device_ids array and explain why in the reply.

Valid device catalog:
%s

User request: %s`, limit, strings.Join(catalogLines, "\n"), strings.TrimSpace(message))

	body := geminiRequest{
		Contents: []geminiContent{{
			Parts: []geminiPart{{Text: prompt}},
		}},
		GenerationConfig: geminiGenerationConfig{
			Temperature:      0.25,
			ResponseMimeType: "application/json",
		},
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal gemini request: %w", err)
	}

	var lastErr error
	for idx, model := range modelsToTry {
		resBody, err := s.callGeminiGenerateContent(ctx, model, payload)
		if err != nil {
			var geminiErr *GeminiAPIError
			if errors.As(err, &geminiErr) && geminiErr.StatusCode == http.StatusNotFound && idx < len(modelsToTry)-1 {
				log.Printf("⚠️ Gemini model %s not found, retrying with fallback model", model)
				lastErr = err
				continue
			}
			return nil, err
		}

		var modelRes geminiResponse
		if err := json.Unmarshal(resBody, &modelRes); err != nil {
			return nil, fmt.Errorf("failed to parse gemini response: %w", err)
		}
		if len(modelRes.Candidates) == 0 || len(modelRes.Candidates[0].Content.Parts) == 0 {
			return nil, fmt.Errorf("gemini response has no content")
		}

		raw := strings.TrimSpace(modelRes.Candidates[0].Content.Parts[0].Text)
		raw = trimCodeFence(raw)

		var parsed geminiStructuredReply
		if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
			return nil, fmt.Errorf("gemini returned invalid JSON: %w", err)
		}

		return &parsed, nil
	}

	if lastErr != nil {
		return nil, lastErr
	}

	return nil, fmt.Errorf("failed to call gemini")
}



func buildGeminiModelCandidates(preferred string) []string {
	fallbacks := []string{"gemini-2.0-flash", "gemini-1.5-flash"}
	result := make([]string, 0, len(fallbacks)+1)
	seen := map[string]struct{}{}

	addModel := func(model string) {
		name := strings.TrimSpace(model)
		if name == "" {
			return
		}
		if _, ok := seen[name]; ok {
			return
		}
		seen[name] = struct{}{}
		result = append(result, name)
	}

	addModel(preferred)
	for _, model := range fallbacks {
		addModel(model)
	}

	if len(result) == 0 {
		return []string{"gemini-2.0-flash"}
	}

	return result
}

func (s *AIChatService) callGeminiGenerateContent(ctx context.Context, model string, payload []byte) ([]byte, error) {
	requestURL := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", model, s.apiKey)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	httpRes, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call gemini: %w", err)
	}
	defer func() {
		if closeErr := httpRes.Body.Close(); closeErr != nil {
			log.Printf("⚠️ Failed to close Gemini response body: %v", closeErr)
		}
	}()

	resBody, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read gemini response: %w", err)
	}

	if httpRes.StatusCode >= 300 {
		trimmedBody := strings.TrimSpace(string(resBody))
		if len(trimmedBody) > 320 {
			trimmedBody = trimmedBody[:320]
		}
		return nil, &GeminiAPIError{StatusCode: httpRes.StatusCode, Body: trimmedBody}
	}

	return resBody, nil
}

func sanitizeYouTubeVideos(raw []geminiVideoReplyRow, limit int) []dto.AIVideoReview {
	if limit <= 0 {
		limit = 3
	}

	results := make([]dto.AIVideoReview, 0, limit)
	seen := make(map[string]struct{}, limit)

	for _, item := range raw {
		cleanURL := normalizeYouTubeURL(item.URL)
		if cleanURL == "" {
			continue
		}
		if _, exists := seen[cleanURL]; exists {
			continue
		}

		title := strings.TrimSpace(item.Title)
		if title == "" {
			title = "YouTube review"
		}

		results = append(results, dto.AIVideoReview{
			Title: title,
			URL:   cleanURL,
		})
		seen[cleanURL] = struct{}{}

		if len(results) >= limit {
			break
		}
	}

	return results
}

func normalizeYouTubeURL(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ""
	}

	parsed, err := url.Parse(trimmed)
	if err != nil {
		return ""
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return ""
	}

	host := strings.ToLower(strings.TrimPrefix(parsed.Hostname(), "www."))
	if host != "youtube.com" && host != "m.youtube.com" && host != "youtu.be" {
		return ""
	}

	if host == "youtube.com" || host == "m.youtube.com" {
		if parsed.Path == "" || parsed.Path == "/" {
			return ""
		}
	}

	return parsed.String()
}

func trimCodeFence(input string) string {
	trimmed := strings.TrimSpace(input)
	if strings.HasPrefix(trimmed, "```") {
		trimmed = strings.TrimPrefix(trimmed, "```json")
		trimmed = strings.TrimPrefix(trimmed, "```")
		trimmed = strings.TrimSuffix(trimmed, "```")
	}
	return strings.TrimSpace(trimmed)
}

func (s *AIChatService) buildCards(ids []string, limit int) []dto.RecommendedDeviceCard {
	cards := make([]dto.RecommendedDeviceCard, 0, limit)
	seen := map[string]struct{}{}

	for _, id := range ids {
		deviceID := strings.TrimSpace(id)
		if deviceID == "" {
			continue
		}
		if _, ok := seen[deviceID]; ok {
			continue
		}
		item, ok := s.deviceIndex[deviceID]
		if !ok {
			continue
		}
		seen[deviceID] = struct{}{}
		cards = append(cards, toCard(item))
		if len(cards) >= limit {
			break
		}
	}

	return cards
}

func toCard(device catalogDevice) dto.RecommendedDeviceCard {
	return dto.RecommendedDeviceCard{
		ID:        device.ID,
		BrandName: device.BrandName,
		ModelName: device.ModelName,
		ImageURL:  device.ImageURL,
		DetailURL: fmt.Sprintf("/devices/%s", device.ID),
		OS:        device.OS,
		Chipset:   device.Chipset,
		Memory:    device.Memory,
		Battery:   device.Battery,
		Price:     device.Price,
	}
}

func buildSingleYouTubeFallbackVideo(deviceName string) []dto.AIVideoReview {
	query := strings.TrimSpace(deviceName)
	if query == "" {
		query = "phone review"
	} else {
		query = query + " review"
	}

	return []dto.AIVideoReview{{
		Title: strings.TrimSpace(deviceName + " - YouTube review search"),
		URL:   "https://www.youtube.com/results?search_query=" + url.QueryEscape(query),
	}}
}
