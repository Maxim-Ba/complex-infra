import { Injectable } from '@angular/core';
import {
  HttpRequest,
  HttpHandler,
  HttpEvent,
  HttpInterceptor
} from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment.development';

@Injectable()
export class BaseUrlInterceptor implements HttpInterceptor {
  private baseUrl = environment.baseUrl;

  intercept(request: HttpRequest<unknown>, next: HttpHandler): Observable<HttpEvent<unknown>> {
    // Исключаем запросы, которые уже имеют абсолютный URL
    if (!request.url.startsWith('http')) {
      const apiReq = request.clone({ url: `${this.baseUrl}/${request.url}`, withCredentials:true });
      return next.handle(apiReq);
    }
    return next.handle(request);
  }
}
