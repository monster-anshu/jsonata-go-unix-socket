<?php

function jsonataEvaluate($jsonInput, $expression) {
    $socketPath = '/tmp/jsonata.sock';
    
    $client = stream_socket_client("unix://$socketPath", $errno, $errstr, 2);

    if (!$client) {
        throw new Exception("Socket connection failed: $errstr ($errno)");
    }

    $payload = [
        "json_input" => $jsonInput,
        "jsonata_expr" => $expression
    ];

    fwrite($client, json_encode($payload) . "\n");

    $response = stream_get_contents($client);
    fclose($client);

    return json_decode($response, true);
}


$jsonInput = json_decode( file_get_contents("input.json"), true);
$jsonataExpr = file_get_contents("input_expr.txt");

$result = jsonataEvaluate(
    $jsonInput,
    $jsonataExpr
);

print_r($result);
