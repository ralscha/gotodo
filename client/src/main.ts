import {
  PreloadAllModules,
  provideRouter,
  RouteReuseStrategy,
  withHashLocation,
  withPreloading,
} from '@angular/router';
import { IonicRouteStrategy, provideIonicAngular } from '@ionic/angular';
import { provideHttpClient, withInterceptors, withXhr } from '@angular/common/http';
import { bootstrapApplication } from '@angular/platform-browser';
import { AppComponent } from './app/app.component';
import { routes } from './app/app-routing';
import { sessionExpiredInterceptor } from './app/service/session-expired.interceptor';

bootstrapApplication(AppComponent, {
  providers: [
    provideIonicAngular(),
    { provide: RouteReuseStrategy, useClass: IonicRouteStrategy },
    provideHttpClient(withXhr(), withInterceptors([sessionExpiredInterceptor])),
    provideRouter(routes, withHashLocation(), withPreloading(PreloadAllModules)),
  ],
}).catch((err) => console.error(err));
