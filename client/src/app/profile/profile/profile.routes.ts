import {Routes} from '@angular/router';
import {ProfilePage} from './profile.page';
import {authGuard} from "../../app-routing";

export const routes: Routes = [
  {
    path: '',
    component: ProfilePage,
    canActivate: [() => authGuard()],
    pathMatch: 'full'
  },
  {
    path: 'password',
    loadComponent: () => import('../password/password.page').then(m => m.PasswordPage),
    canActivate: [() => authGuard()]
  },
  {
    path: 'email',
    loadComponent: () => import('../email/email.page').then(m => m.EmailPage),
    canActivate: [() => authGuard()]
  },
  {
    path: 'email-confirm',
    loadChildren: () => import('../email-change-confirm/email-change-confirm.routes').then(m => m.routes),
    canActivate: [() => authGuard()]
  },
  {
    path: 'account',
    loadComponent: () => import('../account/account.page').then(m => m.AccountPage),
    canActivate: [() => authGuard()]
  },
  {
    path: '**',
    redirectTo: '/profile'
  }
];
