import { create } from "zustand";
import { WorkService } from "@/services/workService";
import type { CreateWorkRequest, UpdateWorkRequest, Work } from "@/types/work";

interface WorkState {
  works: Work[];
  isLoading: boolean;
  error: string | null;
}

interface WorkActions {
  fetchWorks: (limit?: number, offset?: number) => Promise<void>;
  fetchWorkById: (id: number) => Promise<Work | null>;
  fetchWorksByAuthor: (authorId: number) => Promise<Work[]>;
  createWork: (params: CreateWorkRequest) => Promise<Work | null>;
  updateWork: (id: number, params: UpdateWorkRequest) => Promise<void>;
  deleteWork: (id: number) => Promise<void>;
  clearError: () => void;
}

export const useWorkStore = create<WorkState & WorkActions>((set) => ({
  works: [],
  isLoading: false,
  error: null,

  fetchWorks: async (limit = 10, offset = 0) => {
    set({ isLoading: true, error: null });
    try {
      const works = await WorkService.list(limit, offset);
      set({ works, isLoading: false });
    } catch (error) {
      set({
        error: error instanceof Error ? error.message : "Failed to fetch works",
        isLoading: false,
      });
    }
  },

  fetchWorkById: async (id: number) => {
    set({ isLoading: true, error: null });
    try {
      const work = await WorkService.getById(id);
      set({ isLoading: false });
      return work;
    } catch (error) {
      set({
        error: error instanceof Error ? error.message : "Failed to fetch work",
        isLoading: false,
      });
      return null;
    }
  },

  fetchWorksByAuthor: async (authorId: number) => {
    set({ isLoading: true, error: null });
    try {
      const works = await WorkService.getByAuthor(authorId);
      set({ isLoading: false });
      return works;
    } catch (error) {
      set({
        error: error instanceof Error ? error.message : "Failed to fetch works by author",
        isLoading: false,
      });
      return [];
    }
  },

  createWork: async (params: CreateWorkRequest) => {
    set({ isLoading: true, error: null });
    try {
      const work = await WorkService.create(params);
      set((state) => ({
        works: [...state.works, work],
        isLoading: false,
      }));
      return work;
    } catch (error) {
      set({
        error: error instanceof Error ? error.message : "Failed to create work",
        isLoading: false,
      });
      return null;
    }
  },

  updateWork: async (id: number, params: UpdateWorkRequest) => {
    set({ isLoading: true, error: null });
    try {
      await WorkService.update(id, params);
      const updatedWork = await WorkService.getById(id);
      set((state) => ({
        works: state.works.map((w) => (w.id === id ? updatedWork : w)),
        isLoading: false,
      }));
    } catch (error) {
      set({
        error: error instanceof Error ? error.message : "Failed to update work",
        isLoading: false,
      });
    }
  },

  deleteWork: async (id: number) => {
    set({ isLoading: true, error: null });
    try {
      await WorkService.delete(id);
      set((state) => ({
        works: state.works.filter((w) => w.id !== id),
        isLoading: false,
      }));
    } catch (error) {
      set({
        error: error instanceof Error ? error.message : "Failed to delete work",
        isLoading: false,
      });
    }
  },

  clearError: () => {
    set({ error: null });
  },
}));
