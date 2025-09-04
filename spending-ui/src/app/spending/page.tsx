'use client'
import { Spending } from "@/models/spending";
import { deleteSpendingAsync, getSpendingListAsync } from "@/services/spending-service";
import React, { useEffect } from "react";

export default function SpendingPage()
{
    const [spendingList, setSpendingList] = React.useState<Spending[]>([]);

    useEffect(() =>
    {
        fetchSpending();
    }, []);

    async function fetchSpending()
    {
        const spendingList = await getSpendingListAsync();
        setSpendingList(spendingList);
    }

    async function onSpendingDeleted(spending: Spending)
    {
        await deleteSpendingAsync(spending.uuid);
        await fetchSpending();
    }

    return (
        <div className="p-2">
            <div className="flex flex-wrap">
                {spendingList.map((spending) => (
                    <div key={spending.uuid} className="w-1/4 p-4">
                        <div className="flex flex-col border rounded p-2 shadow-md bg-gray-100">
                            <h1 className="self-center text-2xl font-semibold">{spending.remark}</h1>
                            <h1 className="self-center text-2xl font-semibold">${spending.amount}</h1>

                            <div className="flex flex-row mt-5 justify-between">
                                <button className="btn btn-danger" onClick={() => onSpendingDeleted(spending)}>Delete</button>
                            </div>
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
}
