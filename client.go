package indexer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type IndexAssetRequest struct {
	Source          string          `json:"source"`
	ProjectMetadata ProjectMetadata `json:"projectMetadata"`
	Tokens          []Token         `json:"tokens"`
}

type ProjectMetadata struct {
	Source              string                 `json:"source"`
	AssetID             string                 `json:"assetID"`
	ThumbnailURL        string                 `json:"thumbnailURL"`
	GalleryThumbnailURL string                 `json:"galleryThumbnailURL"`
	PreviewURL          string                 `json:"previewURL"`
	ArtworkMetadata     map[string]interface{} `json:"artworkMetadata"`
}

type Token struct {
	ID              string    `json:"id"`
	Blockchain      string    `json:"blockchain"`
	ContractType    string    `json:"contractType"`
	MintedAt        time.Time `json:"mintedAt"`
	ContractAddress string    `json:"contractAddress"`
	Owner           string    `json:"owner"`
}

type Client struct {
	host     string
	APIToken string

	client *http.Client
}

func NewIndexerClient(host string, APIToken string) *Client {
	return &Client{
		APIToken: APIToken,
		host:     host,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// IndexAsset updates an asset
func (c *Client) IndexAsset(asset IndexAssetRequest, assetID string) error {
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   fmt.Sprintf("/asset/%s", assetID),
	}

	body, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, u.String(), bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("API-TOKEN", c.APIToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	// for test
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("indexer returned status code %d with message: %s", resp.StatusCode, resp.Body)
	}

	return nil
}
