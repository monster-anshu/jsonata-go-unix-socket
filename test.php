<?php

class JsonataTransformer
{
    private $socketPath;

    public function __construct(string $socketPath = "/tmp/jsonata.sock")
    {
        $this->socketPath = $socketPath;
    }

    /**
     * @return array{
     *     success: bool,
     *     data: array,
     *     error: mixed,
     *     message: string
     * }
     */
    public function transform($data, string $expression)
    {
        $socket = socket_create(AF_UNIX, SOCK_STREAM, 0);

        if (!$socket) {
            throw new \RuntimeException("Failed to create socket");
        }

        if (!socket_connect($socket, $this->socketPath)) {
            throw new \RuntimeException("Failed to connect to JSONata service");
        }

        $request = json_encode([
            "data" => $data,
            "expression" => $expression,
        ]);

        socket_write($socket, $request . "\n");

        $response = "";
        while ($buf = socket_read($socket, 1024)) {
            $response .= $buf;
        }

        socket_close($socket);

        $result = json_decode($response, true);

        return $result;
    }
}

$jsonInput = json_decode(file_get_contents("input.json"), true);
$jsonataExpr = file_get_contents("input_expr.txt");

$transformer = new JsonataTransformer();
$result = $transformer->transform($jsonInput, $jsonataExpr);

print_r(json_encode($result, JSON_PRETTY_PRINT));
