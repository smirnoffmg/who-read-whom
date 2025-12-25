export interface Opinion {
  writer_id: number;
  work_id: number;
  sentiment: boolean;
  quote: string;
  source: string;
  page: string | null;
  statement_year: number | null;
}

export interface CreateOpinionRequest {
  writer_id: number;
  work_id: number;
  sentiment: boolean;
  quote: string;
  source: string;
  page?: string | null;
  statement_year?: number | null;
}

export interface UpdateOpinionRequest {
  sentiment: boolean;
  quote: string;
  source: string;
  page?: string | null;
  statement_year?: number | null;
}
