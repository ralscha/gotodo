import {Component, OnInit} from '@angular/core';
import {NavController} from '@ionic/angular';
import {AuthService} from '../service/auth.service';
import {MessagesService} from '../service/messages.service';
import {take} from 'rxjs';
import {HttpErrorResponse} from '@angular/common/http';
import {NgForm} from '@angular/forms';
import {FormErrorResponse} from '../model/form-error-response';
import {displayFieldErrors} from '../util';

@Component({
  selector: 'app-login',
  templateUrl: './login.page.html',
  styleUrls: ['./login.page.scss'],
})
export class LoginPage implements OnInit {

  capslockOn = false;

  constructor(private readonly navCtrl: NavController,
              private readonly authService: AuthService,
              private readonly messagesService: MessagesService) {
  }

  ngOnInit(): void {
    // is the user already authenticated
    this.authService.authority$.pipe(take(1)).subscribe(authority => {
      if (authority !== null) {
        this.navCtrl.navigateRoot('home');
      }
    });
  }

  async login(form: NgForm, email: string, password: string): Promise<void> {
    this.authService
      .login(email, password)
      .subscribe({
        next: () => this.navCtrl.navigateRoot('home'),
        error: this.handleErrorResponse(form)
      });
  }

  private handleErrorResponse(form: NgForm) {
    return (errorResponse: HttpErrorResponse) => {
      const response: FormErrorResponse = errorResponse.error;
      if (response && response.fieldErrors) {
        displayFieldErrors(form, response.fieldErrors)
      } else {
        this.messagesService.showErrorToast('Login failed');
      }
    };
  }

}
