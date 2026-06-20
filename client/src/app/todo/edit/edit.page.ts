import { Component, inject, OnInit, signal } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { MessagesService } from '../../service/messages.service';
import {
  AlertController,
  IonBackButton,
  IonButton,
  IonButtons,
  IonCol,
  IonContent,
  IonFab,
  IonFabButton,
  IonGrid,
  IonHeader,
  IonIcon,
  IonInput,
  IonItem,
  IonRow,
  IonText,
  IonTextarea,
  IonTitle,
  IonToolbar,
} from '@ionic/angular/standalone';
import { TodoService } from '../todo.service';
import { HttpErrorResponse } from '@angular/common/http';
import { Errors, Todo } from '../../api/types';
import { addIcons } from 'ionicons';
import { trash } from 'ionicons/icons';
import { FormField, form, required, schema } from '@angular/forms/signals';

@Component({
  selector: 'app-edit',
  templateUrl: './edit.page.html',
  styleUrls: ['./edit.page.scss'],
  imports: [
    FormField,
    IonContent,
    IonGrid,
    IonRow,
    IonCol,
    IonText,
    IonButton,
    IonHeader,
    IonToolbar,
    IonTitle,
    IonItem,
    IonInput,
    IonTextarea,
    IonFab,
    IonFabButton,
    IonIcon,
    IonBackButton,
    IonButtons,
  ],
})
export class EditPage implements OnInit {
  readonly selectedTodo = signal<Todo | undefined>(undefined);
  readonly submitted = signal(false);
  readonly submitError = signal<string | null>(null);
  readonly todoModel = signal({ subject: '', description: '' });
  readonly todoForm = form(
    this.todoModel,
    schema((path) => {
      required(path.subject);
    }),
  );
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly messagesService = inject(MessagesService);
  private readonly alertController = inject(AlertController);
  private readonly todoService = inject(TodoService);

  constructor() {
    addIcons({ trash });
  }

  ngOnInit(): void {
    const todoIdString = this.route.snapshot.paramMap.get('id');
    if (todoIdString) {
      const todo = this.todoService.getTodo(parseInt(todoIdString, 10));
      this.selectedTodo.set(todo);
      this.todoForm().reset({
        subject: todo?.subject ?? '',
        description: todo?.description ?? '',
      });
    } else {
      this.selectedTodo.set({
        id: 0,
        subject: '',
        description: undefined,
      });
      this.todoForm().reset({ subject: '', description: '' });
    }
  }

  async deleteTodo(): Promise<void> {
    if (this.selectedTodo()) {
      const alert = await this.alertController.create({
        header: 'Delete Todo',
        message: 'Do you really want to delete this todo?',
        buttons: [
          {
            text: 'Cancel',
            role: 'cancel',
          },
          {
            text: 'Delete Todo',
            handler: async () => this.reallyDeleteTodo(),
          },
        ],
      });
      await alert.present();
    }
  }

  async save(): Promise<void> {
    this.submitted.set(true);
    this.todoForm().markAsTouched();

    const selectedTodo = this.selectedTodo();
    if (!selectedTodo || this.todoForm().invalid()) {
      return;
    }

    const formValue = this.todoModel();
    selectedTodo.subject = formValue.subject;
    selectedTodo.description = formValue.description || undefined;

    this.todoService.save(selectedTodo).subscribe({
      next: async () => {
        await this.messagesService.showSuccessToast('Todo successfully saved', 500);
        await this.router.navigate(['/todo']);
      },
      error: (err: HttpErrorResponse) => this.handleErrorResponse(err.error),
    });
  }

  private handleErrorResponse(response: Errors | undefined): void {
    if (response?.errors?.['subject']) {
      this.submitError.set('subjectInvalid');
    } else {
      this.messagesService.showErrorToast('Saving Todo failed');
    }
  }

  private async reallyDeleteTodo(): Promise<void> {
    const selectedTodo = this.selectedTodo();
    if (selectedTodo) {
      this.todoService.delete(selectedTodo).subscribe({
        next: async () => {
          await this.router.navigate(['/todo']);
          await this.messagesService.showSuccessToast('Todo successfully deleted', 500);
        },
        error: () => this.messagesService.showErrorToast('Deleting Todo failed'),
      });
    }
  }
}
