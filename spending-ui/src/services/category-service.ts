import { CreateCategoryDto, UpdateCategoryDto, Category, CategoryDto, mapCategoryFromDto } from "@/models/category";

export async function getCategoryListAsync(): Promise<Category[]> {
    const response = await fetch("http://localhost:8001/categories");
    if (!response.ok) {
        throw new Error("Failed to fetch categories");
    }
    const categoryDtos: CategoryDto[] = await response.json();
    var categories = categoryDtos.map(mapCategoryFromDto);
    return categories;
}

export async function createCategoryAsync(requestData: CreateCategoryDto): Promise<void> {
    var response = await fetch("http://localhost:8001/categories", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(requestData),
    });

    if (!response.ok) {
        throw new Error("Failed to create category");
    }
}

export async function updateCategoryAsync(id: string, requestData: UpdateCategoryDto): Promise<void> {
    var response = await fetch(`http://localhost:8001/categories/${id}`, {
        method: "PUT",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(requestData),
    });

    if (!response.ok) {
        throw new Error("Failed to update category");
    }
}

export async function deleteCategoryAsync(id: string): Promise<void> {
    var response = await fetch(`http://localhost:8001/categories/${id}`, {
        method: "DELETE",
    });

    if (!response.ok) {
        throw new Error("Failed to delete category");
    }
}