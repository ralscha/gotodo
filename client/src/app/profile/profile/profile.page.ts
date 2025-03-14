import {Component, OnInit} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {environment} from '../../../environments/environment';
import {AppVersionOutput} from '../../api/types';
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
  IonToolbar
} from "@ionic/angular/standalone";
import {RouterLink} from "@angular/router";

@Component({
  selector: 'app-profile',
  templateUrl: './profile.page.html',
  styleUrls: ['./profile.page.scss'],
  imports: [RouterLink, IonRouterLink, IonContent, IonButton, IonHeader, IonToolbar, IonTitle, IonButtons, IonMenuButton, IonGrid, IonRow, IonCol, IonFooter, IonNote]
})
export class ProfilePage implements OnInit {

  buildInfo: {
    serverVersion: string | null,
    serverBuildTime: string | null,
    clientVersion: string | null,
    clientBuildTime: string | null
  } = {serverVersion: null, serverBuildTime: null, clientVersion: null, clientBuildTime: null};

  constructor(private readonly httpClient: HttpClient) {
  }

  ngOnInit(): void {
    this.buildInfo.clientVersion = environment.version;
    this.buildInfo.clientBuildTime = new Date(environment.buildTimestamp * 1000).toISOString();

    this.httpClient.get<AppVersionOutput>('/v1/profile/build-info')
      .subscribe(response => {
        this.buildInfo.serverVersion = response.version;
        this.buildInfo.serverBuildTime = response.buildTime;
      });
  }

}
