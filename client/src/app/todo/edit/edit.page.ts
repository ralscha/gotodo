import {Component, OnInit} from '@angular/core';
import {ActivatedRoute, Router} from '@angular/router';
import {MessagesService} from '../../service/messages.service';
import {AlertController} from '@ionic/angular';
import {NgForm} from '@angular/forms';
import {Todo} from '../todo';
import {TodoService} from '../todo.service';

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
      if (this.selectedTodo.id > 0) {
        this.todoService.update(this.selectedTodo).subscribe(response => {
          if (response.success) {
            this.messagesService.showSuccessToast('Todo successfully updated', 500);
          } else {
            this.messagesService.showErrorToast('Updating Todo failed');
          }
          this.router.navigate(['/todos']);
        });
      } else {
        this.todoService.insert(this.selectedTodo).subscribe(response => {
          if (response.success) {
            this.messagesService.showSuccessToast('Todo successfully inserted', 500);
          } else {
            this.messagesService.showErrorToast('Inserting Todo failed');
          }
          this.router.navigate(['/todos']);
        });
      }
    }
  }

  private async reallyDeleteTodo(): Promise<void> {
    if (this.selectedTodo) {
      this.todoService.delete(this.selectedTodo).subscribe(success => {
        this.router.navigate(['/todos']);
        if (success) {
          this.messagesService.showSuccessToast('Todo successfully deleted', 500);
        } else {
          this.messagesService.showErrorToast('Deleting Todo failed');
        }
      })
    }
  }
}
