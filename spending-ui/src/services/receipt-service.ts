import { CreateReceiptRequest, Receipt, ReceiptDto } from "@/models/receipt";
import { ReceiptOcr, ReceiptOcrDto } from "@/models/receipt_ocr";

export async function getReceiptsAsync(): Promise<Receipt[]>
{
    const response = await fetch("http://localhost:8001/api/receipts");
    if (!response.ok)
    {
        throw new Error("Failed to fetch receipts");
    }
    const receiptDtos: ReceiptDto[] = await response.json();
    return receiptDtos.map(dto => new Receipt(dto));
}

export async function createReceiptAsync(request: CreateReceiptRequest): Promise<Receipt>
{
    const response = await fetch("http://localhost:8001/api/receipts", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(request),
    });

    if (!response.ok)
    {
        throw new Error("Failed to create receipt");
    }

    const receiptDto: ReceiptDto = await response.json();
    return new Receipt(receiptDto);
}

export async function uploadReceiptAsync(imageFile: File): Promise<ReceiptOcr>
{
    // Send a post request with multipart form data
    const formData = new FormData();
    formData.append("file", imageFile);

    const response = await fetch("http://localhost:8001/api/receipts/upload", {
        method: "POST",
        body: formData,
    });

    if (!response.ok)
    {
        throw new Error("Failed to upload receipt");
    }

    const receiptDto: ReceiptOcrDto = await response.json();
    return new ReceiptOcr(receiptDto);
}