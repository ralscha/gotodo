import { Component, inject, signal } from '@angular/core';
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
import { AuthService } from '../service/auth.service';
import { MessagesService } from '../service/messages.service';
import {
  email,
  FormField,
  FormRoot,
  form,
  minLength,
  required,
  schema,
} from '@angular/forms/signals';
import { HttpErrorResponse } from '@angular/common/http';
import { Errors } from '../api/types';

@Component({
  selector: 'app-signup',
  templateUrl: './signup.page.html',
  styleUrls: ['./signup.page.scss'],
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
    IonBackButton,
    IonButtons,
  ],
})
export class SignupPage {
  readonly signUpSent = signal(false);
  readonly submitError = signal<string | null>(null);
  readonly submitted = signal(false);
  readonly signupModel = signal({ email: '', password: '' });
  readonly signupForm = form(
    this.signupModel,
    schema((path) => {
      required(path.email);
      email(path.email);
      required(path.password);
      minLength(path.password, 8);
    }),
  );
  private readonly authService = inject(AuthService);
  private readonly messagesService = inject(MessagesService);

  async signup(): Promise<void> {
    this.submitted.set(true);
    this.signupForm().markAsTouched();

    if (this.signupForm().invalid()) {
      return;
    }

    const loading = await this.messagesService.showLoading('Signing up');
    this.submitError.set(null);
    const signup = this.signupModel();

    this.authService.signup(signup.email, signup.password).subscribe({
      next: async () => {
        await loading.dismiss();
        this.signUpSent.set(true);
        await this.messagesService.showSuccessToast('Sign-up confirmation successfully sent');
      },
      error: async (err: HttpErrorResponse) => {
        await loading.dismiss();
        this.setSubmitError(err.error);
      },
    });
  }

  private setSubmitError(response: Errors | undefined): void {
    const errors = response?.errors;
    if (errors?.['email']?.includes('exists')) {
      this.submitError.set('emailExists');
    }
    if (errors?.['password']?.includes('weak')) {
      this.submitError.set('weakPassword');
    }
  }
}
