import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {FormsModule} from '@angular/forms';
import {IonicModule} from '@ionic/angular';
import {RouterModule, Routes} from '@angular/router';
import {ListPage} from './list/list.page';
import {EditPage} from './edit/edit.page';
import {TodoService} from './todo.service';
import {authGuard} from "../app-routing.module";

const routes: Routes = [
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

@NgModule({
  declarations: [ListPage, EditPage],
  imports: [
    CommonModule,
    FormsModule,
    IonicModule,
    RouterModule.forChild(routes)
  ],
  providers: [TodoService]
})
export class TodoModule {
}
