export interface Writer {
  id: number;
  name: string;
  birth_year: number;
  death_year: number | null;
  bio: string | null;
}

export interface CreateWriterRequest {
  name: string;
  birth_year: number;
  death_year?: number | null;
  bio?: string | null;
}

export interface UpdateWriterRequest {
  name: string;
  birth_year: number;
  death_year?: number | null;
  bio?: string | null;
}
