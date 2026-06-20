import { computed, inject, Service, signal } from '@angular/core';
import { Observable, of } from 'rxjs';
import { catchError, share, tap } from 'rxjs/operators';
import { HttpClient } from '@angular/common/http';
import { toObservable } from '@angular/core/rxjs-interop';
import {
  Errors,
  LoginInput,
  LoginOutput,
  PasswordResetInput,
  PasswordResetRequestInput,
  SignUpInput,
  TokenInput,
} from '../api/types';

@Service()
export class AuthService {
  private readonly httpClient = inject(HttpClient);

  private readonly authority = signal<string | null>(null);
  readonly authenticated = computed(() => this.authority() != null);
  readonly authority$ = toObservable(this.authority);
  private readonly authorityCall$: Observable<{ authority: string } | null>;

  constructor() {
    this.authorityCall$ = this.httpClient
      .post<{ authority: string }>('/v1/authenticate', null, {
        withCredentials: true,
      })
      .pipe(
        tap((response) => this.authority.set(response.authority)),
        catchError(() => of(null)),
        share(),
      );
  }

  authenticate(): Observable<{ authority: string } | null> {
    return this.authorityCall$;
  }

  isAuthenticated(): boolean {
    return this.authenticated();
  }

  login(email: string, password: string): Observable<LoginOutput | null> {
    const request: LoginInput = { email, password };
    return this.httpClient
      .post<LoginOutput>('/v1/login', request, { withCredentials: true })
      .pipe(tap((response) => this.authority.set(response?.authority)));
  }

  logout(): Observable<void> {
    return this.httpClient
      .post<void>('/v1/logout', null, { withCredentials: true })
      .pipe(tap(() => this.authority.set(null)));
  }

  logoutClient(): void {
    this.authority.set(null);
  }

  signup(email: string, password: string): Observable<Errors | void> {
    const request: SignUpInput = { email, password };
    return this.httpClient.post<Errors | void>('/v1/signup', request);
  }

  confirmSignup(token: string): Observable<void> {
    const request: TokenInput = { token };
    return this.httpClient.post<void>('/v1/signup-confirm', request);
  }

  resetPasswordRequest(email: string): Observable<void> {
    const request: PasswordResetRequestInput = { email };
    return this.httpClient.post<void>('/v1/password-reset-request', request);
  }

  resetPassword(resetToken: string, password: string): Observable<Errors | void> {
    const request: PasswordResetInput = { password, resetToken };
    return this.httpClient.post<Errors | void>('/v1/password-reset', request);
  }
}
