export interface SaveResponse {
  id: number;
  success: boolean;
  fieldErrors?: { [key: string]: string };
  globalError?: string;
}
