package mappers

import (
	"spending/models"
	"spending/responses"
)

func MapCategory(category models.Category) responses.CategoryResponse {
	response := responses.CategoryResponse{
		UUId:      category.UUId,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}

	return response
}

func MapCategoryList(categoryList []*models.Category) []responses.CategoryResponse {
	var responseList []responses.CategoryResponse = make([]responses.CategoryResponse, 0)

	for _, category := range categoryList {
		response := MapCategory(*category)
		responseList = append(responseList, response)
	}
	return responseList
}
