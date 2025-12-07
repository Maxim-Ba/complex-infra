import { ChangeDetectionStrategy, Component } from '@angular/core';

@Component({
  selector: 'app-splitter',
  imports: [],
  template: `<div class="card-shadow">
    <ng-content select="card-title"></ng-content>
    <div class="card-divider"></div>
    <ng-content select="card-body"></ng-content>
  </div>`,
  styleUrl: './splitter.css',
  changeDetection: ChangeDetectionStrategy.OnPush,
  standalone: true,
})
export class Splitter {}
