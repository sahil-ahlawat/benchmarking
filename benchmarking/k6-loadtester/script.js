import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  vus: 500,  // Number of virtual users
  duration: '30s', // Test duration
};

export default function () {
  let res = http.get('http://192.168.31.157:8082/test.php');

  check(res, {
    'status is 200': (r) => r.status === 200,
    'response time < 500ms': (r) => r.timings.duration < 500,
  });

  sleep(1); // Simulate user wait time
}
