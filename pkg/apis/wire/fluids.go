package wire

import (
	"fmt"
	"github.com/cocktailrobots/openbar-server/pkg/db"
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

// ToDbFluids converts a list of Fluids to a list of db.Fluids.
func (f Fluids) ToDbFluids() []db.Fluid {
	fluids := make([]db.Fluid, len(f))
	for i := range f {
		fluids[i] = db.Fluid{
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

// FromDbFluids converts a list of db.Fluids to a list of Fluids.
func FromDbFluids(fluids []db.Fluid) Fluids {
	f := make(Fluids, len(fluids))
	for _, fluid := range fluids {
		if fluid.Idx < 0 || fluid.Idx >= len(fluids) {
			continue
		}

		name := util.ReplaceChars(*fluid.Fluid, map[rune]rune{'_': ' ', '-': ' '})
		f[fluid.Idx] = Fluid{
			ID:   *fluid.Fluid,
			Name: util.TitleCase(name),
		}
	}

	return f
}
