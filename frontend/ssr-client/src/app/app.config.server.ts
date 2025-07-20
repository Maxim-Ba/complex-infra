import { mergeApplicationConfig, ApplicationConfig } from '@angular/core';
import { provideServerRendering, withRoutes } from '@angular/ssr';
import { appConfig } from './app.config';
import { serverRoutes } from './app.routes.server';
import { UNIVERSAL_PROVIDERS } from '@ng-web-apis/universal';
import { provideHttpClient } from '@angular/common/http';

const serverConfig: ApplicationConfig = {
  providers: [
    provideServerRendering(withRoutes(serverRoutes)),
    UNIVERSAL_PROVIDERS,
    provideHttpClient(),
  ],
};

export const config = mergeApplicationConfig(appConfig, serverConfig);
