import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {BehaviorSubject, Observable} from 'rxjs';
import {Todo} from './todo';

@Injectable({
  providedIn: 'root'
})
export class TodoService {

  private readonly todosSubject = new BehaviorSubject<Todo[]>([]);
  private readonly todos$ = this.todosSubject.asObservable();

  constructor(private readonly httpClient: HttpClient) {
  }

  loadTodos(): void {
    this.httpClient.get<Todo[]>('/v1/todo').subscribe(data => this.todosSubject.next(data));
  }

  getTodos(): Observable<Todo[]> {
    return this.todos$;
  }

  async getTodo(id: number): Promise<Todo | undefined> {
    console.log('get', id)
    return undefined;
  }

  async delete(todo: Todo): Promise<void> {
    console.log('delete', todo)
  }

  async save(todo: Todo): Promise<void> {
    console.log('save', todo)
  }

}
