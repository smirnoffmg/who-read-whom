import { Button } from "@/components/common/Button";
import React from "react";

interface DeleteConfirmDialogProps {
    isOpen: boolean;
    onClose: () => void;
    onConfirm: () => void;
    entityName: string;
    entityDetails: string;
    isLoading?: boolean;
    error?: string | null;
}

export const DeleteConfirmDialog: React.FC<DeleteConfirmDialogProps> = ({
    isOpen,
    onClose,
    onConfirm,
    entityName,
    entityDetails,
    isLoading = false,
    error,
}): React.JSX.Element | null => {
    if (!isOpen) {
        return null;
    }

    const handleBackdropClick = (e: React.MouseEvent<HTMLDivElement>): void => {
        if (e.target === e.currentTarget) {
            onClose();
        }
    };

    const handleKeyDown = (e: React.KeyboardEvent<HTMLDivElement>): void => {
        if (e.key === "Escape") {
            onClose();
        }
    };

    return (
        <div
            className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50"
            onClick={handleBackdropClick}
            onKeyDown={handleKeyDown}
            role="dialog"
            aria-modal="true"
            aria-labelledby="delete-dialog-title"
            aria-describedby="delete-dialog-description"
        >
            <div
                className="bg-white rounded-lg shadow-xl max-w-md w-full mx-4 p-6"
                onClick={(e) => e.stopPropagation()}
            >
                <h2 id="delete-dialog-title" className="text-xl font-semibold text-gray-900 mb-4">
                    Confirm Delete
                </h2>
                <p id="delete-dialog-description" className="text-gray-600 mb-2">
                    Are you sure you want to delete this {entityName}?
                </p>
                <p className="text-sm text-gray-500 mb-4">{entityDetails}</p>
                {error && (
                    <div className="mb-4 rounded-md bg-red-50 border border-red-200 p-3">
                        <p className="text-sm text-red-800">{error}</p>
                    </div>
                )}
                <div className="flex justify-end gap-3">
                    <Button variant="secondary" onClick={onClose} disabled={isLoading}>
                        Cancel
                    </Button>
                    <Button variant="danger" onClick={onConfirm} disabled={isLoading}>
                        {isLoading ? "Deleting..." : "Delete"}
                    </Button>
                </div>
            </div>
        </div>
    );
};
