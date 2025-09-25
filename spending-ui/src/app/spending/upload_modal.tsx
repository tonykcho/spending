'use client'

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

    function onFileChange(event: React.ChangeEvent<HTMLInputElement>)
    {
        const file = event.target.files?.[0];
        // Allow Jpeg and Png only
        if (file && (file.type === "image/jpeg" || file.type === "image/png"))
        {
            setImage(file);
        }
        else
        {
            alert("Please select a valid image file (JPG or PNG)");
        }
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

    return (
        <div className={`flex flex-col items-center spending-modal fixed w-full left-0 h-5/6 ${isOpen ? 'top-1/6' : 'top-full'} `}>
            <div className="bg-white border rounded-t border-gray-600 w-9/10 h-full flex flex-col">
                <div className="flex-1 flex flex-col items-center justify-center">
                    {renderImagePreview()}
                </div>
                <button className="btn btn-primary h-12 mx-8 mt-4">Submit</button>
                <button className="btn btn-secondary h-12 mx-8 my-4" onClick={close}>Close</button>
            </div>
        </div>
    )
});

UploadModal.displayName = "UploadModal";

export default UploadModal;