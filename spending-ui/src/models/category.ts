export interface CategoryDto {
    UUId: string;
    Name: string;
    CreatedAt: string;
    UpdatedAt: string;
}

export interface Category {
    uuid: string;
    name: string;
    createdAt: Date;
    updatedAt: Date;
}

export function mapCategoryFromDto(dto: CategoryDto): Category {
    return {
        uuid: dto.UUId,
        name: dto.Name,
        createdAt: new Date(dto.CreatedAt),
        updatedAt: new Date(dto.UpdatedAt),
    };
}

export interface CreateCategoryDto {
    Name: string;
}

export interface UpdateCategoryDto {
    UUId: string;
    Name: string;
}