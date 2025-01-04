import {Routes} from '@angular/router';
import {ListPage} from './list/list.page';
import {EditPage} from './edit/edit.page';
import {authGuard} from "../app-routing";

export const routes: Routes = [
  {
    path: '',
    component: ListPage,
    canActivate: [() => authGuard()],
  },
  {
    path: 'edit',
    canActivate: [() => authGuard()],
    children: [
      {
        path: ':id',
        component: EditPage
      },
      {
        path: '',
        component: EditPage
      }
    ]
  }
];
