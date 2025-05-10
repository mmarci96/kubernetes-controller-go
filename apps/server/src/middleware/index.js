export function requestLogger(req, res, next) {
    const start = Date.now();
    res.on("finish", () => {
        const duration = Date.now() - start;
        const time = new Date().toLocaleTimeString("hu-HU");
        console.log(
            `[${time}] - [${req.method} -> ${req.originalUrl} - ${duration}ms | Status:${res.statusCode}] - IP: ${req.ip}`,
        );
    });
    next();
}
