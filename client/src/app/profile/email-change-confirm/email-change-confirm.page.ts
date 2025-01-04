import {Component, OnInit} from '@angular/core';
import {ActivatedRoute, RouterLink} from '@angular/router';
import {MessagesService} from '../../service/messages.service';
import {ProfileService} from '../profile/profile.service';
import {
  IonButton,
  IonContent,
  IonHeader,
  IonList,
  IonRouterLink,
  IonText,
  IonTitle,
  IonToolbar
} from '@ionic/angular/standalone';

@Component({
  selector: 'app-email-change-confirm',
  templateUrl: './email-change-confirm.page.html',
  imports: [RouterLink, IonRouterLink, IonContent, IonList, IonText, IonButton, IonHeader, IonToolbar, IonTitle]
})
export class EmailChangeConfirmPage implements OnInit {

  success: boolean | null = null;

  constructor(private readonly profileService: ProfileService,
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

    this.profileService.confirmEmailChange(token)
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
