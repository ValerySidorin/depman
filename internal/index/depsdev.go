package index

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Resource struct {
	Name string `json:"name"`
}

type Response struct {
	Results []Resource `json:"results"`
}

type DepsDev struct{}

func (d *DepsDev) Search(query string) ([]string, error) {
	url := "https://deps.dev/_/search/suggest?system=go&kind=package&q=" + query
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("get packages: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read data: %w", err)
	}

	var res Response
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	packages := make([]string, 0, len(res.Results))
	for _, p := range res.Results {
		packages = append(packages, p.Name)
	}

	return packages, nil
}
