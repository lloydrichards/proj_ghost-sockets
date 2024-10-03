import React, { ReactNode } from "react";
import { useMousePosition } from "../../hooks/use-mouse-position";
import { useThrottle } from "../../hooks/use-throttle";
import { cva, VariantProps } from "class-variance-authority";
import { motion } from "framer-motion";

const cursorVariants = cva("absolute size-6 rounded-full", {
  variants: {
    client: {
      true: "stroke-2 stroke-primary",
      false: "",
    },
    color: {
      default: "fill-transparent",
      0: "fill-[#1E90FF]",
      1: "fill-[#9B59B6]",
      2: "fill-[#00CED1]",
      3: "fill-[#2ECC71]",
      4: "fill-[#E67E22]",
      5: "fill-[#E74C3C]",
      6: "fill-[#F1C40F]",
      7: "fill-[#95A5A6]",
      8: "fill-[#D87093]",
      9: "fill-[#3CB371]",
    },
  },
  defaultVariants: {
    color: "default",
    client: false,
  },
});

type CursorProps = {
  onMove?: (x: number, y: number) => void;
  client?: boolean;
  x?: number;
  y?: number;
  vx?: number;
  vy?: number;
  children?: ReactNode;
  className?: string;
} & VariantProps<typeof cursorVariants>;

export const Cursor: React.FC<CursorProps> = ({
  onMove,
  client = false,
  x,
  y,
  vx,
  vy,
  color,
  children,
}) => {
  const mousePosition = useMousePosition();
  const throttlePosition = useThrottle(mousePosition, 100);

  React.useEffect(() => {
    if (!client) return;
    onMove?.(throttlePosition.x ?? 0, throttlePosition.y ?? 0);
  }, [throttlePosition, client, onMove]);

  const xPos = x ?? throttlePosition.x ?? 0;
  const yPos = y ?? throttlePosition.y ?? 0;

  return (
    <motion.g
      initial={{
        x: xPos,
        y: yPos,
      }}
      animate={{
        x: xPos,
        y: yPos,
      }}
    >
      <motion.line
        initial={{
          x2: (vx ?? 0) * -50,
          y2: (vy ?? 0) * -50,
        }}
        animate={{
          x2: (vx ?? 0) * -50,
          y2: (vy ?? 0) * -50,
        }}
        x1={0}
        y1={0}
        className="stroke-2 stroke-primary"
      ></motion.line>
      <circle r={12} className={cursorVariants({ color, client })}></circle>
      <foreignObject x={16} y={-12} width="100" height="100">
        {children}
      </foreignObject>
    </motion.g>
  );
};
