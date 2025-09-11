export interface StoreDto {
    Id: string;
    Name: string;
    CreatedAt: string;
    UpdatedAt: string;
}

export interface Store {
    id: string;
    name: string;
    createdAt: Date;
    updatedAt: Date;
}

export function mapStoreFromDto(dto: StoreDto): Store {
    return {
        id: dto.Id,
        name: dto.Name,
        createdAt: new Date(dto.CreatedAt),
        updatedAt: new Date(dto.UpdatedAt),
    };
}

export interface CreateStoreDto {
    name: string;
}

export interface UpdateStoreDto {
    id: string;
    name: string;
}