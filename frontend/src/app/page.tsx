"use client";

import { LiteraryGraph } from "@/components/graph/LiteraryGraph";
import { SearchBar } from "@/components/graph/SearchBar";
import type { Work } from "@/types/work";
import type { Writer } from "@/types/writer";
import React, { useState } from "react";

export default function Home(): React.JSX.Element {
  const [selectedWriter, setSelectedWriter] = useState<Writer | null>(null);
  const [selectedWork, setSelectedWork] = useState<Work | null>(null);

  const handleWriterSelect = (writer: Writer | null): void => {
    setSelectedWriter(writer);
    // SearchBar already handles clearing selectedWork
  };

  const handleWorkSelect = (work: Work | null): void => {
    setSelectedWork(work);
    // SearchBar already handles clearing selectedWriter
  };

  return (
    <div className="flex h-screen flex-col overflow-hidden bg-gray-50">
      <header className="flex-shrink-0 bg-white shadow-sm">
        <div className="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
          <h1 className="text-3xl font-bold text-gray-900">Literary Opinions Graph</h1>
          <p className="mt-2 text-sm text-gray-600">
            Explore relationships between writers and their opinions about literary works
          </p>
        </div>
      </header>

      <main className="flex flex-1 flex-col overflow-hidden">
        <div className="mx-auto flex h-full w-full max-w-7xl flex-col px-4 py-6 sm:px-6 lg:px-8">
          <div className="mb-6 flex-shrink-0">
            <SearchBar
              onWriterSelect={handleWriterSelect}
              onWorkSelect={handleWorkSelect}
              selectedWriter={selectedWriter}
              selectedWork={selectedWork}
            />
          </div>

          <div className="flex-1 rounded-lg border border-gray-200 bg-white shadow-sm">
            <LiteraryGraph selectedWriter={selectedWriter} selectedWork={selectedWork} />
          </div>
        </div>
      </main>
    </div>
  );
}
