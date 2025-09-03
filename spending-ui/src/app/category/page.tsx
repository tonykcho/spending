'use client'

import { Category, CategoryDto, mapCategoryFromDto } from "@/models/category";
import React, { useEffect } from "react"
import CategoryModal, { CategoryModalRef } from "./category-modal";

export default function CategoryPage() {
    const [categories, setCategories] = React.useState<Category[]>([]);
    const modalRef = React.useRef<CategoryModalRef>(null);

    useEffect(() => {
        fetchCategories();
    }, []);

    async function fetchCategories() {
        const response = await fetch("http://localhost:8001/categories");
        if (!response.ok) {
            throw new Error("Failed to fetch categories");
        }
        const categoryDtos: CategoryDto[] = await response.json();
        var categories = categoryDtos.map(mapCategoryFromDto);
        setCategories(categories);
    }

    function OnCategoryEdit(category: Category) {
        modalRef.current?.open(category);
    }

    return (
        <div className="p-2">
            <div className="flex flex-wrap">
                {categories.map((category) => (
                    <div key={category.uuid} className="w-1/4 p-4">
                        <div className="flex flex-col border rounded p-2 shadow-md bg-gray-100">
                            <h1 className="self-center text-2xl font-semibold">{category.name}</h1>

                            <div className="flex flex-row mt-5 justify-between">
                                <button className="btn btn-danger">Delete</button>
                                <button className="btn btn-primary" onClick={() => OnCategoryEdit(category)}>Edit</button>
                            </div>
                            {/* <button className="mt-2 p-1 bg-blue-500 text-white rounded hover:bg-blue-600">View Details</button> */}
                        </div>
                    </div>
                ))}
            </div>

            <CategoryModal ref={modalRef} onCategoryChanged={() => fetchCategories()} />
        </div>
    )
}