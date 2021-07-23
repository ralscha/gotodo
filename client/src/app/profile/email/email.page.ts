import {Component} from '@angular/core';
import {MessagesService} from '../../service/messages.service';
import {ProfileService} from '../profile/profile.service';
import {NgForm} from '@angular/forms';
import {HttpErrorResponse} from '@angular/common/http';
import {FormErrorResponse} from '../../model/form-error-response';
import {displayFieldErrors} from '../../util';
import {NavController} from '@ionic/angular';

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
    const response: FormErrorResponse = errorResponse.error;
    if (response && response.fieldErrors) {
      displayFieldErrors(form, response.fieldErrors)
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
