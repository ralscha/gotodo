export interface SaveResponse {
  id: number;
  fieldErrors?: { [key: string]: string };
  globalError?: string;
}
