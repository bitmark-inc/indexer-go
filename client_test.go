package indexer

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/spf13/viper"

	"github.com/stretchr/testify/assert"

	"github.com/bitmark-inc/config-loader"
)

// TestIndexAsset is a test function for IndexAsset
func TestIndexAsset(t *testing.T) {
	config.LoadConfig("AUTONOMY_POSTCARD")
	var client = NewIndexerClient("indexer.test.autonomy.io", viper.GetString("indexer.api_token"))

	contract := "KT1MeB8Wntrx4fjksZkCWUwmGDQTGs6DsMwp"
	var tokenID int64 = 1679834315868

	assetID := fmt.Sprintf("%s-%d", contract, tokenID)

	asset := IndexAssetRequest{
		Source: "autonomy-postcard",
		ProjectMetadata: ProjectMetadata{
			Source:              "autonomy-postcard",
			AssetID:             "<contract>-<token_id>",
			ThumbnailURL:        "https://cdn.test.feralfileassets.com/thumbnails/5d7ef388-c2ea-417b-8288-6128be4e5dca/1669015364",
			GalleryThumbnailURL: "https://cdn.test.feralfileassets.com/gallery-thumbnails/5d7ef388-c2ea-417b-8288-6128be4e5dca/1669015364",
			PreviewURL:          "https://cdn.test.feralfileassets.com/previews/5d7ef388-c2ea-417b-8288-6128be4e5dca/1669015365/preview.jpeg",
			ArtworkMetadata: ArtworkMetadata{
				LastOwner: "this is the last owner",
				IsStamped: true,
				LocationInformation: []LocationInformation{
					{
						ClaimedLocation: Geolocation{
							Longitude: 12.34,
							Latitude:  56.78,
						},
						StampedLocation: Geolocation{
							Longitude: 12.34,
							Latitude:  56.78,
						},
					},
					{
						ClaimedLocation: Geolocation{
							Longitude: 123.34,
							Latitude:  456.78,
						},
						StampedLocation: Geolocation{
							Longitude: 112.34,
							Latitude:  56.78,
						},
					},
				},
			},
		},
		Tokens: []Token{
			{
				ID:              strconv.FormatInt(tokenID, 10),
				AssetID:         assetID,
				IndexID:         fmt.Sprintf("tez-%s", assetID),
				Blockchain:      "tezos",
				ContractType:    "fa2",
				MintedAt:        time.Now(),
				ContractAddress: contract,
				Owner:           "<test>",
				Balance:         10,
			},
		},
	}

	err := client.IndexAsset(asset, assetID)

	assert.NoError(t, err)
}

func TestGetAssetInfo(t *testing.T) {
	config.LoadConfig("AUTONOMY_POSTCARD")
	var client = NewIndexerClient("indexer.test.autonomy.io", viper.GetString("indexer.api_token"))

	contract := "KT1MeB8Wntrx4fjksZkCWUwmGDQTGs6DsMwp"
	indexID := "1679815146233"
	tokenID := fmt.Sprintf("tez-%s-%s", contract, indexID)
	assetInfo, err := client.GetAssetInfo(tokenID)
	fmt.Printf("assetInfo%+v\n\n", assetInfo[0])

	assert.NoError(t, err)
}
