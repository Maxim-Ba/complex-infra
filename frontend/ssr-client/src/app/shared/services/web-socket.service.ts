import { Injectable } from '@angular/core';
import { MessageDTO } from '../models/message';

type THandler = (msg: MessageDTO) => void;

@Injectable({ providedIn: 'root' })
export class WebSocketService {
  private pid: number = 0;
  private socket!: WebSocket;
  private providerID!: string;
  private handlers: THandler[] = [];

  connect(providerID: string) {
    this.providerID = providerID;
    if (this.socket) {
      this.disconnect();
    }
    this.socket = new WebSocket(`ws://localhost:8089/ws/${providerID}`);
    this.socket.onopen = this.onOpen;
    this.socket.onerror = this.onError;
    this.socket.onmessage = this.onMessage;
    this.socket.onclose = this.onClose;
  }
  ping() {}
  send(messageData: { payload: string; group: string; action: string }) {
    this.socket.send(JSON.stringify(this.createMessage(messageData)));
  }
  disconnect() {
    this.socket.close();
  }
  private createMessage({
    action,
    group,
    payload,
  }: {
    payload: string;
    group: string;
    action: string;
  }): MessageDTO {
    this.pid += 1;
    return {
      action,
      group,
      payload,
      pid: this.pid.toString(),
      producer: this.providerID,
    };
  }
  private onClose = (event: CloseEvent) => {
    console.log('WS connection closed');
  };
  private onMessage = (event: MessageEvent) => {
    console.log({ event });
    const data = JSON.parse(event.data);
    console.log(data);
    console.log(this.handlers);

    this.handlers.forEach((h) => h(data));
  };
  private onError = (event: Event) => {
    console.error({ event });
  };
  private onOpen = (event: Event) => {
    console.log('WS connection open');
  };
  registerHandler(h: THandler) {
    if (this.handlers.includes(h)) {
      return;
    }
    this.handlers.push(h);
  }
  removeHandler(h: THandler) {
    this.handlers = this.handlers.filter((hdl) => h !== hdl);
  }
}
