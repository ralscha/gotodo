import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {RouterModule, Routes} from '@angular/router';
import {IonicModule} from '@ionic/angular';
import {ProfilePage} from './profile.page';
import {AuthGuard} from '../../service/auth.guard';
import {ProfileService} from './profile.service';

const routes: Routes = [
  {
    path: '',
    component: ProfilePage,
    canActivate: [AuthGuard],
    pathMatch: 'full'
  },
  {
    path: 'password',
    loadChildren: () => import('../password/password.module').then(m => m.PasswordPageModule),
    canActivate: [AuthGuard]
  },
  {
    path: 'email',
    loadChildren: () => import('../email/email.module').then(m => m.EmailPageModule),
    canActivate: [AuthGuard]
  },
  {
    path: 'account',
    loadChildren: () => import('../account/account.module').then(m => m.AccountPageModule),
    canActivate: [AuthGuard]
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
