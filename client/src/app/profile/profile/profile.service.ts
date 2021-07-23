import {Injectable} from '@angular/core';
import {HttpClient, HttpParams} from '@angular/common/http';
import {Observable} from 'rxjs';
import {FormErrorResponse} from '../../model/form-error-response';

@Injectable({
  providedIn: 'root'
})
export class ProfileService {

  constructor(private readonly httpClient: HttpClient) {
  }

  deleteAccount(password: string): Observable<FormErrorResponse | void> {
    return this.httpClient.post<FormErrorResponse | void>('/v1/profile/delete-account', password);
  }

  changePassword(oldPassword: string, newPassword: string): Observable<FormErrorResponse | void> {
    const body = new HttpParams().set('oldPassword', oldPassword).set('newPassword', newPassword);
    return this.httpClient.post<FormErrorResponse | void>('/v1/profile/change-password', body);
  }

  changeEmail(newEmail: string, password: string): Observable<FormErrorResponse> {
    const body = new HttpParams().set('newEmail', newEmail).set('password', password);
    return this.httpClient.post<FormErrorResponse>('/v1/profile/change-email', body);
  }

  confirmEmailChange(token: string): Observable<void> {
    return this.httpClient.post<void>('/v1/profile/email-change-confirm', token);
  }
}