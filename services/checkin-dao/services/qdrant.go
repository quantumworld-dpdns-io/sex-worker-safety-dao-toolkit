package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type QdrantClient struct {
	baseURL  string
	httpC    *http.Client
	mockMode bool
	mu       sync.RWMutex
	mockData []MockPoint
}

type MockPoint struct {
	ID      string                 `json:"id"`
	Vector  []float32              `json:"vector"`
	Payload map[string]interface{} `json:"payload"`
}

func NewQdrantClient(host, port string) *QdrantClient {
	return &QdrantClient{
		baseURL: fmt.Sprintf("http://%s:%s", host, port),
		httpC:   &http.Client{Timeout: 5 * time.Second},
	}
}

type qdrantUpsertRequest struct {
	Points []qdrantPoint `json:"points"`
}

type qdrantPoint struct {
	ID      string                 `json:"id"`
	Vector  []float32              `json:"vector"`
	Payload map[string]interface{} `json:"payload"`
}

type qdrantSearchRequest struct {
	Vector  []float32 `json:"vector"`
	Limit   int       `json:"limit"`
	WithPayload bool `json:"with_payload"`
}

type qdrantSearchResponse struct {
	Result []qdrantScoredPoint `json:"result"`
}

type qdrantScoredPoint struct {
	ID      string                 `json:"id"`
	Score   float64                `json:"score"`
	Payload map[string]interface{} `json:"payload"`
}

func (c *QdrantClient) StoreEmbedding(collection string, id string, vector []float32, payload map[string]interface{}) error {
	if c.mockMode {
		c.mu.Lock()
		defer c.mu.Unlock()
		c.mockData = append(c.mockData, MockPoint{ID: id, Vector: vector, Payload: payload})
		log.Debug().Str("collection", collection).Str("id", id).Msg("qdrant mock: embedding stored")
		return nil
	}

	body := qdrantUpsertRequest{
		Points: []qdrantPoint{
			{ID: id, Vector: vector, Payload: payload},
		},
	}

	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshaling upsert request: %w", err)
	}

	req, err := http.NewRequest("PUT",
		fmt.Sprintf("%s/collections/%s/points", c.baseURL, collection),
		bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("creating upsert request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpC.Do(req)
	if err != nil {
		log.Warn().Err(err).Msg("qdrant unavailable, falling back to mock mode")
		c.mockMode = true
		return c.StoreEmbedding(collection, id, vector, payload)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("qdrant upsert error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	log.Debug().Str("collection", collection).Str("id", id).Msg("qdrant embedding stored")
	return nil
}

func (c *QdrantClient) SearchSimilar(collection string, vector []float32, limit int) ([]MockPoint, error) {
	if c.mockMode {
		c.mu.RLock()
		defer c.mu.RUnlock()

		var results []MockPoint
		for i, p := range c.mockData {
			sim := cosineSimilarity(vector, p.Vector)
			if sim > 0.7 {
				results = append(results, c.mockData[i])
			}
		}
		if len(results) > limit {
			results = results[:limit]
		}
		return results, nil
	}

	body := qdrantSearchRequest{
		Vector:      vector,
		Limit:       limit,
		WithPayload: true,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshaling search request: %w", err)
	}

	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s/collections/%s/points/search", c.baseURL, collection),
		bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("creating search request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpC.Do(req)
	if err != nil {
		log.Warn().Err(err).Msg("qdrant unavailable, falling back to mock mode")
		c.mockMode = true
		return c.SearchSimilar(collection, vector, limit)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("qdrant search error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var searchResp qdrantSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("decoding search response: %w", err)
	}

	var results []MockPoint
	for _, sp := range searchResp.Result {
		results = append(results, MockPoint{
			ID:      sp.ID,
			Payload: sp.Payload,
		})
	}

	return results, nil
}

func cosineSimilarity(a, b []float32) float64 {
	if len(a) != len(b) {
		return 0
	}
	var dot, normA, normB float64
	for i := range a {
		dot += float64(a[i] * b[i])
		normA += float64(a[i] * a[i])
		normB += float64(b[i] * b[i])
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (sqrt(normA) * sqrt(normB))
}

func sqrt(x float64) float64 {
	z := x / 2.0
	for i := 0; i < 10; i++ {
		z -= (z*z - x) / (2 * z)
	}
	return z
}
