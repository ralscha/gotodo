import { inject, Service } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { BehaviorSubject, map, Observable, of } from 'rxjs';
import { tap } from 'rxjs/operators';
import { Todo } from '../api/types';

@Service({ autoProvided: false })
export class TodoService {
  private readonly httpClient = inject(HttpClient);

  private todosMap = new Map<number, Todo>();

  private readonly todosSubject = new BehaviorSubject<Todo[]>([]);
  private readonly todos$ = this.todosSubject.asObservable();

  loadTodos(): Observable<Todo[]> {
    return this.httpClient.get<Todo[]>('/v1/todo').pipe(
      tap((todos) => {
        this.todosMap.clear();
        for (const todo of todos) {
          this.todosMap.set(todo.id, todo);
        }
        this.publish();
      }),
    );
  }

  getTodos(): Observable<Todo[]> {
    return this.todos$;
  }

  getTodo(id: number): Observable<Todo | undefined> {
    const cached = this.todosMap.get(id);
    if (cached) {
      return of(cached);
    }

    return this.loadTodos().pipe(map(() => this.todosMap.get(id)));
  }

  delete(todo: Todo): Observable<void> {
    return this.httpClient.delete<void>(`/v1/todo/${todo.id}`).pipe(
      tap(() => {
        this.todosMap.delete(todo.id);
        this.publish();
      }),
    );
  }

  save(todo: Todo): Observable<Pick<Todo, 'id'> | void> {
    const savedTodo = { ...todo };
    return this.httpClient.post<Pick<Todo, 'id'> | void>('/v1/todo', savedTodo).pipe(
      tap((pickTodo) => {
        if (pickTodo && 'id' in pickTodo) {
          savedTodo.id = pickTodo.id;
        }
        this.todosMap.set(savedTodo.id, savedTodo);
        this.publish();
      }),
    );
  }

  private publish(): void {
    this.todosSubject.next([...this.todosMap.values()]);
  }
}
