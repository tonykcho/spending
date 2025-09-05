package mappers

import (
	"spending/dto"
	"spending/models"
)

func MapStore(store *models.Store) *dto.StoreDto {
	if store == nil {
		return nil
	}

	dto := &dto.StoreDto{
		UUId:      store.UUId,
		Name:      store.Name,
		CreatedAt: store.CreatedAt,
		UpdatedAt: store.UpdatedAt,
	}

	return dto
}

func MapStoreList(stores []*models.Store) []*dto.StoreDto {
	var dtoList []*dto.StoreDto = make([]*dto.StoreDto, 0)

	for _, store := range stores {
		dto := MapStore(store)
		dtoList = append(dtoList, dto)
	}
	return dtoList
}
