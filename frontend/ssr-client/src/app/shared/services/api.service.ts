import { inject, Injectable } from '@angular/core';
import { environment } from '../../../environments/environment.development';
import { HttpClient } from '@angular/common/http';

@Injectable({ providedIn: 'root' })
export class ApiService {
  http = inject(HttpClient);
  login(data: { login: string; password: string }) {
    return this.http.post('login', data);
  }
  register(data: { login: string; password: string }) {
    return this.http.post('register', data);
  }
  logout() {
    return this.http.get('logout');
  }
}
