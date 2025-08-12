import { Routes } from '@angular/router';
import { Auth } from './pages/auth/auth.component';
import { TestChatComponent } from './pages/test-chat/test-chat.component';

export const routes: Routes = [
  {
    component: Auth,
    path: 'auth',
  },
  {
    component: TestChatComponent,
    path: 'test-chat',
  },
  {
    redirectTo: 'auth',
    path: '**',
  },
];
