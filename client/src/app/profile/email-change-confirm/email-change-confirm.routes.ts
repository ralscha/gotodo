import {Routes} from '@angular/router';
import {EmailChangeConfirmPage} from './email-change-confirm.page';

export const routes: Routes = [
  {
    path: ':token',
    component: EmailChangeConfirmPage
  },
  {
    path: '',
    component: EmailChangeConfirmPage
  }
];
