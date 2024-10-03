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
  children?: ReactNode;
  className?: string;
} & VariantProps<typeof cursorVariants>;

export const Cursor: React.FC<CursorProps> = ({
  onMove,
  client = false,
  x,
  y,
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
    <>
      <motion.circle
        initial={{
          cx: xPos - 24,
          cy: yPos - 16,
        }}
        animate={{
          cx: xPos - 24,
          cy: yPos - 16,
        }}
        r={12}
        className={cursorVariants({ color, client })}
      ></motion.circle>
      <motion.foreignObject
        initial={{
          x: xPos - 8,
          y: yPos - 28,
        }}
        animate={{
          x: xPos - 8,
          y: yPos - 28,
        }}
        width="100"
        height="100"
      >
        {children}
      </motion.foreignObject>
    </>
  );
};
