import { WebTransportService } from './../../shared/services/web-transport.service';
import {
  ChangeDetectionStrategy,
  Component,
  inject,
  signal,
} from '@angular/core';
import { TuiButton } from '@taiga-ui/core';
import { WebSocketService } from '../../shared/services/web-socket.service';

@Component({
  standalone: true,
  selector: 'app-chat',
  templateUrl: './chat.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [TuiButton],
})
export class ChatComponent {
  webSocketService = inject(WebSocketService);
  isLoading = signal(false);
  connectionExist = signal(false);
  async connect() {
    try {
      this.isLoading.set(true);
      await this.webSocketService.connect("TODO-provider-ID");
      this.onConnect();
    } catch (err) {
      console.error(err);
    } finally {
      this.isLoading.set(false);
    }
  }
  send() {
    this.webSocketService.send({action:"send", group:"TODO-get-group-id",payload:'Test msg from browser'});
  }
  async disconnect() {
    await this.webSocketService.disconnect();
    this.onDisconnect();
  }
  onConnect() {
    this.connectionExist.set(true);
  }
  onDisconnect() {
    this.connectionExist.set(false);
  }
}
