import { io } from "socket.io-client";

let socket;

document.getElementById("connectBtn").addEventListener("click", async () => {
    const playerId = document.getElementById("playerId").value.trim();
    const gameId = document.getElementById("gameId").value.trim();

    if (!playerId || !gameId) {
        alert("Please enter both Player ID and Game ID.");
        return;
    }

    try {
        socket = io("/", {
            path: "/socket.io",
            query: { gameId, playerId },
        });
        socket.on("connect", () => {
            console.log("Connected with ID:", socket.id);
            document.getElementById("status").innerText =
                `Connected as ${socket.id}`;

            socket.emit("test", { playerId, gameId });
            console.log("Emitted test event with:", { playerId, gameId });
        });

        socket.on("disconnect", () => {
            console.log("Disconnected");
            document.getElementById("status").innerText = "Disconnected";
        });
    } catch (error) {
        console.error(error);
    }
});
