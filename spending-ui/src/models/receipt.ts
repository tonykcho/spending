export interface ReceiptDto
{
    Id: string;
    StoreName: string;
    Date: string;
    Total: number;
    Items: ReceiptItemDto[];
    CreatedAt: string;
    UpdatedAt: string;
}

export interface ReceiptItemDto
{
    Id: string;
    Name: string;
    Price: number;
    CreatedAt: string;
    UpdatedAt: string;
}

export class Receipt
{
    id: string;
    storeName: string;
    date: Date;
    totalAmount: number;
    items: ReceiptItem[];
    createdAt: Date;
    updatedAt: Date;

    constructor(receiptDto: ReceiptDto)
    {
        this.id = receiptDto.Id;
        this.storeName = receiptDto.StoreName;
        this.date = new Date(receiptDto.Date);
        this.totalAmount = receiptDto.Total;
        this.items = receiptDto.Items.map(itemDto => new ReceiptItem(itemDto));
        this.createdAt = new Date(receiptDto.CreatedAt);
        this.updatedAt = new Date(receiptDto.UpdatedAt);
    }
}

export class ReceiptItem
{
    id: string;
    name: string;
    price: number;
    createdAt: Date;
    updatedAt: Date;

    constructor(itemDto: ReceiptItemDto)
    {
        this.id = itemDto.Id;
        this.name = itemDto.Name;
        this.price = itemDto.Price;
        this.createdAt = new Date(itemDto.CreatedAt);
        this.updatedAt = new Date(itemDto.UpdatedAt);
    }
}

export interface CreateReceiptRequest
{
    storeName: string;
    date: Date;
    totalAmount: number;
    items: CreateReceiptItemRequest[];
}

export interface CreateReceiptItemRequest
{
    name: string;
    price: number;
}