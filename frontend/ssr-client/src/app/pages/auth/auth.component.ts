import { AsyncPipe, JsonPipe } from '@angular/common';
import { ChangeDetectorRef, Component, inject } from '@angular/core';
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
import {
  TUI_VALIDATION_ERRORS,
  TuiFieldErrorPipe,
  TuiSegmented,
} from '@taiga-ui/kit';
import { TuiCardLarge, TuiForm, TuiHeader } from '@taiga-ui/layout';
import { ApiService } from '../../shared/services/api.service';
import { catchError, of } from 'rxjs';
import { errorsToFormErrors } from '../../shared/utils/errors/errorsToFormErrors';
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
    TuiSegmented,
    TuiTextfield,
    TuiTitle,
  ],
  providers: [
    {
      provide: TUI_VALIDATION_ERRORS,
      useValue: {
        serviceError: (error: string) => error,
      },
    },
  ],
})
export class Auth {
  private cdr = inject(ChangeDetectorRef);
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
    })
      .pipe(
        catchError((err) => {
          this.setFieldErrors(errorsToFormErrors(err));
          return of(err);
        })
      )
      .subscribe();
  }
  private setFieldErrors(errors: Record<string, string>): void {
    Object.keys(errors).forEach((key) => {
      const control = this.form.get(key);
      if (control) {
        control.setErrors({ serviceError: errors[key] });
        control.markAsTouched();
      }
    });
  }
  readonly form = new FormGroup({
    login: new FormControl('', Validators.required),
    password: new FormControl('', Validators.required),
    basic: new FormControl(true),
  });
}
