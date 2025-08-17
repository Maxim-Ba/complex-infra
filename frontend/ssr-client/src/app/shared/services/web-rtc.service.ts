import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

@Injectable({ providedIn: 'root' })
export class WebRTCService {
  peerConnection!: RTCPeerConnection;
  dataChannel!: RTCDataChannel;
  connectionState$ = new BehaviorSubject<string>('disconnected');
  iceGatheringState$ = new BehaviorSubject<string>('new');

  async initPeer(): Promise<RTCSessionDescriptionInit> {
    const peerConnection = new RTCPeerConnection({
      iceServers: [
        { urls: 'stun:stun1.l.google.com:19302' },
        { urls: 'stun:stun2.l.google.com:19302' },
      ],
    });
    this.monitorConnectionStates();

    this.peerConnection = peerConnection;
    peerConnection.onconnectionstatechange = () => {
      this.connectionState$.next(peerConnection.connectionState);
      console.log('Connection state:', peerConnection.connectionState);
    };

    peerConnection.onicegatheringstatechange = () => {
      this.iceGatheringState$.next(peerConnection.iceGatheringState);
      console.log('ICE gathering state:', peerConnection.iceGatheringState);
    };

    peerConnection.onicecandidate = (event) => {
      if (event.candidate) {
        console.log('New ICE candidate:', event.candidate);
      } else {
        console.log('ICE gathering complete');
      }
    };

    // Создаем канал для игровых данных
    this.dataChannel = peerConnection.createDataChannel('gameData', {
      ordered: true, // Гарантирует порядок доставки сообщений
    });

    this.dataChannel.onopen = () => {
      console.log('Data channel opened!');
      this.connectionState$.next('connected');
    };

    this.dataChannel.onclose = () => {
      console.log('Data channel closed!');
      this.connectionState$.next('disconnected');
    };

    this.dataChannel.onmessage = (e) => {
      console.log('Message received:', e.data);
    };
    const offer = await peerConnection.createOffer();
    await peerConnection.setLocalDescription(offer);
    return offer;
  }
  async setRemoteDescription(sdp: string) {
    if (!this.peerConnection) {
      throw new Error('PeerConnection not initialized');
    }
    try {
      console.log('Setting remote description with SDP:', sdp);
      await this.peerConnection.setRemoteDescription(
        new RTCSessionDescription({ type: 'answer', sdp })
      );
      console.log('Remote description set successfully');
      console.log(
        'Current signaling state:',
        this.peerConnection.signalingState
      );
    } catch (error) {
      console.error('Failed to set remote description:', error);
      throw error;
    }
  }
  async addICECandidate(candidate: RTCIceCandidateInit) {
    return (
      this.peerConnection &&
      (await this.peerConnection.addIceCandidate(
        new RTCIceCandidate(candidate)
      ))
    );
  }

  sendCommand(command: any): boolean {
    console.log(this.dataChannel);

    if (!this.dataChannel || this.dataChannel.readyState !== 'open') {
      console.error('Data channel not ready');
      return false;
    }

    try {
      const data =
        typeof command === 'string' ? command : JSON.stringify(command);
      this.dataChannel.send(data);
      console.log('Command sent:', data);
      return true;
    } catch (error) {
      console.error('Error sending command:', error);
      return false;
    }
  }
  monitorConnectionStates() {
    if (!this.peerConnection) return;

    this.peerConnection.oniceconnectionstatechange = () => {
      console.log(
        'ICE connection state:',
        this.peerConnection.iceConnectionState
      );
    };

    this.peerConnection.onsignalingstatechange = () => {
      console.log('Signaling state:', this.peerConnection.signalingState);
    };

    this.peerConnection.onconnectionstatechange = () => {
      console.log('Connection state:', this.peerConnection.connectionState);
    };
  }
}
