import {Component, ViewChild} from '@angular/core';
import {MessagesService} from '../../service/messages.service';
import {FormsModule, NgForm} from '@angular/forms';
import {
  IonBackButton,
  IonButton,
  IonButtons,
  IonContent,
  IonHeader,
  IonInput,
  IonItem,
  IonList,
  IonText,
  IonTitle,
  IonToolbar,
  NavController
} from '@ionic/angular/standalone';
import {ProfileService} from '../profile/profile.service';
import {HttpErrorResponse} from '@angular/common/http';
import {displayFieldErrors} from '../../util';
import {Errors} from '../../api/types';

@Component({
  selector: 'app-password',
  templateUrl: './password.page.html',
  styleUrls: ['./password.page.scss'],
  imports: [FormsModule, IonContent, IonList, IonText, IonButton, IonHeader, IonToolbar, IonTitle, IonItem, IonInput, IonButtons, IonBackButton]
})
export class PasswordPage {

  @ViewChild('changeForm')
  changeForm!: NgForm;

  constructor(private readonly profileService: ProfileService,
              private readonly navCtrl: NavController,
              private readonly messagesService: MessagesService) {
  }

  async changePassword(oldPassword: string, newPassword: string): Promise<void> {

    const loading = await this.messagesService.showLoading('Changing password');

    this.profileService.changePassword(oldPassword, newPassword)
      .subscribe({
        next: () => {
          loading.dismiss();
          this.changeForm.resetForm();
          this.messagesService.showSuccessToast('Password successfully changed');
          this.navCtrl.back();
        },
        error: err => {
          loading.dismiss();
          this.handleErrorResponse(err);
        }
      });

  }

  private handleErrorResponse(errorResponse: HttpErrorResponse) {
    const response: Errors = errorResponse.error;
    if (response?.errors) {
      displayFieldErrors(this.changeForm, response.errors)
    } else {
      this.changeForm.resetForm();
      this.messagesService.showErrorToast('Changing password failed');
    }
  }

}
