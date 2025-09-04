import { mapSpendingFromDto, Spending, SpendingDto } from "@/models/spending";

export async function getSpendingListAsync(): Promise<Spending[]>
{
    const response = await fetch("http://localhost:8001/spending");
    if (!response.ok)
    {
        throw new Error("Failed to fetch spending");
    }
    const spendingDtos: SpendingDto[] = await response.json();
    const spending = spendingDtos.map(mapSpendingFromDto);
    return spending;
}

export async function deleteSpendingAsync(uuid: string): Promise<void>
{
    const response = await fetch(`http://localhost:8001/spending/${uuid}`, {
        method: "DELETE",
    });
    if (!response.ok)
    {
        throw new Error("Failed to delete spending");
    }
}
