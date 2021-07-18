import {Injectable} from '@angular/core';
import {BehaviorSubject, Observable, of} from 'rxjs';
import {catchError, share, tap} from 'rxjs/operators';
import {HttpClient, HttpParams} from '@angular/common/http';
import {FormErrorResponse} from '../model/form-error-response';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private readonly authoritySubject = new BehaviorSubject<string | null>(null);
  readonly authority$ = this.authoritySubject.asObservable();
  private readonly authorityCall$: Observable<{ authority: string } | null>;

  constructor(private readonly httpClient: HttpClient) {
    this.authorityCall$ = this.httpClient.post<{ authority: string }>('/v1/authenticate', null, {
      withCredentials: true
    })
      .pipe(
        tap(response => this.authoritySubject.next(response.authority)),
        catchError(() => of(null)),
        share()
      );
  }

  authenticate(): Observable<{ authority: string } | null> {
    return this.authorityCall$;
  }

  isAuthenticated(): boolean {
    return this.authoritySubject.getValue() != null;
  }

  login(email: string, password: string): Observable<string | null> {
    const body = new HttpParams().set('email', email).set('password', password);
    return this.httpClient.post<string>('/v1/login', body, {withCredentials: true})
      .pipe(
        tap(authority => this.authoritySubject.next(authority)),
      );
  }

  logout(): Observable<void> {
    return this.httpClient.post<void>('/v1/logout', null, {withCredentials: true})
      .pipe(
        tap(() => this.authoritySubject.next(null))
      );
  }

  signup(email: string, password: string): Observable<FormErrorResponse | void> {
    const body = new HttpParams().set('email', email).set('password', password);
    return this.httpClient.post<FormErrorResponse | void>('/v1/signup', body);
  }

  confirmSignup(token: string): Observable<void> {
    return this.httpClient.post<void>('/v1/signup-confirm', token)
  }

  resetPasswordRequest(email: string): Observable<void> {
    return this.httpClient.post<void>('/v1/reset-password-request', email);
  }

  resetPassword(resetToken: string, password: string): Observable<'WEAK' | void> {
    const body = new HttpParams().set('resetToken', resetToken).set('password', password);
    return this.httpClient.post<'WEAK' | void>('/v1/reset-password', body);
  }

}
