import { CreateStoreDto, mapStoreFromDto, Store, StoreDto, UpdateStoreDto } from "./store";

export interface CategoryDto {
    Id: string;
    Name: string;
    Stores: StoreDto[];
    CreatedAt: string;
    UpdatedAt: string;
}

export interface Category {
    id: string;
    name: string;
    stores: Store[];
    createdAt: Date;
    updatedAt: Date;
}

export function mapCategoryFromDto(dto: CategoryDto): Category {
    return {
        id: dto.Id,
        name: dto.Name,
        stores: dto.Stores ? dto.Stores.map(mapStoreFromDto) : [],
        createdAt: new Date(dto.CreatedAt),
        updatedAt: new Date(dto.UpdatedAt),
    };
}

export interface CreateCategoryDto {
    name: string;
    stores: CreateStoreDto[];
}

export interface UpdateCategoryDto {
    id: string;
    name: string;
    addedStores: CreateStoreDto[];
    editedStores: UpdateStoreDto[];
    deletedStores: string[];
}