import {Routes} from '@angular/router';
import {SignupConfirmPage} from './signup-confirm.page';

export const routes: Routes = [
  {
    path: ':token',
    component: SignupConfirmPage
  },
  {
    path: '',
    component: SignupConfirmPage
  }
];
