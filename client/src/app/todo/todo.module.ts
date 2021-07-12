import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import {FormsModule} from '@angular/forms';
import {IonicModule} from '@ionic/angular';
import {RouterModule, Routes} from '@angular/router';
import {ListPage} from './list/list.page';
import {EditPage} from './edit/edit.page';
import {AuthGuard} from '../service/auth.guard';
import {TodoService} from './todo.service';

const routes: Routes = [
  {
    path: '',
    component: ListPage,
    canActivate: [AuthGuard],
  },
  {
    path: 'edit',
    component: EditPage,
    canActivate: [AuthGuard],
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
export class TodoModule { }
