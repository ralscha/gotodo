import {Component} from '@angular/core';
import {MessagesService} from '../../service/messages.service';
import {ProfileService} from '../profile/profile.service';
import {NgForm} from '@angular/forms';
import {HttpErrorResponse} from '@angular/common/http';
import {displayFieldErrors} from '../../util';
import {NavController} from '@ionic/angular';
import {Errors} from '../../api/types';

@Component({
  selector: 'app-email',
  templateUrl: './email.page.html',
  styleUrls: ['./email.page.scss'],
})
export class EmailPage {

  changeSent = false;

  constructor(private readonly navCtrl: NavController,
              private readonly profileService: ProfileService,
              private readonly messagesService: MessagesService) {
  }

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
