import {Component, OnInit} from '@angular/core';
import {Observable} from 'rxjs';
import {TodoService} from '../todo.service';
import {Todo} from '../todo';

@Component({
  selector: 'app-list',
  templateUrl: './list.page.html',
  styleUrls: ['./list.page.scss'],
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

  trackById(index: number, item: any): number {
    return item.id;
  }

}
