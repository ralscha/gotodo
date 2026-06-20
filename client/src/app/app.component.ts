import { Component, computed, inject } from '@angular/core';
import { AuthService } from './service/auth.service';
import {
  IonApp,
  IonContent,
  IonHeader,
  IonIcon,
  IonItem,
  IonLabel,
  IonList,
  IonMenu,
  IonMenuToggle,
  IonRouterLink,
  IonRouterOutlet,
  IonSplitPane,
  IonTitle,
  IonToolbar,
} from '@ionic/angular/standalone';
import { addIcons } from 'ionicons';
import { checkmarkCircleOutline, logOutOutline, person } from 'ionicons/icons';
import { RouterLink, RouterLinkActive } from '@angular/router';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
  imports: [
    IonContent,
    IonList,
    IonItem,
    IonMenu,
    IonSplitPane,
    IonMenuToggle,
    IonIcon,
    IonLabel,
    IonRouterOutlet,
    IonApp,
    IonHeader,
    IonTitle,
    IonToolbar,
    RouterLink,
    IonRouterLink,
    RouterLinkActive,
  ],
})
export class AppComponent {
  private readonly authService = inject(AuthService);
  readonly authenticated = this.authService.authenticated;
  readonly appPages = computed(() =>
    this.authenticated()
      ? [
          { title: 'My Todos', url: '/todo', icon: 'checkmark-circle-outline' },
          { title: 'Profile', url: '/profile', icon: 'person' },
          { title: 'Log out', url: '/logout', icon: 'log-out-outline' },
        ]
      : [],
  );

  constructor() {
    addIcons({ checkmarkCircleOutline, person, logOutOutline });
  }
}
