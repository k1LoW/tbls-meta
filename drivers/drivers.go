package drivers

import (
	"context"

	"github.com/k1LoW/tbls/schema"
)

type Diff struct {
	From        string
	To          string
	FromContent string
	ToContent   string
}

type Diffs []Diff

// Driver is the common interface for database drivers
type Driver interface {
	Plan(ctx context.Context, from, to *schema.Schema) (Diffs, error)
	Apply(ctx context.Context, from, to *schema.Schema) error
}
