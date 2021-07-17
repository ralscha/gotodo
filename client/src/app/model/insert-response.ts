export interface InsertResponse {
  id: number;
  success: boolean;
  fieldErrors?: { [key: string]: string };
  globalError?: string;
}
