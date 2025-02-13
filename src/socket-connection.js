const wsUrl =  process.env.REACT_APP_WS_URL

class SocketConnection {
  constructor() {
    this.socket = new WebSocket(wsUrl);
    this.reconnectInterval = 5000; // Retry connection every 5 seconds
    this.username = null;
  }

  connect = cb => {
    console.log('Connecting to WebSocket:', this.socket.url);

    this.socket.onopen = () => {
      console.log('Successfully Connected!');
      if (this.username) {
        // If the username exists, send the bootup message immediately
        this.sendMsg({ type: 'bootup', user: this.username });
      }
    };

    this.socket.onmessage = msg => {
      cb(JSON.parse(msg.data)); // Parse the incoming message
    };

    this.socket.onclose = event => {
      console.log('Socket Closed Connection: ', event);
      this.reconnect(); // Attempt reconnect
    };

    this.socket.onerror = error => {
      console.log('Socket Error: ', error);
    };
  };

  sendMsg = msg => {
    if (this.socket.readyState === WebSocket.OPEN) {
      console.log('Sending message:', msg);
      this.socket.send(JSON.stringify(msg));
    } else {
      console.log('WebSocket is not open. Message not sent.');
    }
  };

  connected = user => {
    this.username = user;  // Store username for later re-registration
    this.socket.onopen = () => {
      console.log('Successfully Connected', user);
      this.mapConnection(user);
    };
  };

  mapConnection = user => {
    console.log('Mapping user:', user);
    this.sendMsg({ type: 'bootup', user: user });
  };

  reconnect = () => {
    console.log(`Attempting to reconnect WebSocket..., ${this.username}`);
    setTimeout(() => {
      this.socket = new WebSocket(wsUrl);
      this.connect(msg => console.log('Reconnected:', msg));
      if (this.username) {
        // Resend the bootup message after reconnection
        this.sendMsg({ type: 'bootup', user: this.username });
      }
    }, this.reconnectInterval);
  };

  sendAckMessage = (username) => {
    this.username = username;
    if (this.socket.readyState === WebSocket.OPEN) {
      const ackMessage = { type: 'bootup', user: username };
      console.log('Sending acknowledgment message:', ackMessage);
      this.socket.send(JSON.stringify(ackMessage));
  
      // Now send the bootup message
      const bootupMessage = { type: 'bootup', user: username };
      console.log('Sending bootup message:', bootupMessage);
      this.socket.send(JSON.stringify(bootupMessage));
    } else {
      console.log('WebSocket is not open. Acknowledgment and bootup messages not sent.');
    }
  };
}

export default SocketConnection;
