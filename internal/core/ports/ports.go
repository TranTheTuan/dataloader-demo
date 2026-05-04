package ports

import (
	"context"

	"github.com/TranTheTuan/dataloader-demo/graph/model"
)

// RedemptionRepository defines the data access interface for redemptions.
// This mirrors the original project's ports.RedemptionRepository.
type RedemptionRepository interface {
	// GetPaginatedRedemptionList returns a paginated list of redemptions.
	// After DataLoader refactor: this method only returns the flat redemption data.
	// Nested fields (drugs, documents) are resolved lazily via DataLoader field resolvers.
	GetPaginatedRedemptionList(ctx context.Context, order string, limit int, offset int, programID *string) (*model.PaginatedRedemptionListResponse, error)

	// GetRedemptionDetails returns details for a single redemption.
	// After DataLoader refactor: nested fields are resolved via field resolvers.
	GetRedemptionDetails(ctx context.Context, redemptionID string) (*model.RedemptionDetailsResponse, error)

	// ── Batch methods for DataLoader ──────────────────────────────────
	// These are the NEW methods added to support DataLoader batching.

	// GetRedemptionDrugsByRedemptionIDs fetches drugs for multiple redemptions in one query.
	// Returns a map of redemption_id -> drugs.
	GetRedemptionDrugsByRedemptionIDs(ctx context.Context, redemptionIDs []string) (map[string][]*model.RedemptionDrugDetail, error)

	// GetAcknowledgeDocumentsByRedemptionIDs fetches documents for multiple redemptions in one query.
	// Returns a map of redemption_id -> documents.
	GetAcknowledgeDocumentsByRedemptionIDs(ctx context.Context, redemptionIDs []string) (map[string][]*model.AcknowledgeDocumentDetail, error)
}
