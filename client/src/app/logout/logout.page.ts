import {Component, OnInit} from '@angular/core';
import {AuthService} from '../service/auth.service';

@Component({
  selector: 'app-logout',
  templateUrl: './logout.page.html',
  styleUrls: ['./logout.page.scss'],
})
export class LogoutPage implements OnInit {

  showMsg = false;

  constructor(private readonly authService: AuthService) {
  }

  ngOnInit(): void {
    this.authService.logout().subscribe(() => this.showMsg = true);
  }

}
