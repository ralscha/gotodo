import {Component} from '@angular/core';
import {AuthService} from './service/auth.service';
import {NavController} from '@ionic/angular';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  dark = false

  constructor(readonly authService: AuthService,
              private readonly navCtrl: NavController,) {
  }

  logout(): void {
    this.authService.logout();
    this.navCtrl.navigateRoot('logout');
  }
}
