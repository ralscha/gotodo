import {Routes} from "@angular/router";
import {PasswordResetPage} from "./password-reset.page";

export const routes: Routes = [
  {
    path: ':token',
    component: PasswordResetPage
  },
  {
    path: '',
    component: PasswordResetPage
  }
];
