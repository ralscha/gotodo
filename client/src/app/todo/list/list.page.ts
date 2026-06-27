import { Component, inject, OnInit } from '@angular/core';
import { Observable } from 'rxjs';
import { TodoService } from '../todo.service';
import { Todo } from '../../api/types';
import { AsyncPipe } from '@angular/common';
import {
  IonButtons,
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
  RefresherCustomEvent,
  IonTitle,
  IonToolbar,
} from '@ionic/angular/standalone';
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
    IonButtons,
    IonMenuButton,
    IonRefresher,
    IonRefresherContent,
    IonLabel,
    IonFab,
    IonFabButton,
    IonIcon,
  ],
})
export class ListPage implements OnInit {
  todos$!: Observable<Todo[]>;
  private readonly todoService = inject(TodoService);

  constructor() {
    addIcons({ add });
  }

  ngOnInit(): void {
    this.todos$ = this.todoService.getTodos();
    this.todoService.loadTodos().subscribe();
  }

  refresh(event: RefresherCustomEvent): void {
    this.todoService.loadTodos().subscribe({
      complete: () => event.target.complete(),
      error: () => event.target.complete(),
    });
  }
}
