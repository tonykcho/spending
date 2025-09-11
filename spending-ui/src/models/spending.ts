import { Category, CategoryDto, mapCategoryFromDto } from "./category";

export interface SpendingDto {
    Id: string;
    Amount: number;
    Remark: string;
    SpendingDate: string;
    Category: CategoryDto | null;
}

export interface Spending {
    id: string;
    amount: number;
    remark: string;
    spendingDate: Date;
    category: Category | null;
}

export function mapSpendingFromDto(dto: SpendingDto): Spending {
    return {
        id: dto.Id,
        amount: dto.Amount,
        remark: dto.Remark,
        spendingDate: new Date(dto.SpendingDate),
        category: dto.Category ? mapCategoryFromDto(dto.Category) : null,
    };
}