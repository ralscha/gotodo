import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {RouterModule, Routes} from '@angular/router';
import {IonicModule} from '@ionic/angular';
import {ProfilePage} from './profile.page';
import {ProfileService} from './profile.service';
import {authGuard} from "../../app-routing.module";

const routes: Routes = [
  {
    path: '',
    component: ProfilePage,
    canActivate: [() => authGuard()],
    pathMatch: 'full'
  },
  {
    path: 'password',
    loadChildren: () => import('../password/password.module').then(m => m.PasswordPageModule),
    canActivate: [() => authGuard()]
  },
  {
    path: 'email',
    loadChildren: () => import('../email/email.module').then(m => m.EmailPageModule),
    canActivate: [() => authGuard()]
  },
  {
    path: 'email-confirm',
    loadChildren: () => import('../email-change-confirm/email-change-confirm.module').then(m => m.EmailChangeConfirmPageModule),
    canActivate: [() => authGuard()]
  },
  {
    path: 'account',
    loadChildren: () => import('../account/account.module').then(m => m.AccountPageModule),
    canActivate: [() => authGuard()]
  },
  {
    path: '**',
    redirectTo: '/profile'
  }
];


@NgModule({
  imports: [
    CommonModule,
    IonicModule,
    RouterModule.forChild(routes)
  ],
  declarations: [ProfilePage],
  providers: [ProfileService]
})
export class ProfilePageModule {
}
