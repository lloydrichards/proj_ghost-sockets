import React, { useCallback } from "react";
import useWebSocket, { ReadyState } from "react-use-websocket";
import { Cursor } from "./Cursor";
import { z } from "zod";

type DashboardProps = {
  username: string;
};

const Payload = z.object({
  type: z.string(),
  payload: z.record(
    z.string(),
    z.object({
      username: z.string(),
      state: z.object({
        x: z.number(),
        y: z.number(),
      }),
    })
  ),
});

export const Dashboard: React.FC<DashboardProps> = ({ username }) => {
  const [otherCursors, setOtherCursors] = React.useState<
    Record<string, { x: number; y: number }>
  >({});
  const { sendMessage, lastJsonMessage, readyState } = useWebSocket(
    "ws://localhost:9000/ws",
    {
      queryParams: { username },
    }
  );

  React.useEffect(() => {
    if (lastJsonMessage) {
      const parsed = Payload.safeParse(lastJsonMessage);
      if (parsed.success) {
        const { payload } = parsed.data;
        setOtherCursors(
          Object.fromEntries(
            Object.values(payload)
              .map(({ username, state }) => [username, state])
              .filter(([id]) => id !== username)
          )
        );
      }
    }
  }, [lastJsonMessage, username]);

  const connectionStatus = {
    [ReadyState.CONNECTING]: "Connecting",
    [ReadyState.OPEN]: "Open",
    [ReadyState.CLOSING]: "Closing",
    [ReadyState.CLOSED]: "Closed",
    [ReadyState.UNINSTANTIATED]: "Uninstantiated",
  }[readyState];

  const handleMessages = useCallback(
    (x: number, y: number) => {
      sendMessage(
        JSON.stringify({ type: "update_position", payload: { x, y } })
      );
    },
    [sendMessage]
  );

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
      <Cursor client onMove={handleMessages} />
      {Object.entries(otherCursors).map(([username, { x, y }]) => (
        <Cursor key={username} x={x} y={y}>
          <p
            style={{
              position: "absolute",
              top: -16,
              left: 32,
              color: "white",
              fontSize: 16,
            }}
          >
            {username}
          </p>
        </Cursor>
      ))}
    </div>
  );
};
