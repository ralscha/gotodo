import {Component, inject, OnInit} from '@angular/core';
import {
  IonButton,
  IonContent,
  IonHeader,
  IonInput,
  IonItem,
  IonList,
  IonRouterLink,
  IonText,
  IonTitle,
  IonToolbar,
  NavController
} from '@ionic/angular/standalone';
import {AuthService} from '../service/auth.service';
import {MessagesService} from '../service/messages.service';
import {take} from 'rxjs';
import {HttpErrorResponse} from '@angular/common/http';
import {FormsModule, NgForm} from '@angular/forms';
import {displayFieldErrors} from '../util';
import {Errors} from '../api/types';
import {RouterLink} from "@angular/router";

@Component({
  selector: 'app-login',
  templateUrl: './login.page.html',
  styleUrls: ['./login.page.scss'],
  imports: [FormsModule, RouterLink, IonRouterLink, IonContent, IonList, IonText, IonButton, IonHeader, IonToolbar, IonTitle, IonItem, IonInput]
})
export class LoginPage implements OnInit {
  private readonly navCtrl = inject(NavController);
  private readonly authService = inject(AuthService);
  private readonly messagesService = inject(MessagesService);


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
      const response: Errors = errorResponse.error;
      if (response?.errors) {
        displayFieldErrors(form, response.errors)
      } else {
        this.messagesService.showErrorToast('Login failed');
      }
    };
  }

}
