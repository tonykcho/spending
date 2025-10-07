export interface ReceiptDto
{
    StoreName: string;
    Date: string;
    Items: ReceiptItemDto[];
}

export interface ReceiptItemDto
{
    Name: string;
    Price: number;
}

export interface Receipt
{
    storeName: string;
    date: Date;
    items: ReceiptItem[];
}

export interface ReceiptItem
{
    name: string;
    price: number;
}

export function mapReceiptFromDto(dto: ReceiptDto): Receipt
{
    return {
        storeName: dto.StoreName,
        date: new Date(dto.Date),
        items: dto.Items.map(item => ({
            name: item.Name,
            price: item.Price,
        })),
    };
}