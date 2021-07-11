import {Component} from '@angular/core';
import {HttpClient} from '@angular/common/http';

@Component({
  selector: 'app-home',
  templateUrl: 'home.page.html',
  styleUrls: ['home.page.scss'],
})
export class HomePage {

  constructor(readonly httpClient: HttpClient) {
    httpClient.get('/v1/secret', {responseType: 'text'}).subscribe(data => console.log(data));
  }

}
