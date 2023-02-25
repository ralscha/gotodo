import {inject, NgModule} from '@angular/core';
import {PreloadAllModules, Router, RouterModule, Routes} from '@angular/router';
import {LoginPage} from './login/login.page';
import {LogoutPage} from './logout/logout.page';
import {map} from "rxjs/operators";
import {AuthService} from "./service/auth.service";

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

const routes: Routes = [
  {path: '', redirectTo: 'todo', pathMatch: 'full'},
  {
    path: 'todo',
    loadChildren: () => import('./todo/todo.module').then(m => m.TodoModule),
    canActivate: [() => authGuard()]
  },
  {
    path: 'profile',
    loadChildren: () => import('./profile/profile/profile.module').then(m => m.ProfilePageModule),
    canActivate: [() => authGuard()]
  },
  {path: 'login', component: LoginPage},
  {path: 'logout', component: LogoutPage},
  {
    path: 'signup',
    loadChildren: () => import('./signup/signup.module').then(m => m.SignupPageModule)
  },
  {
    path: 'signup-confirm',
    loadChildren: () => import('./signup-confirm/signup-confirm.module').then(m => m.SignupConfirmPageModule)
  },
  {
    path: 'password-reset-request',
    loadChildren: () => import('./password-reset-request/password-reset-request.module').then(m => m.PasswordResetRequestPageModule)
  },
  {
    path: 'password-reset',
    loadChildren: () => import('./password-reset/password-reset.module').then(m => m.PasswordResetPageModule)
  },
  {path: '**', redirectTo: 'todo'}
];

@NgModule({
  imports: [
    RouterModule.forRoot(routes, {preloadingStrategy: PreloadAllModules, useHash: true})
  ],
  exports: [RouterModule]
})
export class AppRoutingModule {
}
