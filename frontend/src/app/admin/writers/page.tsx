"use client";

import { type Column, DataTable } from "@/components/admin/DataTable";
import { DeleteConfirmDialog } from "@/components/admin/DeleteConfirmDialog";
import { Button } from "@/components/common/Button";
import { ErrorMessage } from "@/components/common/ErrorMessage";
import { useWriterStore } from "@/stores/writerStore";
import type { CreateWriterRequest, UpdateWriterRequest, Writer } from "@/types/writer";
import {
  type CSVImportResult,
  downloadCSV,
  exportWritersToCSV,
  importWritersFromCSV,
} from "@/utils/csvUtils";
import React, { useEffect, useRef, useState } from "react";

interface WriterFormData {
  name: string;
  birth_year: string;
  death_year: string;
  bio: string;
}

const initialFormData: WriterFormData = {
  name: "",
  birth_year: "",
  death_year: "",
  bio: "",
};

export default function WritersAdminPage(): React.JSX.Element {
  const {
    writers,
    isLoading,
    error,
    fetchWriters,
    createWriter,
    updateWriter,
    deleteWriter,
    clearError,
  } = useWriterStore();

  const [editingId, setEditingId] = useState<number | "new" | null>(null);
  const [formData, setFormData] = useState<WriterFormData>(initialFormData);
  const [formErrors, setFormErrors] = useState<Partial<Record<keyof WriterFormData, string>>>({});
  const [deleteConfirmOpen, setDeleteConfirmOpen] = useState<boolean>(false);
  const [writerToDelete, setWriterToDelete] = useState<Writer | null>(null);
  const [importPreview, setImportPreview] = useState<CSVImportResult<Writer> | null>(null);
  const [showImportPreview, setShowImportPreview] = useState<boolean>(false);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const firstInputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    void fetchWriters(1000, 0); // Fetch a large number for admin view
  }, [fetchWriters]);

  useEffect(() => {
    if (editingId && editingId !== "new" && firstInputRef.current) {
      firstInputRef.current.focus();
    }
  }, [editingId]);

  const handleEdit = (id: number | string): void => {
    const writer = writers.find((w) => w.id === Number(id));
    if (writer) {
      setEditingId(writer.id);
      setFormData({
        name: writer.name,
        birth_year: writer.birth_year.toString(),
        death_year: writer.death_year?.toString() ?? "",
        bio: writer.bio ?? "",
      });
      setFormErrors({});
    }
  };

  const handleCreate = (): void => {
    setEditingId("new");
    setFormData(initialFormData);
    setFormErrors({});
  };

  const handleCancel = (): void => {
    setEditingId(null);
    setFormData(initialFormData);
    setFormErrors({});
  };

  const validateForm = (): boolean => {
    const errors: Partial<Record<keyof WriterFormData, string>> = {};

    if (!formData.name.trim()) {
      errors.name = "Name is required";
    }

    if (!formData.birth_year.trim()) {
      errors.birth_year = "Birth year is required";
    } else {
      const birthYear = parseInt(formData.birth_year, 10);
      if (isNaN(birthYear)) {
        errors.birth_year = "Birth year must be a valid number";
      }
    }

    if (formData.death_year.trim()) {
      const deathYear = parseInt(formData.death_year, 10);
      if (isNaN(deathYear)) {
        errors.death_year = "Death year must be a valid number";
      }
    }

    setFormErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const handleSave = async (id: number | string): Promise<void> => {
    if (!validateForm()) {
      return;
    }

    const birthYear = parseInt(formData.birth_year, 10);
    const deathYear = formData.death_year.trim() ? parseInt(formData.death_year, 10) : null;

    // Check both editingId and id parameter to handle "new" case
    if (editingId === "new" || id === "new") {
      const createData: CreateWriterRequest = {
        name: formData.name.trim(),
        birth_year: birthYear,
        death_year: deathYear,
        bio: formData.bio.trim() || null,
      };

      const result = await createWriter(createData);
      if (result) {
        await fetchWriters(1000, 0);
        // Exit edit mode after successful creation
        setEditingId(null);
        setFormData(initialFormData);
        setFormErrors({});
      }
    } else {
      const updateData: UpdateWriterRequest = {
        name: formData.name.trim(),
        birth_year: birthYear,
        death_year: deathYear,
        bio: formData.bio.trim() || null,
      };

      await updateWriter(Number(id), updateData);
      setEditingId(null);
      setFormData(initialFormData);
      await fetchWriters(1000, 0);
    }
  };

  const handleDeleteClick = (id: number | string): void => {
    // Ensure ID is a valid number
    const numericId = typeof id === "string" ? parseInt(id, 10) : id;
    if (isNaN(numericId) || numericId <= 0) {
      console.error("Invalid writer ID for deletion:", id);
      return;
    }

    const writer = writers.find((w) => w.id === numericId);
    if (writer) {
      setWriterToDelete(writer);
      setDeleteConfirmOpen(true);
    } else {
      console.error("Writer not found for deletion:", id, "numericId:", numericId);
    }
  };

  const handleDeleteConfirm = async (): Promise<void> => {
    if (writerToDelete) {
      // Clear any previous errors
      clearError();

      await deleteWriter(writerToDelete.id);
      // The store will set error state if deletion fails
      // We'll use useEffect to handle dialog closing on success
    }
  };

  // Close dialog and refresh on successful deletion (when error clears after delete attempt)
  useEffect(() => {
    if (deleteConfirmOpen && writerToDelete && !isLoading && !error) {
      // Check if the writer was actually deleted (not in the list anymore)
      const stillExists = writers.some((w) => w.id === writerToDelete.id);
      if (!stillExists) {
        setDeleteConfirmOpen(false);
        setWriterToDelete(null);
        void fetchWriters(1000, 0);
      }
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [deleteConfirmOpen, writerToDelete, isLoading, error, writers]);

  const handleExportCSV = (): void => {
    const csvContent = exportWritersToCSV(writers);
    const filename = `writers-${new Date().toISOString().split("T")[0]}.csv`;
    downloadCSV(csvContent, filename);
  };

  const handleImportCSV = (): void => {
    fileInputRef.current?.click();
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>): void => {
    const file = e.target.files?.[0];
    if (!file) {
      return;
    }

    const reader = new FileReader();
    reader.onload = (event) => {
      const csvString = event.target?.result as string;
      const result = importWritersFromCSV(csvString);
      setImportPreview(result);
      setShowImportPreview(true);
    };
    reader.readAsText(file);

    // Reset file input
    if (fileInputRef.current) {
      fileInputRef.current.value = "";
    }
  };

  const handleImportConfirm = async (): Promise<void> => {
    if (!importPreview) {
      return;
    }

    let successCount = 0;
    let errorCount = 0;

    for (const writer of importPreview.data) {
      try {
        const createData: CreateWriterRequest = {
          name: writer.name,
          birth_year: writer.birth_year,
          death_year: writer.death_year,
          bio: writer.bio,
        };
        const result = await createWriter(createData);
        if (result) {
          successCount++;
        } else {
          errorCount++;
        }
      } catch {
        errorCount++;
      }
    }

    setShowImportPreview(false);
    setImportPreview(null);
    void fetchWriters(1000, 0);

    alert(`Import completed: ${successCount} successful, ${errorCount} failed`);
  };

  const truncateText = (text: string | null, maxLength: number): string => {
    if (!text) {
      return "";
    }
    return text.length > maxLength ? `${text.substring(0, maxLength)}...` : text;
  };

  const columns: Column<Writer>[] = [
    {
      key: "id",
      label: "ID",
      sortable: true,
      render: (writer, isEditing) => {
        if (isEditing) {
          return <span className="text-gray-500">{writer.id}</span>;
        }
        return <span>{writer.id}</span>;
      },
    },
    {
      key: "name",
      label: "Name",
      sortable: true,
      render: (writer, isEditing) => {
        if (isEditing) {
          return (
            <input
              ref={
                (editingId === writer.id || (editingId === "new" && writer.id === -1))
                  ? firstInputRef
                  : undefined
              }
              type="text"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              className="w-full px-2 py-1 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="Name"
            />
          );
        }
        return <span>{writer.name}</span>;
      },
    },
    {
      key: "birth_year",
      label: "Birth Year",
      sortable: true,
      render: (writer, isEditing) => {
        if (isEditing) {
          return (
            <input
              type="number"
              value={formData.birth_year}
              onChange={(e) => setFormData({ ...formData, birth_year: e.target.value })}
              className="w-full px-2 py-1 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="Birth Year"
            />
          );
        }
        return <span>{writer.birth_year}</span>;
      },
    },
    {
      key: "death_year",
      label: "Death Year",
      sortable: true,
      render: (writer, isEditing) => {
        if (isEditing) {
          return (
            <input
              type="number"
              value={formData.death_year}
              onChange={(e) => setFormData({ ...formData, death_year: e.target.value })}
              className="w-full px-2 py-1 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="Death Year (optional)"
            />
          );
        }
        return <span>{writer.death_year ?? "—"}</span>;
      },
    },
    {
      key: "bio",
      label: "Bio",
      sortable: false,
      render: (writer, isEditing) => {
        if (isEditing) {
          return (
            <textarea
              value={formData.bio}
              onChange={(e) => setFormData({ ...formData, bio: e.target.value })}
              className="w-full px-2 py-1 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="Bio (optional)"
              rows={2}
            />
          );
        }
        return <span className="text-gray-600">{truncateText(writer.bio, 50)}</span>;
      },
    },
  ];

  // Add new row at top when creating
  const displayData = editingId === "new" ? [{ id: -1 } as Writer, ...writers] : writers;

  return (
    <div className="flex flex-col h-full">
      <div className="p-6 bg-white border-b border-gray-200">
        <h1 className="text-2xl font-bold text-gray-900">Writers</h1>
        <p className="text-sm text-gray-600 mt-1">Manage literary writers and their information</p>
      </div>

      <div className="flex-1 p-6 overflow-auto">
        {error && (
          <div className="mb-4">
            <ErrorMessage message={error} onDismiss={clearError} />
          </div>
        )}

        <div className="mb-4 flex gap-3">
          <Button onClick={handleExportCSV} disabled={isLoading || writers.length === 0}>
            Export CSV
          </Button>
          <Button variant="secondary" onClick={handleImportCSV} disabled={isLoading}>
            Import CSV
          </Button>
          <input
            ref={fileInputRef}
            type="file"
            accept=".csv"
            onChange={handleFileChange}
            className="hidden"
          />
        </div>

        {showImportPreview && importPreview && (
          <div className="mb-4 bg-white p-4 rounded-lg shadow-sm border border-gray-200">
            <h3 className="text-lg font-semibold mb-2">Import Preview</h3>
            <p className="text-sm text-gray-600 mb-4">
              {importPreview.data.length} rows to import
              {importPreview.errors.length > 0 && (
                <span className="text-red-600 ml-2">({importPreview.errors.length} errors)</span>
              )}
            </p>
            {importPreview.errors.length > 0 && (
              <div className="mb-4 max-h-40 overflow-y-auto">
                <p className="text-sm font-medium text-red-600 mb-2">Validation Errors:</p>
                <ul className="list-disc list-inside text-sm text-red-600">
                  {importPreview.errors.map((error, index) => (
                    <li key={index}>
                      Row {error.row}, {error.field}: {error.message}
                    </li>
                  ))}
                </ul>
              </div>
            )}
            {importPreview.data.length > 0 && (
              <div className="mb-4 max-h-60 overflow-y-auto border border-gray-200 rounded">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Name
                      </th>
                      <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Birth Year
                      </th>
                      <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Death Year
                      </th>
                      <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Bio
                      </th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {importPreview.data.slice(0, 20).map((writer, index) => (
                      <tr key={index}>
                        <td className="px-4 py-2 text-sm text-gray-900">{writer.name}</td>
                        <td className="px-4 py-2 text-sm text-gray-900">{writer.birth_year}</td>
                        <td className="px-4 py-2 text-sm text-gray-900">
                          {writer.death_year ?? "—"}
                        </td>
                        <td className="px-4 py-2 text-sm text-gray-600">
                          {writer.bio ? (writer.bio.length > 50 ? `${writer.bio.substring(0, 50)}...` : writer.bio) : "—"}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
                {importPreview.data.length > 20 && (
                  <p className="px-4 py-2 text-sm text-gray-500 bg-gray-50">
                    ... and {importPreview.data.length - 20} more rows
                  </p>
                )}
              </div>
            )}
            {importPreview.data.length === 0 && importPreview.errors.length === 0 && (
              <p className="text-sm text-gray-500 mb-4">No valid data found in CSV file.</p>
            )}
            <div className="flex gap-3">
              <Button onClick={handleImportConfirm} disabled={!importPreview.isValid}>
                Confirm Import
              </Button>
              <Button variant="secondary" onClick={() => setShowImportPreview(false)}>
                Cancel
              </Button>
            </div>
          </div>
        )}

        {formErrors.name && (
          <div className="mb-2 text-sm text-red-600">Name: {formErrors.name}</div>
        )}
        {formErrors.birth_year && (
          <div className="mb-2 text-sm text-red-600">Birth Year: {formErrors.birth_year}</div>
        )}
        {formErrors.death_year && (
          <div className="mb-2 text-sm text-red-600">Death Year: {formErrors.death_year}</div>
        )}

        <DataTable
          columns={columns}
          data={displayData}
          editingId={editingId}
          onEdit={handleEdit}
          onSave={handleSave}
          onCancel={handleCancel}
          onDelete={handleDeleteClick}
          onCreate={handleCreate}
          isLoading={isLoading}
          getRowId={(writer) => writer.id}
        />

        {editingId === "new" && (
          <div className="mt-4 bg-blue-50 p-4 rounded-lg border border-blue-200">
            <div className="flex gap-2">
              <Button
                onClick={() => {
                  void handleSave("new");
                }}
                disabled={isLoading}
              >
                Save New Writer
              </Button>
              <Button variant="secondary" onClick={handleCancel} disabled={isLoading}>
                Cancel
              </Button>
            </div>
          </div>
        )}
      </div>

      <DeleteConfirmDialog
        isOpen={deleteConfirmOpen}
        onClose={() => {
          setDeleteConfirmOpen(false);
          setWriterToDelete(null);
          clearError();
        }}
        onConfirm={handleDeleteConfirm}
        entityName="writer"
        entityDetails={writerToDelete ? `${writerToDelete.name} (ID: ${writerToDelete.id})` : ""}
        isLoading={isLoading}
        error={error}
      />
    </div>
  );
}
