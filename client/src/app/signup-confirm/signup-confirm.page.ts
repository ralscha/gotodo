import {Component, OnInit} from '@angular/core';
import {AuthService} from '../service/auth.service';
import {MessagesService} from '../service/messages.service';
import {ActivatedRoute, RouterLink} from '@angular/router';
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

@Component({
  selector: 'app-signup-confirm',
  templateUrl: './signup-confirm.page.html',
  imports: [RouterLink, IonRouterLink, IonContent, IonList, IonText, IonButton, IonHeader, IonToolbar, IonTitle]
})
export class SignupConfirmPage implements OnInit {

  success: boolean | null = null;

  constructor(private readonly authService: AuthService,
              private readonly route: ActivatedRoute,
              private readonly messagesService: MessagesService) {
  }

  async ngOnInit(): Promise<void> {
    const token = this.route.snapshot.paramMap.get('token');

    if (!token) {
      this.success = false;
      return;
    }

    const loading = await this.messagesService.showLoading('Processing confirmation...');

    this.authService.confirmSignup(token)
      .subscribe({
        next: () => {
          loading.dismiss();
          this.success = true;
        },
        error: () => {
          loading.dismiss();
          this.success = false;
        }
      });
  }

}
