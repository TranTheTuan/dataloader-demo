package repositories

import (
	"context"
	"fmt"
	"log"

	"github.com/example/dpc-dataloader-demo/graph/model"
	"github.com/example/dpc-dataloader-demo/internal/core/ports"
)

// Compile-time check that RedemptionPortImpl implements RedemptionRepository.
var _ ports.RedemptionRepository = (*RedemptionPortImpl)(nil)

type RedemptionPortImpl struct {
	// In the real project this holds *sqlx.DB and *postgres.DBConfig.
	// Here we use in-memory fake data for demonstration.
}

func NewRedemptionPort() *RedemptionPortImpl {
	return &RedemptionPortImpl{}
}

// ─── GetPaginatedRedemptionList ────────────────────────────────────
// AFTER DataLoader refactor: this method NO LONGER fetches drugs or documents.
// It only returns the flat redemption rows. The nested fields are nil here
// and will be populated by gqlgen field resolvers using DataLoader.
func (r *RedemptionPortImpl) GetPaginatedRedemptionList(ctx context.Context, order string, limit int, offset int, programID *string) (*model.PaginatedRedemptionListResponse, error) {
	log.Printf("[repo] GetPaginatedRedemptionList called (order=%s, limit=%d, offset=%d)", order, limit, offset)

	// Simulate a DB query that returns redemption rows WITHOUT nested data.
	// In the real project this is the main SQL query with JOINs on redemption,
	// patient_enrolment, patient, program, etc. — but NOT on redemption_drug.
	items := make([]*model.RedemptionInfo, 0)
	for i := offset; i < offset+limit && i < 6; i++ {
		rid := fmt.Sprintf("RD-%04d", i+1)
		pid := fmt.Sprintf("P-%04d", i+1)
		eid := fmt.Sprintf("E-%04d", i+1)
		pname := fmt.Sprintf("Patient %d", i+1)
		progID := "PROG-001"
		progName := "Demo Program"
		status := "Approved"
		cycle := i + 1

		items = append(items, &model.RedemptionInfo{
			ID:                  intPtr(i + 1),
			RedemptionID:        &rid,
			PatientID:           &pid,
			EnrolmentID:         &eid,
			PatientName:         &pname,
			ProgramID:           &progID,
			ProgramName:         &progName,
			RedemptionStatus:    &status,
			PurchaseCycle:       &cycle,
			// NOTE: RedemptionDrugs and AcknowledgeDocuments are LEFT NIL.
			// They will be resolved by the DataLoader field resolvers.
		})
	}

	return &model.PaginatedRedemptionListResponse{
		Count: 6,
		Items: items,
	}, nil
}

// ─── GetRedemptionDetails ──────────────────────────────────────────
func (r *RedemptionPortImpl) GetRedemptionDetails(ctx context.Context, redemptionID string) (*model.RedemptionDetailsResponse, error) {
	log.Printf("[repo] GetRedemptionDetails called (id=%s)", redemptionID)

	pname := "Patient Detail"
	progName := "Demo Program"
	status := "Approved"

	return &model.RedemptionDetailsResponse{
		RedemptionID:     &redemptionID,
		PatientName:      &pname,
		ProgramName:      &progName,
		RedemptionStatus: &status,
		// NOTE: RedemptionDrugs and AcknowledgeDocuments are LEFT NIL.
	}, nil
}

// ─── Batch Methods for DataLoader ──────────────────────────────────
// These are the NEW methods. Each executes a SINGLE SQL query for ALL
// requested IDs, replacing the old per-item or inline batch approach.

// GetRedemptionDrugsByRedemptionIDs fetches drugs for multiple redemptions
// in a single query using WHERE redemption_id = ANY($1).
func (r *RedemptionPortImpl) GetRedemptionDrugsByRedemptionIDs(ctx context.Context, redemptionIDs []string) (map[string][]*model.RedemptionDrugDetail, error) {
	log.Printf("[repo] GetRedemptionDrugsByRedemptionIDs called with %d IDs: %v", len(redemptionIDs), redemptionIDs)

	// In the real project, this would be:
	//
	//   query := `SELECT dr.id as drug_id, dr.drug_name, dr.drug_dosage, dr.drug_unit,
	//             sku.id as drug_sku_code_id, sku.sku_name, sku.sku_description,
	//             rd.redemption_quantity, rd.redemption_id
	//             FROM schema.redemption_drug rd
	//             JOIN schema.drug dr ON rd.drug_id = dr.id
	//             JOIN schema.ml_drug_sku sku ON rd.drug_sku_id = sku.id
	//             WHERE rd.redemption_id = any($1)`
	//   err := db.SelectContext(ctx, &drugs, query, pq.Array(redemptionIDs))

	result := make(map[string][]*model.RedemptionDrugDetail)
	for _, rid := range redemptionIDs {
		// Simulate 2 drugs per redemption
		result[rid] = []*model.RedemptionDrugDetail{
			{
				RedemptionID:       rid,
				DrugID:             intPtr(1),
				DrugName:           strPtr("Aspirin 100mg"),
				DrugDosage:         strPtr("100mg"),
				DrugUnit:           strPtr("tablet"),
				DrugSkuCodeID:      intPtr(101),
				SkuName:            strPtr("ASP-100"),
				SkuDescription:     strPtr("Aspirin 100mg Tablet"),
				RedemptionQuantity: intPtr(30),
			},
			{
				RedemptionID:       rid,
				DrugID:             intPtr(2),
				DrugName:           strPtr("Metformin 500mg"),
				DrugDosage:         strPtr("500mg"),
				DrugUnit:           strPtr("tablet"),
				DrugSkuCodeID:      intPtr(102),
				SkuName:            strPtr("MET-500"),
				SkuDescription:     strPtr("Metformin 500mg Tablet"),
				RedemptionQuantity: intPtr(60),
			},
		}
	}

	return result, nil
}

// GetAcknowledgeDocumentsByRedemptionIDs fetches acknowledge documents for
// multiple redemptions in a single query.
func (r *RedemptionPortImpl) GetAcknowledgeDocumentsByRedemptionIDs(ctx context.Context, redemptionIDs []string) (map[string][]*model.AcknowledgeDocumentDetail, error) {
	log.Printf("[repo] GetAcknowledgeDocumentsByRedemptionIDs called with %d IDs: %v", len(redemptionIDs), redemptionIDs)

	// In the real project, this would be:
	//
	//   query := `SELECT doc.id, doc.file_name, doc.file_url, doc.remark AS remarks,
	//             u.fullname as created_by, doc.created_at, doc.redemption_id
	//             FROM schema.redemption_acknowledge_document doc
	//             LEFT JOIN schema.user_info u ON doc.created_by = u.id
	//             WHERE doc.redemption_id = any($1)`
	//   err := db.SelectContext(ctx, &docs, query, pq.Array(redemptionIDs))

	result := make(map[string][]*model.AcknowledgeDocumentDetail)
	for _, rid := range redemptionIDs {
		result[rid] = []*model.AcknowledgeDocumentDetail{
			{
				ID:           intPtr(1),
				FileName:     "acknowledgement.pdf",
				FileURL:      "https://storage.example.com/docs/" + rid + "/ack.pdf",
				Remarks:      strPtr("Signed by patient"),
				CreatedBy:    strPtr("Admin User"),
				CreatedAt:    strPtr("2024-01-15T10:30:00Z"),
				RedemptionID: rid,
			},
		}
	}

	return result, nil
}

// ─── Helpers ───────────────────────────────────────────────────────

func intPtr(i int) *int       { return &i }
func strPtr(s string) *string { return &s }
