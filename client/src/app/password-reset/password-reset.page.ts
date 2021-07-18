import {Component, OnInit} from '@angular/core';
import {AuthService} from '../service/auth.service';
import {ActivatedRoute} from '@angular/router';
import {MessagesService} from '../service/messages.service';
import {NgForm} from '@angular/forms';
import {HttpErrorResponse} from '@angular/common/http';
import {FormErrorResponse} from '../model/form-error-response';
import {displayFieldErrors} from '../util';

@Component({
  selector: 'app-password-reset',
  templateUrl: './password-reset.page.html',
  styleUrls: ['./password-reset.page.scss'],
})
export class PasswordResetPage implements OnInit {

  success: boolean | null = null;
  resetToken: string | null = null;
  showForm = true;

  constructor(private readonly authService: AuthService,
              private readonly route: ActivatedRoute,
              private readonly messagesService: MessagesService) {
  }

  async ngOnInit(): Promise<void> {
    this.resetToken = this.route.snapshot.paramMap.get('token');

    if (!this.resetToken) {
      this.showForm = false;
    }
  }

  async reset(form: NgForm, password: string): Promise<void> {
    const loading = await this.messagesService.showLoading('Changing Password...');
    if (this.resetToken !== null) {
      this.authService.resetPassword(this.resetToken, password)
        .subscribe({
          next: () => {
            loading.dismiss();
            this.messagesService.showSuccessToast('Password successfully changed');
            this.success = true;
          },
          error: err => {
            loading.dismiss();
            this.handleErrorResponse(form, err);
          }
        })
    }
  }

  private handleErrorResponse(form: NgForm, errorResponse: HttpErrorResponse) {
    const response: FormErrorResponse = errorResponse.error;
    if (response && response.fieldErrors) {
      displayFieldErrors(form, response.fieldErrors)
      this.success = false;
    } else {
      this.showForm = false;
    }
  }

}
