import http from 'k6/http';

export const options = {
    discardResponseBodies: true,
    scenarios: {
        contacts: {
            executor: 'constant-arrival-rate',

            // How long the test lasts
            duration: '6000s',

            // How many iterations per timeUnit
            rate: 5000,

            // Start `rate` iterations per second
            timeUnit: '1s',

            // Pre-allocate VUs
            preAllocatedVUs: 5000,
        },
    },
};

export default function () {
    http.get('http://localhost:9000/json');
}
