import { CreateSpendingDto, mapSpendingFromDto, Spending, SpendingDto } from "@/models/spending";

export async function getSpendingListAsync(): Promise<Spending[]> {
    const response = await fetch("http://localhost:8001/spending");
    if (!response.ok) {
        throw new Error("Failed to fetch spending");
    }
    const spendingDtos: SpendingDto[] = await response.json();
    const spending = spendingDtos.map(mapSpendingFromDto);
    return spending;
}

export async function createSpendingAsync(requestData: CreateSpendingDto): Promise<Spending> {
    const response = await fetch("http://localhost:8001/spending", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(requestData),
    });
    if (!response.ok) {
        throw new Error("Failed to create spending");
    }
    const spendingDto: SpendingDto = await response.json();
    return mapSpendingFromDto(spendingDto);
}

export async function deleteSpendingAsync(uuid: string): Promise<void> {
    const response = await fetch(`http://localhost:8001/spending/${uuid}`, {
        method: "DELETE",
    });
    if (!response.ok) {
        throw new Error("Failed to delete spending");
    }
}
