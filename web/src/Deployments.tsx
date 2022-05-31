import { useEffect, useState } from 'react';
import { w3cwebsocket as WebSocket } from 'websocket';
import { TitleDeploy } from './util/TitleDeploy';

type Deployed = {

}

export default function Deployments() {
  const [connected, setConnected] = useState(false);
  const [queue, setQueue] = useState<Deployed[]>([]);

  useEffect(() => {
    const socket = new WebSocket("ws://127.0.0.1:9999/sound/deployment");

    socket.onopen = () => setConnected(true);
    
    socket.onmessage = message => {
      console.log(message.data.toString().replace("\n", ""));
    };

    socket.onclose = () => setConnected(false);
  }, []);

  return (
    <TitleDeploy title="Deployments">
      {connected ? (
        <p>what and what not</p>
      ) : (
        <p style={{ 
          fontFamily: "Arial, sans-serif", 
          fontSize: "24px",
          color: "#E4E4E4"
        }}>NOT CONNECTED</p>
      )}
    </TitleDeploy>
  );
}