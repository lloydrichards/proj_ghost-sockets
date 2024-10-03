import React from "react";
import { Label } from "./ui/label";
import { Input } from "./ui/input";
import { Button } from "./ui/button";

type LoginProps = {
  onSubmit: (v: string | null) => void;
};

export const Login: React.FC<LoginProps> = ({ onSubmit }) => {
  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const username = (
      e.currentTarget.elements.namedItem("username") as HTMLInputElement
    ).value;
    onSubmit(username);
  };
  return (
    <form onSubmit={handleSubmit} className="flex items-center gap-4 w-full">
      <Label htmlFor="username">Name</Label>
      <Input type="text" id="username" placeholder="Enter display name..." />
      <Button type="submit">Submit</Button>
    </form>
  );
};
