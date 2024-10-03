import React from "react";
import { useMousePosition } from "../hooks/use-mouse-position";
import { useThrottle } from "../hooks/use-trottle";

type CursorProps = {
  onMove: (x: number, y: number) => void;
  client?: boolean;
  x?: number;
  y?: number;
};

export const Cursor: React.FC<CursorProps> = ({
  onMove,
  client = false,
  x,
  y,
}) => {
  const mousePosition = useMousePosition();
  const throttlePosition = useThrottle(mousePosition, 500);

  React.useEffect(() => {
    if (!client) return;
    onMove(throttlePosition.x ?? 0, throttlePosition.y ?? 0);
  }, [throttlePosition, client, onMove]);

  return (
    <div style={{ position: "fixed", top: 0, left: 0, pointerEvents: "none" }}>
      <div
        style={{
          position: "absolute",
          left: x ?? mousePosition.x ?? 0,
          top: y ?? mousePosition.y ?? 0,
          width: 10,
          height: 10,
          backgroundColor: "red",
        }}
      />
    </div>
  );
};
