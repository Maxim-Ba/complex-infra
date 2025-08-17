import {
  ChangeDetectionStrategy,
  Component,
  inject,
  signal,
} from '@angular/core';
import { TuiButton } from '@taiga-ui/core';
import { WebSocketService } from '../../shared/services/web-socket.service';
import { WebRTCService } from '../../shared/services/web-rtc.service';
import { MessageDTO } from '../../shared/models/message';

@Component({
  standalone: true,
  selector: 'app-chat',
  templateUrl: './chat.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [TuiButton],
})
export class ChatComponent {
  webSocketService = inject(WebSocketService);
  webRTCService = inject(WebRTCService);
  isLoading = signal(false);
  connectionExist = signal(false);

  messages = signal<string[]>([]);

  async connect() {
    try {
      this.isLoading.set(true);
      await this.webSocketService.connect('TODO-provider-ID');
      this.onConnect();
      this.webSocketService.registerHandler(this.onRTCAnswer);
      this.webSocketService.registerHandler(this.onRTCCandidate);
    } catch (err) {
      console.error(err);
    } finally {
      this.isLoading.set(false);
    }
  }
  send() {
    this.webSocketService.send({
      action: 'message',
      group: 'TODO-get-group-id',
      payload: 'Test msg from browser',
    });
  }
  async webRTCStart() {
    const offer = await this.webRTCService.initPeer();
    this.webSocketService.send({
      action: 'webrtc',
      group: 'TODO-get-group-id',
      payload: JSON.stringify({
        payload: {
          sdp: offer.sdp,
          player_id: 'TODO-provider-ID',
          game_id: 'test-game123',
          session_id: 'test-session456',
        },
        type: 'offer',
      }),
    });
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
  onRTCAnswer = (msg: MessageDTO) => {
    console.log('onRTCAnswer');

    if (msg.action === 'answer') {
      this.webRTCService.setRemoteDescription(JSON.parse(msg.payload).sdp);
    }
  };
  onRTCCandidate = (msg: MessageDTO) => {
    if (msg.action === 'candidate') {
      this.webRTCService.addICECandidate(JSON.parse(msg.payload));
    }
  };
  sendPing() {
    const success = this.webRTCService.sendCommand({
      type: 'ping',
      timestamp: Date.now(),
    });

    if (success) {
      this.messages.set([
        ...this.messages(),
        'Ping sent at ' + new Date().toLocaleTimeString(),
      ]);
    }
  }
}
