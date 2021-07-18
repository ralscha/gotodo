import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {FormsModule} from '@angular/forms';
import {RouterModule, Routes} from '@angular/router';
import {IonicModule} from '@ionic/angular';
import {EmailChangeConfirmPage} from './email-change-confirm.page';

const routes: Routes = [
  {
    path: ':token',
    component: EmailChangeConfirmPage
  },
  {
    path: '',
    component: EmailChangeConfirmPage
  }
];

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    IonicModule,
    RouterModule.forChild(routes)
  ],
  declarations: [EmailChangeConfirmPage]
})
export class EmailChangeConfirmPageModule {
}
