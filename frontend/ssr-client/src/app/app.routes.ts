import { Routes } from '@angular/router';
import { Auth } from './pages/auth/auth.component';

export const routes: Routes = [
  {
    component: Auth,
    path: 'auth',
  },
  {
    redirectTo: 'auth',
    path: '**',
  },
];
