import { TestBed } from '@angular/core/testing';
import { ApiService } from './api.service';
import {
  HttpClient,
  provideHttpClient,
  withInterceptorsFromDi,
} from '@angular/common/http';
import { of } from 'rxjs';
import { provideZonelessChangeDetection } from '@angular/core';
import { provideHttpClientTesting } from '@angular/common/http/testing';

describe('ApiService', () => {
  let service: ApiService;
  let httpClient: HttpClient;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        ApiService,
        provideZonelessChangeDetection(),
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
      ]
    });

    service = TestBed.inject(ApiService);
    httpClient = TestBed.inject(HttpClient);

    spyOn(httpClient, 'post').and.callFake((url: string, body: any) =>
      of({ success: true } as any)
    );
    spyOn(httpClient, 'get').and.callFake((url: string) =>
      of({ success: true } as any)
    );
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should call login with correct parameters', () => {
    const testData = { login: 'test', password: '123' };

    service.login(testData).subscribe();

    expect(httpClient.post).toHaveBeenCalledWith('login', testData);
  });

  it('should return expected login response', () => {
    const testResponse = { token: 'abc123' };
    (httpClient.post as jasmine.Spy).and.returnValue(of(testResponse));

    service.login({ login: 'test', password: '123' }).subscribe(response => {
      expect(response).toEqual(testResponse);
    });
  });
});
