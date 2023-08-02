// TemperatureService.js
async function fetchTemperatures() {
    try {
        const response = await fetch("http://localhost:8080/temperatures");
        const data = await response.json();
        return data.temperatures;
    } catch (error) {
        console.error("Error al obtener las temperaturas:", error);
        return [];
    }
}

export { fetchTemperatures };
