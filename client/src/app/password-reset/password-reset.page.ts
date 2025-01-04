import {Component, OnInit} from '@angular/core';
import {AuthService} from '../service/auth.service';
import {ActivatedRoute, RouterLink} from '@angular/router';
import {MessagesService} from '../service/messages.service';
import {FormsModule, NgForm} from '@angular/forms';
import {HttpErrorResponse} from '@angular/common/http';
import {displayFieldErrors} from '../util';
import {Errors} from '../api/types';
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
  IonToolbar
} from "@ionic/angular/standalone";

@Component({
  selector: 'app-password-reset',
  templateUrl: './password-reset.page.html',
  imports: [FormsModule, RouterLink, IonRouterLink, IonContent, IonList, IonText, IonButton, IonHeader, IonToolbar, IonTitle, IonItem, IonInput]
})
export class PasswordResetPage implements OnInit {

  success: boolean | null = null;
  resetToken: string | null = null;
  showForm = true;

  constructor(private readonly authService: AuthService,
              private readonly route: ActivatedRoute,
              private readonly messagesService: MessagesService) {
  }

  async ngOnInit(): Promise<void> {
    this.resetToken = this.route.snapshot.paramMap.get('token');

    if (!this.resetToken) {
      this.showForm = false;
    }
  }

  async reset(form: NgForm, password: string): Promise<void> {
    const loading = await this.messagesService.showLoading('Changing Password...');
    if (this.resetToken !== null) {
      this.authService.resetPassword(this.resetToken, password)
        .subscribe({
          next: () => {
            loading.dismiss();
            this.messagesService.showSuccessToast('Password successfully changed');
            this.success = true;
          },
          error: err => {
            loading.dismiss();
            this.handleErrorResponse(form, err);
          }
        })
    }
  }

  private handleErrorResponse(form: NgForm, errorResponse: HttpErrorResponse) {
    const response: Errors = errorResponse.error;
    if (response?.errors) {
      displayFieldErrors(form, response.errors)
      this.success = false;
    } else {
      this.showForm = false;
    }
  }

}
