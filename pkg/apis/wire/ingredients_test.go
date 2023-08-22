package wire

import (
	"encoding/json"
	"github.com/cocktailrobots/openbar-server/pkg/db/cocktailsdb"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIngredients(t *testing.T) {
	ingredients := []cocktailsdb.Ingredient{
		{
			Name:        "gin",
			DisplayName: "Gin",
			Description: ptr("It's Gin. You know what Gin is."),
		},
		{
			Name:        "sweet_vermouth",
			DisplayName: "Sweet Vermouth",
			Description: ptr("An aromatized fortified wine, flavoured with various botanicals and sometimes colored"),
		},
	}

	ingredientsWire := FromDbIngredients(ingredients)
	data, err := json.Marshal(ingredientsWire)
	require.NoError(t, err)
	require.NoError(t, ingredientsWire.Validate())

	var decoded Ingredients
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	ingredients2 := decoded.ToDbIngredients()
	require.Equal(t, ingredients, ingredients2)
}
