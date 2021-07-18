import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {RouterModule, Routes} from '@angular/router';
import {IonicModule} from '@ionic/angular';
import {SignupConfirmPage} from './signup-confirm.page';

const routes: Routes = [
  {
    path: ':token',
    component: SignupConfirmPage
  },
  {
    path: '',
    component: SignupConfirmPage
  }
];

@NgModule({
  imports: [
    CommonModule,
    IonicModule,
    RouterModule.forChild(routes)
  ],
  declarations: [SignupConfirmPage]
})
export class SignupConfirmPageModule {
}
