export interface FormErrorResponse {
  fieldErrors?: { [key: string]: string };
  globalError?: string;
}
