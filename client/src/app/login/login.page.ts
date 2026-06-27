import { Component, inject, OnInit, signal } from '@angular/core';
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
  NavController,
} from '@ionic/angular/standalone';
import { AuthService } from '../service/auth.service';
import { MessagesService } from '../service/messages.service';
import { take } from 'rxjs';
import {
  email,
  FormField,
  FormRoot,
  form,
  minLength,
  required,
  schema,
} from '@angular/forms/signals';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'app-login',
  templateUrl: './login.page.html',
  styleUrls: ['./login.page.scss'],
  imports: [
    FormField,
    FormRoot,
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
export class LoginPage implements OnInit {
  readonly submitted = signal(false);
  readonly loginModel = signal({ email: '', password: '' });
  readonly loginForm = form(
    this.loginModel,
    schema((path) => {
      required(path.email);
      email(path.email);
      required(path.password);
      minLength(path.password, 8);
    }),
  );
  private readonly navCtrl = inject(NavController);
  private readonly authService = inject(AuthService);
  private readonly messagesService = inject(MessagesService);

  ngOnInit(): void {
    this.authService.authority$.pipe(take(1)).subscribe((authority) => {
      if (authority !== null) {
        this.navCtrl.navigateRoot('/todo');
      }
    });
  }

  async login(): Promise<void> {
    this.submitted.set(true);
    this.loginForm().markAsTouched();

    if (this.loginForm().invalid()) {
      return;
    }

    const loading = await this.messagesService.showLoading('Logging in');
    const credentials = this.loginModel();

    this.authService.login(credentials.email, credentials.password).subscribe({
      next: async () => {
        await loading.dismiss();
        await this.navCtrl.navigateRoot('/todo');
      },
      error: async () => {
        await loading.dismiss();
        await this.messagesService.showErrorToast('Login failed');
      },
    });
  }
}
