"use client";

import { useEffect, useMemo, useState } from "react";
import { WriterService } from "@/services/writerService";
import { WorkService } from "@/services/workService";
import type { Writer } from "@/types/writer";
import type { Work } from "@/types/work";

interface SearchBarProps {
  onWriterSelect: (writer: Writer | null) => void;
  onWorkSelect: (work: Work | null) => void;
  selectedWriter: Writer | null;
  selectedWork: Work | null;
}

export const SearchBar: React.FC<SearchBarProps> = ({
  onWriterSelect,
  onWorkSelect,
  selectedWriter,
  selectedWork,
}): React.JSX.Element => {
  const [searchQuery, setSearchQuery] = useState<string>("");
  const [searchType, setSearchType] = useState<"writer" | "work">("writer");
  const [results, setResults] = useState<(Writer | Work)[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!searchQuery.trim()) {
      setResults([]);
      return;
    }

    const debounceTimer = setTimeout(() => {
      const performSearch = async (): Promise<void> => {
        setIsLoading(true);
        setError(null);
        try {
          if (searchType === "writer") {
            const writersData = await WriterService.search(searchQuery.trim(), 20, 0);
            setResults(writersData);
          } else {
            const worksData = await WorkService.search(searchQuery.trim(), 20, 0);
            setResults(worksData);
          }
        } catch (err) {
          setError(err instanceof Error ? err.message : "Search failed");
          setResults([]);
        } finally {
          setIsLoading(false);
        }
      };

      void performSearch();
    }, 300); // 300ms debounce

    return () => {
      clearTimeout(debounceTimer);
    };
  }, [searchQuery, searchType]);

  const showResults = useMemo(
    () => results.length > 0 && searchQuery.trim().length > 0 && !isLoading,
    [results, searchQuery, isLoading]
  );

  const handleSelect = (item: Writer | Work): void => {
    if (searchType === "writer") {
      onWriterSelect(item as Writer);
      onWorkSelect(null);
    } else {
      onWorkSelect(item as Work);
      onWriterSelect(null);
    }
    setSearchQuery("");
  };

  const handleClear = (): void => {
    setSearchQuery("");
    onWriterSelect(null);
    onWorkSelect(null);
  };

  return (
    <div className="relative">
      <div className="flex gap-2">
        <div className="flex-1">
          <div className="relative">
            <input
              type="text"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              onFocus={() => {
                // Focus handled by showResults computed value
              }}
              placeholder={`Search ${searchType === "writer" ? "writers" : "works"}...`}
              className="w-full rounded-md border border-gray-300 px-4 py-2 pl-10 focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <div className="absolute left-3 top-1/2 -translate-y-1/2">
              <svg
                className="h-5 w-5 text-gray-400"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                />
              </svg>
            </div>
          </div>

          {isLoading && searchQuery.trim().length > 0 && (
            <div className="absolute z-10 mt-1 w-full rounded-md border border-gray-200 bg-white px-4 py-2 text-sm text-gray-500">
              Searching...
            </div>
          )}
          {error && searchQuery.trim().length > 0 && (
            <div className="absolute z-10 mt-1 w-full rounded-md border border-red-200 bg-red-50 px-4 py-2 text-sm text-red-600">
              {error}
            </div>
          )}
          {showResults && (
            <div className="absolute z-10 mt-1 max-h-60 w-full overflow-auto rounded-md border border-gray-200 bg-white shadow-lg">
              {results.map((item) => (
                <button
                  key={searchType === "writer" ? (item as Writer).id : (item as Work).id}
                  type="button"
                  onClick={() => handleSelect(item)}
                  className="w-full px-4 py-2 text-left hover:bg-gray-100"
                >
                  {searchType === "writer" ? (
                    <div>
                      <div className="font-medium text-gray-900">{(item as Writer).name}</div>
                      {(item as Writer).birth_year && (
                        <div className="text-sm text-gray-500">
                          {(item as Writer).birth_year}
                          {(item as Writer).death_year ? ` - ${(item as Writer).death_year}` : ""}
                        </div>
                      )}
                    </div>
                  ) : (
                    <div className="font-medium text-gray-900">{(item as Work).title}</div>
                  )}
                </button>
              ))}
            </div>
          )}
        </div>

        <div className="flex gap-2">
          <button
            type="button"
            onClick={() => {
              setSearchType("writer");
              setSearchQuery("");
            }}
            className={`rounded-md px-4 py-2 text-sm font-medium ${
              searchType === "writer"
                ? "bg-blue-600 text-white"
                : "bg-gray-200 text-gray-700 hover:bg-gray-300"
            }`}
          >
            Writers
          </button>
          <button
            type="button"
            onClick={() => {
              setSearchType("work");
              setSearchQuery("");
            }}
            className={`rounded-md px-4 py-2 text-sm font-medium ${
              searchType === "work"
                ? "bg-blue-600 text-white"
                : "bg-gray-200 text-gray-700 hover:bg-gray-300"
            }`}
          >
            Works
          </button>
        </div>

        {(selectedWriter || selectedWork) && (
          <button
            type="button"
            onClick={handleClear}
            className="rounded-md bg-red-600 px-4 py-2 text-sm font-medium text-white hover:bg-red-700"
          >
            Clear
          </button>
        )}
      </div>

      {(selectedWriter || selectedWork) && (
        <div className="mt-2 rounded-md bg-blue-50 px-4 py-2">
          <span className="text-sm text-gray-600">Selected: </span>
          <span className="text-sm font-medium text-blue-900">
            {selectedWriter ? selectedWriter.name : selectedWork?.title}
          </span>
        </div>
      )}
    </div>
  );
};
