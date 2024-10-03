import React, { useCallback } from "react";
import useWebSocket, { ReadyState } from "react-use-websocket";
import { Cursor } from "./ui/cursor";
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
    `ws://${import.meta.env.SERVER_HOST ?? "localhost"}:9000/ws`,
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
    <div className="w-full flex flex-col justify-between mt-4">
      <section className="w-full flex justify-between">
        <h1>Welcome, {username}</h1>
        <h2>{connectionStatus}</h2>
      </section>
      <section className="w-full flex flex-row-reverse mb-4">
        <ul>
          {Object.entries(otherCursors).map(([username, { x, y }]) => (
            <li key={username}>
              <p>
                {username} x: {x}, y: {y}
              </p>
            </li>
          ))}
        </ul>
      </section>
      <Cursor client onMove={handleMessages} className="z-50" />
      {Object.entries(otherCursors).map(([username, { x, y }], idx) => (
        <Cursor key={username} color={idx as 0} x={x} y={y}>
          <p className="absolute top-0 left-8 text-foreground">{username}</p>
        </Cursor>
      ))}
    </div>
  );
};
