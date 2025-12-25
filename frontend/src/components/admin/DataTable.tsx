import React, { useMemo, useState } from "react";

export interface Column<T> {
    key: string;
    label: string;
    render: (item: T, isEditing: boolean) => React.ReactNode;
    sortable?: boolean;
}

interface DataTableProps<T> {
    columns: Column<T>[];
    data: T[];
    editingId: string | number | null;
    onEdit: (id: string | number) => void;
    onSave: (id: string | number, data: Partial<T>) => void;
    onCancel: () => void;
    onDelete: (id: string | number) => void;
    onCreate?: () => void;
    isLoading?: boolean;
    getRowId: (item: T) => string | number;
    itemsPerPage?: number;
}

export const DataTable = <T,>({
    columns,
    data,
    editingId,
    onEdit,
    onSave,
    onCancel,
    onDelete,
    onCreate,
    isLoading = false,
    getRowId,
    itemsPerPage = 10,
}: DataTableProps<T>): React.JSX.Element => {
    const [sortColumn, setSortColumn] = useState<string | null>(null);
    const [sortDirection, setSortDirection] = useState<"asc" | "desc">("asc");
    const [currentPage, setCurrentPage] = useState<number>(1);

    const handleSort = (columnKey: string): void => {
        if (sortColumn === columnKey) {
            setSortDirection(sortDirection === "asc" ? "desc" : "asc");
        } else {
            setSortColumn(columnKey);
            setSortDirection("asc");
        }
    };

    const sortedData = useMemo(() => {
        if (!sortColumn) {
            return data;
        }

        return [...data].sort((a, b) => {
            const aValue = (a as Record<string, unknown>)[sortColumn];
            const bValue = (b as Record<string, unknown>)[sortColumn];

            if (aValue === null || aValue === undefined) {
                return 1;
            }
            if (bValue === null || bValue === undefined) {
                return -1;
            }

            if (typeof aValue === "string" && typeof bValue === "string") {
                return sortDirection === "asc"
                    ? aValue.localeCompare(bValue)
                    : bValue.localeCompare(aValue);
            }

            if (typeof aValue === "number" && typeof bValue === "number") {
                return sortDirection === "asc" ? aValue - bValue : bValue - aValue;
            }

            return 0;
        });
    }, [data, sortColumn, sortDirection]);

    const paginatedData = useMemo(() => {
        const startIndex = (currentPage - 1) * itemsPerPage;
        const endIndex = startIndex + itemsPerPage;
        return sortedData.slice(startIndex, endIndex);
    }, [sortedData, currentPage, itemsPerPage]);

    const totalPages = Math.ceil(sortedData.length / itemsPerPage);

    const SortIcon = ({ columnKey }: { columnKey: string }): React.JSX.Element => {
        if (sortColumn !== columnKey) {
            return (
                <span className="ml-1 text-gray-400">
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"
                        />
                    </svg>
                </span>
            );
        }

        return (
            <span className="ml-1 text-blue-600">
                {sortDirection === "asc" ? (
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 15l7-7 7 7" />
                    </svg>
                ) : (
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                    </svg>
                )}
            </span>
        );
    };

    return (
        <div className="bg-white shadow-sm rounded-lg overflow-hidden">
            {onCreate && (
                <div className="px-6 py-4 border-b border-gray-200">
                    <button
                        onClick={onCreate}
                        disabled={isLoading || editingId !== null}
                        className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:bg-gray-300 disabled:cursor-not-allowed text-sm font-medium"
                    >
                        Add New
                    </button>
                </div>
            )}

            <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                    <thead className="bg-gray-50">
                        <tr>
                            {columns.map((column) => (
                                <th
                                    key={column.key}
                                    className={`px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider ${column.sortable !== false ? "cursor-pointer hover:bg-gray-100" : ""
                                        }`}
                                    onClick={() => column.sortable !== false && handleSort(column.key)}
                                >
                                    <div className="flex items-center">
                                        {column.label}
                                        {column.sortable !== false && <SortIcon columnKey={column.key} />}
                                    </div>
                                </th>
                            ))}
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Actions
                            </th>
                        </tr>
                    </thead>
                    <tbody className="bg-white divide-y divide-gray-200">
                        {paginatedData.length === 0 && !isLoading ? (
                            <tr>
                                <td colSpan={columns.length + 1} className="px-6 py-8 text-center text-gray-500">
                                    No data available
                                </td>
                            </tr>
                        ) : (
                            paginatedData.map((item) => {
                                const itemId = getRowId(item);
                                // Handle "new" editing case: if editingId is "new" and itemId is -1 (placeholder), consider it editing
                                const isEditing =
                                    editingId === itemId ||
                                    (editingId === "new" && itemId === -1);

                                return (
                                    <tr
                                        key={String(itemId)}
                                        className={isEditing ? "bg-blue-50" : "hover:bg-gray-50"}
                                    >
                                        {columns.map((column) => (
                                            <td
                                                key={column.key}
                                                className="px-6 py-4 whitespace-nowrap text-sm text-gray-900"
                                            >
                                                {column.render(item, isEditing)}
                                            </td>
                                        ))}
                                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                                            {isEditing ? (
                                                <div className="flex gap-2">
                                                    <button
                                                        onClick={async () => {
                                                            // For "new" case, pass "new" as the ID, otherwise use itemId
                                                            const saveId =
                                                                editingId === "new" && itemId === -1
                                                                    ? "new"
                                                                    : itemId;
                                                            await onSave(saveId, {});
                                                        }}
                                                        className="text-green-600 hover:text-green-900 text-sm font-medium"
                                                    >
                                                        Save
                                                    </button>
                                                    <button
                                                        onClick={onCancel}
                                                        className="text-gray-600 hover:text-gray-900 text-sm font-medium"
                                                    >
                                                        Cancel
                                                    </button>
                                                </div>
                                            ) : (
                                                <div className="flex gap-2">
                                                    <button
                                                        onClick={(e) => {
                                                            e.stopPropagation();
                                                            onEdit(itemId);
                                                        }}
                                                        disabled={isLoading || editingId !== null}
                                                        className="text-blue-600 hover:text-blue-900 disabled:text-gray-400 disabled:cursor-not-allowed text-sm font-medium"
                                                    >
                                                        Edit
                                                    </button>
                                                    <button
                                                        onClick={(e) => {
                                                            e.stopPropagation();
                                                            onDelete(itemId);
                                                        }}
                                                        disabled={isLoading || editingId !== null}
                                                        className="text-red-600 hover:text-red-900 disabled:text-gray-400 disabled:cursor-not-allowed text-sm font-medium"
                                                    >
                                                        Delete
                                                    </button>
                                                </div>
                                            )}
                                        </td>
                                    </tr>
                                );
                            })
                        )}
                        {isLoading && (
                            <tr>
                                <td colSpan={columns.length + 1} className="px-6 py-8 text-center">
                                    <div className="flex justify-center">
                                        <div className="w-8 h-8 border-4 border-gray-200 border-t-blue-600 rounded-full animate-spin" />
                                    </div>
                                </td>
                            </tr>
                        )}
                    </tbody>
                </table>
            </div>

            {totalPages > 1 && (
                <div className="px-6 py-4 border-t border-gray-200 flex items-center justify-between">
                    <div className="text-sm text-gray-700">
                        Showing {(currentPage - 1) * itemsPerPage + 1} to{" "}
                        {Math.min(currentPage * itemsPerPage, sortedData.length)} of {sortedData.length} results
                    </div>
                    <div className="flex gap-2">
                        <button
                            onClick={() => setCurrentPage((prev) => Math.max(1, prev - 1))}
                            disabled={currentPage === 1 || isLoading}
                            className="px-3 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 disabled:bg-gray-100 disabled:text-gray-400 disabled:cursor-not-allowed"
                        >
                            Previous
                        </button>
                        <span className="px-3 py-2 text-sm text-gray-700">
                            Page {currentPage} of {totalPages}
                        </span>
                        <button
                            onClick={() => setCurrentPage((prev) => Math.min(totalPages, prev + 1))}
                            disabled={currentPage === totalPages || isLoading}
                            className="px-3 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 disabled:bg-gray-100 disabled:text-gray-400 disabled:cursor-not-allowed"
                        >
                            Next
                        </button>
                    </div>
                </div>
            )}
        </div>
    );
};
