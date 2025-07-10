import {Component, inject} from '@angular/core';
import {MessagesService} from '../../service/messages.service';
import {ProfileService} from '../profile/profile.service';
import {FormsModule, NgForm} from '@angular/forms';
import {HttpErrorResponse} from '@angular/common/http';
import {displayFieldErrors} from '../../util';
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
import {Errors} from '../../api/types';

@Component({
  selector: 'app-email',
  templateUrl: './email.page.html',
  styleUrls: ['./email.page.scss'],
  imports: [FormsModule, IonContent, IonList, IonText, IonButton, IonHeader, IonToolbar, IonTitle, IonItem, IonInput, IonButtons, IonBackButton]
})
export class EmailPage {
  changeSent = false;
  private readonly navCtrl = inject(NavController);
  private readonly profileService = inject(ProfileService);
  private readonly messagesService = inject(MessagesService);

  private static handleErrorResponse(form: NgForm, errorResponse: HttpErrorResponse) {
    const response: Errors = errorResponse.error;
    if (response?.errors) {
      displayFieldErrors(form, response.errors)
    }
  }

  async changeEmail(form: NgForm, newEmail: string, password: string): Promise<void> {
    const loading = await this.messagesService.showLoading('Changing email...');

    this.profileService.changeEmail(newEmail, password)
      .subscribe({
        next: () => {
          loading.dismiss();
          this.changeSent = true;
          this.messagesService.showSuccessToast('Email change confirmation successfully sent');
          this.navCtrl.navigateBack("profile");
        },
        error: err => {
          loading.dismiss();
          EmailPage.handleErrorResponse(form, err);
        }
      });
  }

}
