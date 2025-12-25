import type { CreateOpinionRequest, Opinion, UpdateOpinionRequest } from "@/types/opinion";

export class OpinionService {
  private static readonly BASE_URL =
    process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

  static async create(params: CreateOpinionRequest): Promise<Opinion> {
    const response = await fetch(`${this.BASE_URL}/opinions`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(params),
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: response.statusText }));
      throw new Error(error.error || `OpinionService.create failed: ${response.statusText}`);
    }

    return response.json();
  }

  static async getByWriter(writerId: number): Promise<Opinion[]> {
    const response = await fetch(`${this.BASE_URL}/opinions/writer/${writerId}`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: response.statusText }));
      throw new Error(error.error || `OpinionService.getByWriter failed: ${response.statusText}`);
    }

    return response.json();
  }

  static async getByWork(workId: number): Promise<Opinion[]> {
    const response = await fetch(`${this.BASE_URL}/opinions/work/${workId}`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: response.statusText }));
      throw new Error(error.error || `OpinionService.getByWork failed: ${response.statusText}`);
    }

    return response.json();
  }

  static async getByWriterAndWork(writerId: number, workId: number): Promise<Opinion> {
    const response = await fetch(`${this.BASE_URL}/opinions/writer/${writerId}/work/${workId}`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: response.statusText }));
      throw new Error(
        error.error || `OpinionService.getByWriterAndWork failed: ${response.statusText}`
      );
    }

    return response.json();
  }

  static async list(limit: number = 10, offset: number = 0): Promise<Opinion[]> {
    const response = await fetch(`${this.BASE_URL}/opinions?limit=${limit}&offset=${offset}`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: response.statusText }));
      throw new Error(error.error || `OpinionService.list failed: ${response.statusText}`);
    }

    return response.json();
  }

  static async update(
    writerId: number,
    workId: number,
    params: UpdateOpinionRequest
  ): Promise<void> {
    const response = await fetch(`${this.BASE_URL}/opinions/writer/${writerId}/work/${workId}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(params),
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: response.statusText }));
      throw new Error(error.error || `OpinionService.update failed: ${response.statusText}`);
    }
  }

  static async delete(writerId: number, workId: number): Promise<void> {
    const response = await fetch(`${this.BASE_URL}/opinions/writer/${writerId}/work/${workId}`, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: response.statusText }));
      throw new Error(error.error || `OpinionService.delete failed: ${response.statusText}`);
    }
  }
}
