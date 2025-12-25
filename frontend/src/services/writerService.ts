import type { CreateWriterRequest, UpdateWriterRequest, Writer } from "@/types/writer";

export class WriterService {
  private static readonly BASE_URL =
    process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

  static async create(params: CreateWriterRequest): Promise<Writer> {
    const response = await fetch(`${this.BASE_URL}/writers`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(params),
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: response.statusText }));
      throw new Error(error.error || `WriterService.create failed: ${response.statusText}`);
    }

    return response.json();
  }

  static async getById(id: number): Promise<Writer> {
    const response = await fetch(`${this.BASE_URL}/writers/${id}`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: response.statusText }));
      throw new Error(error.error || `WriterService.getById failed: ${response.statusText}`);
    }

    return response.json();
  }

  static async list(limit: number = 10, offset: number = 0): Promise<Writer[]> {
    const response = await fetch(`${this.BASE_URL}/writers?limit=${limit}&offset=${offset}`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: response.statusText }));
      throw new Error(error.error || `WriterService.list failed: ${response.statusText}`);
    }

    return response.json();
  }

  static async search(query: string, limit: number = 20, offset: number = 0): Promise<Writer[]> {
    const searchParam = encodeURIComponent(query);
    const response = await fetch(
      `${this.BASE_URL}/writers?search=${searchParam}&limit=${limit}&offset=${offset}`,
      {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      }
    );

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: response.statusText }));
      throw new Error(error.error || `WriterService.search failed: ${response.statusText}`);
    }

    return response.json();
  }

  static async update(id: number, params: UpdateWriterRequest): Promise<void> {
    const response = await fetch(`${this.BASE_URL}/writers/${id}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(params),
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: response.statusText }));
      throw new Error(error.error || `WriterService.update failed: ${response.statusText}`);
    }
  }

  static async delete(id: number): Promise<void> {
    const response = await fetch(`${this.BASE_URL}/writers/${id}`, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: response.statusText }));
      throw new Error(error.error || `WriterService.delete failed: ${response.statusText}`);
    }
  }
}
