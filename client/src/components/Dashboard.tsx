import React, { useCallback, useMemo } from "react";
import useWebSocket, { ReadyState } from "react-use-websocket";
import { Cursor } from "./ui/cursor";
import { z } from "zod";

type DashboardProps = {
  username: string;
};

const User = z.object({
  username: z.string(),
  color: z.number(),
  mood: z.string(),
  state: z.object({
    x: z.number(),
    y: z.number(),
    vx: z.number(),
    vy: z.number(),
    spd: z.number(),
    acc: z.number(),
    ang: z.number(),
  }),
});

const Payload = z.object({
  type: z.string(),
  payload: z.record(z.string(), User),
});

export const Dashboard: React.FC<DashboardProps> = ({ username }) => {
  const [players, setPlayers] = React.useState<
    Array<z.infer<typeof User> & { id: string }>
  >([]);

  const { sendMessage, lastJsonMessage, readyState } = useWebSocket(
    `ws://${import.meta.env.VITE_SERVER_HOST ?? "localhost"}:9000/ws`,
    {
      queryParams: { username },
    }
  );

  React.useEffect(() => {
    if (lastJsonMessage) {
      const parsed = Payload.safeParse(lastJsonMessage);
      if (parsed.success) {
        const { payload } = parsed.data;
        setPlayers(
          Object.entries(payload).map(([id, data]) => ({ id, ...data }))
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
        JSON.stringify({
          type: "update_position",
          payload: { x, y, delta: 100 },
        })
      );
    },
    [sendMessage]
  );

  const currentPlayer = useMemo(
    () => players.find((d) => d.username === username),
    [players, username]
  );
  return (
    <div className="w-full h-[100dvh] flex flex-col justify-between pt-4">
      <section className="w-full flex justify-between">
        <h1>Welcome, {username}</h1>
        <h2>{connectionStatus}</h2>
      </section>
      <section className="w-full flex flex-row-reverse pb-4">
        <table className="w-full table-fixed text-left">
          <thead>
            <tr>
              <th>Name</th>
              {Object.keys(players[0]?.state || {}).map((key) => (
                <th className="text-right" key={key}>
                  {key}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {players.map(({ username, state }) => (
              <tr key={username}>
                <td>{username}</td>
                {Object.entries(state).map(([key, value]) => (
                  <td className="text-right" key={key}>
                    {value.toFixed(2)}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </section>
      <svg className="absolute left-0 top-0 size-full overflow-visible">
        <g className="translate-x-1/2 translate-y-1/2">
          <Cursor client onMove={handleMessages} className="z-50" />
          {players
            .filter((d) => d.username != username)
            .map((d) => {
              // check if the cursor is near the current player
              const isNear =
                currentPlayer &&
                Math.hypot(
                  d.state.x - currentPlayer.state.x,
                  d.state.y - currentPlayer.state.y
                ) < 100;
              return (
                <Cursor
                  key={d.id}
                  color={d.color as 0}
                  mood={isNear ? d.mood : undefined}
                  {...d.state}
                >
                  <p>{d.username}</p>
                </Cursor>
              );
            })}
        </g>
      </svg>
    </div>
  );
};
