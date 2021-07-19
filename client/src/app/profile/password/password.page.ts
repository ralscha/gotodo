import {Component, ViewChild} from '@angular/core';
import {MessagesService} from '../../service/messages.service';
import {NgForm} from '@angular/forms';
import {AuthService} from '../../service/auth.service';
import {NavController} from '@ionic/angular';
import {ProfileService} from '../profile/profile.service';

@Component({
  selector: 'app-password',
  templateUrl: './password.page.html',
  styleUrls: ['./password.page.scss'],
})
export class PasswordPage {

  submitError: string | null = null;

  @ViewChild('changeForm')
  changeForm!: NgForm;

  constructor(private readonly profileService: ProfileService,
              private readonly authService: AuthService,
              private readonly navCtrl: NavController,
              private readonly messagesService: MessagesService) {
  }

  async changePassword(oldPassword: string, newPassword: string): Promise<void> {

    const loading = await this.messagesService.showLoading('Changing password');
    this.submitError = null;

    this.profileService.changePassword(oldPassword, newPassword)
      .subscribe(async (response) => {
        await loading.dismiss();
        if (!response) {
          this.changeForm.resetForm();
          await this.messagesService.showSuccessToast('Password successfully changed');
          // await this.authService.deleteTokens();
          await this.navCtrl.navigateRoot('/login');
        } else if (response === 'INVALID') {
          this.submitError = 'passwordInvalid';
          await this.messagesService.showErrorToast('Old Password invalid');
        } else if (response === 'WEAK_PASSWORD') {
          this.submitError = 'weakPassword';
          await this.messagesService.showErrorToast('Weak password');
        }
      }, () => {
        loading.dismiss();
        this.messagesService.showErrorToast('Changing password failed');
      });

  }

}
