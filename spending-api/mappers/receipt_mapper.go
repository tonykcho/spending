package mappers

import (
	"spending/dto"
	"spending/models"
)

func MapReceipt(receipt *models.Receipt) *dto.ReceiptDto {
	if receipt == nil {
		return nil
	}

	dto := &dto.ReceiptDto{
		Id:        receipt.UUId,
		StoreName: receipt.StoreName,
		Total:     receipt.Total,
		Date:      receipt.Date,
		CreatedAt: receipt.CreatedAt,
		UpdatedAt: receipt.UpdatedAt,
	}

	if receipt.Items != nil {
		dto.Items = MapReceiptItems(receipt.Items)
	}

	return dto
}

func MapReceipts(receiptList []*models.Receipt) []*dto.ReceiptDto {
	var dtoList []*dto.ReceiptDto = make([]*dto.ReceiptDto, 0)

	for _, receipt := range receiptList {
		dto := MapReceipt(receipt)
		dtoList = append(dtoList, dto)
	}
	return dtoList
}
