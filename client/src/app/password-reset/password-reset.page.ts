import { Component, inject, OnInit, signal } from '@angular/core';
import { AuthService } from '../service/auth.service';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { MessagesService } from '../service/messages.service';
import { FormField, form, minLength, required, schema } from '@angular/forms/signals';
import { HttpErrorResponse } from '@angular/common/http';
import { Errors } from '../api/types';
import {
  IonButton,
  IonCol,
  IonContent,
  IonGrid,
  IonHeader,
  IonInput,
  IonItem,
  IonRouterLink,
  IonRow,
  IonText,
  IonTitle,
  IonToolbar,
} from '@ionic/angular/standalone';

@Component({
  selector: 'app-password-reset',
  templateUrl: './password-reset.page.html',
  imports: [
    FormField,
    RouterLink,
    IonRouterLink,
    IonContent,
    IonGrid,
    IonRow,
    IonCol,
    IonText,
    IonButton,
    IonHeader,
    IonToolbar,
    IonTitle,
    IonItem,
    IonInput,
  ],
})
export class PasswordResetPage implements OnInit {
  readonly success = signal<boolean | null>(null);
  readonly submitError = signal<string | null>(null);
  readonly submitted = signal(false);
  readonly resetModel = signal({ password: '' });
  readonly resetForm = form(
    this.resetModel,
    schema((path) => {
      required(path.password);
      minLength(path.password, 8);
    }),
  );
  private resetToken: string | null = null;
  private readonly authService = inject(AuthService);
  private readonly route = inject(ActivatedRoute);
  private readonly messagesService = inject(MessagesService);

  ngOnInit(): void {
    this.resetToken = this.route.snapshot.paramMap.get('token');

    if (!this.resetToken) {
      this.success.set(false);
    }
  }

  async reset(): Promise<void> {
    this.submitted.set(true);
    this.resetForm().markAsTouched();

    if (this.resetForm().invalid() || this.resetToken === null) {
      return;
    }

    const loading = await this.messagesService.showLoading('Changing Password');
    this.submitError.set(null);

    this.authService.resetPassword(this.resetToken, this.resetModel().password).subscribe({
      next: async () => {
        await loading.dismiss();
        await this.messagesService.showSuccessToast('Password successfully changed');
        this.success.set(true);
      },
      error: async (err: HttpErrorResponse) => {
        await loading.dismiss();
        this.handleErrorResponse(err.error);
      },
    });
  }

  private handleErrorResponse(response: Errors | undefined): void {
    const errors = response?.errors;
    if (errors?.['password']?.includes('weak')) {
      this.submitError.set('weakPassword');
    } else {
      this.success.set(false);
    }
  }
}
