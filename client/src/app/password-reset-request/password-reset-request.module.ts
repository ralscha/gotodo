import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {FormsModule} from '@angular/forms';
import {RouterModule, Routes} from '@angular/router';

import {IonicModule} from '@ionic/angular';

import {PasswordResetRequestPage} from './password-reset-request.page';

const routes: Routes = [
  {
    path: '',
    component: PasswordResetRequestPage
  }
];

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    IonicModule,
    RouterModule.forChild(routes)
  ],
  declarations: [PasswordResetRequestPage]
})
export class PasswordResetRequestPageModule {
}
