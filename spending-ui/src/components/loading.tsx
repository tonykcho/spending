'use client'

import React from "react"

interface LoadingContextType
{
    showLoading: () => void;
    hideLoading: () => void;
}

const LoadingContext = React.createContext<LoadingContextType | undefined>(undefined);

export function useLoading()
{
    const context = React.useContext(LoadingContext);
    if (!context)
    {
        throw new Error("useLoading must be used within a LoadingProvider");
    }
    return context;
}

export function LoadingProvider({ children }: { children: React.ReactNode })
{
    const [isLoading, setIsLoading] = React.useState(false);

    const showLoading = () =>
    {
        setIsLoading(true);
    };

    const hideLoading = () =>
    {
        setIsLoading(false);
    };

    return (
        <LoadingContext.Provider value={{ showLoading, hideLoading }}>
            {children}
            {isLoading && (
                <div className="fixed h-full w-full top-0 left-0 flex items-center justify-center bg-gray-500/10">
                    <div className="loader ease-linear rounded-full border-8 border-t-8 border-gray-200 h-16 w-16"></div>
                </div>
            )}
        </LoadingContext.Provider>
    );
}