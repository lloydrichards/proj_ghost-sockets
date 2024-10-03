import React from "react";
import useWebSocket, { ReadyState } from "react-use-websocket";
import { Cursor } from "./Cursor";

type DashboardProps = {
  username: string;
};

export const Dashboard: React.FC<DashboardProps> = ({ username }) => {
  const { sendMessage, lastJsonMessage, readyState } = useWebSocket(
    "ws://localhost:9000/ws",
    {
      queryParams: { username },
    }
  );

  const connectionStatus = {
    [ReadyState.CONNECTING]: "Connecting",
    [ReadyState.OPEN]: "Open",
    [ReadyState.CLOSING]: "Closing",
    [ReadyState.CLOSED]: "Closed",
    [ReadyState.UNINSTANTIATED]: "Uninstantiated",
  }[readyState];

  return (
    <div>
      <div>
        <pre>
          {lastJsonMessage
            ? JSON.stringify(lastJsonMessage, null, 2)
            : "No messages yet"}
        </pre>
      </div>
      Welcome, {username}
      <p>{connectionStatus}</p>
      <Cursor client onMove={(x, y) => sendMessage(JSON.stringify({ x, y }))} />
    </div>
  );
};
