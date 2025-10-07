export interface ReceiptOcrDto
{
    StoreName: string;
    Date: string;
    Items: ReceiptItemOcrDto[];
}

export interface ReceiptItemOcrDto
{
    Name: string;
    Price: number;
}

export class ReceiptOcr
{
    storeName: string;
    date: Date;
    items: ReceiptItemOcr[];

    constructor(receiptOcrDto: ReceiptOcrDto)
    {
        this.storeName = receiptOcrDto.StoreName;
        this.date = new Date(receiptOcrDto.Date);
        this.items = receiptOcrDto.Items.map(itemDto => ({
            name: itemDto.Name,
            price: itemDto.Price,
        }));
    }
}

export class ReceiptItemOcr
{
    name: string;
    price: number;

    constructor(itemDto: ReceiptItemOcrDto)
    {
        this.name = itemDto.Name;
        this.price = itemDto.Price;
    }
}
