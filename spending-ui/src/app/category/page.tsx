'use client'

import { Category } from "@/models/category";
import React, { useEffect } from "react"
import CategoryModal, { CategoryModalRef } from "./category-modal";
import { deleteCategoryAsync, getCategoryListAsync } from "@/services/category-service";

export default function CategoryPage() {
    const [categories, setCategories] = React.useState<Category[]>([]);
    const modalRef = React.useRef<CategoryModalRef>(null);

    useEffect(() => {
        fetchCategories();
    }, []);

    async function fetchCategories() {
        const categories = await getCategoryListAsync();
        setCategories(categories);
    }

    function OnCategoryEdit(category: Category) {
        modalRef.current?.open(category);
    }

    async function OnCategoryDeleted(category: Category) {
        await deleteCategoryAsync(category.uuid);
        fetchCategories();
    }

    return (
        <div className="p-2">
            <div className="flex flex-wrap">
                {categories.map((category) => (
                    <div key={category.uuid} className="w-1/4 p-4">
                        <div className="flex flex-col border rounded p-2 shadow-md bg-gray-100">
                            <h1 className="self-center text-2xl font-semibold">{category.name}</h1>

                            <div className="flex flex-row mt-5 justify-between">
                                <button className="btn btn-danger" onClick={() => OnCategoryDeleted(category)}>Delete</button>
                                <button className="btn btn-primary" onClick={() => OnCategoryEdit(category)}>Edit</button>
                            </div>
                        </div>
                    </div>
                ))}
            </div>

            <CategoryModal ref={modalRef} onCategoryChanged={() => fetchCategories()} />
        </div>
    )
}