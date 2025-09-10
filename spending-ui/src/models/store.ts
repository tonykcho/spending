export interface StoreDto {
    UUId: string;
    Name: string;
    CreatedAt: string;
    UpdatedAt: string;
}

export interface Store {
    uuid: string;
    name: string;
    createdAt: Date;
    updatedAt: Date;
}

export function mapStoreFromDto(dto: StoreDto): Store {
    return {
        uuid: dto.UUId,
        name: dto.Name,
        createdAt: new Date(dto.CreatedAt),
        updatedAt: new Date(dto.UpdatedAt),
    };
}

export interface CreateStoreDto {
    Name: string;
}

export interface UpdateStoreDto {
    Name: string;
}