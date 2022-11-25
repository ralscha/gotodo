/* Do not change, this code is generated from Golang structs */


export interface Todo {
    id: number;
    subject: string;
    description?: string;
}
export interface LoginInput {
    password: string;
    email: string;
}
export interface PasswordResetInput {
    password: string;
    resetToken: string;
}
export interface PasswordResetRequestInput {
    email: string;
}
export interface EmailChangeInput {
    password: string;
    newEmail: string;
}
export interface PasswordInput {
    password: string;
}
export interface TokenInput {
    token: string;
}
export interface PasswordChangeInput {
    oldPassword: string;
    newPassword: string;
}
export interface SignUpInput {
    email: string;
    password: string;
}
export interface LoginOutput {
    authority: string;
}
export interface Errors {
    errors: {[key: string]: string[]};
}
export interface AppVersionOutput {
    buildTime: string;
    version: string;
}