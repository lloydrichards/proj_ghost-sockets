import React, { ReactNode, useMemo } from "react";
import { useMousePosition } from "../hooks/use-mouse-position";
import { useThrottle } from "../hooks/use-trottle";

type CursorProps = {
  onMove?: (x: number, y: number) => void;
  client?: boolean;
  x?: number;
  y?: number;
  children?: ReactNode;
};

const colors = [
  "#1E90FF",
  "#9B59B6",
  "#00CED1",
  "#2ECC71",
  "#E67E22",
  "#E74C3C",
  "#F1C40F",
  "#95A5A6",
  "#D87093",
  "#3CB371",
];

const randomColor = () => {
  const rndIdx = Math.floor(Math.random() * colors.length);
  return colors[rndIdx];
};

export const Cursor: React.FC<CursorProps> = ({
  onMove,
  client = false,
  x,
  y,
  children,
}) => {
  const mousePosition = useMousePosition();
  const throttlePosition = useThrottle(mousePosition, 100);

  React.useEffect(() => {
    if (!client) return;
    onMove?.(throttlePosition.x ?? 0, throttlePosition.y ?? 0);
  }, [throttlePosition, client, onMove]);

  const size = 24;
  const xPos = x ?? throttlePosition.x ?? 0;
  const yPos = y ?? throttlePosition.y ?? 0;
  const color = useMemo(() => randomColor(), []);

  return (
    <div style={{ position: "fixed", top: 0, left: 0, pointerEvents: "none" }}>
      <div
        style={{
          position: "absolute",
          left: xPos - size / 2,
          top: yPos - size / 2,
          width: size,
          height: size,
          borderRadius: 99,
          backgroundColor: color,
        }}
      >
        {children}
      </div>
    </div>
  );
};
