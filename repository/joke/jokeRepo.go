package joke

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"instagram/util/httpx"
)

const endpoint = "https://api.api-ninjas.com/v1/jokes?limit=1"

type Repo interface {
	FetchJoke(ctx context.Context) (string, error)
}

type repo struct {
	apiKey string
	client *http.Client
}

func New(apiKey string) Repo {
	return &repo{
		apiKey: apiKey,
		client: httpx.Client(),
	}
}

func (r *repo) FetchJoke(ctx context.Context) (string, error) {
	if r.apiKey == "" {
		return "", errors.New("API_NINJAS_KEY is empty")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("X-Api-Key", r.apiKey)

	resp, err := r.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("joke api non-200 response")
	}

	var arr []struct {
		Joke string `json:"joke"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&arr); err != nil {
		return "", err
	}
	if len(arr) == 0 || arr[0].Joke == "" {
		return "", errors.New("no joke")
	}
	return arr[0].Joke, nil
}
