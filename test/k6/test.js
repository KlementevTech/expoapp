import grpc from 'k6/net/grpc';
import { check, sleep } from 'k6';

const GRPC_ADDRESS = __ENV.GRPC_ADDRESS || "expo-service:50051"
const DURATION = __ENV.DURATION || "1m"
const VUS = __ENV.VUS || "50"

const client = new grpc.Client();
client.load(['/proto'], 'expo/v1/expo.proto');

export const options = {
    stages: [
        { duration: '30s', target: 10 },
        { duration: DURATION, target: VUS },
        { duration: '30s', target: 0 },
    ],
    thresholds: {
        'checks': ['rate>0.9'],
        'grpc_req_duration': ['p(95)<500'],
    },
};

export default () => {
    client.connect(GRPC_ADDRESS, {
        plaintext: true
    });

    const response = client.invoke('expo.v1.ExpoService/GetVersion', {});

    check(response, {
        'status is OK': (r) => r && r.status === grpc.StatusOK,
    });

    client.close();
    sleep(1);
};

