import type { CreateWorkRequest, UpdateWorkRequest, Work } from "@/types/work";

export class WorkService {
  private static readonly BASE_URL =
    process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

  static async create(params: CreateWorkRequest): Promise<Work> {
    const response = await fetch(`${this.BASE_URL}/works`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(params),
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: response.statusText }));
      throw new Error(error.error || `WorkService.create failed: ${response.statusText}`);
    }

    return response.json();
  }

  static async getById(id: number): Promise<Work> {
    const numericId = Number(id);
    if (isNaN(numericId) || numericId <= 0) {
      throw new Error(`Invalid work ID: ${id}`);
    }
    const response = await fetch(`${this.BASE_URL}/works/${numericId}`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: response.statusText }));
      throw new Error(error.error || `WorkService.getById failed: ${response.statusText}`);
    }

    return response.json();
  }

  static async getByAuthor(authorId: number): Promise<Work[]> {
    const response = await fetch(`${this.BASE_URL}/works/author/${authorId}`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: response.statusText }));
      throw new Error(error.error || `WorkService.getByAuthor failed: ${response.statusText}`);
    }

    return response.json();
  }

  static async list(limit: number = 10, offset: number = 0): Promise<Work[]> {
    const response = await fetch(`${this.BASE_URL}/works?limit=${limit}&offset=${offset}`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: response.statusText }));
      throw new Error(error.error || `WorkService.list failed: ${response.statusText}`);
    }

    return response.json();
  }

  static async search(query: string, limit: number = 20, offset: number = 0): Promise<Work[]> {
    const searchParam = encodeURIComponent(query);
    const response = await fetch(
      `${this.BASE_URL}/works?search=${searchParam}&limit=${limit}&offset=${offset}`,
      {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      }
    );

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: response.statusText }));
      throw new Error(error.error || `WorkService.search failed: ${response.statusText}`);
    }

    return response.json();
  }

  static async update(id: number, params: UpdateWorkRequest): Promise<void> {
    const numericId = Number(id);
    if (isNaN(numericId) || numericId <= 0) {
      throw new Error(`Invalid work ID: ${id}`);
    }
    const response = await fetch(`${this.BASE_URL}/works/${numericId}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(params),
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: response.statusText }));
      throw new Error(error.error || `WorkService.update failed: ${response.statusText}`);
    }
  }

  static async delete(id: number): Promise<void> {
    const numericId = Number(id);
    if (isNaN(numericId) || numericId <= 0) {
      throw new Error(`Invalid work ID: ${id}`);
    }
    const response = await fetch(`${this.BASE_URL}/works/${numericId}`, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: response.statusText }));
      throw new Error(error.error || `WorkService.delete failed: ${response.statusText}`);
    }
  }
}
