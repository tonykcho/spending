'use client'

import { Category } from "@/models/category";
import { CreateSpendingDto } from "@/models/spending";
import { getCategoryListAsync } from "@/services/category-service";
import { createSpendingAsync } from "@/services/spending-service";
import React, { useEffect, useImperativeHandle } from "react";
import { forwardRef } from "react"

export interface SpendingModalRef {
    open: () => void;
}

export interface SpendingModalProps {
    onSpendingChanged?: () => void
}

const SpendingModal = forwardRef<SpendingModalRef, SpendingModalProps>((props, ref) => {
    const [categories, setCategories] = React.useState<Category[]>([]);
    const [isOpen, setIsOpen] = React.useState(false);
    const [formData, setFormData] = React.useState({ amount: "", categoryId: "", storeId: "", spendingDate: "", remark: "" });

    useEffect(() => {
        fetchCategories();
    }, []);

    useImperativeHandle(ref, () => ({
        open: open,
    }));

    async function fetchCategories() {
        const categories = await getCategoryListAsync()
        setCategories(categories)
    }

    function open() {
        setFormData({ amount: "", categoryId: "", storeId: "", spendingDate: "", remark: "" });
        setIsOpen(!isOpen);
    }

    function close() {
        setIsOpen(false);
    }

    async function submit() {
        var requestData: CreateSpendingDto = {
            amount: parseFloat(formData.amount),
            remark: formData.remark,
            spendingDate: formData.spendingDate === "" ? null : new Date(formData.spendingDate),
            categoryId: formData.categoryId === "" ? null : formData.categoryId,
            storeId: formData.storeId === "" ? null : formData.storeId,
        };
        await createSpendingAsync(requestData);
        props.onSpendingChanged?.();
        close()
    }

    function renderStores() {
        if (formData.categoryId === "") {
            return null;
        }
        else {
            const category = categories.find(c => c.id === formData.categoryId);
            if (!category) {
                return null;
            }
            return (
                <div className="flex flex-col mt-4">
                    <label htmlFor="store" className="text-sm text-gray-600 mb-1">Store</label>
                    <select id="store" className="form-control" value={formData.storeId} onChange={(e) => setFormData({ ...formData, storeId: e.target.value })}>
                        <option value="">Select Store</option>
                        {category.stores.map((store) => (
                            <option key={store.id} value={store.id}>{store.name}</option>
                        ))}
                    </select>
                </div>
            );
        }
    }

    return (
        <div className={`flex flex-col items-center spending-modal fixed w-full left-0 h-4/5 ${isOpen ? 'top-1/5' : 'top-full'} `}>
            <div className="bg-white border rounded-t border-gray-600 w-9/10 h-full flex flex-col">
                <div className="flex-1 flex flex-col px-4 pt-4">
                    <div className="flex flex-col">
                        <label htmlFor="amount" className="text-sm text-gray-600 mb-1">Amount</label>
                        <input
                            id="amount"
                            type="number"
                            placeholder="Amount"
                            className="form-control"
                            value={formData.amount}
                            onChange={(e) => setFormData({ ...formData, amount: e.target.value })}
                        />
                    </div>

                    <div className="flex flex-col mt-4">
                        <label htmlFor="date" className="text-sm text-gray-600 mb-1">Date</label>
                        <input
                            id="date"
                            type="date"
                            placeholder="Date"
                            className="form-control"
                            value={formData.spendingDate}
                            onChange={(e) => setFormData({ ...formData, spendingDate: e.target.value })}
                        />
                    </div>

                    <div className="flex flex-col mt-4">
                        <label htmlFor="category" className="text-sm text-gray-600 mb-1">Category</label>
                        <select id="category" className="form-control" value={formData.categoryId} onChange={(e) => setFormData({ ...formData, categoryId: e.target.value, storeId: "" })}>
                            <option value="">Select Category</option>
                            {categories.map((category) => (
                                <option key={category.id} value={category.id}>{category.name}</option>
                            ))}
                        </select>
                    </div>

                    {renderStores()}

                    <div className="flex flex-col mt-4">
                        <label htmlFor="remark" className="text-sm text-gray-600 mb-1">Remark</label>
                        <input
                            id="remark"
                            type="text"
                            placeholder="Remark"
                            className="form-control"
                            value={formData.remark}
                            onChange={(e) => setFormData({ ...formData, remark: e.target.value })}
                        />
                    </div>
                </div>
                <button className="btn btn-primary h-12 mx-8 mt-4" onClick={submit}>Submit</button>
                <button className="btn btn-secondary h-12 mx-8 my-4" onClick={close}>Close</button>
            </div>
        </div>
    )
});

SpendingModal.displayName = "SpendingModal"

export default SpendingModal;