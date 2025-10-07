import { mapReceiptFromDto, Receipt, ReceiptDto } from "@/models/receipt";

export async function uploadReceiptAsync(imageFile: File): Promise<Receipt>
{
    // Send a post request with multipart form data
    const formData = new FormData();
    formData.append("file", imageFile);

    const response = await fetch("http://localhost:8001/api/upload", {
        method: "POST",
        body: formData,
    });

    if (!response.ok)
    {
        throw new Error("Failed to upload receipt");
    }

    const receiptDto: ReceiptDto = await response.json();
    return mapReceiptFromDto(receiptDto);
}