import Axios, { AxiosResponse } from 'axios';
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

// TODO: set up the layout box for deployments
export default function Deployments() {
  const [connected, setConnected] = useState(false);

  let alertContent = "";
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

    const alertContainer = document.getElementById("alertContainer");
    if (!alertContainer) {
      isPlaying = false;
      return;
    }
    
    const child = document.createElement("div");
    child.id = "alert-main";
    child.innerHTML = alertContent
      .replaceAll("{{deployer}}", next.state ? next.state["display-name"] : next.tester)
      .replaceAll("{{id}}", next.file_name); // TODO: fetch name instead of file name for this

    window["deploymentStart"](alertContainer, child).then(() => {
      const audio = new Audio(`sounds/${next.file_name}`);
      audio.load();
      audio.loop = false;
      audio.onended = () => { 
        audio.remove();
        window["deploymentEnd"](alertContainer, child).then(() => {
          child.remove();
          isPlaying = false;
        });
      };
      audio.play();
    });
  };

  useEffect(() => {
    Axios
      .get("/struct/alert.html")
      .catch(console.log)
      .then(it => {
        const res = it as AxiosResponse<any, any>;
        alertContent = res.data.toString();
      });

    if (document) {
      const alertCss = document.createElement("link");
      alertCss.rel = "stylesheet";
      alertCss.href = "/struct/alert.css";

      const alertScript = document.createElement("script");
      alertScript.type = "text/javascript";
      alertScript.src = "/struct/alert.js";

      document.head.appendChild(alertCss);
      document.body.appendChild(alertScript);
    }

    const socket = new WebSocket("ws://127.0.0.1:9999/sound/deployment");
    socket.onopen = () => setConnected(true);
    socket.onclose = () => setConnected(false);
    socket.onmessage = message => {
      const obj = JSON.parse(message.data.toString());
      
      if (obj === undefined || !obj.file_name) {
        return;
      }
      
      queue.push(obj);
    };

    setInterval(handleNext, 500); // attempt to handle every 500ms
  }, []);

  return (
    <TitleDeploy title="Deployments">
      {connected ? (
        <div 
          id="alertContainer" 
          style={{
            display: "flex", 
            width: "100%", 
            height: "100vh"
          }}
        />
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