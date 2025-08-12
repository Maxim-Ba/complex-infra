import { Injectable } from '@angular/core';

@Injectable({providedIn:"root"})
export class WebTransportService {
  url = 'https://localhost:8087/webtransport';
  transport!: WebTransport | null;
  writer!: WritableStreamDefaultWriter<any> | null;
  pingInterval!: NodeJS.Timeout;
  async connect() {
    try {
      this.transport = new WebTransport(this.url);
      await this.transport.ready;
      console.log('WebTransport connection established');

      const stream = await this.transport.createBidirectionalStream();
      this.writer = stream.writable.getWriter();

      // Запускаем чтение сообщений без блокировки event loop
      this.readMessages(stream.readable.getReader());

      // Периодическая отправка сообщений
      this.pingInterval = setInterval(() => this.send('Ping'), 3000);
    } catch (e) {
      console.error('Connection error:', e);
    }
  }

  async readMessages(reader: ReadableStreamDefaultReader<any>) {
    try {
      while (true) {
        const { value, done } = await reader.read();
        if (done) break;

        // Обработка сообщения без блокировки UI
        this.handleMessage(new TextDecoder().decode(value));
      }
    } catch (e) {
      console.error('Read error:', e);
    }
  }

  handleMessage(msg: string) {
    console.log({ msg });
  }

  async send(message: string) {
    if (!this.writer) return;

    try {
      await this.writer.write(new TextEncoder().encode(message));
    } catch (e) {
      console.error('Send error:', e);
    }
  }

  async disconnect() {
    clearInterval(this.pingInterval);
    if (this.writer) {
      this.writer.close();
      this.writer = null;
    }
    if (this.transport) {
      this.transport.close();
      this.transport = null;
    }
  }
}
