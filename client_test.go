package indexer

import (
	"fmt"
	"github.com/spf13/viper"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/bitmark-inc/config-loader"
)

// TestIndexAsset is a test function for IndexAsset
func TestIndexAsset(t *testing.T) {
	config.LoadConfig("NFT_INDEXER")

	client := NewIndexerClient("indexer.test.autonomy.io", viper.GetString("indexer.api_token"))
	contract := "KT1MeB8Wntrx4fjksZkCWUwmGDQTGs6DsMwp"
	var tokenID int64 = 1

	assetID := fmt.Sprintf("%s-%d", contract, tokenID)

	request := IndexAssetRequest{
		Source: "autonomy-postcard",
		ProjectMetadata: ProjectMetadata{
			Source:              "autonomy-postcard",
			AssetID:             "<contract>-<token_id>",
			ThumbnailURL:        "https://cdn.test.feralfileassets.com/thumbnails/5d7ef388-c2ea-417b-8288-6128be4e5dca/1669015364",
			GalleryThumbnailURL: "https://cdn.test.feralfileassets.com/gallery-thumbnails/5d7ef388-c2ea-417b-8288-6128be4e5dca/1669015364",
			PreviewURL:          "https://cdn.test.feralfileassets.com/previews/5d7ef388-c2ea-417b-8288-6128be4e5dca/1669015365/preview.jpeg",
		},
		Tokens: []Token{
			{
				ID:              strconv.FormatInt(tokenID, 10),
				Blockchain:      "tezos",
				ContractType:    "fa2",
				MintedAt:        time.Now(),
				ContractAddress: contract,
				Owner:           "<owner>",
			},
		},
	}

	err := client.IndexAsset(request, assetID)

	assert.NoError(t, err)
}
