import {Component, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {MessagesService} from '../../service/messages.service';
import {ProfileService} from '../profile/profile.service';

@Component({
    selector: 'app-email-change-confirm',
    templateUrl: './email-change-confirm.page.html',
    standalone: false
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
