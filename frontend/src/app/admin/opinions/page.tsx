"use client";

import React, { useEffect, useRef, useState } from "react";
import { type Column, DataTable } from "@/components/admin/DataTable";
import { DeleteConfirmDialog } from "@/components/admin/DeleteConfirmDialog";
import { Button } from "@/components/common/Button";
import { ErrorMessage } from "@/components/common/ErrorMessage";
import { WriterService } from "@/services/writerService";
import { WorkService } from "@/services/workService";
import { useOpinionStore } from "@/stores/opinionStore";
import type { CreateOpinionRequest, Opinion, UpdateOpinionRequest } from "@/types/opinion";
import type { Work } from "@/types/work";
import type { Writer } from "@/types/writer";
import {
  type CSVImportResult,
  downloadCSV,
  exportOpinionsToCSV,
  importOpinionsFromCSV,
} from "@/utils/csvUtils";

interface OpinionFormData {
  writer_id: string;
  work_id: string;
  sentiment: string;
  quote: string;
  source: string;
  page: string;
  statement_year: string;
}

const initialFormData: OpinionFormData = {
  writer_id: "",
  work_id: "",
  sentiment: "true",
  quote: "",
  source: "",
  page: "",
  statement_year: "",
};

const getOpinionKey = (opinion: Opinion): string => {
  return `${opinion.writer_id}-${opinion.work_id}`;
};

export default function OpinionsAdminPage(): React.JSX.Element {
  const {
    opinions,
    isLoading,
    error,
    fetchOpinions,
    createOpinion,
    updateOpinion,
    deleteOpinion,
    clearError,
  } = useOpinionStore();

  const [writers, setWriters] = useState<Writer[]>([]);
  const [works, setWorks] = useState<Work[]>([]);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [formData, setFormData] = useState<OpinionFormData>(initialFormData);
  const [formErrors, setFormErrors] = useState<Partial<Record<keyof OpinionFormData, string>>>({});
  const [deleteConfirmOpen, setDeleteConfirmOpen] = useState<boolean>(false);
  const [opinionToDelete, setOpinionToDelete] = useState<Opinion | null>(null);
  const [importPreview, setImportPreview] = useState<CSVImportResult<Opinion> | null>(null);
  const [showImportPreview, setShowImportPreview] = useState<boolean>(false);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const firstInputRef = useRef<HTMLSelectElement>(null);

  useEffect(() => {
    void fetchOpinions(1000, 0);
  }, [fetchOpinions]);

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
    const loadWorks = async (): Promise<void> => {
      try {
        const worksData = await WorkService.list(1000, 0);
        setWorks(worksData);
      } catch (error) {
        console.error("Failed to load works:", error);
      }
    };
    void loadWorks();
  }, []);

  useEffect(() => {
    if (editingId && editingId !== "new" && firstInputRef.current) {
      firstInputRef.current.focus();
    }
  }, [editingId]);

  const getWriterName = (writerId: number): string => {
    const writer = writers.find((w) => w.id === writerId);
    return writer ? writer.name : `ID: ${writerId}`;
  };

  const getWorkTitle = (workId: number): string => {
    const work = works.find((w) => w.id === workId);
    return work ? work.title : `ID: ${workId}`;
  };

  const getWorkAuthorId = (workId: number): number | null => {
    const work = works.find((w) => w.id === workId);
    return work ? work.author_id : null;
  };

  const handleEdit = (id: string | number): void => {
    const [writerIdStr, workIdStr] = String(id).split("-");
    const writerId = parseInt(writerIdStr, 10);
    const workId = parseInt(workIdStr, 10);

    const opinion = opinions.find((o) => o.writer_id === writerId && o.work_id === workId);
    if (opinion) {
      setEditingId(getOpinionKey(opinion));
      setFormData({
        writer_id: opinion.writer_id.toString(),
        work_id: opinion.work_id.toString(),
        sentiment: opinion.sentiment ? "true" : "false",
        quote: opinion.quote,
        source: opinion.source,
        page: opinion.page ?? "",
        statement_year: opinion.statement_year?.toString() ?? "",
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
    const errors: Partial<Record<keyof OpinionFormData, string>> = {};

    if (!formData.writer_id.trim()) {
      errors.writer_id = "Writer is required";
    }

    if (!formData.work_id.trim()) {
      errors.work_id = "Work is required";
    }

    // Validate writer is not the work's author
    if (formData.writer_id && formData.work_id) {
      const writerId = parseInt(formData.writer_id, 10);
      const workId = parseInt(formData.work_id, 10);
      const workAuthorId = getWorkAuthorId(workId);

      if (workAuthorId === writerId) {
        errors.writer_id = "Writer cannot express an opinion about their own work";
      }
    }

    if (!formData.quote.trim()) {
      errors.quote = "Quote is required";
    }

    if (!formData.source.trim()) {
      errors.source = "Source is required";
    }

    if (formData.statement_year.trim()) {
      const year = parseInt(formData.statement_year, 10);
      if (isNaN(year)) {
        errors.statement_year = "Statement year must be a valid number";
      }
    }

    setFormErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const handleSave = async (id: string | number): Promise<void> => {
    if (!validateForm()) {
      return;
    }

    const writerId = parseInt(formData.writer_id, 10);
    const workId = parseInt(formData.work_id, 10);
    const sentiment = formData.sentiment === "true";
    const statementYear = formData.statement_year.trim()
      ? parseInt(formData.statement_year, 10)
      : null;

    // Check both editingId and id parameter to handle "new" case
    if (editingId === "new" || id === "new") {
      const createData: CreateOpinionRequest = {
        writer_id: writerId,
        work_id: workId,
        sentiment,
        quote: formData.quote.trim(),
        source: formData.source.trim(),
        page: formData.page.trim() || null,
        statement_year: statementYear,
      };

      const result = await createOpinion(createData);
      if (result) {
        await fetchOpinions(1000, 0);
        // Exit edit mode after successful creation
        setEditingId(null);
        setFormData(initialFormData);
        setFormErrors({});
      }
    } else {
      const [oldWriterIdStr, oldWorkIdStr] = String(editingId).split("-");
      const oldWriterId = parseInt(oldWriterIdStr, 10);
      const oldWorkId = parseInt(oldWorkIdStr, 10);

      const updateData: UpdateOpinionRequest = {
        sentiment,
        quote: formData.quote.trim(),
        source: formData.source.trim(),
        page: formData.page.trim() || null,
        statement_year: statementYear,
      };

      await updateOpinion(oldWriterId, oldWorkId, updateData);
      setEditingId(null);
      setFormData(initialFormData);
      await fetchOpinions(1000, 0);
    }
  };

  const handleDeleteClick = (id: string | number): void => {
    const [writerIdStr, workIdStr] = String(id).split("-");
    const writerId = parseInt(writerIdStr, 10);
    const workId = parseInt(workIdStr, 10);

    const opinion = opinions.find((o) => o.writer_id === writerId && o.work_id === workId);
    if (opinion) {
      setOpinionToDelete(opinion);
      setDeleteConfirmOpen(true);
    }
  };

  const handleDeleteConfirm = async (): Promise<void> => {
    if (opinionToDelete) {
      await deleteOpinion(opinionToDelete.writer_id, opinionToDelete.work_id);
      setDeleteConfirmOpen(false);
      setOpinionToDelete(null);
      void fetchOpinions(1000, 0);
    }
  };

  const handleExportCSV = (): void => {
    const csvContent = exportOpinionsToCSV(opinions);
    const filename = `opinions-${new Date().toISOString().split("T")[0]}.csv`;
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
      const result = importOpinionsFromCSV(csvString);
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

    for (const opinion of importPreview.data) {
      try {
        const createData: CreateOpinionRequest = {
          writer_id: opinion.writer_id,
          work_id: opinion.work_id,
          sentiment: opinion.sentiment,
          quote: opinion.quote,
          source: opinion.source,
          page: opinion.page,
          statement_year: opinion.statement_year,
        };
        const result = await createOpinion(createData);
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
    void fetchOpinions(1000, 0);

    alert(`Import completed: ${successCount} successful, ${errorCount} failed`);
  };

  const truncateText = (text: string, maxLength: number): string => {
    return text.length > maxLength ? `${text.substring(0, maxLength)}...` : text;
  };

  const columns: Column<Opinion>[] = [
    {
      key: "writer_id",
      label: "Writer ID",
      sortable: true,
      render: (opinion, isEditing) => {
        if (isEditing) {
          return (
            <select
              ref={editingId === getOpinionKey(opinion) ? firstInputRef : undefined}
              value={formData.writer_id}
              onChange={(e) => setFormData({ ...formData, writer_id: e.target.value })}
              className="w-full px-2 py-1 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="">Select Writer</option>
              {writers.map((writer) => (
                <option key={writer.id} value={writer.id.toString()}>
                  {writer.name} (ID: {writer.id})
                </option>
              ))}
            </select>
          );
        }
        return <span>{opinion.writer_id}</span>;
      },
    },
    {
      key: "writer_name",
      label: "Writer Name",
      sortable: false,
      render: (opinion) => {
        return <span className="text-gray-600">{getWriterName(opinion.writer_id)}</span>;
      },
    },
    {
      key: "work_id",
      label: "Work ID",
      sortable: true,
      render: (opinion, isEditing) => {
        if (isEditing) {
          return (
            <select
              value={formData.work_id}
              onChange={(e) => setFormData({ ...formData, work_id: e.target.value })}
              className="w-full px-2 py-1 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="">Select Work</option>
              {works.map((work) => (
                <option key={work.id} value={work.id.toString()}>
                  {work.title} (ID: {work.id})
                </option>
              ))}
            </select>
          );
        }
        return <span>{opinion.work_id}</span>;
      },
    },
    {
      key: "work_title",
      label: "Work Title",
      sortable: false,
      render: (opinion) => {
        return <span className="text-gray-600">{getWorkTitle(opinion.work_id)}</span>;
      },
    },
    {
      key: "sentiment",
      label: "Sentiment",
      sortable: true,
      render: (opinion, isEditing) => {
        if (isEditing) {
          return (
            <div className="flex gap-4">
              <label className="flex items-center">
                <input
                  type="radio"
                  name={`sentiment-${getOpinionKey(opinion)}`}
                  value="true"
                  checked={formData.sentiment === "true"}
                  onChange={(e) => setFormData({ ...formData, sentiment: e.target.value })}
                  className="mr-2"
                />
                <span className="text-green-600">Positive</span>
              </label>
              <label className="flex items-center">
                <input
                  type="radio"
                  name={`sentiment-${getOpinionKey(opinion)}`}
                  value="false"
                  checked={formData.sentiment === "false"}
                  onChange={(e) => setFormData({ ...formData, sentiment: e.target.value })}
                  className="mr-2"
                />
                <span className="text-red-600">Negative</span>
              </label>
            </div>
          );
        }
        return (
          <span className={`font-medium ${opinion.sentiment ? "text-green-600" : "text-red-600"}`}>
            {opinion.sentiment ? "Positive" : "Negative"}
          </span>
        );
      },
    },
    {
      key: "quote",
      label: "Quote",
      sortable: false,
      render: (opinion, isEditing) => {
        if (isEditing) {
          return (
            <textarea
              value={formData.quote}
              onChange={(e) => setFormData({ ...formData, quote: e.target.value })}
              className="w-full px-2 py-1 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="Quote"
              rows={2}
            />
          );
        }
        return <span className="text-gray-600">{truncateText(opinion.quote, 50)}</span>;
      },
    },
    {
      key: "source",
      label: "Source",
      sortable: true,
      render: (opinion, isEditing) => {
        if (isEditing) {
          return (
            <input
              type="text"
              value={formData.source}
              onChange={(e) => setFormData({ ...formData, source: e.target.value })}
              className="w-full px-2 py-1 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="Source"
            />
          );
        }
        return <span>{truncateText(opinion.source, 30)}</span>;
      },
    },
    {
      key: "page",
      label: "Page",
      sortable: false,
      render: (opinion, isEditing) => {
        if (isEditing) {
          return (
            <input
              type="text"
              value={formData.page}
              onChange={(e) => setFormData({ ...formData, page: e.target.value })}
              className="w-full px-2 py-1 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="Page (optional)"
            />
          );
        }
        return <span>{opinion.page ?? "—"}</span>;
      },
    },
    {
      key: "statement_year",
      label: "Statement Year",
      sortable: true,
      render: (opinion, isEditing) => {
        if (isEditing) {
          return (
            <input
              type="number"
              value={formData.statement_year}
              onChange={(e) => setFormData({ ...formData, statement_year: e.target.value })}
              className="w-full px-2 py-1 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="Year (optional)"
            />
          );
        }
        return <span>{opinion.statement_year ?? "—"}</span>;
      },
    },
  ];

  const displayData =
    editingId === "new" ? [{ writer_id: -1, work_id: -1 } as Opinion, ...opinions] : opinions;

  return (
    <div className="flex flex-col h-full">
      <div className="p-6 bg-white border-b border-gray-200">
        <h1 className="text-2xl font-bold text-gray-900">Opinions</h1>
        <p className="text-sm text-gray-600 mt-1">Document opinions with quotes and sources</p>
      </div>

      <div className="flex-1 p-6 overflow-auto">
        {error && (
          <div className="mb-4">
            <ErrorMessage message={error} onDismiss={clearError} />
          </div>
        )}

        <div className="mb-4 flex gap-3">
          <Button onClick={handleExportCSV} disabled={isLoading || opinions.length === 0}>
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

        {formErrors.writer_id && (
          <div className="mb-2 text-sm text-red-600">Writer: {formErrors.writer_id}</div>
        )}
        {formErrors.work_id && (
          <div className="mb-2 text-sm text-red-600">Work: {formErrors.work_id}</div>
        )}
        {formErrors.quote && (
          <div className="mb-2 text-sm text-red-600">Quote: {formErrors.quote}</div>
        )}
        {formErrors.source && (
          <div className="mb-2 text-sm text-red-600">Source: {formErrors.source}</div>
        )}
        {formErrors.statement_year && (
          <div className="mb-2 text-sm text-red-600">
            Statement Year: {formErrors.statement_year}
          </div>
        )}

        <div className="overflow-x-auto">
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
            getRowId={(opinion) => getOpinionKey(opinion)}
          />
        </div>

        {editingId === "new" && (
          <div className="mt-4 bg-blue-50 p-4 rounded-lg border border-blue-200">
            <div className="flex gap-2">
              <Button
                onClick={() => {
                  void handleSave("new");
                }}
                disabled={isLoading}
              >
                Save New Opinion
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
          setOpinionToDelete(null);
        }}
        onConfirm={handleDeleteConfirm}
        entityName="opinion"
        entityDetails={
          opinionToDelete
            ? `Writer ${getWriterName(opinionToDelete.writer_id)} about "${getWorkTitle(opinionToDelete.work_id)}"`
            : ""
        }
        isLoading={isLoading}
      />
    </div>
  );
}
