package service

import (
	"context"
	"log"

	"github.com/example/dpc-dataloader-demo/graph/model"
	"github.com/example/dpc-dataloader-demo/internal/core/ports"
)

// RedemptionSvc defines the service interface.
// This mirrors the original project's service.RedemptionSvc.
type RedemptionSvc interface {
	GetPaginatedRedemptionList(ctx context.Context, input model.PaginatedRedemptionListInput) (*model.PaginatedRedemptionListResponse, error)
	GetRedemptionDetails(ctx context.Context, redemptionID string) (*model.RedemptionDetailsResponse, error)
}

// Compile-time check.
var _ RedemptionSvc = (*RedemptionSvcImpl)(nil)

type RedemptionSvcImpl struct {
	redemptionRepo ports.RedemptionRepository
}

func NewRedemptionSvc(redemptionRepo ports.RedemptionRepository) *RedemptionSvcImpl {
	return &RedemptionSvcImpl{
		redemptionRepo: redemptionRepo,
	}
}

// GetPaginatedRedemptionList handles business logic for the paginated list.
//
// BEFORE DataLoader: The repository method fetched redemptions AND their drugs
// and documents all in one call (3 SQL queries inside one method).
//
// AFTER DataLoader: The repository only fetches the flat redemption rows.
// Drugs and documents are resolved lazily by gqlgen field resolvers via DataLoader.
// This means if the client doesn't request `redemption_drugs` in their query,
// the drug SQL query is NEVER executed.
func (s *RedemptionSvcImpl) GetPaginatedRedemptionList(ctx context.Context, input model.PaginatedRedemptionListInput) (*model.PaginatedRedemptionListResponse, error) {
	log.Println("[service] GetPaginatedRedemptionList called")

	order := "created_at DESC"
	if input.Order != nil {
		order = *input.Order
	}
	limit := 10
	if input.Limit != nil {
		limit = *input.Limit
	}
	offset := 0
	if input.Offset != nil {
		offset = *input.Offset
	}

	response, err := s.redemptionRepo.GetPaginatedRedemptionList(ctx, order, limit, offset, input.ProgramID)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetRedemptionDetails returns details for a single redemption.
func (s *RedemptionSvcImpl) GetRedemptionDetails(ctx context.Context, redemptionID string) (*model.RedemptionDetailsResponse, error) {
	log.Println("[service] GetRedemptionDetails called")

	return s.redemptionRepo.GetRedemptionDetails(ctx, redemptionID)
}
