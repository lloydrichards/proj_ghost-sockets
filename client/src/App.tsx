import { useState } from "react";
import { Login } from "./components/Login";
import { Dashboard } from "./components/Dashboard";

function App() {
  const [username, setUsername] = useState<string | null>(null);

  return (
    <main className="h-screen">
      {!username ? (
        <>
          <Login onSubmit={setUsername} />
        </>
      ) : (
        <Dashboard username={username} />
      )}
    </main>
  );
}

export default App;
