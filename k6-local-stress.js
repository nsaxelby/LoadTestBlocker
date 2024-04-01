import http from 'k6/http';

export const options = {
    discardResponseBodies: true,

    scenarios: {
        contacts: {
            executor: 'ramping-arrival-rate',

            // Start iterations per `timeUnit`
            startRate: 3000,

            // Start `startRate` iterations per minute
            timeUnit: '1s',

            // Pre-allocate necessary VUs.
            preAllocatedVUs: 10000,

            stages: [
                // Start 300 iterations per `timeUnit` for the first minute.
                { target: 6000, duration: '30s' },

                // Linearly ramp-up to starting 600 iterations per `timeUnit` over the following two minutes.
                { target: 9000, duration: '30s' },

                // Continue starting 600 iterations per `timeUnit` for the following four minutes.
                { target: 12000, duration: '30s' },

                // Linearly ramp-down to starting 60 iterations per `timeUnit` over the last two minutes.
                { target: 15000, duration: '30s' },
            ],
        },
    },
};

export default function () {
    http.get('http://localhost:9000/json');
}