'use client'

import { useLoading } from "@/components/loading";
import { Receipt } from "@/models/receipt";
import { CreateSpendingDto } from "@/models/spending";
import { uploadReceiptAsync } from "@/services/receipt-service";
import { createSpendingAsync } from "@/services/spending-service";
import Image from "next/image";
import React from "react";
import { forwardRef, useImperativeHandle } from "react";

export interface UploadModalRef
{
    open: () => void;
}

export interface UploadModalProps
{
    onUploadCompleted?: () => void
}

const UploadModal = forwardRef<UploadModalRef, UploadModalProps>((props, ref) =>
{
    const [isOpen, setIsOpen] = React.useState(false);
    const [image, setImage] = React.useState<File | null>(null);
    const [receipt, setReceipt] = React.useState<Receipt | null>(null);
    const { showLoading, hideLoading } = useLoading();

    useImperativeHandle(ref, () => ({
        open: open,
    }));

    function open()
    {
        setIsOpen(!isOpen);
    }

    function close()
    {
        setIsOpen(false);
    }

    async function onFileChange(event: React.ChangeEvent<HTMLInputElement>)
    {
        const file = event.target.files?.[0];
        // Allow Jpeg and Png only
        if (file && (file.type === "image/jpeg" || file.type === "image/png"))
        {
            setImage(file);
            showLoading();
            const result = await uploadReceiptAsync(file)
            hideLoading();
            setReceipt(result);
            console.log("Upload result:", result);
        }
        else
        {
            alert("Please select a valid image file (JPG or PNG)");
        }
    }

    async function submitReceipt(receipt: Receipt)
    {
        const requestData: CreateSpendingDto = {
            amount: parseFloat(receipt.items.reduce((sum, item) => sum + item.price, 0).toFixed(2)),
            remark: receipt.storeName,
            spendingDate: receipt.date,
            categoryId: null,
            storeId: null,
        };
        await createSpendingAsync(requestData);
        props.onUploadCompleted?.();
        close()
    }

    function renderImagePreview()
    {
        if (image)
        {
            return <Image src={URL.createObjectURL(image)} alt="Selected" width={0} height={0} style={{ width: '275px', height: 'auto' }} className="mb-4" />
        }
        else
        {
            return (
                <div className="flex flex-col items-center">
                    <button className="btn btn-secondary h-12 mx-8 mb-4" onClick={() => document.getElementById('fileInput')?.click()}>Upload Image</button>
                    <input type="file" className="hidden" id="fileInput" onChange={onFileChange} />
                </div>
            )
        }
    }

    function renderReceiptDetails()
    {
        if (receipt)
        {
            return (<div className="text-center w-full">
                <h2 className="text-lg font-semibold">Receipt Details</h2>

                <div className="flex flex-row">
                    <div className="flex-1 text-right pr-2">Store:</div>
                    <div className="flex-1 text-left pl-2">{receipt.storeName}</div>
                </div>

                <div className="flex flex-row">
                    <div className="flex-1 text-right pr-2">Date:</div>
                    <div className="flex-1 text-left pl-2">{receipt.date.toDateString()}</div>
                </div>

                {receipt.items.map((item, index) => (
                    <div key={index} className="flex flex-row">
                        <div className="flex-1 text-right pr-2">{item.name}:</div>
                        <div className="flex-1 text-left pl-2">${item.price.toFixed(2)}</div>
                    </div>
                ))}
                <div className="flex flex-row">
                    <div className="flex-1 text-right pr-2">Total:</div>
                    <div className="flex-1 text-left pl-2">${receipt.items.reduce((sum, item) => sum + item.price, 0).toFixed(2)}</div>
                </div>
            </div>)
        }
    }

    return (
        <div className={`flex flex-col items-center spending-modal fixed w-full left-0 h-5/6 ${isOpen ? 'top-1/6' : 'top-full'} `}>
            <div className="bg-white border border-b-0 rounded-t border-gray-600 w-9/10 h-full flex flex-col overflow-y-auto">
                <div className="flex-1 flex flex-col items-center justify-center">
                    {renderImagePreview()}
                    {renderReceiptDetails()}
                </div>
                {receipt && (
                    <button className="btn btn-primary h-12 mx-8 mt-4" onClick={() => submitReceipt(receipt)}>Submit</button>
                )}
                <button className="btn btn-secondary h-12 mx-8 my-4" onClick={close}>Close</button>
            </div>
        </div>
    )
});

UploadModal.displayName = "UploadModal";

export default UploadModal;