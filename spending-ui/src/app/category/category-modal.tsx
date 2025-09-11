'use client'

import { Category, CreateCategoryDto, UpdateCategoryDto } from "@/models/category";
import React, { forwardRef, useImperativeHandle } from "react";
import { createCategoryAsync, updateCategoryAsync } from "@/services/category-service";

interface CategoryFormData {
    id: string | null;
    name: string;
    stores: StoreFormData[];
}

interface StoreFormData {
    id: string | null;
    name: string;
    isDeleted: boolean;
}

export interface CategoryModalRef {
    open: (category: Category | null) => void;
}

export interface CategoryModalProps {
    onCategoryChanged?: () => void
}

const CategoryModal = forwardRef<CategoryModalRef, CategoryModalProps>((props, ref) => {
    const [isOpen, setIsOpen] = React.useState(false);
    const [formData, setFormData] = React.useState<CategoryFormData>({ id: null, name: "", stores: [] });

    function openModal(category: Category | null = null) {
        setIsOpen(true);
        if (category != null) {
            setFormData({
                id: category.id,
                name: category.name,
                stores: category.stores.map(store => ({ id: store.id, name: store.name, isDeleted: false }))
            });
        }
    }

    function closeModal() {
        setIsOpen(false);
        setFormData({ id: null, name: "", stores: [] });
    }

    useImperativeHandle(ref, () => ({
        open: openModal,
    }));

    function onStoreDeleted(store: StoreFormData) {
        if (store.id == null) {
            setFormData({
                ...formData,
                stores: formData.stores.filter(s => s !== store)
            });
        }
        else {
            const updatedStores = formData.stores.map(s => {
                if (s === store) {
                    return { ...s, isDeleted: true };
                }
                return s;
            });
            setFormData({
                ...formData,
                stores: updatedStores
            });
        }
    }

    async function submit() {
        if (formData.name.trim() === "") {
            return;
        }

        if (formData.id != null) {
            await updateCategory();
        } else {
            await createCategory();
        }

        props.onCategoryChanged?.();

        closeModal();
    }

    async function createCategory() {
        const requestData: CreateCategoryDto = {
            name: formData.name,
            stores: formData.stores.map(store => {
                return { name: store.name }
            })
        }

        await createCategoryAsync(requestData);
    }

    async function updateCategory() {
        if (!formData.id) {
            throw new Error("Category ID is required for update.");
        }

        const requestData: UpdateCategoryDto = {
            id: formData.id,
            name: formData.name,
            addedStores: formData.stores
                .filter(store => store.id == null)
                .map(store => {
                    return { name: store.name }
                }),
            editedStores: formData.stores
                .filter(store => store.id != null)
                .map(store => {
                    return { id: store.id!, name: store.name }
                }),
            deletedStores: formData.stores
                .filter(store => store.isDeleted && store.id != null)
                .map(store => store.id!)
        }

        await updateCategoryAsync(formData.id, requestData)
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

                        <div className="flex flex-col">
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

                        <div className="flex-1 flex flex-col mt-4">
                            <label className="text-sm text-gray-600 mb-1">Stores</label>
                            <div className="space-y-2 max-h-72 overflow-y-auto">
                                {formData.stores.map((store, index) => (
                                    !store.isDeleted && (
                                        <div key={index} className="flex items-center space-x-2">
                                            <input
                                                type="text"
                                                placeholder="Store Name"
                                                className="form-control flex-1"
                                                value={store.name}
                                                onChange={(e) => {
                                                    const newStores = [...formData.stores];
                                                    newStores[index].name = e.target.value;
                                                    setFormData({ ...formData, stores: newStores });
                                                }}
                                            />
                                            <button
                                                className="btn btn-sm btn-danger"
                                                onClick={() => onStoreDeleted(store)}
                                            >
                                                Delete
                                            </button>
                                        </div>
                                    )
                                ))}
                            </div>
                            <button
                                className="btn btn-sm btn-secondary mt-2"
                                onClick={() => setFormData({ ...formData, stores: [...formData.stores, { id: null, name: "", isDeleted: false }] })}
                            >
                                Add Store
                            </button>
                            <div className="border-b my-2 border-gray-300"></div>
                        </div>


                        <div className="flex justify-between mt-2">
                            <button className="btn btn-secondary" onClick={closeModal}>Close</button>
                            <button className="btn btn-primary" onClick={() => { submit() }}>Submit</button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
});

CategoryModal.displayName = "CategoryModal";

export default CategoryModal;