import {Component} from '@angular/core';
import {NavController} from '@ionic/angular';
import {AuthService} from '../service/auth.service';
import {MessagesService} from '../service/messages.service';
import {NgForm} from '@angular/forms';
import { HttpErrorResponse } from '@angular/common/http';
import {displayFieldErrors} from '../util';
import {Errors} from '../api/types';

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
    const response: Errors = errorResponse.error;
    if (response?.errors) {
      displayFieldErrors(form, response.errors)
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
