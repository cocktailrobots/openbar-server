package wire

import (
	"encoding/json"
	"github.com/cocktailrobots/openbar-server/pkg/db"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCocktails(t *testing.T) {
	cocktails := []db.Cocktail{
		{
			Name:        "americano",
			DisplayName: "Americano",
			Description: ptr("A classic cocktail"),
		},
		{
			Name:        "negroni",
			DisplayName: "Negroni",
			Description: ptr("A classic cocktail"),
		},
	}

	cocktailsWire := FromDbCocktails(cocktails)
	data, err := json.Marshal(cocktailsWire)
	require.NoError(t, err)
	require.NoError(t, cocktailsWire.Validate())

	var decoded Cocktails
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	cocktails2 := decoded.ToDbCocktails()
	require.Equal(t, cocktails, cocktails2)
}
