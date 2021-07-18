import {Injectable} from '@angular/core';
import {BehaviorSubject, Observable, of} from 'rxjs';
import {map} from 'rxjs/operators';
import {HttpClient, HttpParams} from '@angular/common/http';
import {catchError, share, tap} from 'rxjs/operators';

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

  signup(email: string, password: string): Observable<'EMAIL_REGISTERED' | 'WEAK_PASSWORD' | null> {
    const body = new HttpParams().set('email', email).set('password', password);
    return this.httpClient.post<'EMAIL_REGISTERED' | 'WEAK_PASSWORD' | null>('/be/signup', body);
  }

  confirmSignup(token: string): Observable<boolean> {
    return this.httpClient.post('/be/confirm-signup', token, {responseType: 'text'})
      .pipe(
        map(response => response === 'true')
      );
  }

  resetPasswordRequest(email: string): Observable<void> {
    return this.httpClient.post<void>('/v1/reset-password-request', email);
  }

  resetPassword(resetToken: string, password: string): Observable<'INVALID' | 'WEAK_PASSWORD' | null> {
    const body = new HttpParams().set('resetToken', resetToken).set('password', password);
    return this.httpClient.post<'INVALID' | 'WEAK_PASSWORD' | null>('/v1/reset-password', body);
  }

}
