import {inject} from '@angular/core';
import {Router, Routes} from '@angular/router';
import {LoginPage} from './login/login.page';
import {LogoutPage} from './logout/logout.page';
import {map} from "rxjs/operators";
import {AuthService} from "./service/auth.service";
import {ProfileService} from "./profile/profile/profile.service";
import {TodoService} from "./todo/todo.service";

export const authGuard = (authService = inject(AuthService), router = inject(Router)) => {
  if (authService.isAuthenticated()) {
    return true;
  }

  return authService.authenticate().pipe(map(authority => {
      if (authority !== null) {
        return true;
      } else {
        return router.createUrlTree(['login']);
      }
    }
  ));
}

export const routes: Routes = [
  {path: '', redirectTo: 'todo', pathMatch: 'full'},
  {
    path: 'todo',
    providers: [TodoService],
    loadChildren: () => import('./todo/todo.routes').then(m => m.routes),
    canActivate: [() => authGuard()]
  },
  {
    path: 'profile',
    providers: [ProfileService],
    loadChildren: () => import('./profile/profile/profile.routes').then(m => m.routes),
    canActivate: [() => authGuard()]
  },
  {path: 'login', component: LoginPage},
  {path: 'logout', component: LogoutPage},
  {
    path: 'signup',
    loadComponent: () => import('./signup/signup.page').then(m => m.SignupPage),
  },
  {
    path: 'signup-confirm',
    loadChildren: () => import('./signup-confirm/signup-confirm.routes').then(m => m.routes)
  },
  {
    path: 'password-reset-request',
    loadComponent: () => import('./password-reset-request/password-reset-request.page').then(m => m.PasswordResetRequestPage)
  },
  {
    path: 'password-reset',
    loadChildren: () => import('./password-reset/password-reset.routes').then(m => m.routes)
  },
  {path: '**', redirectTo: 'todo'}
];
