import {NgModule} from '@angular/core';
import {PreloadAllModules, RouterModule, Routes} from '@angular/router';
import {AuthGuard} from './service/auth.guard';
import {LoginPage} from './login/login.page';
import {LogoutPage} from './logout/logout.page';

const routes: Routes = [
  {path: '', redirectTo: 'home', pathMatch: 'full'},
  {path: 'home', loadChildren: () => import('./home/home.module').then(m => m.HomePageModule), canActivate: [AuthGuard]},
  {path: 'login', component: LoginPage},
  {path: 'logout', component: LogoutPage},
  {path: '**', redirectTo: 'home'}
];

@NgModule({
  imports: [
    RouterModule.forRoot(routes, { preloadingStrategy: PreloadAllModules, useHash: true })
  ],
  exports: [RouterModule]
})
export class AppRoutingModule {
}
