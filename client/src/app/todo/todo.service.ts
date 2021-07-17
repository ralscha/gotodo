import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {BehaviorSubject, Observable} from 'rxjs';
import {map, tap} from 'rxjs/operators';
import {Todo} from './todo';
import {InsertResponse} from '../model/insert-response';
import {UpdateResponse} from '../model/update-response';
import {DeleteResponse} from '../model/delete-response';

@Injectable({
  providedIn: 'root'
})
export class TodoService {
  private todosMap: Map<number, Todo> = new Map();

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

  delete(todo: Todo): Observable<boolean> {
    return this.httpClient.delete<DeleteResponse>(`/v1/todo/${todo.id}`)
      .pipe(
        map(response => {
          if (response.success) {
            this.todosMap.delete(todo.id);
            this.publish();
            return true;
          } else {
            return false
          }
        }));
  }

  insert(todo: Todo): Observable<InsertResponse> {
    return this.httpClient.post<InsertResponse>('/v1/todo', todo)
      .pipe(
        tap(response => {
          if (response.success) {
            todo.id = response.id;
            this.todosMap.set(todo.id, todo)
            this.publish();
          }
        })
      )
  }

  update(todo: Todo): Observable<UpdateResponse> {
    return this.httpClient.put<UpdateResponse>(`/v1/todo/${todo.id}`, todo)
      .pipe(
        tap(response => {
          if (response.success) {
            this.todosMap.set(todo.id, todo)
            this.publish();
          }
        })
      )
  }

  private publish(): void {
    this.todosSubject.next([...this.todosMap.values()])
  }

}
