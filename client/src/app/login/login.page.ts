import {Component, OnInit} from '@angular/core';
import {NavController} from '@ionic/angular';
import {AuthService} from '../service/auth.service';
import {MessagesService} from '../service/messages.service';
import {take} from 'rxjs';
import {HttpErrorResponse} from '@angular/common/http';

@Component({
  selector: 'app-login',
  templateUrl: './login.page.html',
  styleUrls: ['./login.page.scss'],
})
export class LoginPage implements OnInit {

  capslockOn = false;

  constructor(private readonly navCtrl: NavController,
              private readonly authService: AuthService,
              private readonly messagesService: MessagesService) {
  }

  ngOnInit(): void {
    // is the user already authenticated
    this.authService.authority$.pipe(take(1)).subscribe(authority => {
      if (authority !== null) {
        this.navCtrl.navigateRoot('home');
      }
    });
  }

  async login(username: string, password: string): Promise<void> {
    this.authService
      .login(username, password)
      .subscribe({
        next: () => this.navCtrl.navigateRoot('home'),
        error: () => this.showLoginFailedToast()
      });
  }

  private showLoginFailedToast(): void {
    this.messagesService.showErrorToast('Login failed');
  }

}
