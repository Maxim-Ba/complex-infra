import { ComponentFixture, TestBed } from '@angular/core/testing';
import { Auth } from './auth.component';
import { ReactiveFormsModule } from '@angular/forms';
import { ApiService } from '../../shared/services/api.service';
import { of } from 'rxjs';
import { provideHttpClientTesting } from '@angular/common/http/testing';
import { TuiRoot } from '@taiga-ui/core';
import {
  provideHttpClient,
  withInterceptorsFromDi,
} from '@angular/common/http';
import { provideZonelessChangeDetection } from '@angular/core';

describe('AuthComponent', () => {
  let component: Auth;
  let fixture: ComponentFixture<Auth>;
  let apiService: ApiService;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ReactiveFormsModule, TuiRoot, Auth],
      providers: [
        provideZonelessChangeDetection(),
        ApiService,
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(Auth);
    component = fixture.componentInstance;
    apiService = TestBed.inject(ApiService);
    fixture.detectChanges();
  });

  it('should create', async() => {
    await fixture.whenStable()
    expect(component).toBeTruthy();
  });

  it('should initialize form with default values', () => {
    expect(component.form.value).toEqual({
      login: '',
      password: '',
      basic: true,
    });
  });

  it('should require login and password fields', () => {
    const loginControl = component.form.get('login');
    const passwordControl = component.form.get('password');

    expect(loginControl?.hasError('required')).toBeTrue();
    expect(passwordControl?.hasError('required')).toBeTrue();

    loginControl?.setValue('test');
    passwordControl?.setValue('password');

    expect(loginControl?.hasError('required')).toBeFalse();
    expect(passwordControl?.hasError('required')).toBeFalse();
  });

  it('should call register when basic is true and form is submitted', () => {
    spyOn(apiService, 'register').and.returnValue(of(null));
    component.form.patchValue({
      login: 'testuser',
      password: 'testpass',
      basic: true,
    });

    component.onSubmit();

    expect(apiService.register).toHaveBeenCalledWith({
      login: 'testuser',
      password: 'testpass',
    });
  });

  it('should call login when basic is false and form is submitted', () => {
    spyOn(apiService, 'login').and.returnValue(of(null));
    component.form.patchValue({
      login: 'testuser',
      password: 'testpass',
      basic: false,
    });

    component.onSubmit();

    expect(apiService.login).toHaveBeenCalledWith({
      login: 'testuser',
      password: 'testpass',
    });
  });

  it('should reset form on cancel', async() => {
    component.form.patchValue({
      login: 'test',
      password: 'test',
      basic: false,
    });

    component.onCancel();
    await fixture.whenStable()

    expect(component.form.value).toEqual({
      login: "",
      password: "",
      basic: true,
    });
  });

  it('should call logout when logout method is called', async () => {
    const logoutSpy = spyOn(apiService, 'logout').and.returnValue(of({}));

    component.logout();


    expect(logoutSpy).toHaveBeenCalled();
  });
});
