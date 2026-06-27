import { Component, inject, signal } from '@angular/core';
import { AuthService } from '../service/auth.service';
import { MessagesService } from '../service/messages.service';
import { email, FormField, FormRoot, form, required, schema } from '@angular/forms/signals';
import {
  IonBackButton,
  IonButton,
  IonButtons,
  IonCol,
  IonContent,
  IonGrid,
  IonHeader,
  IonInput,
  IonItem,
  IonRow,
  IonText,
  IonTitle,
  IonToolbar,
} from '@ionic/angular/standalone';

@Component({
  selector: 'app-password-reset-request',
  templateUrl: './password-reset-request.page.html',
  imports: [
    FormField,
    FormRoot,
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
    IonButtons,
    IonBackButton,
  ],
})
export class PasswordResetRequestPage {
  readonly resetSent = signal(false);
  readonly submitted = signal(false);
  readonly resetRequestModel = signal({ email: '' });
  readonly resetRequestForm = form(
    this.resetRequestModel,
    schema((path) => {
      required(path.email);
      email(path.email);
    }),
  );
  private readonly authService = inject(AuthService);
  private readonly messagesService = inject(MessagesService);

  async resetRequest(): Promise<void> {
    this.submitted.set(true);
    this.resetRequestForm().markAsTouched();

    if (this.resetRequestForm().invalid()) {
      return;
    }

    const loading = await this.messagesService.showLoading('Sending email');

    this.authService.resetPasswordRequest(this.resetRequestModel().email).subscribe({
      next: async () => {
        await loading.dismiss();
        await this.messagesService.showSuccessToast('Email successfully sent');
        this.resetSent.set(true);
      },
      error: async () => {
        await loading.dismiss();
        await this.messagesService.showErrorToast('Sending email failed');
      },
    });
  }
}
