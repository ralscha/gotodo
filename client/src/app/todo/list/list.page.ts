import { Component, inject, OnInit, signal } from '@angular/core';
import { Observable } from 'rxjs';
import { TodoService } from '../todo.service';
import { Todo } from '../../api/types';
import { AsyncPipe } from '@angular/common';
import {
  IonButtons,
  IonButton,
  IonContent,
  IonFab,
  IonFabButton,
  IonHeader,
  IonIcon,
  IonItem,
  IonLabel,
  IonList,
  IonMenuButton,
  IonRefresher,
  IonRefresherContent,
  IonRouterLink,
  IonSpinner,
  IonText,
  RefresherCustomEvent,
  IonTitle,
  IonToolbar,
} from '@ionic/angular';
import { addIcons } from 'ionicons';
import { add } from 'ionicons/icons';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'app-list',
  templateUrl: './list.page.html',
  imports: [
    RouterLink,
    IonRouterLink,
    AsyncPipe,
    IonContent,
    IonList,
    IonHeader,
    IonToolbar,
    IonTitle,
    IonItem,
    IonButton,
    IonButtons,
    IonMenuButton,
    IonRefresher,
    IonRefresherContent,
    IonLabel,
    IonFab,
    IonFabButton,
    IonIcon,
    IonSpinner,
    IonText,
  ],
})
export class ListPage implements OnInit {
  todos$!: Observable<Todo[]>;
  readonly loading = signal(true);
  readonly loadFailed = signal(false);
  private readonly todoService = inject(TodoService);

  constructor() {
    addIcons({ add });
  }

  ngOnInit(): void {
    this.todos$ = this.todoService.getTodos();
    this.loadTodos();
  }

  refresh(event: RefresherCustomEvent): void {
    this.loadFailed.set(false);
    this.todoService.loadTodos().subscribe({
      complete: () => event.target.complete(),
      error: () => {
        this.loadFailed.set(true);
        event.target.complete();
      },
    });
  }

  retry(): void {
    this.loading.set(true);
    this.loadTodos();
  }

  private loadTodos(): void {
    this.loadFailed.set(false);
    this.todoService.loadTodos().subscribe({
      error: () => {
        this.loading.set(false);
        this.loadFailed.set(true);
      },
      complete: () => this.loading.set(false),
    });
  }
}
