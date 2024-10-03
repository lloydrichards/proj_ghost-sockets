import React from "react";

type LoginProps = {
  onSubmit: (v: string | null) => void;
};

export const Login: React.FC<LoginProps> = ({ onSubmit }) => {
  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    const username = formData.get("username") as string;
    onSubmit(username);
  };
  return (
    <form onSubmit={handleSubmit}>
      <label htmlFor="username">Username</label>
      <input type="text" id="username" name="username" />
      <button type="submit">Submit</button>
    </form>
  );
};
