package mappers

import (
	"spending/models"
	"spending/responses"
)

func MapSpending(spending models.SpendingRecord) responses.SpendingResponse {
	response := responses.SpendingResponse{
		UUId:         spending.UUId,
		Amount:       spending.Amount,
		Remark:       spending.Remark,
		SpendingDate: spending.SpendingDate,
	}

	if spending.Category != nil {
		category := MapCategory(*spending.Category)
		response.Category = &category
	}

	return response
}

func MapSpendingList(spendingList []*models.SpendingRecord) []responses.SpendingResponse {
	var responseList []responses.SpendingResponse = make([]responses.SpendingResponse, 0)

	for _, spending := range spendingList {
		response := MapSpending(*spending)
		responseList = append(responseList, response)
	}

	return responseList
}
