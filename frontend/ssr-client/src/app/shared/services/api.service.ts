import { inject, Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { plainToClass } from 'class-transformer';
import { Auth } from '../../models/auth/auth';
import { validate } from 'class-validator';
import { from, switchMap, throwError } from 'rxjs';

@Injectable({ providedIn: 'root' })
export class ApiService {
  http = inject(HttpClient);
  login(data: { login: string; password: string }) {
    const authData = plainToClass(Auth, data);
    // TODO Either
    return from(validate(authData)).pipe(
      switchMap((errors) =>
        !!errors.length
          ? throwError(() => errors)
          : this.http.post<null>('login', data)
      )
    );
  }
  register(data: { login: string; password: string }) {
    // TODO Either
    const authData = plainToClass(Auth, data);
    // TODO Either
    return from(validate(authData)).pipe(
      switchMap((errors) =>
        !!errors.length
          ? throwError(() => errors)
          : this.http.post<null>('register', data)
      )
    );
  }
  logout() {
    return this.http.get('logout');
  }
}
