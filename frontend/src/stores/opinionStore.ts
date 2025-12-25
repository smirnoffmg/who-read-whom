import { create } from "zustand";
import { OpinionService } from "@/services/opinionService";
import type { CreateOpinionRequest, Opinion, UpdateOpinionRequest } from "@/types/opinion";

interface OpinionState {
  opinions: Opinion[];
  isLoading: boolean;
  error: string | null;
}

interface OpinionActions {
  fetchOpinions: (limit?: number, offset?: number) => Promise<void>;
  fetchOpinionByWriterAndWork: (writerId: number, workId: number) => Promise<Opinion | null>;
  fetchOpinionsByWriter: (writerId: number) => Promise<Opinion[]>;
  fetchOpinionsByWork: (workId: number) => Promise<Opinion[]>;
  createOpinion: (params: CreateOpinionRequest) => Promise<Opinion | null>;
  updateOpinion: (writerId: number, workId: number, params: UpdateOpinionRequest) => Promise<void>;
  deleteOpinion: (writerId: number, workId: number) => Promise<void>;
  clearError: () => void;
}

export const useOpinionStore = create<OpinionState & OpinionActions>((set) => ({
  opinions: [],
  isLoading: false,
  error: null,

  fetchOpinions: async (limit = 10, offset = 0) => {
    set({ isLoading: true, error: null });
    try {
      const opinions = await OpinionService.list(limit, offset);
      set({ opinions, isLoading: false });
    } catch (error) {
      set({
        error: error instanceof Error ? error.message : "Failed to fetch opinions",
        isLoading: false,
      });
    }
  },

  fetchOpinionByWriterAndWork: async (writerId: number, workId: number) => {
    set({ isLoading: true, error: null });
    try {
      const opinion = await OpinionService.getByWriterAndWork(writerId, workId);
      set({ isLoading: false });
      return opinion;
    } catch (error) {
      set({
        error: error instanceof Error ? error.message : "Failed to fetch opinion",
        isLoading: false,
      });
      return null;
    }
  },

  fetchOpinionsByWriter: async (writerId: number) => {
    set({ isLoading: true, error: null });
    try {
      const opinions = await OpinionService.getByWriter(writerId);
      set({ isLoading: false });
      return opinions;
    } catch (error) {
      set({
        error: error instanceof Error ? error.message : "Failed to fetch opinions by writer",
        isLoading: false,
      });
      return [];
    }
  },

  fetchOpinionsByWork: async (workId: number) => {
    set({ isLoading: true, error: null });
    try {
      const opinions = await OpinionService.getByWork(workId);
      set({ isLoading: false });
      return opinions;
    } catch (error) {
      set({
        error: error instanceof Error ? error.message : "Failed to fetch opinions by work",
        isLoading: false,
      });
      return [];
    }
  },

  createOpinion: async (params: CreateOpinionRequest) => {
    set({ isLoading: true, error: null });
    try {
      const opinion = await OpinionService.create(params);
      set((state) => ({
        opinions: [...state.opinions, opinion],
        isLoading: false,
      }));
      return opinion;
    } catch (error) {
      set({
        error: error instanceof Error ? error.message : "Failed to create opinion",
        isLoading: false,
      });
      return null;
    }
  },

  updateOpinion: async (writerId: number, workId: number, params: UpdateOpinionRequest) => {
    set({ isLoading: true, error: null });
    try {
      await OpinionService.update(writerId, workId, params);
      const updatedOpinion = await OpinionService.getByWriterAndWork(writerId, workId);
      set((state) => ({
        opinions: state.opinions.map((o) =>
          o.writer_id === writerId && o.work_id === workId ? updatedOpinion : o
        ),
        isLoading: false,
      }));
    } catch (error) {
      set({
        error: error instanceof Error ? error.message : "Failed to update opinion",
        isLoading: false,
      });
    }
  },

  deleteOpinion: async (writerId: number, workId: number) => {
    set({ isLoading: true, error: null });
    try {
      await OpinionService.delete(writerId, workId);
      set((state) => ({
        opinions: state.opinions.filter((o) => !(o.writer_id === writerId && o.work_id === workId)),
        isLoading: false,
      }));
    } catch (error) {
      set({
        error: error instanceof Error ? error.message : "Failed to delete opinion",
        isLoading: false,
      });
    }
  },

  clearError: () => {
    set({ error: null });
  },
}));
