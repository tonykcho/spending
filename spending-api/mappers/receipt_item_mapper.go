package mappers

import (
	"spending/dto"
	"spending/models"
)

func MapReceiptItem(receiptItem *models.ReceiptItem) *dto.ReceiptItemDto {
	if receiptItem == nil {
		return nil
	}

	dto := &dto.ReceiptItemDto{
		Id:        receiptItem.UUId,
		Name:      receiptItem.Name,
		Price:     receiptItem.Price,
		CreatedAt: receiptItem.CreatedAt,
		UpdatedAt: receiptItem.UpdatedAt,
	}

	return dto
}

func MapReceiptItems(receiptItems []*models.ReceiptItem) []*dto.ReceiptItemDto {
	var dtoList []*dto.ReceiptItemDto = make([]*dto.ReceiptItemDto, 0)

	for _, item := range receiptItems {
		dto := MapReceiptItem(item)
		dtoList = append(dtoList, dto)
	}
	return dtoList
}
