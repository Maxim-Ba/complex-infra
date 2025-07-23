import { AsyncPipe, JsonPipe } from '@angular/common';
import { Component, inject } from '@angular/core';
import {
  FormControl,
  FormGroup,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';
import {
  TuiAppearance,
  TuiButton,
  TuiError,
  TuiNotification,
  TuiTextfield,
  TuiTitle,
} from '@taiga-ui/core';
import { TuiFieldErrorPipe, TuiSegmented } from '@taiga-ui/kit';
import { TuiCardLarge, TuiForm, TuiHeader } from '@taiga-ui/layout';
import { ApiService } from '../../shared/services/api.service';
@Component({
  selector: 'app-auth',
  templateUrl: './auth.component.html',
  imports: [
    AsyncPipe,
    ReactiveFormsModule,
    TuiAppearance,
    TuiButton,
    TuiCardLarge,
    TuiError,
    TuiFieldErrorPipe,
    TuiForm,
    TuiHeader,
    TuiNotification,
    TuiSegmented,
    TuiTextfield,
    TuiTitle,
    JsonPipe,
  ],
})
export class Auth {
  logout() {
    this.apiService.logout().subscribe();
  }
  apiService = inject(ApiService);
  onCancel() {
    this.form.reset({
      login: '',
      password: '',
      basic: true,
    });
  }
  onSubmit() {
    this.apiService[this.form.value.basic ? 'register' : 'login']({
      login: this.form.value.login || '',
      password: this.form.value.password || '',
    }).subscribe();
  }
  readonly form = new FormGroup({
    login: new FormControl('', Validators.required),
    password: new FormControl('', Validators.required),
    basic: new FormControl(true),
  });
}
