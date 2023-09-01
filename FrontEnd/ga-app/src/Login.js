import React, { useState } from "react";
import { useNavigate } from "react-router-dom";

function Login() {
    const navigate = useNavigate();
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");

    const handleLogin = async () => {
        try {
            const response = await fetch("/api/login", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ username, password }),
            });

            if (response.ok) {
                // Login successful, navigate to Temperatures page
                navigate("/temperatures");
            } else {
                // Handle login error
                console.log(response)
                console.error("Login failed");
            }
        } catch (error) {
            console.error("Error during login:", error);
        }
    };

    return (
        <div>
            <h1>Login</h1>
            <div>
                <label htmlFor="username">Username:</label>
                <input
                    type="text"
                    id="username"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                />
            </div>
            <div>
                <label htmlFor="password">Password:</label>
                <input
                    type="password"
                    id="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                />
            </div>
            <button onClick={handleLogin}>Login</button>
            <p>Don't have an account? <a href="/register">Register</a></p>
        </div>
    );
}

export default Login;
