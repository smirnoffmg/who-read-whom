import { create } from "zustand";
import { WriterService } from "@/services/writerService";
import type { CreateWriterRequest, UpdateWriterRequest, Writer } from "@/types/writer";

interface WriterState {
  writers: Writer[];
  isLoading: boolean;
  error: string | null;
}

interface WriterActions {
  fetchWriters: (limit?: number, offset?: number) => Promise<void>;
  fetchWriterById: (id: number) => Promise<Writer | null>;
  createWriter: (params: CreateWriterRequest) => Promise<Writer | null>;
  updateWriter: (id: number, params: UpdateWriterRequest) => Promise<void>;
  deleteWriter: (id: number) => Promise<void>;
  clearError: () => void;
}

export const useWriterStore = create<WriterState & WriterActions>((set) => ({
  writers: [],
  isLoading: false,
  error: null,

  fetchWriters: async (limit = 10, offset = 0) => {
    set({ isLoading: true, error: null });
    try {
      const writers = await WriterService.list(limit, offset);
      set({ writers, isLoading: false });
    } catch (error) {
      set({
        error: error instanceof Error ? error.message : "Failed to fetch writers",
        isLoading: false,
      });
    }
  },

  fetchWriterById: async (id: number) => {
    set({ isLoading: true, error: null });
    try {
      const writer = await WriterService.getById(id);
      set({ isLoading: false });
      return writer;
    } catch (error) {
      set({
        error: error instanceof Error ? error.message : "Failed to fetch writer",
        isLoading: false,
      });
      return null;
    }
  },

  createWriter: async (params: CreateWriterRequest) => {
    set({ isLoading: true, error: null });
    try {
      const writer = await WriterService.create(params);
      set((state) => ({
        writers: [...state.writers, writer],
        isLoading: false,
      }));
      return writer;
    } catch (error) {
      set({
        error: error instanceof Error ? error.message : "Failed to create writer",
        isLoading: false,
      });
      return null;
    }
  },

  updateWriter: async (id: number, params: UpdateWriterRequest) => {
    set({ isLoading: true, error: null });
    try {
      await WriterService.update(id, params);
      const updatedWriter = await WriterService.getById(id);
      set((state) => ({
        writers: state.writers.map((w) => (w.id === id ? updatedWriter : w)),
        isLoading: false,
      }));
    } catch (error) {
      set({
        error: error instanceof Error ? error.message : "Failed to update writer",
        isLoading: false,
      });
    }
  },

  deleteWriter: async (id: number) => {
    set({ isLoading: true, error: null });
    try {
      await WriterService.delete(id);
      set((state) => ({
        writers: state.writers.filter((w) => w.id !== id),
        isLoading: false,
      }));
    } catch (error) {
      set({
        error: error instanceof Error ? error.message : "Failed to delete writer",
        isLoading: false,
      });
    }
  },

  clearError: () => {
    set({ error: null });
  },
}));
