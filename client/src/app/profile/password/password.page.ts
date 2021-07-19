import {Component, ViewChild} from '@angular/core';
import {MessagesService} from '../../service/messages.service';
import {NgForm} from '@angular/forms';
import {AuthService} from '../../service/auth.service';
import {NavController} from '@ionic/angular';
import {ProfileService} from '../profile/profile.service';
import {HttpErrorResponse} from '@angular/common/http';
import {FormErrorResponse} from '../../model/form-error-response';
import {displayFieldErrors} from '../../util';

@Component({
  selector: 'app-password',
  templateUrl: './password.page.html',
  styleUrls: ['./password.page.scss'],
})
export class PasswordPage {

  @ViewChild('changeForm')
  changeForm!: NgForm;

  constructor(private readonly profileService: ProfileService,
              private readonly authService: AuthService,
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
    const response: FormErrorResponse = errorResponse.error;
    if (response && response.fieldErrors) {
      displayFieldErrors(this.changeForm, response.fieldErrors)
    } else {
      this.changeForm.resetForm();
      this.messagesService.showErrorToast('Changing password failed');
    }
  }

}
