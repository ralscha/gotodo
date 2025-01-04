import {Component} from '@angular/core';
import {AuthService} from './service/auth.service';
import {
  IonApp,
  IonContent,
  IonIcon, IonItem, IonLabel,
  IonList, IonListHeader, IonMenu, IonMenuToggle, IonRouterLink, IonRouterOutlet, IonSplitPane,
  NavController
} from '@ionic/angular/standalone';
import {addIcons} from "ionicons";
import {logOutOutline, person} from "ionicons/icons";
import {RouterLink, RouterLinkActive} from "@angular/router";

@Component({
    selector: 'app-root',
    templateUrl: './app.component.html',
    styleUrls: ['./app.component.scss'],
  imports: [IonContent, IonList, IonItem, IonMenu, IonSplitPane, IonListHeader, IonMenuToggle, IonIcon, IonLabel, IonRouterOutlet, IonApp, RouterLink, IonRouterLink, RouterLinkActive]
})
export class AppComponent {
  dark = false

  constructor(readonly authService: AuthService,
              private readonly navCtrl: NavController) {
    addIcons({ person, logOutOutline });
  }

  logout(): void {
    this.authService.logout().subscribe(() => this.navCtrl.navigateRoot('logout'));
  }
}
