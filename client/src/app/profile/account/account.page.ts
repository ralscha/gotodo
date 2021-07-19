import {Component, ViewChild} from '@angular/core';
import {MessagesService} from '../../service/messages.service';
import {AlertController, NavController} from '@ionic/angular';
import {NgForm} from '@angular/forms';
import {ProfileService} from '../profile/profile.service';
import {HttpErrorResponse} from '@angular/common/http';
import {FormErrorResponse} from '../../model/form-error-response';
import {displayFieldErrors} from '../../util';
import {AuthService} from '../../service/auth.service';

@Component({
  selector: 'app-account',
  templateUrl: './account.page.html',
  styleUrls: ['./account.page.scss'],
})
export class AccountPage {

  submitError: string | null = null;

  @ViewChild('deleteForm')
  deleteForm!: NgForm;

  constructor(private readonly navCtrl: NavController,
              private readonly authService: AuthService,
              private readonly profileService: ProfileService,
              private readonly messagesService: MessagesService,
              private readonly alertController: AlertController) {
  }

  async deleteAccount(password: string): Promise<void> {
    this.submitError = null;

    const alert = await this.alertController.create({
      header: 'Delete Account',
      message: 'Do you really want to delete your account? This action is <strong>irreversible!</strong>',
      buttons: [
        {
          text: 'Cancel',
          role: 'cancel',
          handler: () => {
            this.deleteForm.resetForm();
          }
        }, {
          text: 'Delete Account',
          handler: () => this.reallyDeleteAccount(password)
        }
      ]
    });

    await alert.present();

  }

  private async reallyDeleteAccount(password: string): Promise<void> {
    const loading = await this.messagesService.showLoading('Deleting account...');

    this.profileService.deleteAccount(password)
      .subscribe({
        next: () => {
          loading.dismiss();
          this.messagesService.showSuccessToast('Account successfully deleted');
          this.authService.logoutClient();
          this.navCtrl.navigateRoot('/login');
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
      displayFieldErrors(this.deleteForm, response.fieldErrors)
    } else {
      this.deleteForm.resetForm();
      this.messagesService.showErrorToast('Deleting account failed');
    }
  }
}
