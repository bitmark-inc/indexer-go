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
	ID              string          `json:"id,omitempty"`
	IndexID         string          `json:"indexID,omitempty"`
	Source          string          `json:"source"`
	ProjectMetadata ProjectMetadata `json:"projectMetadata"`
	Tokens          []Token         `json:"tokens,omitempty"`
}

type AssetInfo struct {
	ID              string                   `json:"id"`
	IndexID         string                   `json:"indexID"`
	Source          string                   `json:"source"`
	ProjectMetadata VersionedProjectMetadata `json:"projectMetadata"`
	Tokens          []Token                  `json:"tokens,omitempty"`
}

type VersionedProjectMetadata struct {
	Origin ProjectMetadata `json:"origin"`
	Latest ProjectMetadata `json:"latest"`
}

type AssetAttributes struct {
	Scrollable bool `json:"scrollable" bson:"scrollable"`
}

type Geolocation struct {
	Longitude float64 `json:"lon"`
	Latitude  float64 `json:"lat"`
}

type LocationInformation struct {
	ClaimedLocation Geolocation `json:"claimedLocation"`
	StampedLocation Geolocation `json:"stampedLocation,omitempty"`
}

type ArtworkMetadata struct {
	LastOwner           string                `json:"lastOwner"`
	IsStamped           bool                  `json:"isStamped"`
	LocationInformation []LocationInformation `json:"locationInformation"`
}

type Medium string

type ProjectMetadata struct {
	// Common attributes
	ArtistID            string `json:"artistID" structs:"artistID" bson:"artistID"`                                  // Artist blockchain address
	ArtistName          string `json:"artistName" structs:"artistName" bson:"artistName"`                            // <creator.user.username>,
	ArtistURL           string `json:"artistURL" structs:"artistURL" bson:"artistURL"`                               // <OpenseaAPI/creator.address>,
	AssetID             string `json:"assetID" structs:"assetID" bson:"assetID"`                                     // <asset_contract.address>,
	Title               string `json:"title" structs:"title" bson:"title"`                                           // <name>,
	Description         string `json:"description" structs:"description" bson:"description"`                         // <description>,
	MIMEType            string `json:"mimeType" structs:"mimeType" bson:"mimeType"`                                  // <mime_type from file extension or metadata>,
	Medium              Medium `json:"medium" structs:"medium" bson:"medium"`                                        // <"image" if image_url is present; "other" if animation_url is present> ,
	MaxEdition          int64  `json:"maxEdition" structs:"maxEdition" bson:"maxEdition"`                            // 0,
	BaseCurrency        string `json:"baseCurrency,omitempty" structs:"baseCurrency" bson:"baseCurrency"`            // null,
	BasePrice           int64  `json:"basePrice,omitempty" structs:"basePrice" bson:"basePrice"`                     // null,
	Source              string `json:"source" structs:"source" bson:"source"`                                        // <Opeasea/Artblock>,
	SourceURL           string `json:"sourceURL" structs:"sourceURL" bson:"sourceURL"`                               // <linktoSourceWebsite>,
	PreviewURL          string `json:"previewURL" structs:"previewURL" bson:"previewURL"`                            // <image_url or animation_url>,
	ThumbnailURL        string `json:"thumbnailURL" structs:"thumbnailURL" bson:"thumbnailURL"`                      // <image_thumbnail_url>,
	GalleryThumbnailURL string `json:"galleryThumbnailURL" structs:"galleryThumbnailURL" bson:"galleryThumbnailURL"` // <image_thumbnail_url>,
	AssetData           string `json:"assetData" structs:"assetData" bson:"assetData"`                               // null,
	AssetURL            string `json:"assetURL" structs:"assetURL" bson:"assetURL"`                                  // <permalink>

	// autonomy customized attributes
	Attributes *AssetAttributes `json:"attributes,omitempty" bson:"attributes,omitempty"`

	// artwork metadata from source. currently on for Feral File
	ArtworkMetadata ArtworkMetadata `json:"artworkMetadata" structs:"artworkMetadata" bson:"artworkMetadata"`

	// Operation attributes
	LastUpdatedAt time.Time `json:"lastUpdatedAt" structs:"lastUpdatedAt" bson:"lastUpdatedAt"`

	// Feral File attributes
	InitialSaleModel string `json:"initialSaleModel" structs:"initialSaleModel" bson:"initialSaleModel"` // airdrop|fix-price|highest-bid-auction|group-auction

	// Deprecated attributes
	OriginalFileURL string `json:"originalFileURL" structs:"-" bson:"-"`
}

type Token struct {
	ID              string    `json:"id"`
	Fungible        bool      `json:"fungible"`
	Blockchain      string    `json:"blockchain"`
	ContractType    string    `json:"contractType"`
	MintedAt        time.Time `json:"mintedAt"`
	ContractAddress string    `json:"contractAddress"`
	Owner           string    `json:"owner"`
	AssetID         string    `json:"assetID"`
	IndexID         string    `json:"indexID"`
	Balance         int64     `json:"balance"`
	Source          string    `json:"source"`
}

type NFTQuery struct {
	IDs []string `json:"ids"`
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

type Test struct {
	ID         string `json:"id"`
	Blockchain string `json:"blockchain"`
}

func (c *Client) GetAssetInfo(tokenID string) ([]AssetInfo, error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   "/v1/nft/query/",
	}

	queryIDs := NFTQuery{
		IDs: []string{tokenID},
	}
	body, err := json.Marshal(queryIDs)
	if err != nil {
		return []AssetInfo{}, err
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(body))
	if err != nil {
		return []AssetInfo{}, err
	}

	req.Header.Set("API-TOKEN", c.APIToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return []AssetInfo{}, err
	}
	defer resp.Body.Close()

	var indexAssetRequest []AssetInfo
	err = json.NewDecoder(resp.Body).Decode(&indexAssetRequest)
	if err != nil {
		return []AssetInfo{}, err
	}

	return indexAssetRequest, nil
}
