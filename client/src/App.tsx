import { useState } from "react";
import { Login } from "./components/Login";
import { Dashboard } from "./components/Dashboard";

function App() {
  const [username, setUsername] = useState<string | null>(null);

  return (
    <main className="h-screen flex justify-center">
      <section className="flex size-full max-w-screen-sm justify-center">
        {!username ? (
          <Login onSubmit={setUsername} />
        ) : (
          <Dashboard username={username} />
        )}
      </section>
    </main>
  );
}

export default App;
