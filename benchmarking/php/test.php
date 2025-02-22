<?php
header('Content-Type: application/json');

// Simulate database query delay
usleep(rand(200000, 500000)); // Sleep for 200-500ms to mimic DB latency

// Function to generate a random string
function randomString($length = 100) {
    $characters = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
    $charactersLength = strlen($characters);
    $randomString = '';
    for ($i = 0; $i < $length; $i++) {
        $randomString .= $characters[rand(0, $charactersLength - 1)];
    }
    return $randomString;
}

// Simulate heavy DB processing (complex operations)
function simulateHeavyComputation() {
    $total = 0;
    for ($i = 0; $i < 1000000; $i++) { // Simulate CPU-bound workload
        $total += sqrt($i) * log($i + 1);
    }
    return $total;
}

// Generate an array of 50 items with random data
$data = [];
for ($i = 0; $i < 50; $i++) {
    $data[] = [
        'id' => $i + 1,
        'name' => randomString(20),
        'description' => randomString(200),
        'price' => rand(100, 9999) / 100,
        'stock' => rand(0, 500),
        'category' => randomString(15),
        'timestamp' => date('Y-m-d H:i:s')
    ];
}

// Perform artificial heavy computation
simulateHeavyComputation();

// Encode to JSON and output
echo json_encode(['status' => 'success', 'data' => $data], JSON_PRETTY_PRINT);
?>
