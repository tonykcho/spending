import { mapStoreFromDto, Store, StoreDto } from "./store";

export interface CategoryDto {
    UUId: string;
    Name: string;
    Stores: StoreDto[];
    CreatedAt: string;
    UpdatedAt: string;
}

export interface Category {
    uuid: string;
    name: string;
    stores: Store[];
    createdAt: Date;
    updatedAt: Date;
}

export function mapCategoryFromDto(dto: CategoryDto): Category {
    return {
        uuid: dto.UUId,
        name: dto.Name,
        stores: dto.Stores ? dto.Stores.map(mapStoreFromDto) : [],
        createdAt: new Date(dto.CreatedAt),
        updatedAt: new Date(dto.UpdatedAt),
    };
}

export interface CreateCategoryDto {
    Name: string;
}

export interface UpdateCategoryDto {
    Name: string;
}