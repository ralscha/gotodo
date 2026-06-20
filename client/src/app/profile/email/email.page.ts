import { Component, inject, signal } from '@angular/core';
import { MessagesService } from '../../service/messages.service';
import { ProfileService } from '../profile/profile.service';
import { email, FormField, form, minLength, required, schema } from '@angular/forms/signals';
import { HttpErrorResponse } from '@angular/common/http';
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
import { Errors } from '../../api/types';

@Component({
  selector: 'app-email',
  templateUrl: './email.page.html',
  styleUrls: ['./email.page.scss'],
  imports: [
    FormField,
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
export class EmailPage {
  readonly submitError = signal<string | null>(null);
  readonly changeSent = signal(false);
  readonly submitted = signal(false);
  readonly emailModel = signal({ password: '', newEmail: '' });
  readonly emailForm = form(
    this.emailModel,
    schema((path) => {
      required(path.password);
      minLength(path.password, 8);
      required(path.newEmail);
      email(path.newEmail);
    }),
  );
  private readonly profileService = inject(ProfileService);
  private readonly messagesService = inject(MessagesService);

  async changeEmail(): Promise<void> {
    this.submitted.set(true);
    this.emailForm().markAsTouched();

    if (this.emailForm().invalid()) {
      return;
    }

    const loading = await this.messagesService.showLoading('Changing email');
    this.submitError.set(null);
    const emailChange = this.emailModel();

    this.profileService.changeEmail(emailChange.newEmail, emailChange.password).subscribe({
      next: async () => {
        await loading.dismiss();
        await this.messagesService.showSuccessToast('Email change confirmation successfully sent');
        this.changeSent.set(true);
      },
      error: async (err: HttpErrorResponse) => {
        await loading.dismiss();
        this.setSubmitError(err.error);
      },
    });
  }

  private setSubmitError(response: Errors | undefined): void {
    const errors = response?.errors;
    if (errors?.['password']?.includes('invalid')) {
      this.submitError.set('passwordInvalid');
    } else if (errors?.['newEmail']?.includes('exists')) {
      this.submitError.set('emailExists');
    } else {
      this.messagesService.showErrorToast('Changing email failed');
    }
  }
}
