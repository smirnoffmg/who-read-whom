"use client";

import { type Column, DataTable } from "@/components/admin/DataTable";
import { DeleteConfirmDialog } from "@/components/admin/DeleteConfirmDialog";
import { Button } from "@/components/common/Button";
import { ErrorMessage } from "@/components/common/ErrorMessage";
import { WriterService } from "@/services/writerService";
import { useWorkStore } from "@/stores/workStore";
import type { CreateWorkRequest, UpdateWorkRequest, Work } from "@/types/work";
import type { Writer } from "@/types/writer";
import {
  type CSVImportResult,
  downloadCSV,
  exportWorksToCSV,
  importWorksFromCSV,
} from "@/utils/csvUtils";
import React, { useEffect, useRef, useState } from "react";

interface WorkFormData {
  title: string;
  author_id: string;
}

const initialFormData: WorkFormData = {
  title: "",
  author_id: "",
};

export default function WorksAdminPage(): React.JSX.Element {
  const { works, isLoading, error, fetchWorks, createWork, updateWork, deleteWork, clearError } =
    useWorkStore();

  const [writers, setWriters] = useState<Writer[]>([]);
  const [editingId, setEditingId] = useState<number | "new" | null>(null);
  const [formData, setFormData] = useState<WorkFormData>(initialFormData);
  const [formErrors, setFormErrors] = useState<Partial<Record<keyof WorkFormData, string>>>({});
  const [deleteConfirmOpen, setDeleteConfirmOpen] = useState<boolean>(false);
  const [workToDelete, setWorkToDelete] = useState<Work | null>(null);
  const [importPreview, setImportPreview] = useState<CSVImportResult<Work> | null>(null);
  const [showImportPreview, setShowImportPreview] = useState<boolean>(false);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const firstInputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    void fetchWorks(1000, 0);
  }, [fetchWorks]);

  useEffect(() => {
    const loadWriters = async (): Promise<void> => {
      try {
        const writersData = await WriterService.list(1000, 0);
        setWriters(writersData);
      } catch (error) {
        console.error("Failed to load writers:", error);
      }
    };
    void loadWriters();
  }, []);

  useEffect(() => {
    if (editingId && editingId !== "new" && firstInputRef.current) {
      firstInputRef.current.focus();
    }
  }, [editingId]);

  const getAuthorName = (authorId: number): string => {
    const writer = writers.find((w) => w.id === authorId);
    return writer ? writer.name : `ID: ${authorId}`;
  };

  const handleEdit = (id: number | string): void => {
    const work = works.find((w) => w.id === Number(id));
    if (work) {
      setEditingId(work.id);
      setFormData({
        title: work.title,
        author_id: work.author_id.toString(),
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
    const errors: Partial<Record<keyof WorkFormData, string>> = {};

    if (!formData.title.trim()) {
      errors.title = "Title is required";
    }

    if (!formData.author_id.trim()) {
      errors.author_id = "Author is required";
    } else {
      const authorId = parseInt(formData.author_id, 10);
      if (isNaN(authorId)) {
        errors.author_id = "Author ID must be a valid number";
      }
    }

    setFormErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const handleSave = async (id: number | string): Promise<void> => {
    if (!validateForm()) {
      return;
    }

    const authorId = parseInt(formData.author_id, 10);

    // Check both editingId and id parameter to handle "new" case
    if (editingId === "new" || id === "new") {
      const createData: CreateWorkRequest = {
        title: formData.title.trim(),
        author_id: authorId,
      };

      const result = await createWork(createData);
      if (result) {
        await fetchWorks(1000, 0);
        // Exit edit mode after successful creation
        setEditingId(null);
        setFormData(initialFormData);
        setFormErrors({});
      }
    } else {
      const updateData: UpdateWorkRequest = {
        title: formData.title.trim(),
        author_id: authorId,
      };

      await updateWork(Number(id), updateData);
      setEditingId(null);
      setFormData(initialFormData);
      await fetchWorks(1000, 0);
    }
  };

  const handleDeleteClick = (id: number | string): void => {
    const work = works.find((w) => w.id === Number(id));
    if (work) {
      setWorkToDelete(work);
      setDeleteConfirmOpen(true);
    }
  };

  const handleDeleteConfirm = async (): Promise<void> => {
    if (workToDelete) {
      await deleteWork(workToDelete.id);
      setDeleteConfirmOpen(false);
      setWorkToDelete(null);
      void fetchWorks(1000, 0);
    }
  };

  const handleExportCSV = (): void => {
    const csvContent = exportWorksToCSV(works);
    const filename = `works-${new Date().toISOString().split("T")[0]}.csv`;
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
      const result = importWorksFromCSV(csvString);
      setImportPreview(result);
      setShowImportPreview(true);
    };
    reader.readAsText(file);

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

    for (const work of importPreview.data) {
      try {
        const createData: CreateWorkRequest = {
          title: work.title,
          author_id: work.author_id,
        };
        const result = await createWork(createData);
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
    void fetchWorks(1000, 0);

    alert(`Import completed: ${successCount} successful, ${errorCount} failed`);
  };

  const columns: Column<Work>[] = [
    {
      key: "id",
      label: "ID",
      sortable: true,
      render: (work, isEditing) => {
        if (isEditing) {
          return <span className="text-gray-500">{work.id}</span>;
        }
        return <span>{work.id}</span>;
      },
    },
    {
      key: "title",
      label: "Title",
      sortable: true,
      render: (work, isEditing) => {
        if (isEditing) {
          return (
            <input
              ref={editingId === work.id ? firstInputRef : undefined}
              type="text"
              value={formData.title}
              onChange={(e) => setFormData({ ...formData, title: e.target.value })}
              className="w-full px-2 py-1 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="Title"
            />
          );
        }
        return <span>{work.title}</span>;
      },
    },
    {
      key: "author_id",
      label: "Author ID",
      sortable: true,
      render: (work, isEditing) => {
        if (isEditing) {
          return (
            <select
              value={formData.author_id}
              onChange={(e) => setFormData({ ...formData, author_id: e.target.value })}
              className="w-full px-2 py-1 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="">Select Author</option>
              {writers.map((writer) => (
                <option key={writer.id} value={writer.id.toString()}>
                  {writer.name} (ID: {writer.id})
                </option>
              ))}
            </select>
          );
        }
        return <span>{work.author_id}</span>;
      },
    },
    {
      key: "author_name",
      label: "Author Name",
      sortable: false,
      render: (work) => {
        return <span className="text-gray-600">{getAuthorName(work.author_id)}</span>;
      },
    },
  ];

  const displayData = editingId === "new" ? [{ id: -1 } as Work, ...works] : works;

  return (
    <div className="flex flex-col h-full">
      <div className="p-6 bg-white border-b border-gray-200">
        <h1 className="text-2xl font-bold text-gray-900">Works</h1>
        <p className="text-sm text-gray-600 mt-1">Track literary works and their authors</p>
      </div>

      <div className="flex-1 p-6 overflow-auto">
        {error && (
          <div className="mb-4">
            <ErrorMessage message={error} onDismiss={clearError} />
          </div>
        )}

        <div className="mb-4 flex gap-3">
          <Button onClick={handleExportCSV} disabled={isLoading || works.length === 0}>
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
            <p className="text-sm text-gray-600 mb-2">
              {importPreview.data.length} rows to import
              {importPreview.errors.length > 0 && (
                <span className="text-red-600 ml-2">({importPreview.errors.length} errors)</span>
              )}
            </p>
            {importPreview.errors.length > 0 && (
              <div className="mb-4 max-h-40 overflow-y-auto">
                <ul className="list-disc list-inside text-sm text-red-600">
                  {importPreview.errors.map((error, index) => (
                    <li key={index}>
                      Row {error.row}, {error.field}: {error.message}
                    </li>
                  ))}
                </ul>
              </div>
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

        {formErrors.title && (
          <div className="mb-2 text-sm text-red-600">Title: {formErrors.title}</div>
        )}
        {formErrors.author_id && (
          <div className="mb-2 text-sm text-red-600">Author: {formErrors.author_id}</div>
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
          getRowId={(work) => work.id}
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
                Save New Work
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
          setWorkToDelete(null);
        }}
        onConfirm={handleDeleteConfirm}
        entityName="work"
        entityDetails={workToDelete ? `${workToDelete.title} (ID: ${workToDelete.id})` : ""}
        isLoading={isLoading}
      />
    </div>
  );
}
