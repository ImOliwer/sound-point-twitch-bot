import { useEffect, useState } from 'react';
import { w3cwebsocket as WebSocket } from 'websocket';
import { Deployed } from './util/shared';
import { TitleDeploy } from './util/TitleDeploy';

type UserState = {
  [key: string]: any;
};

type Deployment = Deployed & {
  state?: UserState;
  tester?: string;
};

export default function Deployments() {
  const [connected, setConnected] = useState(false);

  let isPlaying = false;
  let queue: Deployment[] = [];

  const handleNext = () => {
    if (queue.length == 0 || isPlaying) {
      return;
    }

    isPlaying = true;
    const next = queue.shift();
    
    if (next === undefined) {
      return;
    }

    const audio = new Audio(`sounds/${next.file_name}`);
    audio.load();
    audio.loop = false;
    audio.onended = () => { audio.remove(); isPlaying = false };
    audio.play();
  };

  useEffect(() => {
    const socket = new WebSocket("ws://127.0.0.1:9999/sound/deployment");

    socket.onopen = () => setConnected(true);
    
    socket.onmessage = message => {
      const obj = JSON.parse(message.data.toString());
      
      if (obj === undefined || !obj.file_name) {
        return;
      }
      
      queue.push(obj);
    };

    socket.onclose = () => setConnected(false);

    setInterval(handleNext, 500); // attempt to handle every 500ms
  }, []);

  return (
    <TitleDeploy title="Deployments">
      {connected ? (
        <>
        </>
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