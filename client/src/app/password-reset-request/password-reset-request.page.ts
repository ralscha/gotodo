import {Component} from '@angular/core';
import {NavController} from '@ionic/angular';
import {AuthService} from '../service/auth.service';
import {MessagesService} from '../service/messages.service';

@Component({
  selector: 'app-password-reset-request',
  templateUrl: './password-reset-request.page.html'
})
export class PasswordResetRequestPage {

  resetSent = false;

  constructor(private readonly navCtrl: NavController,
              private readonly authService: AuthService,
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
