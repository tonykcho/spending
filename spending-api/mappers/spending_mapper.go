package mappers

import (
	"spending/dto"
	"spending/models"
)

func MapSpending(spending *models.SpendingRecord) *dto.SpendingDto {
	dto := &dto.SpendingDto{
		UUId:         spending.UUId,
		Amount:       spending.Amount,
		Remark:       spending.Remark,
		SpendingDate: spending.SpendingDate,
	}

	if spending.Category != nil {
		category := MapCategory(spending.Category)
		dto.Category = category
	}

	return dto
}

func MapSpendingList(spendingList []*models.SpendingRecord) []*dto.SpendingDto {
	var dtoList []*dto.SpendingDto = make([]*dto.SpendingDto, 0)

	for _, spending := range spendingList {
		dto := MapSpending(spending)
		dtoList = append(dtoList, dto)
	}

	return dtoList
}
