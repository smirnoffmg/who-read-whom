import Papa from "papaparse";
import type { Writer } from "@/types/writer";
import type { Work } from "@/types/work";
import type { Opinion } from "@/types/opinion";

export interface CSVValidationError {
  row: number;
  field: string;
  message: string;
}

export interface CSVImportResult<T> {
  data: T[];
  errors: CSVValidationError[];
  isValid: boolean;
}

// Writer CSV functions
export const exportWritersToCSV = (writers: Writer[]): string => {
  const csvData = writers.map((writer) => ({
    id: writer.id.toString(),
    name: writer.name,
    birth_year: writer.birth_year.toString(),
    death_year: writer.death_year?.toString() ?? "",
    bio: writer.bio ?? "",
  }));

  return Papa.unparse(csvData, {
    header: true,
    columns: ["id", "name", "birth_year", "death_year", "bio"],
  });
};

export const importWritersFromCSV = (csvString: string): CSVImportResult<Writer> => {
  const result = Papa.parse(csvString, {
    header: true,
    skipEmptyLines: true,
    transformHeader: (header: string) => header.trim().toLowerCase(),
  });

  const errors: CSVValidationError[] = [];
  const writers: Writer[] = [];

  result.data.forEach((row: unknown, index: number) => {
    const rowData = row as Record<string, string>;
    const rowNumber = index + 2; // +2 because index is 0-based and we skip header

    // Validate required fields
    if (!rowData.name || rowData.name.trim() === "") {
      errors.push({
        row: rowNumber,
        field: "name",
        message: "Name is required",
      });
    }

    if (!rowData.birth_year || rowData.birth_year.trim() === "") {
      errors.push({
        row: rowNumber,
        field: "birth_year",
        message: "Birth year is required",
      });
    }

    // Validate birth_year is a number
    const birthYear = parseInt(rowData.birth_year ?? "", 10);
    if (isNaN(birthYear)) {
      errors.push({
        row: rowNumber,
        field: "birth_year",
        message: "Birth year must be a valid number",
      });
    }

    // Validate death_year if provided
    let deathYear: number | null = null;
    if (rowData.death_year && rowData.death_year.trim() !== "") {
      const parsed = parseInt(rowData.death_year, 10);
      if (isNaN(parsed)) {
        errors.push({
          row: rowNumber,
          field: "death_year",
          message: "Death year must be a valid number",
        });
      } else {
        deathYear = parsed;
      }
    }

    // Parse ID if provided (for updates, but we'll ignore it for imports)
    const id = rowData.id && rowData.id.trim() !== "" ? parseInt(rowData.id, 10) : 0;

    if (errors.filter((e) => e.row === rowNumber).length === 0) {
      writers.push({
        id,
        name: rowData.name.trim(),
        birth_year: birthYear,
        death_year: deathYear,
        bio: rowData.bio && rowData.bio.trim() !== "" ? rowData.bio.trim() : null,
      });
    }
  });

  // Add parsing errors from PapaParse
  result.errors.forEach((error) => {
    errors.push({
      row: error.row ?? 0,
      field: error.code ?? "unknown",
      message: error.message ?? "CSV parsing error",
    });
  });

  return {
    data: writers,
    errors,
    isValid: errors.length === 0,
  };
};

// Work CSV functions
export const exportWorksToCSV = (works: Work[]): string => {
  const csvData = works.map((work) => ({
    id: work.id.toString(),
    title: work.title,
    author_id: work.author_id.toString(),
  }));

  return Papa.unparse(csvData, {
    header: true,
    columns: ["id", "title", "author_id"],
  });
};

export const importWorksFromCSV = (csvString: string): CSVImportResult<Work> => {
  const result = Papa.parse(csvString, {
    header: true,
    skipEmptyLines: true,
    transformHeader: (header: string) => header.trim().toLowerCase(),
  });

  const errors: CSVValidationError[] = [];
  const works: Work[] = [];

  result.data.forEach((row: unknown, index: number) => {
    const rowData = row as Record<string, string>;
    const rowNumber = index + 2;

    if (!rowData.title || rowData.title.trim() === "") {
      errors.push({
        row: rowNumber,
        field: "title",
        message: "Title is required",
      });
    }

    if (!rowData.author_id || rowData.author_id.trim() === "") {
      errors.push({
        row: rowNumber,
        field: "author_id",
        message: "Author ID is required",
      });
    }

    const authorId = parseInt(rowData.author_id ?? "", 10);
    if (isNaN(authorId)) {
      errors.push({
        row: rowNumber,
        field: "author_id",
        message: "Author ID must be a valid number",
      });
    }

    const id = rowData.id && rowData.id.trim() !== "" ? parseInt(rowData.id, 10) : 0;

    if (errors.filter((e) => e.row === rowNumber).length === 0) {
      works.push({
        id,
        title: rowData.title.trim(),
        author_id: authorId,
      });
    }
  });

  result.errors.forEach((error) => {
    errors.push({
      row: error.row ?? 0,
      field: error.code ?? "unknown",
      message: error.message ?? "CSV parsing error",
    });
  });

  return {
    data: works,
    errors,
    isValid: errors.length === 0,
  };
};

// Opinion CSV functions
export const exportOpinionsToCSV = (opinions: Opinion[]): string => {
  const csvData = opinions.map((opinion) => ({
    writer_id: opinion.writer_id.toString(),
    work_id: opinion.work_id.toString(),
    sentiment: opinion.sentiment ? "true" : "false",
    quote: opinion.quote,
    source: opinion.source,
    page: opinion.page ?? "",
    statement_year: opinion.statement_year?.toString() ?? "",
  }));

  return Papa.unparse(csvData, {
    header: true,
    columns: ["writer_id", "work_id", "sentiment", "quote", "source", "page", "statement_year"],
  });
};

export const importOpinionsFromCSV = (csvString: string): CSVImportResult<Opinion> => {
  const result = Papa.parse(csvString, {
    header: true,
    skipEmptyLines: true,
    transformHeader: (header: string) => header.trim().toLowerCase(),
  });

  const errors: CSVValidationError[] = [];
  const opinions: Opinion[] = [];

  result.data.forEach((row: unknown, index: number) => {
    const rowData = row as Record<string, string>;
    const rowNumber = index + 2;

    if (!rowData.writer_id || rowData.writer_id.trim() === "") {
      errors.push({
        row: rowNumber,
        field: "writer_id",
        message: "Writer ID is required",
      });
    }

    if (!rowData.work_id || rowData.work_id.trim() === "") {
      errors.push({
        row: rowNumber,
        field: "work_id",
        message: "Work ID is required",
      });
    }

    if (!rowData.sentiment || rowData.sentiment.trim() === "") {
      errors.push({
        row: rowNumber,
        field: "sentiment",
        message: "Sentiment is required",
      });
    }

    if (!rowData.quote || rowData.quote.trim() === "") {
      errors.push({
        row: rowNumber,
        field: "quote",
        message: "Quote is required",
      });
    }

    if (!rowData.source || rowData.source.trim() === "") {
      errors.push({
        row: rowNumber,
        field: "source",
        message: "Source is required",
      });
    }

    const writerId = parseInt(rowData.writer_id ?? "", 10);
    if (isNaN(writerId)) {
      errors.push({
        row: rowNumber,
        field: "writer_id",
        message: "Writer ID must be a valid number",
      });
    }

    const workId = parseInt(rowData.work_id ?? "", 10);
    if (isNaN(workId)) {
      errors.push({
        row: rowNumber,
        field: "work_id",
        message: "Work ID must be a valid number",
      });
    }

    const sentiment = rowData.sentiment?.toLowerCase().trim();
    if (sentiment !== "true" && sentiment !== "false" && sentiment !== "1" && sentiment !== "0") {
      errors.push({
        row: rowNumber,
        field: "sentiment",
        message: "Sentiment must be true/false or 1/0",
      });
    }

    let statementYear: number | null = null;
    if (rowData.statement_year && rowData.statement_year.trim() !== "") {
      const parsed = parseInt(rowData.statement_year, 10);
      if (isNaN(parsed)) {
        errors.push({
          row: rowNumber,
          field: "statement_year",
          message: "Statement year must be a valid number",
        });
      } else {
        statementYear = parsed;
      }
    }

    if (errors.filter((e) => e.row === rowNumber).length === 0) {
      const sentimentValue =
        sentiment === "true" || sentiment === "1"
          ? true
          : sentiment === "false" || sentiment === "0"
            ? false
            : false;

      opinions.push({
        writer_id: writerId,
        work_id: workId,
        sentiment: sentimentValue,
        quote: rowData.quote.trim(),
        source: rowData.source.trim(),
        page: rowData.page && rowData.page.trim() !== "" ? rowData.page.trim() : null,
        statement_year: statementYear,
      });
    }
  });

  result.errors.forEach((error) => {
    errors.push({
      row: error.row ?? 0,
      field: error.code ?? "unknown",
      message: error.message ?? "CSV parsing error",
    });
  });

  return {
    data: opinions,
    errors,
    isValid: errors.length === 0,
  };
};

// Generic CSV download helper
export const downloadCSV = (csvContent: string, filename: string): void => {
  const blob = new Blob([csvContent], { type: "text/csv;charset=utf-8;" });
  const link = document.createElement("a");
  const url = URL.createObjectURL(blob);

  link.setAttribute("href", url);
  link.setAttribute("download", filename);
  link.style.visibility = "hidden";
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
};
