export interface UpdateResponse {
  success: boolean;
  fieldErrors?: { [key: string]: string };
  globalError?: string;
}
