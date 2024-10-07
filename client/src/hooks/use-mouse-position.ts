import { useEffect, useState } from "react";

export const useMousePosition = () => {
  const [mousePosition, setMousePosition] = useState<{
    x: number | null;
    y: number | null;
  }>({ x: null, y: null });
  useEffect(() => {
    const updateMousePosition = (ev: MouseEvent) => {
      // calculate the position relative to the center of the screen
      const x = ev.clientX - window.innerWidth / 2;
      const y = ev.clientY - window.innerHeight / 2;
      setMousePosition({ x, y });
    };
    window.addEventListener("mousemove", updateMousePosition);
    return () => {
      window.removeEventListener("mousemove", updateMousePosition);
    };
  }, []);
  return mousePosition;
};
