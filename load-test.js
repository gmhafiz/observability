// k6 run k6.js

import http from 'k6/http';
import { sleep } from 'k6';

export const options = {
    vus: 8,
    // duration: '10s',

   stages: [
       { duration: '2m', target: 4 }, // traffic ramp-up from 1 to N users over X minutes.
       { duration: '6m', target: 8 }, // stay at N users for X minute
       { duration: '2m', target: 0 }, // ramp-down to 0 users
   ],
};

export default function () {
    http.get('http://localhost:3080/uuid');
}
