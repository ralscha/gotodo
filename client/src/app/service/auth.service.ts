import {Injectable} from '@angular/core';
import {BehaviorSubject, Observable, of} from 'rxjs';
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

  login(username: string, password: string): Observable<string | null> {
    const body = new HttpParams().set('username', username).set('password', password);
    return this.httpClient.post<string>('/v1/login', body, {withCredentials: true})
      .pipe(
        tap(authority => this.authoritySubject.next(authority)),
        catchError(() => of(null))
      );
  }

  logout(): Observable<void> {
    return this.httpClient.post<void>('/v1/logout', {withCredentials: true})
      .pipe(
        tap(() => this.authoritySubject.next(null))
      );
  }

}
