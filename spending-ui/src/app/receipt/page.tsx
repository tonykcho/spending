'use client'

import { Receipt } from "@/models/receipt";
import { getReceiptsAsync } from "@/services/receipt-service";
import React, { useEffect } from "react";
import UploadModal, { UploadModalRef } from "./upload_modal";

export default function ReceiptPage()
{
    const [receipts, setReceipts] = React.useState<Receipt[]>([]);
    const modalRef = React.useRef<UploadModalRef>(null);

    useEffect(() =>
    {
        fetchReceipts();
    }, []);

    async function onUpload()
    {
        modalRef.current?.open();
    }

    async function fetchReceipts()
    {
        const fetchedReceipts = await getReceiptsAsync();
        setReceipts(fetchedReceipts);
        console.log(fetchedReceipts)
    }

    return (
        <div className="p-2">
            <div className="flex flex-col">
                <div className="flex flex-col overflow-y-auto space-y-4 h-[60vh]">
                    {receipts.map((receipt) => (
                        <div key={receipt.id} className="border rounded p-2 shadow-md bg-gray-100 flex justify-between items-center">
                            <div>
                                <h1 className="text-2xl font-semibold">{receipt.storeName}</h1>
                                <p className="text-gray-600">{new Date(receipt.date).toDateString()}</p>
                            </div>
                            <div className="flex flex-col items-end">
                                <span className="text-2xl font-bold">${receipt.totalAmount}</span>
                                <p className="text-gray-600">{receipt.items.length} items</p>
                            </div>
                        </div>
                    ))}
                </div>
            </div>

            <div className="fixed flex-col left-0 bottom-20 flex items-center w-full">
                <button className="btn btn-primary h-12 w-72 mt-8" onClick={onUpload}>Upload</button>
            </div>
            <UploadModal ref={modalRef} onUploadCompleted={() => fetchReceipts()} />
        </div>
    )
}