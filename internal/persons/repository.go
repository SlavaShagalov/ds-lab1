package persons

import (
	"context"

	"github.com/SlavaShagalov/ds-lab1/internal/models"
)

type CreateParams struct {
	Name    string
	Age     int
	Address string
	Work    string
}

type PartialUpdateParams struct {
	ID      int64
	Name    *string
	Age     *int
	Address *string
	Work    *string
}

type Repository interface {
	//HealthCheck(ctx context.Context) error
	Create(ctx context.Context, params *CreateParams) (*models.Person, error)
	Get(ctx context.Context, personID int64) (*models.Person, error)
	List(ctx context.Context, offset, limit int64) ([]models.Person, error)
	PartialUpdate(ctx context.Context, person *PartialUpdateParams) (*models.Person, error)
	Delete(ctx context.Context, personID int64) error
}
