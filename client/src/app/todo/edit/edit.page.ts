import {Component, inject, OnInit} from '@angular/core';
import {ActivatedRoute, Router} from '@angular/router';
import {MessagesService} from '../../service/messages.service';
import {
  AlertController,
  IonBackButton,
  IonButton,
  IonButtons,
  IonContent,
  IonFab,
  IonFabButton,
  IonHeader,
  IonIcon,
  IonInput,
  IonItem,
  IonList,
  IonText,
  IonTextarea,
  IonTitle,
  IonToolbar
} from '@ionic/angular/standalone';
import {FormsModule, NgForm} from '@angular/forms';
import {TodoService} from '../todo.service';
import {displayFieldErrors} from '../../util';
import {HttpErrorResponse} from '@angular/common/http';
import {Errors, Todo} from '../../api/types';
import {addIcons} from "ionicons";
import {trash} from "ionicons/icons";

@Component({
  selector: 'app-edit',
  templateUrl: './edit.page.html',
  styleUrls: ['./edit.page.scss'],
  imports: [FormsModule, IonContent, IonList, IonText, IonButton, IonHeader, IonToolbar, IonTitle, IonItem, IonInput, IonTextarea, IonFab, IonFabButton, IonIcon, IonBackButton, IonButtons]
})
export class EditPage implements OnInit {
  selectedTodo?: Todo;
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly messagesService = inject(MessagesService);
  private readonly alertController = inject(AlertController);
  private readonly todoService = inject(TodoService);

  constructor() {
    addIcons({trash});
  }

  async ngOnInit(): Promise<void> {
    const todoIdString = this.route.snapshot.paramMap.get('id');
    if (todoIdString) {
      this.selectedTodo = await this.todoService.getTodo(parseInt(todoIdString, 10));
    } else {
      this.selectedTodo = {
        id: 0,
        subject: "",
        description: undefined,
      };
    }
  }

  async deleteTodo(): Promise<void> {
    if (this.selectedTodo) {
      const alert = await this.alertController.create({
        header: 'Delete Todo',
        message: 'Do you really want to delete this todo?',
        buttons: [
          {
            text: 'Cancel',
            role: 'cancel'
          }, {
            text: 'Delete Todo',
            handler: async () => this.reallyDeleteTodo()
          }
        ]
      });
      await alert.present();
    }
  }

  async save(form: NgForm): Promise<void> {
    if (this.selectedTodo) {
      this.selectedTodo.subject = form.value.subject;
      this.selectedTodo.description = form.value.description;
      this.todoService.save(this.selectedTodo).subscribe({
        next: () => this.handleSuccessResponse(),
        error: this.handleErrorResponse(form)
      });
    }
  }

  private handleSuccessResponse(): void {
    this.messagesService.showSuccessToast('Todo successfully saved', 500);
    this.router.navigate(['/todos']);
  }

  private handleErrorResponse(form: NgForm) {
    return (errorResponse: HttpErrorResponse) => {
      const response: Errors = errorResponse.error;
      if (response?.errors) {
        displayFieldErrors(form, response.errors)
      } else {
        this.messagesService.showErrorToast('Saving Todo failed');
      }
    };
  }

  private async reallyDeleteTodo(): Promise<void> {
    if (this.selectedTodo) {
      this.todoService.delete(this.selectedTodo).subscribe({
        next: () => {
          this.router.navigate(['/todos']);
          this.messagesService.showSuccessToast('Todo successfully deleted', 500);
        },
        error: () => this.messagesService.showErrorToast('Deleting Todo failed')
      });
    }
  }
}
