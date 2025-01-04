import {Component} from '@angular/core';
import {AuthService} from '../service/auth.service';
import {MessagesService} from '../service/messages.service';
import {FormsModule} from '@angular/forms';
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
  IonToolbar
} from "@ionic/angular/standalone";

@Component({
  selector: 'app-password-reset-request',
  templateUrl: './password-reset-request.page.html',
  imports: [FormsModule, IonContent, IonList, IonText, IonButton, IonHeader, IonToolbar, IonTitle, IonItem, IonInput, IonButtons, IonBackButton]
})
export class PasswordResetRequestPage {

  resetSent = false;

  constructor(private readonly authService: AuthService,
              private readonly messagesService: MessagesService) {
  }

  async resetRequest(email: string): Promise<void> {
    const loading = await this.messagesService.showLoading('Sending email...');

    this.authService.resetPasswordRequest(email)
      .subscribe({
        next: () => {
          this.messagesService.showSuccessToast('Email successfully sent');
          this.resetSent = true;
          loading.dismiss();
        },
        error: () => {
          this.messagesService.showErrorToast('Sending email failed');
          loading.dismiss();
        }
      });
  }


}
