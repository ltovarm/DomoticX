import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const navigate = useNavigate();

  // Simulate authentication logic
  const login = () => {
    setIsAuthenticated(true);
    navigate("/temperatures");
  };

  useEffect(() => {
    if (isAuthenticated) {
      // Connect to WebSocket and fetch data
    }
  }, [isAuthenticated]);

  return (
    <div>
      {isAuthenticated ? (
        <h1>Authenticated</h1>
      ) : (
        <div>
          <h1>Login</h1>
          <button onClick={login}>Login</button>
        </div>
      )}
    </div>
  );
}

export default App;
