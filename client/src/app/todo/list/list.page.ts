import {Component, OnInit} from '@angular/core';
import {Observable} from 'rxjs';
import {TodoService} from '../todo.service';
import {Todo} from '../../api/types';

@Component({
    selector: 'app-list',
    templateUrl: './list.page.html',
    standalone: false
})
export class ListPage implements OnInit {

  todos$!: Observable<Todo[]>;

  constructor(private readonly todoService: TodoService) {
  }

  ngOnInit(): void {
    this.todoService.loadTodos();
    this.todos$ = this.todoService.getTodos();
  }

  refresh(event: Event): void {
    this.todoService.loadTodos();
    (event as CustomEvent).detail.complete();
  }

}
