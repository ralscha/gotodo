import {Component, OnInit} from '@angular/core';
import {ActivatedRoute, Router} from '@angular/router';
import {MessagesService} from '../../service/messages.service';
import {AlertController} from '@ionic/angular';
import {NgForm} from '@angular/forms';
import {Todo} from '../todo';
import {TodoService} from '../todo.service';
import {displayFieldErrors} from '../../util';
import {FormErrorResponse} from '../../model/form-error-response';
import {HttpErrorResponse} from '@angular/common/http';

@Component({
  selector: 'app-edit',
  templateUrl: './edit.page.html',
  styleUrls: ['./edit.page.scss'],
})
export class EditPage implements OnInit {

  selectedTodo: Todo | undefined;

  constructor(private readonly route: ActivatedRoute,
              private readonly router: Router,
              private readonly messagesService: MessagesService,
              private readonly alertController: AlertController,
              private readonly todoService: TodoService) {
  }

  async ngOnInit(): Promise<void> {
    const todoIdString = this.route.snapshot.paramMap.get('id');
    if (todoIdString) {
      this.selectedTodo = await this.todoService.getTodo(parseInt(todoIdString, 10));
    } else {
      this.selectedTodo = {
        id: 0,
        subject: null,
        description: null,
      };
    }
  }

  async deleteTodo(): Promise<void> {
    if (this.selectedTodo) {
      const alert = await this.alertController.create({
        header: 'Delete Todo',
        message: 'Do you really want to delete this todo?</strong>',
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
      const response: FormErrorResponse = errorResponse.error;
      if (response) {
        if (response.fieldErrors) {
          displayFieldErrors(form, response.fieldErrors)
        } else {
          this.messagesService.showErrorToast('Saving Todo failed');
        }
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
