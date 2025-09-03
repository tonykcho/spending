'use client'

import { Category, CreateCategoryDto, UpdateCategoryDto } from "@/models/category";
import React, { forwardRef, useImperativeHandle } from "react";
import { createCategoryAsync, updateCategoryAsync } from "@/services/category-service";

interface CategoryFormData {
    uuid: string | null;
    name: string;
}

export interface CategoryModalRef {
    open: (category: Category | null) => void;
}

export interface CategoryModalProps {
    onCategoryChanged?: () => void
}

const CategoryModal = forwardRef<CategoryModalRef, CategoryModalProps>((props, ref) => {
    const [isOpen, setIsOpen] = React.useState(false);
    const [formData, setFormData] = React.useState<CategoryFormData>({ uuid: null, name: "" });

    function openModal(category: Category | null = null) {
        setIsOpen(true);
        if (category != null) {
            setFormData({ uuid: category.uuid, name: category.name });
        }
    }

    function closeModal() {
        setIsOpen(false);
    }

    useImperativeHandle(ref, () => ({
        open: openModal,
    }));

    async function submit() {
        if (formData.name.trim() === "") {
            return;
        }

        if (formData.uuid != null) {
            await updateCategory();
        } else {
            await createCategory();
        }

        props.onCategoryChanged?.();

        closeModal();
    }

    async function createCategory() {
        var requestData: CreateCategoryDto = {
            Name: formData.name
        }

        await createCategoryAsync(requestData);
    }

    async function updateCategory() {
        if (!formData.uuid) {
            throw new Error("Category UUID is required for update.");
        }

        var requestData: UpdateCategoryDto = {
            Name: formData.name
        }

        await updateCategoryAsync(formData.uuid, requestData)
    }

    return (
        <div>
            <button className="btn btn-primary fixed bottom-8 right-8" onClick={() => openModal()}>
                Create
            </button>

            {isOpen && (
                <div className="modal-bg">
                    <div className="modal">
                        <div className="flex justify-between items-center mb-4">
                            <h2 className="text-2xl">Create Category</h2>
                            <button className="close-button" onClick={closeModal}></button>
                        </div>

                        <div className="flex-1 flex flex-col">
                            <label htmlFor="category-name" className="text-sm text-gray-600 mb-1">Category Name</label>
                            <input
                                id="category-name"
                                type="text"
                                placeholder="Name"
                                className="form-control"
                                value={formData.name}
                                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                            />
                        </div>

                        <div className="flex justify-between">
                            <button className="btn btn-secondary" onClick={closeModal}>Close</button>
                            <button className="btn btn-primary" onClick={() => { submit() }}>Submit</button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
});

export default CategoryModal;