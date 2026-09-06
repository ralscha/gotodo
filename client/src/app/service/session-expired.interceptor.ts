import { HttpErrorResponse, HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { Router } from '@angular/router';
import { catchError, throwError } from 'rxjs';
import { AuthService } from './auth.service';

const authenticationEndpoints = new Set(['/v1/authenticate', '/v1/login']);

export const sessionExpiredInterceptor: HttpInterceptorFn = (request, next) => {
  const authService = inject(AuthService);
  const router = inject(Router);

  return next(request).pipe(
    catchError((error: HttpErrorResponse) => {
      if (error.status === 401 && !authenticationEndpoints.has(request.url)) {
        authService.logoutClient();
        void router.navigateByUrl('/login');
      }
      return throwError(() => error);
    }),
  );
};
