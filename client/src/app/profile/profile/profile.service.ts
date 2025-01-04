import {Injectable} from '@angular/core';
import { HttpClient } from '@angular/common/http';
import {Observable} from 'rxjs';
import {EmailChangeInput, Errors, PasswordChangeInput, PasswordInput, TokenInput} from '../../api/types';

@Injectable()
export class ProfileService {

  constructor(private readonly httpClient: HttpClient) {
  }

  deleteAccount(password: string): Observable<Errors | void> {
    const request: PasswordInput = {password};
    return this.httpClient.post<Errors | void>('/v1/profile/account-delete', request);
  }

  changePassword(oldPassword: string, newPassword: string): Observable<Errors | void> {
    const request: PasswordChangeInput = {oldPassword, newPassword};
    return this.httpClient.post<Errors | void>('/v1/profile/password-change', request);
  }

  changeEmail(newEmail: string, password: string): Observable<Errors> {
    const request: EmailChangeInput = {newEmail, password};
    return this.httpClient.post<Errors>('/v1/profile/email-change', request);
  }

  confirmEmailChange(token: string): Observable<void> {
    const request: TokenInput = {token};
    return this.httpClient.post<void>('/v1/profile/email-change-confirm', request);
  }
}
