import { TestBed } from '@angular/core/testing';
import { ApiService } from './api.service';
import {
  HttpClient,
  provideHttpClient,
  withInterceptorsFromDi,
} from '@angular/common/http';
import { of, switchMap, throwError } from 'rxjs';
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
      ],
    });

    service = TestBed.inject(ApiService);
    httpClient = TestBed.inject(HttpClient);

    spyOn(httpClient, 'post').and.returnValue(of(null));
    spyOn(httpClient, 'get').and.returnValue(of(null));
  });

  afterEach(() => {
    TestBed.resetTestingModule();
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  // Login tests
  it('should call login with correct parameters when validation passes', (done) => {
    const testData = { login: 'testtest', password: '123123' };

    spyOn(service, 'login').and.callFake((data: any) => {
      return of(null).pipe(
        switchMap(() => httpClient.post<null>('login', data))
      );
    });

    service.login(testData).subscribe(() => {
      expect(httpClient.post).toHaveBeenCalledWith('login', testData);
      done();
    });
  });

  it('should not call http.post when login validation fails', (done) => {
    const testData = { login: 'short', password: '123' };
    const mockErrors = [{ constraints: { minLength: 'too short' } }];

    spyOn(service, 'login').and.returnValue(throwError(() => mockErrors));

    service.login(testData).subscribe({
      error: () => {
        expect(httpClient.post).not.toHaveBeenCalled();
        done();
      },
    });
  });

  it('should return expected login response', (done) => {
    const testResponse = null;
    const testData = { login: 'testtest', password: '123123' };

    (httpClient.post as jasmine.Spy).and.returnValue(of(testResponse));

    spyOn(service, 'login').and.callFake((data: any) => {
      return of(null).pipe(
        switchMap(() => httpClient.post<null>('login', data))
      );
    });

    service.login(testData).subscribe((response) => {
      expect(response).toEqual(testResponse);
      done();
    });
  });

  // Register tests
  it('should call register with correct parameters when validation passes', (done) => {
    const testData = { login: 'testuser', password: 'password123' };

    spyOn(service, 'register').and.callFake((data: any) => {
      return of(null).pipe(
        switchMap(() => httpClient.post<null>('register', data))
      );
    });

    service.register(testData).subscribe(() => {
      expect(httpClient.post).toHaveBeenCalledWith('register', testData);
      done();
    });
  });

  it('should not call http.post when register validation fails', (done) => {
    const testData = { login: 'usr', password: 'pwd' };
    const mockErrors = [{ constraints: { minLength: 'invalid credentials' } }];

    spyOn(service, 'register').and.returnValue(throwError(() => mockErrors));

    service.register(testData).subscribe({
      error: () => {
        expect(httpClient.post).not.toHaveBeenCalled();
        done();
      },
    });
  });

  it('should return expected register response', (done) => {
    const testResponse = null;
    const testData = { login: 'newuser', password: 'secure123' };

    (httpClient.post as jasmine.Spy).and.returnValue(of(testResponse));

    spyOn(service, 'register').and.callFake((data: any) => {
      return of(null).pipe(
        switchMap(() => httpClient.post<null>('register', data))
      );
    });

    service.register(testData).subscribe((response) => {
      expect(response).toEqual(testResponse);
      done();
    });
  });

  // Logout tests
  it('should call logout endpoint', (done) => {
    service.logout().subscribe(() => {
      expect(httpClient.get).toHaveBeenCalledWith('logout');
      done();
    });
  });

  it('should return expected logout response', (done) => {
    const testResponse = {};
    (httpClient.get as jasmine.Spy).and.returnValue(of(testResponse));

    service.logout().subscribe((response) => {
      expect(response).toEqual(testResponse);
      done();
    });
  });

  // Error handling tests
  it('should propagate validation errors for login', (done) => {
    const testData = { login: 'x', password: 'y' };
    const mockErrors = [
      {
        property: 'login',
        constraints: { minLength: 'Login must be at least 3 characters' },
      },
      {
        property: 'password',
        constraints: { minLength: 'Password must be at least 6 characters' },
      },
    ];

    spyOn(service, 'login').and.returnValue(throwError(() => mockErrors));

    service.login(testData).subscribe({
      error: (errors) => {
        expect(errors).toEqual(mockErrors);
        done();
      },
    });
  });

  it('should propagate validation errors for register', (done) => {
    const testData = { login: 'a', password: 'b' };
    const mockErrors = [
      {
        property: 'login',
        constraints: { minLength: 'Login must be at least 3 characters' },
      },
    ];

    spyOn(service, 'register').and.returnValue(throwError(() => mockErrors));

    service.register(testData).subscribe({
      error: (errors) => {
        expect(errors).toEqual(mockErrors);
        done();
      },
    });
  });
});
