import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {RouterModule, Routes} from '@angular/router';

import {IonicModule} from '@ionic/angular';

import {PasswordResetPage} from './password-reset.page';
import {FormsModule} from '@angular/forms';

const routes: Routes = [
  {
    path: ':token',
    component: PasswordResetPage
  },
  {
    path: '',
    component: PasswordResetPage
  }
];

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    IonicModule,
    RouterModule.forChild(routes)
  ],
  declarations: [PasswordResetPage]
})
export class PasswordResetPageModule {
}
