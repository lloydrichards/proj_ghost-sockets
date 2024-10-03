import { useState } from "react";
import "./App.css";
import { Login } from "./components/Login";
import { Dashboard } from "./components/Dashboard";

function App() {
  const [username, setUsername] = useState<string | null>(null);

  return !username ? (
    <>
      <Login onSubmit={setUsername} />
    </>
  ) : (
    <Dashboard username={username} />
  );
}

export default App;
