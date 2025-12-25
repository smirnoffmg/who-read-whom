export interface Work {
  id: number;
  title: string;
  author_id: number;
}

export interface CreateWorkRequest {
  title: string;
  author_id: number;
}

export interface UpdateWorkRequest {
  title: string;
  author_id: number;
}
