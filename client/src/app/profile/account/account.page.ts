import { Component, inject, signal } from '@angular/core';
import { MessagesService } from '../../service/messages.service';
import {
  AlertController,
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
import { FormField, form, minLength, required, schema } from '@angular/forms/signals';
import { ProfileService } from '../profile/profile.service';
import { HttpErrorResponse } from '@angular/common/http';
import { AuthService } from '../../service/auth.service';
import { Errors } from '../../api/types';

@Component({
  selector: 'app-account',
  templateUrl: './account.page.html',
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
export class AccountPage {
  readonly submitError = signal<string | null>(null);
  readonly submitted = signal(false);
  readonly deleteModel = signal({ password: '' });
  readonly deleteForm = form(
    this.deleteModel,
    schema((path) => {
      required(path.password);
      minLength(path.password, 8);
    }),
  );
  private readonly navCtrl = inject(NavController);
  private readonly authService = inject(AuthService);
  private readonly profileService = inject(ProfileService);
  private readonly messagesService = inject(MessagesService);
  private readonly alertController = inject(AlertController);

  async deleteAccount(): Promise<void> {
    this.submitted.set(true);
    this.deleteForm().markAsTouched();

    if (this.deleteForm().invalid()) {
      return;
    }

    this.submitError.set(null);
    const password = this.deleteModel().password;

    const alert = await this.alertController.create({
      header: 'Delete Account',
      message:
        'Do you really want to delete your account? This action is <strong>irreversible!</strong>',
      buttons: [
        {
          text: 'Cancel',
          role: 'cancel',
          handler: () => this.resetDeleteForm(),
        },
        {
          text: 'Delete Account',
          handler: () => this.reallyDeleteAccount(password),
        },
      ],
    });

    await alert.present();
  }

  private async reallyDeleteAccount(password: string): Promise<void> {
    const loading = await this.messagesService.showLoading('Deleting account');

    this.profileService.deleteAccount(password).subscribe({
      next: async () => {
        await loading.dismiss();
        await this.messagesService.showSuccessToast('Account successfully deleted');
        this.authService.logoutClient();
        await this.navCtrl.navigateRoot('/login');
      },
      error: async (err: HttpErrorResponse) => {
        await loading.dismiss();
        this.handleErrorResponse(err.error);
      },
    });
  }

  private handleErrorResponse(response: Errors | undefined): void {
    if (response?.errors?.['password']?.includes('invalid')) {
      this.submitError.set('passwordInvalid');
    } else {
      this.resetDeleteForm();
      this.messagesService.showErrorToast('Deleting account failed');
    }
  }

  private resetDeleteForm(): void {
    this.deleteForm().reset({ password: '' });
    this.submitted.set(false);
  }
}
