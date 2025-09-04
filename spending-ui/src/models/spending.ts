import { Category, CategoryDto, mapCategoryFromDto } from "./category";

export interface SpendingDto
{
    UUId: string;
    Amount: number;
    Remark: string;
    SpendingDate: string;
    Category: CategoryDto | null;
}

export interface Spending
{
    uuid: string;
    amount: number;
    remark: string;
    spendingDate: Date;
    category: Category | null;
}

export function mapSpendingFromDto(dto: SpendingDto): Spending
{
    return {
        uuid: dto.UUId,
        amount: dto.Amount,
        remark: dto.Remark,
        spendingDate: new Date(dto.SpendingDate),
        category: dto.Category ? mapCategoryFromDto(dto.Category) : null,
    };
}