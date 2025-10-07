package receipt_repo

import (
	"context"
	"database/sql"
	"spending/models"

	"go.opentelemetry.io/otel"
)

func (repo *receiptRepository) LoadReceiptItems(ctx context.Context, tx *sql.Tx, receipt *models.Receipt) error {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(ctx, "DB:LoadReceiptItems")
	defer span.End()

	if receipt == nil {
		return nil
	}

	items, err := repo.receipt_item_repo.GetItemsByReceiptId(ctx, tx, receipt.Id)
	if err != nil {
		return err
	}

	receipt.Items = items
	return nil
}

func (repo *receiptRepository) LoadReceiptsItems(ctx context.Context, tx *sql.Tx, receipts []*models.Receipt) error {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(ctx, "DB:LoadReceiptsItems")
	defer span.End()

	if len(receipts) == 0 {
		return nil
	}

	receiptIds := make([]int, 0, len(receipts))
	receiptById := make(map[int]*models.Receipt, len(receipts))
	for _, receipt := range receipts {
		receiptIds = append(receiptIds, receipt.Id)
		receiptById[receipt.Id] = receipt
	}

	itemsByReceiptId, err := repo.receipt_item_repo.GetItemsByReceiptIds(ctx, tx, receiptIds)
	if err != nil {
		return err
	}

	for receiptId, items := range itemsByReceiptId {
		if receipt, ok := receiptById[receiptId]; ok {
			receipt.Items = items
		}
	}

	return nil
}
