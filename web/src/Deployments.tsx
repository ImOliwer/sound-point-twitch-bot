import { useEffect, useState } from 'react';
import { w3cwebsocket as WebSocket } from 'websocket';

interface Deployed {

}

export default function Deployments() {
  const [connected, setConnected] = useState(false);
  const [queue, setQueue] = useState<Deployed[]>([]);

  useEffect(() => {
    const socket = new WebSocket("ws://127.0.0.1:9999");

    socket.onopen = () => setConnected(true);
    
    socket.onmessage = message => {
      console.log(message.data.toString().replace("\n", ""));
    };

    socket.onclose = () => setConnected(false);
  }, []);

  if (!connected) {
    return (
      <p style={{ 
        fontFamily: "Arial, sans-serif", 
        fontSize: "24px",
        color: "#E4E4E4"
      }}>NOT CONNECTED</p>
    );
  }

  return (
    <p>what and what not</p>
  );
}