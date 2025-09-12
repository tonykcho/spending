'use client'

import React from "react";

interface MessageContextType {
    showMessage: (message: string) => void;
    hideMessage: () => void;
}

const MessageContext = React.createContext<MessageContextType | undefined>(undefined);

export function useMessage() {
    const context = React.useContext(MessageContext);
    if (!context) {
        throw new Error("useMessage must be used within a MessageProvider");
    }
    return context;
}

export function MessageProvider({ children }: { children: React.ReactNode }) {
    const [message, setMessage] = React.useState<string | null>(null);

    const showMessage = (msg: string) => {
        setMessage(msg);
    };

    const hideMessage = () => {
        setMessage(null);
    };

    return (
        <MessageContext.Provider value={{ showMessage, hideMessage }}>
            {children}
            {message && (
                <div className="fixed h-full w-full top-0 left-0 flex items-center justify-center bg-gray-500/10" onClick={hideMessage}>
                    <div className="flex flex-col w-80 h-40 bg-white text-black px-4 py-4 rounded border border-gray-300 shadow-lg">
                        <p className="flex-1">{message}</p>
                        <div className="flex justify-end mt-4">
                            <button className="btn btn-primary" onClick={hideMessage}>OK</button>
                        </div>
                    </div>
                </div>
            )}
        </MessageContext.Provider>
    );
}