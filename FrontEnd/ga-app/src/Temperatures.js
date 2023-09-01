import React, { useState, useEffect } from "react";

function Temperatures() {
    const [temperatures, setTemperatures] = useState([]);
    const [message, setMessage] = useState({
        Id: 0,
        Ndata: 0,
        Data: [],
    });

    useEffect(() => {
        // Connect to WebSocket
        const socket = new WebSocket("ws://localhost:8080/ws");

        // Listen for messages from the server
        socket.onmessage = (event) => {
            const messageData = JSON.parse(event.data);
            console.log("messageData from WebSocket:", messageData);
            setMessage(messageData);
        };

        // Close the WebSocket connection when unmounting
        return () => {
            socket.close();
        };
    }, []);

    useEffect(() => {
        // Extract the array of objects "Data" from the message
        const data = message.datajson;
        console.log("data from WebSocket:", data);

        // Check if "data" is defined before trying to map over it
        if (Array.isArray(data)) {
            // Extract the floating-point numbers from the "Data" objects and update the temperatures state
            const temperatureData = [];
            for (const nestedJson of data) {
                console.log("nestedJson:", nestedJson); // Agregamos el log aquí
                const floatData = parseFloat(nestedJson.data);
                temperatureData.push(floatData);
            }
            console.log("temperatureData:", temperatureData); // Agregamos el log aquí
            setTemperatures(temperatureData);
        }
    }, [message]);

    return (
        <div>
            <h1>Temperaturas</h1>
            {temperatures.map((temp, index) => (
                <p key={index}>Temperatura {index + 1}: {temp.toFixed(1)} °C</p>
            ))}
        </div>
    );
}

export default Temperatures;
