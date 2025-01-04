import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {BehaviorSubject, Observable} from 'rxjs';
import {tap} from 'rxjs/operators';
import {Errors, Todo} from '../api/types';

@Injectable()
export class TodoService {
  private todosMap = new Map<number, Todo>();

  private readonly todosSubject = new BehaviorSubject<Todo[]>([]);
  private readonly todos$ = this.todosSubject.asObservable();

  constructor(private readonly httpClient: HttpClient) {
  }

  loadTodos(): void {
    this.httpClient.get<Todo[]>('/v1/todo').subscribe(todos => {
      this.todosMap.clear();
      for (const todo of todos) {
        this.todosMap.set(todo.id, todo);
      }
      this.publish();
    });
  }

  getTodos(): Observable<Todo[]> {
    return this.todos$;
  }

  getTodo(id: number): Todo | undefined {
    return this.todosMap.get(id);
  }

  delete(todo: Todo): Observable<void> {
    return this.httpClient.delete<void>(`/v1/todo/${todo.id}`)
      .pipe(
        tap(() => {
          this.todosMap.delete(todo.id);
          this.publish();
        }));
  }

  save(todo: Todo): Observable<Errors | Pick<Todo, "id"> | void> {
    return this.httpClient.post<Errors | Pick<Todo, "id"> | void>('/v1/todo', todo)
      .pipe(
        tap(pickTodo => {
          if (pickTodo && 'id' in pickTodo) {
            todo.id = pickTodo.id;
          }
          this.todosMap.set(todo.id, todo)
          this.publish();
        })
      )
  }

  private publish(): void {
    this.todosSubject.next([...this.todosMap.values()])
  }

}
