'use client'
import { Spending } from "@/models/spending";
import { deleteSpendingAsync, getSpendingListAsync } from "@/services/spending-service";
import React, { useEffect } from "react";
import SpendingModal, { SpendingModalRef } from "./spending_modal";
import UploadModal, { UploadModalRef } from "./upload_modal";

export default function SpendingPage()
{
    const [spendingList, setSpendingList] = React.useState<Spending[]>([]);
    const modalRef = React.useRef<SpendingModalRef>(null);
    const uploadModelRef = React.useRef<UploadModalRef>(null);

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
        await deleteSpendingAsync(spending.id);
        await fetchSpending();
    }

    async function onAddSpending()
    {
        modalRef.current?.open();
    }

    async function onUpload()
    {
        uploadModelRef.current?.open();
    }

    function renderSpendingName(spending: Spending)
    {
        if (spending.category != null)
        {
            return <h1 className="text-2xl font-semibold">{spending.category.name}</h1>
        }
        else
        {
            return <h1 className="text-2xl font-semibold">{spending.remark}</h1>
        }
    }

    return (
        <div className="p-2">
            <div className="flex flex-col">
                <p className="text-3xl font-bold mb-4 text-center">{(new Date()).toDateString()}</p>
                <div className="flex flex-col overflow-y-auto space-y-4 h-[60vh]">
                    {spendingList.map((spending) => (
                        <div key={spending.id} className="border rounded p-2 shadow-md bg-gray-100 flex justify-between items-center">
                            <div>
                                {renderSpendingName(spending)}
                                <p className="text-gray-600">{new Date(spending.spendingDate).toDateString()}</p>
                            </div>
                            <div className="flex items-center space-x-4">
                                <span className="text-2xl font-bold">${spending.amount}</span>
                                <button className="btn btn-danger" onClick={() => onSpendingDeleted(spending)}>Delete</button>
                            </div>
                        </div>
                    ))}
                </div>
            </div>

            <div className="fixed flex-col left-0 bottom-20 flex items-center w-full">
                <button className="btn btn-primary h-12 w-72 mt-8" onClick={onUpload}>Upload</button>

                <button className="btn btn-primary h-12 w-72 mt-8" onClick={onAddSpending}>Add Spending</button>
            </div>
            <UploadModal ref={uploadModelRef} />
            <SpendingModal ref={modalRef} onSpendingChanged={() => fetchSpending()} />
        </div>
    );
}
