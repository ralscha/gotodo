import {Component} from '@angular/core';
import {NavController} from '@ionic/angular';
import {AuthService} from '../service/auth.service';
import {MessagesService} from '../service/messages.service';
import {NgForm} from '@angular/forms';
import {HttpErrorResponse} from '@angular/common/http';
import {FormErrorResponse} from '../model/form-error-response';
import {displayFieldErrors} from '../util';

@Component({
  selector: 'app-signup',
  templateUrl: './signup.page.html',
  styleUrls: ['./signup.page.scss'],
})
export class SignupPage {

  signUpSent = false;

  constructor(private readonly navCtrl: NavController,
              private readonly authService: AuthService,
              private readonly messagesService: MessagesService) {
  }

  private static handleErrorResponse(form: NgForm, errorResponse: HttpErrorResponse) {
    const response: FormErrorResponse = errorResponse.error;
    if (response && response.fieldErrors) {
      displayFieldErrors(form, response.fieldErrors)
    }
  }

  async signup(form: NgForm, email: string, password: string): Promise<void> {
    const loading = await this.messagesService.showLoading('Signing up...');
    this.authService.signup(email, password)
      .subscribe({
        next: () => {
          loading.dismiss();
          this.signUpSent = true;
          this.messagesService.showSuccessToast('Sign-up confirmation successful sent');
          this.navCtrl.navigateRoot('login');
        },
        error: err => {
          loading.dismiss();
          SignupPage.handleErrorResponse(form, err);
        }
      });
  }


}
