import {inject, Injectable} from '@angular/core';
import {LoadingController, ToastController} from '@ionic/angular/standalone';

@Injectable({
  providedIn: 'root'
})
export class MessagesService {
  private readonly toastCtrl = inject(ToastController);
  private readonly loadingCtrl = inject(LoadingController);


  async showLoading(message = 'Working'): Promise<HTMLIonLoadingElement> {
    const loading = await this.loadingCtrl.create({
      spinner: 'bubbles',
      message
    });
    await loading.present();
    return loading;
  }

  async showSuccessToast(message: string, duration = 4000): Promise<void> {
    const toast = await this.toastCtrl.create({
      message,
      duration,
      position: 'bottom',
      color: 'success',
    });
    await toast.present();
  }

  async showErrorToast(message = 'Unexpected error occurred'): Promise<void> {
    const toast = await this.toastCtrl.create({
      message,
      duration: 4000,
      position: 'bottom',
      color: 'danger'
    });
    await toast.present();
  }

}
