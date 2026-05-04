package dataloader

import (
	"context"
	"log"
	"time"

	"github.com/vikstrous/dataloadgen"

	"github.com/TranTheTuan/dataloader-demo/graph/model"
	"github.com/TranTheTuan/dataloader-demo/internal/core/ports"
)

// Loaders holds all DataLoader instances for a single request.
//
// DataLoaders MUST be request-scoped because:
//  1. Their internal cache is only valid for the lifetime of a single request.
//  2. Different requests may have different auth contexts / data visibility.
//  3. Stale cache across requests would cause data inconsistency.
//
// Each Loader batches calls that happen within a short time window (WithWait)
// and deduplicates keys, then calls the batch function ONCE with all collected keys.
type Loaders struct {
	// RedemptionDrugsByRedemptionID batches and caches drug lookups by redemption_id.
	// When gqlgen resolves `redemption_drugs` for 10 RedemptionInfo items concurrently,
	// instead of 10 separate SQL queries, this loader collects all 10 redemption_ids
	// and calls GetRedemptionDrugsByRedemptionIDs ONCE.
	RedemptionDrugsByRedemptionID *dataloadgen.Loader[string, []*model.RedemptionDrugDetail]

	// AcknowledgeDocsByRedemptionID batches and caches document lookups by redemption_id.
	AcknowledgeDocsByRedemptionID *dataloadgen.Loader[string, []*model.AcknowledgeDocumentDetail]
}

// NewLoaders creates a new Loaders instance with all DataLoaders configured.
// Call this once per HTTP request (done automatically by the Middleware).
func NewLoaders(redemptionRepo ports.RedemptionRepository) *Loaders {
	return &Loaders{
		RedemptionDrugsByRedemptionID: dataloadgen.NewMappedLoader(
			func(ctx context.Context, keys []string) (map[string][]*model.RedemptionDrugDetail, error) {
				log.Printf("[dataloader] Batching drug fetch for %d redemption IDs", len(keys))
				return redemptionRepo.GetRedemptionDrugsByRedemptionIDs(ctx, keys)
			},
			// WithWait controls how long the loader waits to collect keys before
			// executing the batch function. 2ms is a good default — long enough
			// to collect concurrent resolver calls, short enough to not add
			// noticeable latency.
			dataloadgen.WithWait(2*time.Millisecond),
		),

		AcknowledgeDocsByRedemptionID: dataloadgen.NewMappedLoader(
			func(ctx context.Context, keys []string) (map[string][]*model.AcknowledgeDocumentDetail, error) {
				log.Printf("[dataloader] Batching document fetch for %d redemption IDs", len(keys))
				return redemptionRepo.GetAcknowledgeDocumentsByRedemptionIDs(ctx, keys)
			},
			dataloadgen.WithWait(2*time.Millisecond),
		),
	}
}
