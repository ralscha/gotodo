import { Component, inject, signal } from '@angular/core';
import { MessagesService } from '../../service/messages.service';
import { FormField, FormRoot, form, minLength, required, schema } from '@angular/forms/signals';
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
  NavController,
} from '@ionic/angular/standalone';
import { ProfileService } from '../profile/profile.service';
import { HttpErrorResponse } from '@angular/common/http';
import { Errors } from '../../api/types';

@Component({
  selector: 'app-password',
  templateUrl: './password.page.html',
  styleUrls: ['./password.page.scss'],
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
export class PasswordPage {
  readonly submitError = signal<string | null>(null);
  readonly submitted = signal(false);
  readonly passwordModel = signal({ oldPassword: '', newPassword: '' });
  readonly passwordForm = form(
    this.passwordModel,
    schema((path) => {
      required(path.oldPassword);
      minLength(path.oldPassword, 8);
      required(path.newPassword);
      minLength(path.newPassword, 8);
    }),
  );
  private readonly profileService = inject(ProfileService);
  private readonly navCtrl = inject(NavController);
  private readonly messagesService = inject(MessagesService);

  async changePassword(): Promise<void> {
    this.submitted.set(true);
    this.passwordForm().markAsTouched();

    if (this.passwordForm().invalid()) {
      return;
    }

    const loading = await this.messagesService.showLoading('Changing password');
    this.submitError.set(null);
    const passwords = this.passwordModel();

    this.profileService.changePassword(passwords.oldPassword, passwords.newPassword).subscribe({
      next: async () => {
        await loading.dismiss();
        this.passwordForm().reset({ oldPassword: '', newPassword: '' });
        this.submitted.set(false);
        await this.messagesService.showSuccessToast('Password successfully changed');
        await this.navCtrl.back();
      },
      error: async (err: HttpErrorResponse) => {
        await loading.dismiss();
        this.setSubmitError(err.error);
      },
    });
  }

  private setSubmitError(response: Errors | undefined): void {
    const errors = response?.errors;
    if (errors?.['oldPassword']?.includes('invalid')) {
      this.submitError.set('passwordInvalid');
    } else if (errors?.['newPassword']?.includes('weak')) {
      this.submitError.set('weakPassword');
    } else {
      this.messagesService.showErrorToast('Changing password failed');
    }
  }
}
