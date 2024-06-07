import {NgModule} from '@angular/core';
import {BrowserModule} from '@angular/platform-browser';
import {RouteReuseStrategy} from '@angular/router';

import {IonicModule, IonicRouteStrategy} from '@ionic/angular';

import {AppComponent} from './app.component';
import {AppRoutingModule} from './app-routing.module';
import { provideHttpClient, withInterceptorsFromDi } from '@angular/common/http';
import {FormsModule} from '@angular/forms';
import {LoginPage} from './login/login.page';
import {LogoutPage} from './logout/logout.page';

@NgModule({ declarations: [AppComponent, LoginPage, LogoutPage],
    bootstrap: [AppComponent], imports: [BrowserModule,
        FormsModule,
        IonicModule.forRoot(),
        AppRoutingModule], providers: [
        { provide: RouteReuseStrategy, useClass: IonicRouteStrategy },
        provideHttpClient(withInterceptorsFromDi())
    ] })
export class AppModule {
}
