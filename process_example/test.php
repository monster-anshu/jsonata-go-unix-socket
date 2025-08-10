<?php
$binary =  __DIR__ . '/main';

$input = json_encode([
    "json_input" => ["name" => "Himanshu"],
    "jsonata_expr" => '$.name'
]);

$descriptors = [
    0 => ["pipe", "r"],
    1 => ["pipe", "w"],
    2 => ["pipe", "w"]
];

$process = proc_open($binary, $descriptors, $pipes);

if (is_resource($process)) {
    fwrite($pipes[0], $input);
    fclose($pipes[0]);

    $output = stream_get_contents($pipes[1]);
    fclose($pipes[1]);

    $errors = stream_get_contents($pipes[2]);
    fclose($pipes[2]);

    $code = proc_close($process);

    if ($code === 0) {
        echo $output;
        // echo "Go output: $output\n";
    } else {
        echo "Error: $errors\n";
    }
}
