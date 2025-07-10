import {Component, inject, OnInit} from '@angular/core';
import {Observable} from 'rxjs';
import {TodoService} from '../todo.service';
import {Todo} from '../../api/types';
import {AsyncPipe} from '@angular/common';
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
  IonTitle,
  IonToolbar
} from "@ionic/angular/standalone";
import {addIcons} from "ionicons";
import {add} from "ionicons/icons";
import {RouterLink} from "@angular/router";

@Component({
  selector: 'app-list',
  templateUrl: './list.page.html',
  imports: [RouterLink, IonRouterLink, AsyncPipe, IonContent, IonList, IonHeader, IonToolbar, IonTitle, IonItem, IonButtons, IonMenuButton, IonRefresher, IonRefresherContent, IonLabel, IonFab, IonFabButton, IonIcon]
})
export class ListPage implements OnInit {
  todos$!: Observable<Todo[]>;
  private readonly todoService = inject(TodoService);

  constructor() {
    addIcons({add});
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
