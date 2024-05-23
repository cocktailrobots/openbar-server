package wire

import (
	"fmt"
	"github.com/cocktailrobots/openbar-server/pkg/db/openbardb"
	"github.com/cocktailrobots/openbar-server/pkg/util"
	"strings"
)

// Fluid is a fluid as it will be written to the wire in HTTP responses.
type Fluid struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Fluids is a slice of fluids.
type Fluids []Fluid

// ToDbFluids converts a list of Fluids to a list of openbardb.Fluids.
func (f Fluids) ToDbFluids() []openbardb.Fluid {
	fluids := make([]openbardb.Fluid, len(f))
	for i := range f {
		fluids[i] = openbardb.Fluid{
			Idx:   i,
			Fluid: &f[i].ID,
		}
	}

	return fluids
}

func (f Fluids) Validate() error {
	for i := range f {
		f[i].ID = strings.TrimSpace(f[i].ID)
		if len(f[i].ID) == 0 {
			return fmt.Errorf("fluid id cannot be empty")
		}

		f[i].Name = strings.TrimSpace(f[i].Name)
		if len(f[i].Name) == 0 {
			return fmt.Errorf("fluid name cannot be empty")
		}
	}

	return nil
}

// FromDbFluids converts a list of openbardb.Fluids to a list of Fluids.
func FromDbFluids(fluids []openbardb.Fluid) Fluids {
	f := make(Fluids, len(fluids))
	for _, fluid := range fluids {
		if fluid.Idx < 0 || fluid.Idx >= len(fluids) {
			continue
		}

		if fluid.Fluid == nil {
			f[fluid.Idx] = Fluid{
				ID:   "empty",
				Name: "Empty",
			}
		} else {
			name := util.ReplaceChars(*fluid.Fluid, map[rune]rune{'_': ' ', '-': ' '})
			f[fluid.Idx] = Fluid{
				ID:   *fluid.Fluid,
				Name: util.TitleCase(name),
			}
		}
	}

	return f
}
