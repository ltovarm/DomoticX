import React, { useState } from "react";
import { useNavigate } from "react-router-dom";

function Register() {
    const navigate = useNavigate();
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");

    const handleRegister = async () => {
        try {
            const response = await fetch("/api/register", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ username, password }),
            });

            if (response.ok) {
                // Registration successful, navigate to the Login page
                navigate("/");
            } else {
                // Handle registration error
                console.error("Registration failed");
            }
        } catch (error) {
            console.error("Error during registration:", error);
        }
    };

    return (
        <div>
            <h1>Register</h1>
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
            <button onClick={handleRegister}>Register</button>
            <p>Already have an account? <a href="/">Login</a></p>
        </div>
    );
}

export default Register;
