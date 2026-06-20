import { Component, inject, OnInit, signal } from '@angular/core';
import {
  IonButton,
  IonCol,
  IonContent,
  IonGrid,
  IonHeader,
  IonRouterLink,
  IonRow,
  IonText,
  IonTitle,
  IonToolbar,
} from '@ionic/angular/standalone';
import { RouterLink } from '@angular/router';
import { AuthService } from '../service/auth.service';

@Component({
  selector: 'app-logout',
  templateUrl: './logout.page.html',
  imports: [
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
  ],
})
export class LogoutPage implements OnInit {
  private readonly authService = inject(AuthService);
  readonly showMsg = signal(false);

  ngOnInit(): void {
    this.authService.logout().subscribe(() => this.showMsg.set(true));
  }
}
