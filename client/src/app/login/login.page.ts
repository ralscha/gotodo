import {Component, OnInit} from '@angular/core';
import {NavController} from '@ionic/angular';
import {AuthService} from '../service/auth.service';
import {MessagesService} from '../service/messages.service';
import {take} from 'rxjs';
import {HttpClient} from '@angular/common/http';

@Component({
  selector: 'app-login',
  templateUrl: './login.page.html',
  styleUrls: ['./login.page.scss'],
})
export class LoginPage implements OnInit {

  capslockOn = false;

  constructor(private readonly navCtrl: NavController,
              private httpClient: HttpClient,
              private readonly authService: AuthService,
              private readonly messagesService: MessagesService) {
    httpClient.get('/v1/secret', {responseType: 'text'}).subscribe(data => console.log(data));
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
        next:  authority => {
          if (authority !== null) {
            this.navCtrl.navigateRoot('home')
          } else {
            this.showLoginFailedToast()
          }
        },
        error:  () => this.showLoginFailedToast()
      });
  }

  private showLoginFailedToast(): void {
    this.messagesService.showErrorToast('Login failed');
  }

}
