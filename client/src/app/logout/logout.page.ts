import {Component} from '@angular/core';
import {
  IonButton,
  IonContent,
  IonHeader,
  IonList,
  IonRouterLink,
  IonText,
  IonTitle,
  IonToolbar
} from "@ionic/angular/standalone";
import {RouterLink} from "@angular/router";

@Component({
    selector: 'app-logout',
    templateUrl: './logout.page.html',
  imports: [RouterLink, IonRouterLink, IonContent, IonList, IonText, IonButton, IonHeader, IonToolbar, IonTitle]
})
export class LogoutPage {

}
