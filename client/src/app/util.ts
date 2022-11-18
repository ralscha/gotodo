import {NgForm} from '@angular/forms';

export function displayFieldErrors(form: NgForm, fieldErrors: { [key: string]: string[] }): void {
  for (const [key, value] of Object.entries(fieldErrors)) {
    const comp = form.form.get(key);
    if (comp) {
      for (const v of value) {
        comp.setErrors({[v]: true});
      }
    }
  }
}
