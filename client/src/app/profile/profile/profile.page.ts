import { Component, computed } from '@angular/core';
import { httpResource } from '@angular/common/http';
import { environment } from '../../../environments/environment';
import { AppVersionOutput } from '../../api/types';
import {
  IonButton,
  IonButtons,
  IonCol,
  IonContent,
  IonFooter,
  IonGrid,
  IonHeader,
  IonMenuButton,
  IonNote,
  IonRouterLink,
  IonRow,
  IonTitle,
  IonToolbar,
} from '@ionic/angular/standalone';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'app-profile',
  templateUrl: './profile.page.html',
  styleUrls: ['./profile.page.scss'],
  imports: [
    RouterLink,
    IonRouterLink,
    IonContent,
    IonButton,
    IonHeader,
    IonToolbar,
    IonTitle,
    IonButtons,
    IonMenuButton,
    IonGrid,
    IonRow,
    IonCol,
    IonFooter,
    IonNote,
  ],
})
export class ProfilePage {
  private readonly serverBuildInfo = httpResource<AppVersionOutput>(() => '/v1/profile/build-info');

  readonly buildInfo = computed(() => ({
    serverVersion: this.serverBuildInfo.value()?.version ?? null,
    serverBuildTime: this.serverBuildInfo.value()?.buildTime ?? null,
    clientVersion: environment.version,
    clientBuildTime: new Date(environment.buildTimestamp * 1000).toISOString(),
  }));
}
