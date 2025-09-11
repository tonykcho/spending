package mappers

import (
	"spending/dto"
	"spending/models"
)

func MapCategory(category *models.Category) *dto.CategoryDto {
	if category == nil {
		return nil
	}

	dto := &dto.CategoryDto{
		Id:        category.UUId,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}

	if category.Stores != nil {
		dto.Stores = MapStoreList(category.Stores)
	}

	return dto
}

func MapCategoryList(categoryList []*models.Category) []*dto.CategoryDto {
	var dtoList []*dto.CategoryDto = make([]*dto.CategoryDto, 0)

	for _, category := range categoryList {
		dto := MapCategory(category)
		dtoList = append(dtoList, dto)
	}
	return dtoList
}
