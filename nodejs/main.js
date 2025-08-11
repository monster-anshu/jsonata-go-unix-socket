const net = require("net");
const fs = require("fs");
const jsonata = require("jsonata");

const socketPath = "/tmp/jsonata.sock";

// Remove socket if it exists
if (fs.existsSync(socketPath)) {
  fs.unlinkSync(socketPath);
}

/**
 * @param {net.Socket} socket
 */
function handleConnection(socket) {
  let buffer = "";

  const send = (data) => {
    console.log(data);
    socket.write(JSON.stringify(data) + "\n");
    socket.end();
  };

  socket.on("data", async (chunk) => {
    let transfomationRequest;
    try {
      buffer += chunk.toString();
      transfomationRequest = JSON.parse(buffer);
      buffer = ""; // clear after parsing
    } catch (error) {
      send({
        success: false,
        error: err.message,
        message: "Invalid JSON",
      });
      return;
    }

    const { data, expression } = transfomationRequest;

    /**@type {jsonata.Expression} */
    let expr;
    try {
      expr = jsonata(expression);
    } catch (err) {
      send({
        success: false,
        error: err.message,
        message: "Invalid expression",
      });
      return;
    }

    try {
      const resp = await expr.evaluate(data);
      send({
        success: true,
        data: resp,
        error: null,
      });
    } catch (err) {
      send({
        success: false,
        error: err.message,
        message: "Evaluation error",
      });
    }
  });

  socket.on("error", (err) => {
    console.error("Socket error:", err);
  });
}

const server = net.createServer(handleConnection);

server.listen(socketPath, () => {
  console.log(`JSONata Unix socket server started: ${socketPath}`);
});
