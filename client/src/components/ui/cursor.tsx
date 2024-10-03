import React, { ReactNode } from "react";
import { useMousePosition } from "../../hooks/use-mouse-position";
import { useThrottle } from "../../hooks/use-throttle";
import { cva, VariantProps } from "class-variance-authority";
import { cn } from "@/lib/utils";

const cursorVariants = cva("absolute size-6 rounded-full", {
  variants: {
    client: {
      true: "border-2 border-primary",
      false: "",
    },
    color: {
      default: "bg-transparent",
      0: "bg-[#1E90FF]",
      1: "bg-[#9B59B6]",
      2: "bg-[#00CED1]",
      3: "bg-[#2ECC71]",
      4: "bg-[#E67E22]",
      5: "bg-[#E74C3C]",
      6: "bg-[#F1C40F]",
      7: "bg-[#95A5A6]",
      8: "bg-[#D87093]",
      9: "bg-[#3CB371]",
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
  className,
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
    <div className={cn("fixed top-0 left-0 pointer-events-none", className)}>
      <div
        className={cursorVariants({ color, client })}
        style={{
          left: xPos - 12,
          top: yPos - 12,
        }}
      >
        {children}
      </div>
    </div>
  );
};
