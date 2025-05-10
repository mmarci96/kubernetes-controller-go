import express from "express";
import http from "http";
import { Server } from "socket.io";
import { requestLogger } from "./middleware/index.js";

const app = express();
app.use(requestLogger);

app.get("/ping", (req, res) => res.status(200).send("PONG"));

app.use(express.json());

const server = http.createServer(app);
const io = new Server(server);

const websocketController = (io) => {
    io.on("connection", (socket) => {
        console.log("Connected: ", socket.id);
        socket.on("test", (data) => {
            console.log("Test socket data: ", data);
        });
        socket.on("disconnect", () => {
            console.log("Disconnected: ", socket.id);
        });
    });
};

websocketController(io);

const main = async () => {
    try {
        server.listen(8080, "0.0.0.0", () => {
            console.log("Server listening on 8080 port");
        });
    } catch (err) {
        console.error(err);
    }
};
main();
